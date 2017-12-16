package UserModel

import (
	"crypto/sha512"
	"errors"
	"fmt"

	"./BaseModel"
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
		return nil, errors.New("Database error")
	}
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.id, &user.name, &user.email, &user.enable); err != nil {
			return users, errors.New("Database error")
		}
		users = append(users, user)
	}
	return users, nil
}

func (m *UserModel) AllEnable() ([]User, error) {
	m.Open()
	defer m.Close()

	var users []User

	query := fmt.Sprintf("SELECT id, name, email, enable FROM %s WHERE enable = 1", m.Table)
	rows, err := m.Connection.Query(query)
	if err != nil {
		return nil, errors.New("Database error")
	}
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.id, &user.name, &user.email, &user.enable); err != nil {
			return users, errors.New("Database error")
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
		return user, errors.New("Database error")
	}
	if err := stmtOut.QueryRow(id).Scan(&user.id, &user.name, &user.email, &user.enable); err != nil {
		return user, errors.New("Database error")
	}
	return user, nil
}

func (m *UserModel) Add(name string, email string, password string) error {
	m.Open()
	defer m.Close()

	hashedPassword := GenerateHashedPassword(password)
	stmtOut, err := m.Connection.Prepare(fmt.Sprintf("INSERT INTO %s (name, email, hashed_password) VALUES(?, ?, ?)", m.Table))
	if err != nil {
		return errors.New("Database : query error")
	}
	// error is here...why err == nil
	if err := stmtOut.QueryRow(name, email, hashedPassword); err == nil {
		fmt.Println(err)
		return errors.New("Database error")
	}
	return nil
}

func (m *UserModel) PasswordCheck(email string, password string) (string, error) {
	m.Open()
	defer m.Close()

	var name string
	hashedPassword := GenerateHashedPassword(password)
	stmtOut, err := m.Connection.Prepare(fmt.Sprintf("SELECT name FROM %s WHERE email = ? AND hashed_password = ?", m.Table))
	if err != nil {
		return "", err
	}

	if err := stmtOut.QueryRow(email, hashedPassword).Scan(&name); err != nil {
		return "", err
	}
	return name, nil
}

func (m *UserModel) Enable(id int) error {
	m.Open()
	defer m.Close()

	stmtOut, err := m.Connection.Prepare(fmt.Sprintf("UPDATE %s SET enable = 1 WHERE id = ?", m.Table))
	if err != nil {
		return errors.New("Database error")
	}
	if err := stmtOut.QueryRow(id); err != nil {
		return errors.New("Database error")
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
	if err := stmtOut.QueryRow(id); err != nil {
		return errors.New("Database error")
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
			return false, errors.New("Database error")
		}
		users = append(users, user)
	}
	return len(users) != 0, nil
}
