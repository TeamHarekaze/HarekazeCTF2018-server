package AnswerModel

import (
	// "fmt"

	"errors"
	"fmt"

	"../BaseModel"
)

const (
	table      = "answer"
	primarykey = "id"
)

type AnswerModel struct {
	BaseModel.Base
}

func New() *AnswerModel {
	base := new(AnswerModel)
	base.Table = table
	base.Primarykey = primarykey
	return base
}

func (m *AnswerModel) CheckFlag(id int, flag string) (bool, error) {
	m.Open()
	defer m.Close()

	query := fmt.Sprintf("SELECT COUNT(*) FROM question WHERE id = ? AND flag = ?")
	stmtOut, err := m.Connection.Prepare(query)
	if err != nil {
		return false, errors.New("Database query error")
	}
	var count int
	if stmtOut.QueryRow(id, flag).Scan(&count) != nil {
		return false, errors.New("Database error")
	}
	return count > 0, nil
}

func (m *AnswerModel) IsCorrected(id int, username string) (bool, error) {
	m.Open()
	defer m.Close()

	query := fmt.Sprintf(`
		SELECT COUNT(*) FROM answer
		INNER JOIN question ON question.flag = answer.flag
		INNER JOIN user ON user.id = answer.user_id
		WHERE answer.question_id = ? AND
			user.team_id IN (SELECT team_id FROM user WHERE name = ?)
	`)
	stmtOut, err := m.Connection.Prepare(query)
	if err != nil {
		return false, errors.New("Database query error")
	}
	var count int
	if stmtOut.QueryRow(id, username).Scan(&count) != nil {
		return false, errors.New("Database error")
	}
	return count > 0, nil
}

func (m *AnswerModel) Insert(QuestionId int, userName string, flag string) error {
	m.Open()
	defer m.Close()

	query := fmt.Sprintf("INSERT INTO %s (user_id, question_id, flag) SELECT id, ?, ? FROM user WHERE name = ?", m.Table)
	stmtOut, err := m.Connection.Prepare(query)
	if err != nil {
		return errors.New("Database : query error")
	}
	if stmtOut.QueryRow(QuestionId, flag, userName) == nil {
		return errors.New("Database error")
	}
	return nil
}
