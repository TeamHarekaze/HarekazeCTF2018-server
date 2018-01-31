package RankingCache

import (
	"io/ioutil"
	"reflect"
	"sort"
	"strconv"
	"time"

	"github.com/HayatoDoi/HarekazeCTF-Competition/app/datamodels/AnswerModel"
	"github.com/HayatoDoi/HarekazeCTF-Competition/app/datamodels/QuestionModel"
	"github.com/HayatoDoi/HarekazeCTF-Competition/app/redismodels/BaseCache"
	"github.com/go-redis/redis"
)

// struct for rank
type Rank struct {
	Rank       int
	Name       string
	Score      int
	UpdateTime int
}

// 構造体のスライス
type R []Rank

func (r R) Len() int {
	return len(r)
}

func (r R) Swap(i int, j int) {
	r[i], r[j] = r[j], r[i]
}

func (r R) Less(i, j int) bool {
	if r[i].Score != r[j].Score {
		return r[i].Score < r[j].Score
	} else {
		return r[i].UpdateTime > r[j].UpdateTime
	}
}

type Cache struct {
	BaseCache.Base
}

func New() *Cache {
	cache := new(Cache)
	cache.Client = BaseCache.NewClient(0)
	return cache
}

func (c *Cache) setData() error {
	answerModel := AnswerModel.New()
	rank, err := answerModel.Ranking()
	if err != nil {
		return err
	}
	for i := 0; i < len(rank); i = i + 1 {
		r := rank[i]
		name := r.Name
		score := r.Score
		updateTime := r.UpdateTime.Unix()
		err := c.Client.HSet(name, "score", score).Err()
		if err != nil {
			return err
		}
		err = c.Client.HSet(name, "update_time", updateTime).Err()
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Cache) Set(teamName string, questionID int) error {
	r, err := c.IsNoSet()
	if err != nil {
		return err
	}
	if r == true {
		err := c.setData()
		if err != nil {
			return err
		}
		return nil
	}
	bonus := 0

	answerModel := AnswerModel.New()
	isFast, err := answerModel.IsFast(questionID)
	if err != nil {
		return err
	}
	if isFast == true {
		bonus = 10
	}

	questionModel := QuestionModel.New()
	score, err := questionModel.GetScore(questionID)
	if err != nil {
		return err
	}

	val, err := c.Client.HGet(teamName, "score").Result()
	if err != nil && err != redis.Nil {
		return err
	}
	if err == redis.Nil {
		err := c.Client.HSet(teamName, "score", score+bonus).Err()
		if err != nil {
			return err
		}
		err = c.Client.HSet(teamName, "update_time", time.Now().Unix()).Err()
		if err != nil {
			return err
		}
	} else {
		old_score, _ := strconv.Atoi(val)
		err := c.Client.HSet(teamName, "score", old_score+score+bonus).Err()
		if err != nil {
			return err
		}
		err = c.Client.HSet(teamName, "update_time", time.Now().Unix()).Err()
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Cache) Rank() (R, error) {
	var rank R
	r, err := c.IsNoSet()
	if err != nil {
		return rank, err
	}
	if r == true {
		err := c.setData()
		if err != nil {
			return rank, err
		}
	}
	// get data
	data, err := ioutil.ReadFile(`./getRankingCacheAllData.lua`)
	if err != nil {
		return rank, err
	}
	IncrByXX := redis.NewScript(string(data))
	dataArrayInterface, err := IncrByXX.Run(c.Client, []string{}).Result()
	if err != nil {
		return rank, err
	}

	dataArrayValue := reflect.ValueOf(dataArrayInterface)
	for i := 0; i < dataArrayValue.Len(); i = i + 1 {
		var r Rank
		dataInterface := dataArrayValue.Index(i).Interface()
		dataValue := reflect.ValueOf(dataInterface)
		r.Name = dataValue.Index(0).Interface().(string)
		scoreStr := dataValue.Index(1).Interface().(string)
		r.Score, _ = strconv.Atoi(scoreStr)
		updateTimeStr := dataValue.Index(2).Interface().(string)
		r.UpdateTime, _ = strconv.Atoi(updateTimeStr)
		rank = append(rank, r)
	}

	// sort
	sort.Sort(sort.Reverse(rank))
	// add Rank
	for i := 0; i < len(rank); i = i + 1 {
		rank[i].Rank = i + 1
	}
	return rank, nil
}
