package models

import (
	"database/sql"
	"errors"
	"github.com/mikelpsv/data-mapping-service/app"
)

type Key struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

func (k *Key) FindById(keyId int64) (*Key, error) {
	sql := "SELECT _id, name FROM keys WHERE _id = $1"
	rows, err := app.Db.Query(sql, keyId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, errors.New("No records found matching the specified conditions")
	}

	err = rows.Scan(&k.Id, &k.Name)
	if err != nil {
		return nil, err
	}
	return k, nil
}

func (k *Key) Delete() error {

	if k.Id == 0 {
		return errors.New("Row ID is empty")
	}

	mappings := new(Mappings)
	count, err := mappings.CountRowsByKey(k.Id)
	if err != nil {
		return err
	}

	if count != 0 {
		return errors.New("Key is used")
	}

	sql := "DELETE FROM keys WHERE _id = $1"
	res, err := app.Db.Exec(sql, k.Id)
	if err != nil {
		return err
	}
	aff, err := res.RowsAffected()
	if err != nil {
		return nil
	}
	if aff != 1 {
		return errors.New("Rows affected 0")
	}
	return nil
}

func (k *Key) Store() (*Key, error) {

	var err error
	var res sql.Result

	if k.Id == 0 {
		sql := "INSERT INTO keys (name) VALUES($1)"
		res, err = app.Db.Exec(sql, k.Name)
	} else {
		sql := "UPDATE keys SET name = $1 WHERE _id = $2"
		res, err = app.Db.Exec(sql, k.Name, k.Id)
	}

	if err != nil {
		return nil, err
	}

	aff, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}
	if aff != 1 {
		return nil, errors.New("Rows affected 0")
	}

	return k, nil
}
