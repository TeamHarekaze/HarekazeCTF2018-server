package TeamModel

import (
	"crypto/sha512"
	"errors"
	"fmt"

	"github.com/TeamHarekaze/HarekazeCTF2018-server/app/datamodels/BaseModel"
)

const (
	table      = "team"
	primarykey = "id"
)

func GenerateHashedPassword(userPassword string) string {
	return fmt.Sprintf("%x", sha512.Sum512([]byte(userPassword)))
}

type Team struct {
	Id     int
	Name   string
	Enable bool
}

type TeamModel struct {
	BaseModel.Base
	Team Team
}

func New() *TeamModel {
	base := new(TeamModel)
	base.Table = table
	base.Primarykey = primarykey
	return base
}

func (m *TeamModel) All() ([]Team, error) {
	m.Open()
	defer m.Close()

	var teams []Team

	query := fmt.Sprintf("SELECT id, name, enable FROM %s", m.Table)
	rows, err := m.Connection.Query(query)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var team Team
		if err := rows.Scan(&team.Id, &team.Name, &team.Enable); err != nil {
			return teams, errors.New("Database scan error")
		}
		teams = append(teams, team)
	}
	return teams, nil
}
func (m *TeamModel) AllEnable() ([]Team, error) {
	m.Open()
	defer m.Close()

	var teams []Team

	query := fmt.Sprintf("SELECT id, name, email, enable FROM %s WHERE enable = 1", m.Table)
	rows, err := m.Connection.Query(query)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var team Team
		if err := rows.Scan(&team.Id, &team.Name, &team.Enable); err != nil {
			return teams, err
		}
		teams = append(teams, team)
	}
	return teams, nil
}

func (m *TeamModel) Enable(id int) error {
	m.Open()
	defer m.Close()

	stmtOut, err := m.Connection.Prepare(fmt.Sprintf("UPDATE %s SET enable = 1 WHERE id = ?", m.Table))
	if err != nil {
		return errors.New("Database query error")
	}
	if stmtOut.QueryRow(id) == nil {
		return errors.New("Database error")
	}
	return nil
}

func (m *TeamModel) Disable(id int) error {
	m.Open()
	defer m.Close()

	stmtOut, err := m.Connection.Prepare(fmt.Sprintf("UPDATE %s SET enable = 0 WHERE id = ?", m.Table))
	if err != nil {
		return errors.New("Database query error")
	}
	if stmtOut.QueryRow(id) == nil {
		return errors.New("Database error(stmtOut.QueryRow(id))")
	}
	return nil
}
func (m *TeamModel) PasswordCheck(name string, password string) (bool, error) {
	m.Open()
	defer m.Close()

	hashedPassword := GenerateHashedPassword(password)
	stmtOut, err := m.Connection.Prepare(fmt.Sprintf("SELECT COUNT(name) FROM %s WHERE name = ? AND hashed_password = ?", m.Table))
	if err != nil {
		return false, err
	}

	var count int
	if err := stmtOut.QueryRow(name, hashedPassword).Scan(&count); err != nil {
		return false, err
	}
	return count == 1, nil
}

func (m *TeamModel) UsedChack(name string) (bool, error) {
	m.Open()
	defer m.Close()

	stmtOut, err := m.Connection.Prepare(fmt.Sprintf("SELECT COUNT(name) FROM %s WHERE name = ?", m.Table))
	if err != nil {
		return false, err
	}

	var count int
	if err := stmtOut.QueryRow(name).Scan(&count); err != nil {
		return false, err
	}
	return count != 0, nil
}
func (m *TeamModel) Add(name string, password string) error {
	m.Open()
	defer m.Close()

	hashedPassword := GenerateHashedPassword(password)
	query := fmt.Sprintf("INSERT INTO %s (name, hashed_password) VALUES(?, ?)", m.Table)
	_, err := m.Connection.Exec(query, name, hashedPassword)
	if err != nil {
		return err
	}
	return nil
}

func (m *TeamModel) GetMember(teamID string) ([]string, error) {
	m.Open()
	defer m.Close()

	var member []string

	query := fmt.Sprintf("SELECT name FROM user WHERE team_id = ?")
	rows, err := m.Connection.Query(query, teamID)
	if err != nil {
		return member, err
	}
	for rows.Next() {
		var m string
		if err := rows.Scan(&m); err != nil {
			return member, err
		}
		member = append(member, m)
	}
	return member, nil
}

// func (m *TeamModel) FindMember(id int) ([]string, error) {
// 	m.Open()
// 	defer m.Close()

// 	query := fmt.Sprintf("SELECT id, name, email, enable FROM %s WHERE enable = 1", m.Table)
// 	rows, err := m.Connection.Query(query)
// 	if err != nil {
// 		return nil, errors.New("Database error")
// 	}
// 	for rows.Next() {
// 		var team Team
// 		if err := rows.Scan(&team.id, &user.name, &team.enable); err != nil {
// 			return teams, errors.New("Database error")
// 		}
// 		teams = append(teams, team)
// 	}
// 	return teams, nil
// }

// func (m *UserModel) PasswordCheck(email string, password string) (string, error) {
// 	m.Open()
// 	defer m.Close()

// 	var name string
// 	hashedPassword := GenerateHashedPassword(password)
// 	stmtOut, err := m.Connection.Prepare(fmt.Sprintf("SELECT name FROM %s WHERE email = ? AND hashed_password = ?", m.Table))
// 	if err != nil {
// 		return "", err
// 	}

// 	if err := stmtOut.QueryRow(email, hashedPassword).Scan(&name); err != nil {
// 		return "", err
// 	}
// 	return name, nil
// }

// func (m *UserModel) Enable(id int) error {
// 	m.Open()
// 	defer m.Close()

// 	stmtOut, err := m.Connection.Prepare(fmt.Sprintf("UPDATE %s SET enable = 1 WHERE id = ?", m.Table))
// 	if err != nil {
// 		return errors.New("Database error")
// 	}
// 	if err := stmtOut.QueryRow(id); err != nil {
// 		return errors.New("Database error")
// 	}
// 	return nil
// }

// func (m *UserModel) Disenable(id int) error {
// 	m.Open()
// 	defer m.Close()

// 	stmtOut, err := m.Connection.Prepare(fmt.Sprintf("UPDATE %s SET enable = 0 WHERE id = ?", m.Table))
// 	if err != nil {
// 		return errors.New("Database : query error")
// 	}
// 	if err := stmtOut.QueryRow(id); err != nil {
// 		return errors.New("Database error")
// 	}
// 	return nil
// }
