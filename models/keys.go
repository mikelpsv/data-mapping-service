package models

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/mikelpsv/data-mapping-service/app"
)

type Key struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
	NsId int64  `json:"-"`
}

func (k *Key) FindById(keyId int64) (*Key, error) {
	k.Clean()

	sql := "SELECT id, name FROM keys WHERE id = $1"
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

	sql := "DELETE FROM keys WHERE id = $1"
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
		sqlIns := "INSERT INTO keys (name, ns_id) VALUES($1, $2) RETURNING id"
		err = app.Db.QueryRow(sqlIns, k.Name, k.NsId).Scan(&k.Id)
		if err != nil {
			return nil, err
		}
	} else {
		sqlUpd := "UPDATE keys SET name = $1, ns_id = $2 WHERE id = $2"
		res, err = app.Db.Exec(sqlUpd, k.Name, k.Id, k.NsId)
		if err != nil {
			return nil, err
		}
		aff, err := res.RowsAffected()
		if err != nil {
			return nil, err
		}
		if aff != 1 {
			return nil, fmt.Errorf("rows affected %d", aff)
		}
	}
	return k, nil
}

func (k *Key) Clean() {
	k.Id = 0
	k.Name = ""
}

func (k *Key) FindByName(keyName string) (*Key, error) {
	k.Clean()

	sql := "SELECT id, name, ns_id FROM keys WHERE name = $1"
	rows, err := app.Db.Query(sql, keyName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, errors.New("no records found")
	}

	err = rows.Scan(&k.Id, &k.Name, &k.NsId)
	if err != nil {
		return nil, err
	}
	return k, nil
}
