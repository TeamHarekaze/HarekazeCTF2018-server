package QuestionModel

import (
	// "fmt"

	"errors"
	"fmt"

	"../BaseModel"
)

const (
	table      = "question"
	primarykey = "id"
)

type Question struct {
	Id               int
	Name             string
	Flag             string
	Score            int
	Sentence         string
	Genre            string
	PublishStartTime string
}

type QuestionModel struct {
	BaseModel.Base
}

func New() *QuestionModel {
	base := new(QuestionModel)
	base.Table = table
	base.Primarykey = primarykey
	return base
}
func (m *QuestionModel) FindAll() ([]Question, error) {
	m.Open()
	defer m.Close()

	var questions []Question

	query := fmt.Sprintf("SELECT id, name, flag, score, sentence, genre, publish_start_time  FROM %s", m.Table)
	rows, err := m.Connection.Query(query)
	if err != nil {
		return nil, errors.New("Database query error")
	}
	for rows.Next() {
		var question Question
		if err := rows.Scan(&question.Id, &question.Name, &question.Flag, &question.Score,
			&question.Sentence, &question.Genre, &question.PublishStartTime); err != nil {
			return questions, err
		}
		questions = append(questions, question)
	}
	return questions, nil
}
func (m *QuestionModel) FindAllEnable() ([]Question, error) {
	m.Open()
	defer m.Close()

	var questions []Question

	query := fmt.Sprintf("SELECT id, name FROM %s WHERE publish_start_time < NOW()", m.Table)
	rows, err := m.Connection.Query(query)
	if err != nil {
		return nil, errors.New("Database error")
	}
	for rows.Next() {
		var question Question
		if err := rows.Scan(&question.Id, &question.Name); err != nil {
			return questions, errors.New("Database error")
		}
		questions = append(questions, question)
	}
	return questions, nil
}

func (m *QuestionModel) FindId(id int) (Question, error) {
	m.Open()
	defer m.Close()

	var question Question
	stmtOut, err := m.Connection.Prepare(fmt.Sprintf("SELECT id, name, flag, score, sentence FROM %s WHERE id = ?", m.Table))
	if err != nil {
		return question, errors.New("Database query error")
	}
	if stmtOut.QueryRow(id).Scan(&question.Id, &question.Name, &question.Flag, &question.Score, &question.Sentence) != nil {
		return question, errors.New("Database error")
	}
	return question, nil
}

func (m *QuestionModel) Save(args map[string]string) error {
	m.Open()
	defer m.Close()

	var query string
	if args["publish_now"] == "on" {
		query = fmt.Sprintf(`
			INSERT INTO %s (name, flag, genre, score, sentence, author_id)
				SELECT ?, ?, ?, ?, ?, id FROM user WHERE name = ?`, m.Table)
	} else {
		query = fmt.Sprintf(`
		INSERT INTO %s (name, flag, genre, score, publish_start_time, sentence, author_id)
			SELECT ?, ?, ?, ?, ?, ?, id FROM user WHERE name = ?`, m.Table)
	}
	stmtOut, err := m.Connection.Prepare(query)
	if err != nil {
		return errors.New("Database : query error")
	}
	if args["publish_now"] == "on" {
		if stmtOut.QueryRow(args["name"], args["flag"], args["genre"], args["score"], args["sentence"], args["auther_name"]) == nil {
			fmt.Println(err)
			return errors.New("Database error")
		}
	} else {
		if stmtOut.QueryRow(args["name"], args["flag"], args["genre"], args["score"], args["publish_start_time"], args["sentence"], args["auther_name"]) == nil {
			fmt.Println(err)
			return errors.New("Database error")
		}
	}
	return nil
}

func (m *QuestionModel) Update(questionId int, name string, flag string, score string, sentence string) error {
	m.Open()
	defer m.Close()

	query := fmt.Sprintf("UPDATE %s SET name = ?, flag = ?, score = ?, sentence = ? WHERE id = ?", m.Table)
	stmtOut, err := m.Connection.Prepare(query)
	if err != nil {
		return errors.New("Database : query error")
	}
	if stmtOut.QueryRow(name, flag, score, sentence, questionId) == nil {
		fmt.Println(err)
		return errors.New("Database error")
	}
	return nil
}

func (m *QuestionModel) Delete(questionId int) error {
	m.Open()
	defer m.Close()

	query := fmt.Sprintf("DELETE FROM %s WHERE id = ?", m.Table)
	stmtOut, err := m.Connection.Prepare(query)
	if err != nil {
		return errors.New("Database : query error")
	}
	if stmtOut.QueryRow(questionId) == nil {
		fmt.Println(err)
		return errors.New("Database error")
	}
	return nil
}
