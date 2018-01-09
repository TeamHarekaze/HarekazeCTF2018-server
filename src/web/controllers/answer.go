package controllers

import (
	"errors"

	"./BaseController"

	"../../redisClient/RankingCache"
	"../models/AnswerModel"
	"../models/QuestionModel"
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
		return mvc.Response{Path: "/user/login"}
	} else if c.IsBeforeCompetition() {
		return mvc.Response{Code: 404}
	} else if !c.IsNowCompetition() {
		questionModel := QuestionModel.New()
		question, _ := questionModel.FindId(questionId)
		return c.answerViewTemplete(context.Map{
			"Title":         question.Name,
			"Sentence":      question.Sentence,
			"IsSubmitBlock": false,
			"Message":       "The competition end.",
			"MessageType":   "danger",
			"Token":         c.MakeToken(),
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
		message = "Corrected answer"
		messageType = "success"
	}

	return c.answerViewTemplete(context.Map{
		"Title":         question.Name,
		"Sentence":      question.Sentence,
		"IsSubmitBlock": isCorrected,
		"Message":       message,
		"MessageType":   messageType,
		"Token":         c.MakeToken(),
	})
}

// PostBy handles GET: http://localhost:8080/answer/<quesion id>.
func (c *AnswerController) PostBy(questionId int) mvc.Result {
	if !c.IsLoggedIn() {
		return mvc.Response{Path: "/user/login"}
	} else if c.IsBeforeCompetition() {
		return mvc.Response{Code: 404}
	} else if !c.IsNowCompetition() {
		questionModel := QuestionModel.New()
		question, _ := questionModel.FindId(questionId)
		return c.answerViewTemplete(context.Map{
			"Title":         question.Name,
			"Sentence":      question.Sentence,
			"IsSubmitBlock": true,
			"Message":       "The competition end.",
			"MessageType":   "danger",
			"Token":         c.MakeToken(),
		})
	}
	var (
		flag  = c.Ctx.FormValue("flag")
		token = c.Ctx.FormValue("csrf_token")
	)
	if !c.CheckTaken(token) {
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
			err := rankingCache.Set(c.GetLoggedTeamName(), questionId)
			if err != nil {
				return mvc.Response{Err: err}
			}
			message = "Correct answer"
			messageType = "success"
		}
	} else {
		message = "Corrected answer"
		messageType = "success"
	}
	questionModel := QuestionModel.New()
	question, _ := questionModel.FindId(questionId)
	return c.answerViewTemplete(context.Map{
		"Title":         question.Name,
		"Sentence":      question.Sentence,
		"IsSubmitBlock": isCorrected,
		"Message":       message,
		"MessageType":   messageType,
		"Token":         c.MakeToken(),
	})
}
