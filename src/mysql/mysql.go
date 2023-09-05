package mysql

import (
	"database/sql"
	"game-server/src/user"
	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
)

type DbWorker struct {
	dsn  string
	db   *sql.DB
	user user.User
}

func InitDB(url string) DbWorker {
	dbw := DbWorker{
		dsn: url,
	}

	dbw.db, _ = sql.Open("mysql", dbw.dsn)

	//连接数据库
	err := dbw.db.Ping()

	if err != nil {
		logrus.Error("The database connection failed.")
	}

	//设置数据库连接池的最大连接数目
	dbw.db.SetMaxOpenConns(20)

	return dbw
}

func CheckUser(user user.User) bool {
	dbw := InitDB("root:zk000000@tcp(127.0.0.1:3306)/sunny_land")

	stmt, _ := dbw.db.Prepare("SELECT * FROM user WHERE username = ? AND password = ?")
	defer stmt.Close()

	rows, err := stmt.Query(user.Username, user.Password)
	defer rows.Close()

	if err != nil {
		logrus.Error("This query failed.")
	}

	if rows != nil {
		return true
	}

	return false
}
