package models

import (
	"errors"
	"fmt"
	"github.com/mikelpsv/data-mapping-service/app"
)

type Namespace struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

func (n *Namespace) Get() ([]Namespace, error) {
	retVal := make([]Namespace, 0)
	sqlGet := "SELECT id, name FROM namespaces"
	rows, err := app.Db.Query(sqlGet)
	if err != nil {
		return retVal, err
	}
	defer rows.Close()

	for rows.Next() {
		ns := Namespace{}
		err := rows.Scan(&ns.Id, &ns.Name)
		if err != nil {
			return retVal, err
		}
		retVal = append(retVal, ns)
	}
	return retVal, nil
}

func (n *Namespace) FindById(nsId int64) (*Namespace, error) {
	sql := "SELECT id, name FROM namespaces WHERE id = $1"
	rows, err := app.Db.Query(sql, nsId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, errors.New("no records found")
	}

	err = rows.Scan(&n.Id, &n.Name)
	if err != nil {
		return nil, err
	}
	return n, nil
}

func (n *Namespace) FindByName(nsName string) (*Namespace, error) {
	sql := "SELECT id, name FROM namespaces WHERE name = $1"
	rows, err := app.Db.Query(sql, nsName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, errors.New("no records found")
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

	sql := "DELETE FROM namespaces WHERE id = $1"
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

	if n.Id == 0 {
		sqlIns := "INSERT INTO namespaces (name) VALUES($1) RETURNING id"
		err = app.Db.QueryRow(sqlIns, n.Name).Scan(&n.Id)
		if err != nil {
			return nil, err
		}
	} else {
		sqlUpd := "UPDATE namespaces SET name = $1 WHERE id = $2"
		res, err := app.Db.Exec(sqlUpd, n.Name, n.Id)
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

	return n, nil
}

func (n *Namespace) Clean() {
	n.Id = 0
	n.Name = ""
}

// GetKeys возвращает список ключей в текущем namespace
func (n *Namespace) GetKeys() ([]Key, error) {
	retVal := make([]Key, 0)

	sqlS := "SELECT id, name, ns_id FROM keys WHERE ns_id = $1"
	rows, err := app.Db.Query(sqlS, n.Id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		k := Key{}
		err = rows.Scan(&k.Id, &k.Name, &k.NsId)
		if err != nil {
			return nil, err
		}
		retVal = append(retVal, k)
	}

	return retVal, nil
}

// CreateKey создает ключ по имени в текущем namespace
func (n *Namespace) CreateKey(keyName string) (*Key, error) {
	k := new(Key)
	k.Name = keyName
	k.NsId = n.Id
	return k.Store()
}

func (n *Namespace) KeyExists(keyName string) bool {
	var count int
	sql := "SELECT COUNT(*) AS count FROM keys WHERE name = $1 AND ns_id = $2"
	row := app.Db.QueryRow(sql, keyName, n.Id)
	_ = row.Scan(&count)
	return count > 0
}
