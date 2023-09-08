package mysql

import (
	"database/sql"
	"game-server/src/user"
	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
	"sync"
)

var (
	db   *sql.DB
	once sync.Once
)

const (
	url = "root:langwudong@tcp(127.0.0.1:3306)/sunnyland"
)

func InitDB() *sql.DB {
	once.Do(func() {
		db, _ = sql.Open("mysql", url)

		//连接数据库
		err := db.Ping()

		if err != nil {
			logrus.Error(err)
		}

		//设置数据库连接池的最大连接数目
		db.SetMaxOpenConns(100)
	})

	return db
}

func CheckUser(user user.User) bool {
	dbw := InitDB()

	stmt, _ := dbw.Prepare("SELECT * FROM user WHERE username = ? AND password = ?")

	defer stmt.Close()

	rows, err := stmt.Query(user.Username, user.Password)
	defer rows.Close()

	if err != nil {
		logrus.Error(err)
	}

	return rows.Next()
}

func AddUser(user user.User) {
	dbw := InitDB()

	stmt, err := dbw.Prepare("INSERT INTO user (username,password) VALUES (?,?)")

	defer stmt.Close()

	_, err = stmt.Exec(user.Username, user.Password)

	if err != nil {
		logrus.Error(err)
	}
}

func UpdateUser(user user.User, newPassword string) bool {
	dbw := InitDB()

	stmt, err := dbw.Prepare("UPDATE user SET password=? WHERE username=? AND password=?")

	defer stmt.Close()

	result, err := stmt.Exec(newPassword, user.Username, user.Password)

	if err != nil {
		logrus.Error(err)
		return false
	}

	rowsAffected, _ := result.RowsAffected()

	if rowsAffected == 0 {
		return false
	} else {
		return true
	}
}
