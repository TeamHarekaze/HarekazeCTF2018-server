package RankingCache

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"time"

	"../../web/models/AnswerModel"
	"../../web/models/QuestionModel"
	"github.com/go-redis/redis"
)

// struct for rank
type Rank struct {
	Rank       int
	Name       string
	Score      int
	UpdateTime time.Time
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
	return r[i].Score < r[j].Score
}

type Cache struct {
	client *redis.Client
}

func New() *Cache {
	cache := new(Cache)
	cache.client = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})
	return cache
}

func (c *Cache) isNoSet() (bool, error) {
	var cursor uint64
	var n int
	for {
		var keys []string
		var err error
		keys, cursor, err = c.client.Scan(cursor, "", 10).Result()
		if err != nil {
			return false, err
		}
		n += len(keys)
		if cursor == 0 {
			break
		}
	}
	return n == 0, nil
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
		// aliveTime := time.Duration(1600000000 - r.UpdateTime.Unix())
		aliveTime := time.Duration(0)
		err := c.client.Set(name, score, aliveTime).Err()
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Cache) Set(teamName string, questionID int) error {
	r, err := c.isNoSet()
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

	val, err := c.client.Get(teamName).Result()
	if err != nil && err != redis.Nil {
		return err
	}
	if err == redis.Nil {
		err := c.client.Set(teamName, score+bonus, 0).Err()
		if err != nil {
			return err
		}
	} else {
		old_score, _ := strconv.Atoi(val)
		err := c.client.Set(teamName, old_score+score+bonus, 0).Err()
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Cache) Rank() (R, error) {
	var rank R
	r, err := c.isNoSet()
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
	var cursor uint64
	for {
		var keys []string
		var err error
		keys, cursor, err = c.client.Scan(cursor, "", 10).Result()
		if err != nil {
			return rank, err
		}
		for _, key := range keys {
			val, err := c.client.Get(key).Result()
			if err != nil {
				return rank, err
			}
			var r Rank
			intVal, _ := strconv.Atoi(val)
			r.Name = key
			r.Score = intVal
			rank = append(rank, r)
		}
		if cursor == 0 {
			break
		}
	}
	// sort
	sort.Sort(sort.Reverse(rank))
	for i := 0; i < len(rank); i = i + 1 {
		rank[i].Rank = i + 1
	}
	return rank, nil
}
