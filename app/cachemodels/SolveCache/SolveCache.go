package SolveCache

import (
	"io/ioutil"
	"reflect"

	"github.com/TeamHarekaze/HarekazeCTF2018-server/app/cachemodels/BaseCache"
	"github.com/TeamHarekaze/HarekazeCTF2018-server/app/datamodels/AnswerModel"
	"github.com/TeamHarekaze/HarekazeCTF2018-server/app/datamodels/QuestionModel"
	"github.com/go-redis/redis"
)

func contains(s []string, e string) bool {
	for _, v := range s {
		if e == v {
			return true
		}
	}
	return false
}

type Cache struct {
	BaseCache.Base
}

func New() *Cache {
	cache := new(Cache)
	cache.Client = BaseCache.NewClient(2)
	return cache
}

func (c *Cache) setData() error {

	answerModel := AnswerModel.New()
	solves, err := answerModel.GetSolve()
	if err != nil {
		return err
	}
	for i := range solves {
		err := c.Client.LPush(solves[i].QuestionName, solves[i].Team).Err()
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Cache) Set(questionID int, teamName string) error {
	questionModel := QuestionModel.New()
	question, err := questionModel.FindId(questionID)
	if err != nil {
		return err
	}
	questionName := question.Name
	err = c.Client.LPush(questionName, teamName).Err()
	if err != nil {
		return err
	}
	return nil
}

func (c *Cache) List(teamName string) ([]QuestionModel.Question, error) {
	r, err := c.IsNoSet()
	if err != nil {
		return nil, err
	}
	if r == true {
		err := c.setData()
		if err != nil {
			return nil, err
		}
	}

	// get question data
	questionModel := QuestionModel.New()
	questions, err := questionModel.FindAllEnable()
	if err != nil {
		return nil, err
	}

	// get solve data
	data, err := ioutil.ReadFile(`./getSolveCacheAllData.lua`)
	if err != nil {
		return nil, err
	}
	IncrByXX := redis.NewScript(string(data))
	row, err := IncrByXX.Run(c.Client, []string{}).Result()
	if err != nil {
		return nil, err
	}

	//cast
	solveData := make(map[string][]string) //---
	rowValue := reflect.ValueOf(row)
	qns := reflect.ValueOf(rowValue.Index(0).Interface())
	tss := reflect.ValueOf(rowValue.Index(1).Interface())
	for i := 0; i < qns.Len(); i = i + 1 {
		qn := qns.Index(i).Interface().(string)
		var solveTeamName []string
		ts := reflect.ValueOf(tss.Index(i).Interface())
		for j := 0; j < ts.Len(); j = j + 1 {
			t := ts.Index(j).Interface().(string)
			solveTeamName = append(solveTeamName, t)
		}
		solveData[qn] = solveTeamName
	}

	// Join question data to solve data
	for i := range questions {
		var solvesCount int
		_, ok := solveData[questions[i].Name]
		if ok == true {
			solvesCount = len(solveData[questions[i].Name])
		}
		questions[i].SolvesCount = solvesCount
		questions[i].IsSolve = contains(solveData[questions[i].Name], teamName)
	}

	return questions, nil
}
