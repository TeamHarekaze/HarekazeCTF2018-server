package AnswerModel

import (
	"time"
	// "fmt"

	"errors"
	"fmt"

	"../BaseModel"
)

const (
	table      = "answer"
	primarykey = "id"
)

type Solve struct {
	QuestionName string
	Team         string
}

type Rank struct {
	Rank       int
	Name       string
	Score      int
	UpdateTime time.Time
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
	_, err := m.Connection.Exec(query, userId, questionId, flag)
	if err != nil {
		return errors.New("Database : query error")
	}
	return nil
}
func (m *AnswerModel) Ranking() ([]Rank, error) {
	m.Open()
	defer m.Close()

	var rs []Rank
	query := fmt.Sprintf(`
		SELECT team_name,
				SUM(CASE WHEN score_table.create_time = (SELECT MIN(answer.create_time) FROM answer
						INNER JOIN question
						ON question.id = answer.question_id AND question.flag = answer.flag
						WHERE question.id = score_table.question_id )
				THEN score+10
				ELSE score
				END ) AS score_sum,
				MAX(create_time) AS update_time
			FROM (
				SELECT team.name AS team_name, question.id AS question_id, IFNULL(question.score, 0) AS score, MIN(answer.create_time) AS create_time
					FROM team
					RIGHT JOIN user
						ON user.team_id = team.id
					RIGHT JOIN answer
						ON answer.user_id = user.id
					INNER JOIN question
						ON question.id = answer.question_id AND question.flag = answer.flag
					GROUP BY team.name, question.id
				) score_table
			GROUP BY score_table.team_name
			ORDER BY score_sum DESC, MAX(create_time)
	`)
	rows, err := m.Connection.Query(query)
	if err != nil {
		return nil, errors.New("Database query error")
	}
	count := 0
	for rows.Next() {
		var r Rank
		if rows.Scan(&r.Name, &r.Score, &r.UpdateTime) != nil {
			return rs, errors.New("Database error")
		}
		count = count + 1
		r.Rank = count
		rs = append(rs, r)
	}
	return rs, nil
}

func (m *AnswerModel) IsFast(questionID int) (bool, error) {
	m.Open()
	defer m.Close()
	var count int
	query := fmt.Sprintf("SELECT COUNT(%s.id) FROM %s LEFT JOIN question ON answer.question_id = question.id AND answer.flag = question.flag WHERE question.id = ?", m.Table, m.Table)
	stmtOut, err := m.Connection.Prepare(query)
	if err != nil {
		return false, errors.New("Database query error")
	}
	if stmtOut.QueryRow(questionID).Scan(&count) != nil {
		return false, errors.New("Database error")
	}
	return count == 1, nil
}

func (m *AnswerModel) GetSolve() ([]Solve, error) {
	m.Open()
	defer m.Close()

	var solves []Solve
	query := fmt.Sprintf(`
		SELECT question.name, team.name FROM answer
		INNER JOIN question ON question.id = answer.question_id AND question.flag = answer.flag
		LEFT JOIN user ON user.id = answer.user_id
		LEFT JOIN team ON team.id = user.team_id
	`)
	rows, err := m.Connection.Query(query)
	if err != nil {
		return solves, err
	}
	for rows.Next() {
		var s Solve
		if rows.Scan(&s.QuestionName, &s.Team) != nil {
			return solves, err
		}
		solves = append(solves, s)
	}
	return solves, nil
}
