package app

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"time"
)

var Db *sql.DB

func InitDb(host, dbname, dbuser, dbpass string) {
	var err error

	connStr := fmt.Sprintf("host=%s dbname=%s user=%s password=%s sslmode=disable", host, dbname, dbuser, dbpass)
	Db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	err = Db.Ping()
	if err != nil {
		log.Fatal(err)
	}

}

func Close() {
	Db.Close()
}

func Install(loadTestData bool) {

	sql := ""
	_, err := Db.Exec(sql)
	if err != nil {
		log.Fatal(err.Error())
	}

	if loadTestData {
		sql = ""
		_, err = Db.Exec(sql, "Основная система", "$2a$10$SyaL6fNLoPplhxqOlmN7MuA/MxXm7/F9AX.NqVDRSb4xi9YrHQg36", "1234567890", 3600, time.Now(), time.Now())
		if err != nil {
			log.Fatal(err.Error())
		}
	}
	sql = ""

	_, err = Db.Exec(sql)
	if err != nil {
		log.Fatal(err.Error())
	}

	if loadTestData {
		sql = ""
		_, err = Db.Exec(sql, "Пользователь1", "$2a$10$/ui7v1gRNVLSRtfHOib/muwP5TAr7e33c9y7LPpfdUHmCIWJSO8ny", "1", time.Now(), time.Now())
		if err != nil {
			log.Fatal(err.Error())
		}

		sql = ""
		_, err = Db.Exec(sql, "Пользователь2", "$2a$10$B2pAjD62tq0QOAswYaXqFe9cxVEgMm8PVTL4SfgIl3CNJUkmNITQm", "1", time.Now(), time.Now())
		if err != nil {
			log.Fatal(err.Error())
		}
	}
}
