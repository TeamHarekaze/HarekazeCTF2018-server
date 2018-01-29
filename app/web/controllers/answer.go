package controllers

import (
	"errors"
	"fmt"

	"github.com/HayatoDoi/HarekazeCTF-Competition/app/datamodels/AnswerModel"
	"github.com/HayatoDoi/HarekazeCTF-Competition/app/datamodels/QuestionModel"
	"github.com/HayatoDoi/HarekazeCTF-Competition/app/redismodels/RankingCache"
	"github.com/HayatoDoi/HarekazeCTF-Competition/app/redismodels/SolveCache"
	"github.com/HayatoDoi/HarekazeCTF-Competition/app/web/controllers/BaseController"
	"github.com/kataras/iris/context"
	"github.com/kataras/iris/mvc"
)

type AnswerController struct {
	BaseController.Base
}

func (c *AnswerController) answerViewTemplete(data context.Map) mvc.View {
	return mvc.View{
		Name: "answer/index.html",
		Data: data,
	}
}

// GetBy handles GET: http://localhost:8080/answer/<quesion id>.
func (c *AnswerController) GetBy(questionId int) mvc.Result {
	if !c.IsLoggedIn() {
		c.SetRedirectPath(fmt.Sprintf("/answer/%d", questionId))
		return mvc.Response{Path: "/user/login"}
	} else if c.IsBeforeCompetition() {
		return mvc.Response{Code: 404}
	} else if !c.IsNowCompetition() {
		questionModel := QuestionModel.New()
		question, _ := questionModel.FindId(questionId)
		return c.answerViewTemplete(context.Map{
			"Title":         question.Name,
			"Sentence":      question.Sentence,
			"Genre":         question.Genre,
			"Score":         question.Score,
			"IsSubmitBlock": false,
			"Message":       "The competition end.",
			"MessageType":   "danger",
			"Token":         c.MakeToken(fmt.Sprintf("/answer/%d", questionId)),
			"IsLoggedIn":    c.IsLoggedIn(),
		})
	}

	questionModel := QuestionModel.New()
	question, err := questionModel.FindId(questionId)
	if err != nil {
		return mvc.Response{Err: err}
	}
	message := ""
	messageType := ""
	answerModel := AnswerModel.New()
	isCorrected, err := answerModel.IsCorrected(questionId, c.GetLoggedTeamID())
	if err != nil {
		return mvc.Response{Err: err}
	}
	if isCorrected {
		message = "Question already solved"
		messageType = "success"
	}

	return c.answerViewTemplete(context.Map{
		"Title":         question.Name,
		"Sentence":      question.Sentence,
		"Genre":         question.Genre,
		"Score":         question.Score,
		"IsSubmitBlock": isCorrected,
		"Message":       message,
		"MessageType":   messageType,
		"Token":         c.MakeToken(fmt.Sprintf("/answer/%d", questionId)),
		"IsLoggedIn":    c.IsLoggedIn(),
	})
}

// PostBy handles GET: http://localhost:8080/answer/<quesion id>.
func (c *AnswerController) PostBy(questionId int) mvc.Result {
	if !c.IsLoggedIn() {
		c.SetRedirectPath(fmt.Sprintf("/answer/%d", questionId))
		return mvc.Response{Path: "/user/login"}
	} else if c.IsBeforeCompetition() {
		return mvc.Response{Code: 404}
	} else if !c.IsNowCompetition() {
		questionModel := QuestionModel.New()
		question, _ := questionModel.FindId(questionId)
		return c.answerViewTemplete(context.Map{
			"Title":         question.Name,
			"Sentence":      question.Sentence,
			"Genre":         question.Genre,
			"Score":         question.Score,
			"IsSubmitBlock": true,
			"Message":       "The competition end.",
			"MessageType":   "danger",
			"Token":         c.MakeToken(fmt.Sprintf("/answer/%d", questionId)),
			"IsLoggedIn":    c.IsLoggedIn(),
		})
	}
	var (
		flag  = c.Ctx.FormValue("flag")
		token = c.Ctx.FormValue("csrf_token")
	)
	if !c.CheckTaken(token, fmt.Sprintf("/answer/%d", questionId)) {
		return mvc.Response{Err: errors.New("token error"), Code: 400}
	}

	message := ""
	messageType := ""

	answerModel := AnswerModel.New()
	isCorrected, err := answerModel.IsCorrected(questionId, c.GetLoggedTeamID())
	if err != nil {
		return mvc.Response{Err: err}
	}
	if isCorrected == false {
		isCorrect, err := answerModel.CheckFlag(questionId, flag)
		if err != nil {
			return mvc.Response{Err: err}
		}
		if answerModel.Insert(questionId, c.GetLoggedUserID(), flag) != nil {
			return mvc.Response{Err: errors.New("db error")}
		}
		message = "Incorrect answer"
		messageType = "danger"
		if isCorrect {
			rankingCache := RankingCache.New()
			defer rankingCache.Close()
			err := rankingCache.Set(c.GetLoggedTeamName(), questionId)
			if err != nil {
				return mvc.Response{Err: err}
			}
			solveCache := SolveCache.New()
			defer solveCache.Close()
			err = solveCache.Set(questionId, c.GetLoggedTeamName())
			if err != nil {
				return mvc.Response{Err: err}
			}
			message = "Correct answer"
			messageType = "success"
		}
	} else {
		message = "Question already solved"
		messageType = "success"
	}
	questionModel := QuestionModel.New()
	question, _ := questionModel.FindId(questionId)
	return c.answerViewTemplete(context.Map{
		"Title":         question.Name,
		"Sentence":      question.Sentence,
		"Genre":         question.Genre,
		"Score":         question.Score,
		"IsSubmitBlock": isCorrected,
		"Message":       message,
		"MessageType":   messageType,
		"Token":         c.MakeToken(fmt.Sprintf("/answer/%d", questionId)),
		"IsLoggedIn":    c.IsLoggedIn(),
	})
}
