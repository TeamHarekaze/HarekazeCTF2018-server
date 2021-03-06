package UserModel

import (
	"crypto/sha512"
	"errors"
	"fmt"

	"github.com/TeamHarekaze/HarekazeCTF2018-server/datamodels/BaseModel"
)

const (
	table      = "user"
	primarykey = "id"
)

func GenerateHashedPassword(userPassword string) string {
	return fmt.Sprintf("%x", sha512.Sum512([]byte(userPassword)))
}

// ValidatePassword will check if passwords are matched.
// func ValidatePassword(userPassword string, hashed []byte) (bool, error) {
//     if err := bcrypt.CompareHashAndPassword(hashed, []byte(userPassword)); err != nil {
//         return false, err
//     }
//     return true, nil
// }

type User struct {
	id     int
	name   string
	email  string
	enable bool
}

type UserModel struct {
	BaseModel.Base
	User User
}

func New() *UserModel {
	base := new(UserModel)
	base.Table = table
	base.Primarykey = primarykey
	return base
}

func (m *UserModel) All() ([]User, error) {
	m.Open()
	defer m.Close()

	var users []User

	query := fmt.Sprintf("SELECT id, name, email, enable FROM %s", m.Table)
	rows, err := m.Connection.Query(query)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.id, &user.name, &user.email, &user.enable); err != nil {
			return users, err
		}
		users = append(users, user)
	}
	return users, nil
}
func (m *UserModel) GetNameFromEmail(email string) (string, error) {
	m.Open()
	defer m.Close()

	var name string
	stmtOut, err := m.Connection.Prepare(fmt.Sprintf("SELECT name FROM %s WHERE email = ?", m.Table))
	if err != nil {
		return "", errors.New("Database : query error")
	}
	if err := stmtOut.QueryRow(email).Scan(&name); err != nil {
		return "", err
	}
	if name == "" {
		return "", errors.New("user not found")
	}
	return name, nil
}

func (m *UserModel) AllEnable() ([]User, error) {
	m.Open()
	defer m.Close()

	var users []User

	query := fmt.Sprintf("SELECT id, name, email, enable FROM %s WHERE enable = 1", m.Table)
	rows, err := m.Connection.Query(query)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.id, &user.name, &user.email, &user.enable); err != nil {
			return users, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (m *UserModel) FindId(id int) (User, error) {
	m.Open()
	defer m.Close()

	var user User
	stmtOut, err := m.Connection.Prepare(fmt.Sprintf("SELECT id, name, email, enable FROM %s WHERE id = ?", m.Table))
	if err != nil {
		return user, err
	}
	if stmtOut.QueryRow(id).Scan(&user.id, &user.name, &user.email, &user.enable) != nil {
		return user, err
	}
	return user, nil
}

func (m *UserModel) Add(name string, email string, password string, teamName string) error {
	m.Open()
	defer m.Close()

	hashedPassword := GenerateHashedPassword(password)
	query := fmt.Sprintf("INSERT INTO %s (name, email, hashed_password, team_id) SELECT ?, ?, ?, id FROM team WHERE name = ?", m.Table)
	_, err := m.Connection.Exec(query, name, email, hashedPassword, teamName)
	if err != nil {
		return err
	}
	return nil
}

func (m *UserModel) PasswordCheck(email string, password string) (bool, error) {
	m.Open()
	defer m.Close()

	hashedPassword := GenerateHashedPassword(password)
	stmtOut, err := m.Connection.Prepare(fmt.Sprintf("SELECT COUNT(name) FROM %s WHERE email = ? AND hashed_password = ?", m.Table))
	if err != nil {
		return false, errors.New("Database : query error")
	}

	var count int
	if stmtOut.QueryRow(email, hashedPassword).Scan(&count) != nil {
		return false, errors.New("email or username is incorrect")
	}
	return count == 1, nil
}

func (m *UserModel) Enable(id int) error {
	m.Open()
	defer m.Close()

	stmtOut, err := m.Connection.Prepare(fmt.Sprintf("UPDATE %s SET enable = 1 WHERE id = ?", m.Table))
	if err != nil {
		return err
	}
	if stmtOut.QueryRow(id) == nil {
		return errors.New("Database error(stmtOut.QueryRow(id) == nil)")
	}
	return nil
}

func (m *UserModel) Disenable(id int) error {
	m.Open()
	defer m.Close()

	stmtOut, err := m.Connection.Prepare(fmt.Sprintf("UPDATE %s SET enable = 0 WHERE id = ?", m.Table))
	if err != nil {
		return errors.New("Database : query error")
	}
	if stmtOut.QueryRow(id) == nil {
		return errors.New("Database error(stmtOut.QueryRow(id))")
	}
	return nil
}

func (m *UserModel) UsedChack(name string, email string) (bool, error) {
	m.Open()
	defer m.Close()

	var users []User

	query := fmt.Sprintf("select id, name, email, enable from %s WHERE name = ? OR email = ?", m.Table)
	rows, err := m.Connection.Query(query, name, email)
	if err != nil {
		return false, errors.New("Database : query error")
	}
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.id, &user.name, &user.email, &user.enable); err != nil {
			return false, err
		}
		users = append(users, user)
	}
	return len(users) != 0, nil
}

func (m *UserModel) GetUserInfo(username string) (string, string, string, error) {
	m.Open()
	defer m.Close()

	query := fmt.Sprintf("SELECT user.id, team.id, team.name FROM %s LEFT JOIN team ON  team.id = user.team_id WHERE user.name = ?", m.Table)
	stmtOut, err := m.Connection.Prepare(query)
	if err != nil {
		return "", "", "", errors.New("Database : query error")
	}

	var userid string
	var teamname string
	var teamid string
	if stmtOut.QueryRow(username).Scan(&userid, &teamid, &teamname) != nil {
		return "", "", "", errors.New("email or username is incorrect")
	}
	return userid, teamname, teamid, nil
}
