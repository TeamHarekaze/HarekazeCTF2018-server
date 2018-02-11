package controllers

import (
	"fmt"

	"github.com/TeamHarekaze/HarekazeCTF2018-server/cachemodels/RankingCache"
	"github.com/TeamHarekaze/HarekazeCTF2018-server/cachemodels/SolveCache"
	"github.com/TeamHarekaze/HarekazeCTF2018-server/datamodels/AnswerModel"
	"github.com/TeamHarekaze/HarekazeCTF2018-server/datamodels/QuestionModel"
	"github.com/TeamHarekaze/HarekazeCTF2018-server/web/controllers/BaseController"
	"github.com/kataras/iris/context"
	"github.com/kataras/iris/mvc"
)

type AnswerController struct {
	BaseController.Base
}

func (c *AnswerController) answerViewTemplete(qID interface{}, msgs ...string) mvc.View {
	questionModel := QuestionModel.New()
	question, err := questionModel.FindId(qID.(int))
	if err != nil {
		c.Error(err)
	}
	if len(msgs) != 2 {
		return mvc.View{
			Name: "answer/index.html",
			Data: context.Map{
				"Title":         question.Name,
				"Sentence":      question.Sentence,
				"Genre":         question.Genre,
				"Score":         question.Score,
				"IsSubmitBlock": false,
				"Token":         c.MakeToken(fmt.Sprintf("/answer/%d", qID.(int))),
				"IsLoggedIn":    c.IsLoggedIn(),
				"CurrentPage":   "questions",
			},
		}
	} else {
		isSubmitBlock := false
		if msgs[0] == "success" {
			isSubmitBlock = true
		}
		return mvc.View{
			Name: "answer/index.html",
			Data: context.Map{
				"Title":         question.Name,
				"Sentence":      question.Sentence,
				"Genre":         question.Genre,
				"Score":         question.Score,
				"IsSubmitBlock": isSubmitBlock,
				"Message":       msgs[1],
				"MessageType":   msgs[0],
				"Token":         c.MakeToken(fmt.Sprintf("/answer/%d", qID.(int))),
				"IsLoggedIn":    c.IsLoggedIn(),
				"CurrentPage":   "questions",
			},
		}
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
		return c.answerViewTemplete(questionId, "danger", "The competition end.")
	}

	questionModel := QuestionModel.New()
	isExit, err := questionModel.ExitByID(questionId)
	if err != nil {
		return c.Error(err)
	}
	if isExit == false {
		return c.Error("Can not found question.", 404)
	}

	answerModel := AnswerModel.New()
	isCorrected, err := answerModel.IsCorrected(questionId, c.GetLoggedTeamID())
	if err != nil {
		return c.Error(err)
	}
	if isCorrected {
		return c.answerViewTemplete(questionId, "success", "Question already solved")
	} else {
		return c.answerViewTemplete(questionId)
	}
}

// PostBy handles GET: http://localhost:8080/answer/<quesion id>.
func (c *AnswerController) PostBy(questionId int) mvc.Result {
	if !c.IsLoggedIn() {
		c.SetRedirectPath(fmt.Sprintf("/answer/%d", questionId))
		return mvc.Response{Path: "/user/login"}
	} else if c.IsBeforeCompetition() {
		return c.Error("Competition is not start.", 404)
	} else if !c.IsNowCompetition() {
		return c.answerViewTemplete(questionId, "danger", "The competition end.")
	}
	questionModel := QuestionModel.New()
	isExit, err := questionModel.ExitByID(questionId)
	if err != nil {
		return c.Error(err)
	}
	if isExit == false {
		return c.Error("Can not found question.", 404)
	}

	var (
		flag  = c.Ctx.FormValue("flag")
		token = c.Ctx.FormValue("csrf_token")
	)
	if !c.CheckTaken(token, fmt.Sprintf("/answer/%d", questionId)) {
		c.answerViewTemplete(questionId, "danger", "Token error.")
	}

	message := ""
	messageType := ""

	answerModel := AnswerModel.New()
	isCorrected, err := answerModel.IsCorrected(questionId, c.GetLoggedTeamID())
	if err != nil {
		return c.Error(err)
	}
	if isCorrected == false {
		isCorrect, err := answerModel.CheckFlag(questionId, flag)
		if err != nil {
			return c.Error(err)
		}
		if err := answerModel.Insert(questionId, c.GetLoggedUserID(), flag); err != nil {
			return c.Error(err)
		}
		message = "Incorrect answer"
		messageType = "danger"
		if isCorrect {
			rankingCache := RankingCache.New()
			defer rankingCache.Close()
			err := rankingCache.Set(c.GetLoggedTeamName(), questionId)
			if err != nil {
				return c.Error(err)
			}
			solveCache := SolveCache.New()
			defer solveCache.Close()
			err = solveCache.Set(questionId, c.GetLoggedTeamName())
			if err != nil {
				return c.Error(err)
			}
			message = "Correct answer"
			messageType = "success"
		}
	} else {
		message = "Question already solved"
		messageType = "success"
	}
	return c.answerViewTemplete(questionId, messageType, message)
}
