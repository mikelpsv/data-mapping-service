package models

import (
	"database/sql"
	"errors"
	"github.com/mikelpsv/data-mapping-service/app"
)

type Namespace struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

func (n *Namespace) FindById(nsId int64) (*Namespace, error) {
	sql := "SELECT _id, name FROM namespaces WHERE _id = $1"
	rows, err := app.Db.Query(sql, nsId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, errors.New("No records found matching the specified conditions")
	}

	err = rows.Scan(&n.Id, &n.Name)
	if err != nil {
		return nil, err
	}
	return n, nil
}

func (n *Namespace) Delete() error {
	if n.Id == 0 {
		return errors.New("Row ID is empty")
	}
	mappings := new(Mappings)
	count, err := mappings.CountRowsByNamespace(n.Id)
	if err != nil {
		return err
	}
	if count != 0 {
		return errors.New("Namespace is used")
	}

	sql := "DELETE FROM namespaces WHERE _id = $1"
	res, err := app.Db.Exec(sql, n.Id)
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

func (n *Namespace) Store() (*Namespace, error) {
	var err error
	var res sql.Result

	if n.Id == 0 {
		sql := "INSERT INTO namespaces (name) VALUES($1)"
		res, err = app.Db.Exec(sql, n.Name)
	} else {
		sql := "UPDATE namespaces SET name = $1 WHERE _id = $2"
		res, err = app.Db.Exec(sql, n.Name, n.Id)
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

	return n, nil
}
