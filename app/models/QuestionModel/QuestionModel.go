package QuestionModel

import (
	"errors"
	"fmt"
	"time"

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
	AutherName       string
	PublishStartTime time.Time
	IsSolve          bool
	SolvesCount      int
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

	query := fmt.Sprintf("SELECT id, name, flag, score, sentence, genre, publish_start_time  FROM %s WHERE publish_start_time < NOW()", m.Table)
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

func (m *QuestionModel) List(teamName string) ([]Question, error) {
	m.Open()
	defer m.Close()

	var questions []Question

	query := fmt.Sprintf(`
		SELECT
			SUM((CASE WHEN Q2.answer_team = ? THEN TRUE
				ELSE  FALSE END)) AS is_solve,
			Q2.question_id,
			Q2.question_name,
			Q2.question_score,
			Q2.question_genre,
			COUNT(Q2.answer_team) AS solves_count,
			Q2.question_author
		FROM( SELECT DISTINCT
					question.id AS question_id,
					question.name AS question_name,
					question.score AS question_score,
					question.genre AS question_genre,
					team.name AS answer_team,
					author.name AS question_author
				FROM %s
				LEFT JOIN answer ON question.id = answer.question_id AND question.flag = answer.flag
				LEFT JOIN user ON user.id = answer.user_id
				LEFT JOIN team ON team.id = user.team_id
				LEFT JOIN user author ON author.id = question.author_id
				WHERE publish_start_time < NOW()
			) Q2
		GROUP BY Q2.question_id
	`, m.Table)
	rows, err := m.Connection.Query(query, teamName)
	if err != nil {
		return nil, errors.New("Database query error")
	}
	for rows.Next() {
		var question Question
		if err := rows.Scan(&question.IsSolve, &question.Id, &question.Name, &question.Score, &question.Genre, &question.SolvesCount, &question.AutherName); err != nil {
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
	stmtOut, err := m.Connection.Prepare(fmt.Sprintf("SELECT id, name, flag, score, sentence, genre, publish_start_time FROM %s WHERE id = ?", m.Table))
	if err != nil {
		return question, errors.New("Database query error")
	}
	if stmtOut.QueryRow(id).Scan(&question.Id, &question.Name, &question.Flag, &question.Score,
		&question.Sentence, &question.Genre, &question.PublishStartTime) != nil {
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
			INSERT INTO %s (name, flag, genre, score, sentence, author_id) VALUES (?, ?, ?, ?, ?, ?)`, m.Table)
	} else {
		query = fmt.Sprintf(`
			INSERT INTO %s (name, flag, genre, score, publish_start_time, sentence, author_id) VALUES (?, ?, ?, ?, ?, ?, ?)`, m.Table)
	}

	if args["publish_now"] == "on" {
		_, err := m.Connection.Exec(query, args["name"], args["flag"], args["genre"], args["score"], args["sentence"], args["auther_id"])
		if err != nil {
			fmt.Println(err)
			return errors.New("Database error")
		}
	} else {
		_, err := m.Connection.Exec(query, args["name"], args["flag"], args["genre"], args["score"], args["publish_start_time"], args["sentence"], args["auther_id"])
		if err != nil {
			fmt.Println(err)
			return errors.New("Database error")
		}
	}
	return nil
}

func (m *QuestionModel) Update(questionId int, args map[string]string) error {
	m.Open()
	defer m.Close()

	var query string
	if args["publish_now"] == "on" {
		query = fmt.Sprintf("UPDATE %s SET name = ?, flag = ?, score = ?, genre = ?, sentence = ? WHERE id = ?", m.Table)
	} else {
		query = fmt.Sprintf("UPDATE %s SET name = ?, flag = ?, score = ?, genre = ?, publish_start_time = ?, sentence = ? WHERE id = ?", m.Table)
	}
	stmtOut, err := m.Connection.Prepare(query)
	if err != nil {
		return errors.New("Database : query error")
	}
	if args["publish_now"] == "on" {
		if stmtOut.QueryRow(args["name"], args["flag"], args["score"], args["genre"], args["sentence"], questionId) == nil {
			fmt.Println(err)
			return errors.New("Database error")
		}
	} else {
		if stmtOut.QueryRow(args["name"], args["flag"], args["score"], args["genre"], args["publish_start_time"], args["sentence"], questionId) == nil {
			fmt.Println(err)
			return errors.New("Database error")
		}
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

func (m *QuestionModel) GetScore(questionID int) (int, error) {
	m.Open()
	defer m.Close()
	var score int
	query := fmt.Sprintf("SELECT score FROM %s WHERE id = ?", m.Table)
	stmtOut, err := m.Connection.Prepare(query)
	if err != nil {
		return 0, errors.New("Database query error")
	}
	if stmtOut.QueryRow(questionID).Scan(&score) != nil {
		return 0, errors.New("Database error")
	}
	return score, nil
}
