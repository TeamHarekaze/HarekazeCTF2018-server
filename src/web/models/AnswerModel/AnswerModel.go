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

type Rank struct {
	Rank  int
	Name  string
	Score int
}

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

func (m *AnswerModel) IsCorrected(questionId int, teamId interface{}) (bool, error) {
	m.Open()
	defer m.Close()

	query := fmt.Sprintf(`
		SELECT COUNT(*) FROM answer
			INNER JOIN question ON question.flag = answer.flag
			INNER JOIN user ON user.id = answer.user_id
			WHERE answer.question_id = ? AND user.team_id = ?
	`)
	stmtOut, err := m.Connection.Prepare(query)
	if err != nil {
		return false, errors.New("Database query error")
	}
	var count int
	if stmtOut.QueryRow(questionId, teamId.(string)).Scan(&count) != nil {
		return false, errors.New("Database error")
	}
	return count > 0, nil
}

func (m *AnswerModel) Insert(questionId int, userId interface{}, flag string) error {
	m.Open()
	defer m.Close()

	query := fmt.Sprintf("INSERT INTO %s (user_id, question_id, flag) VALUES( ?, ?, ?)", m.Table)
	stmtOut, err := m.Connection.Prepare(query)
	if err != nil {
		return errors.New("Database : query error")
	}
	if stmtOut.QueryRow(userId, questionId, flag) == nil {
		return errors.New("Database error")
	}
	return nil
}
func (m *AnswerModel) Ranking() ([]Rank, error) {
	m.Open()
	defer m.Close()

	var rs []Rank
	query := fmt.Sprintf(`
		SELECT team.name as team_name, IFNULL(SUM(answer.score), 0) as score, IFNULL(MAX(update_time), '') as update_time
			FROM team
			LEFT JOIN user
				ON user.team_id = team.id
			LEFT JOIN (
				SELECT DISTINCT answer.user_id, question.score as score, answer.create_time as update_time, NULL
					FROM answer
						INNER JOIN question
						ON question.id = answer.question_id
					UNION (
						SELECT answer.user_id, 10, NULL, MIN(answer.create_time)
							FROM answer
							INNER JOIN question
								ON question.id = answer.question_id
								GROUP BY question.id, answer.user_id
					)
				) answer
				ON answer.user_id = user.id
			GROUP BY team.id
			ORDER BY score DESC, update_time
	`)
	rows, err := m.Connection.Query(query)
	if err != nil {
		return nil, errors.New("Database query error")
	}
	count := 0
	for rows.Next() {
		var r Rank
		var tmp string
		if rows.Scan(&r.Name, &r.Score, &tmp) != nil {
			return rs, errors.New("Database error")
		}
		count = count + 1
		r.Rank = count
		rs = append(rs, r)
	}
	return rs, nil
}
