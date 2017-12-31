package controllers

import (
	"errors"

	"./BaseController"

	"../models/AnswerModel"
	"../models/QuestionModel"
	"github.com/kataras/iris/context"
	"github.com/kataras/iris/mvc"
)

type AnswerController struct {
	BaseController.Base
}

// GetBy handles GET: http://localhost:8080/answer/<quesion id>.
func (c *AnswerController) GetBy(questionId int) mvc.Result {
	if !c.IsLoggedIn() {
		return mvc.Response{Path: "/user/login"}
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

	return mvc.View{
		Name: "answer/index.html",
		Data: context.Map{
			"Title":       question.Name,
			"Sentence":    question.Sentence,
			"IsCorrected": isCorrected,
			"Message":     message,
			"MessageType": messageType,
			"Token":       c.MakeToken(),
		},
	}
}

// PostBy handles GET: http://localhost:8080/answer/<quesion id>.
func (c *AnswerController) PostBy(questionId int) mvc.Result {
	if !c.IsLoggedIn() {
		return mvc.Response{Path: "/user/login"}
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
	isCorrected, err := answerModel.IsCorrected(questionId, c.GetLoggedUserName())
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
			message = "Correct answer"
			messageType = "success"
		}
	}
	questionModel := QuestionModel.New()
	question, _ := questionModel.FindId(questionId)
	return mvc.View{
		Name: "answer/index.html",
		Data: context.Map{
			"Title":       question.Name,
			"Sentence":    question.Sentence,
			"IsCorrected": isCorrected,
			"Message":     message,
			"MessageType": messageType,
			"Token":       c.MakeToken(),
		},
	}
}
