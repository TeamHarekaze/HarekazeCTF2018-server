package BaseModel

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	// "errors"
)

type Base struct {
	Table      string
	Connection *sql.DB
	Primarykey string
}

func (m *Base) Open() {
	adder := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)
	var err error
	m.Connection, err = sql.Open("mysql", adder)
	if err != nil {
		panic(err.Error())
	}
}

func (m *Base) Close() {
	m.Connection.Close()
}
