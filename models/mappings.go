package models

import (
	"errors"
	"github.com/mikelpsv/data-mapping-service/app"
)

type Mapping struct {
	Id          int64  `json:"id"`
	NamespaceId int64  `json:"namespace_id"`
	KeyId       int64  `json:"key_id"`
	ValExt      string `json:"val_ext"`
	ValInt      string `json:"val_int"`
	Payload     string `json:"payload"`
}

type Mappings []Mapping

func (m *Mapping) FindById(mapId int64) (*Mapping, error) {
	sql := "SELECT _id, ns_id, key_id, val_ext, val_int, payload FROM mappings WHERE _id = $1"
	rows, err := app.Db.Query(sql, mapId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, errors.New("No records found matching the specified conditions")
	}

	err = rows.Scan(&m.Id, &m.NamespaceId, &m.KeyId, &m.ValExt, &m.ValInt, &m.Payload)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (m *Mapping) Delete() error {

	if m.Id == 0 {
		return errors.New("Row ID is empty")
	}

	sql := "DELETE FROM mappings WHERE _id = $1"
	res, err := app.Db.Exec(sql, m.Id)
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

func (m *Mapping) Store() (*Mapping, error) {
	if m.Id == 0 {
		sql := "INSERT INTO mappings (ns_id, key_id, val_ext, val_int, payload) VALUES($1, $2, $3, $4, $5)"
		_, err := app.Db.Exec(sql, m.NamespaceId, m.KeyId, m.ValExt, m.ValInt, m.Payload)
		if err != nil {
			return nil, err
		}

	}
	return m, nil
}

func (m *Mappings) GetByExtValue(nsId int64, keyId int64, extVal string) (*Mappings, error) {
	sql := "SELECT _id, ns_id, key_id, val_ext, val_int, payload FROM mappings WHERE ns_id = $1 AND key_id = $2 AND val_ext = $3"
	rows, err := app.Db.Query(sql, nsId, keyId, extVal)
	if err != nil {
		return nil, err
	}

	*m = make(Mappings, 0, 5)
	for rows.Next() {
		mapRow := new(Mapping)
		err = rows.Scan(&mapRow.Id, &mapRow.NamespaceId, &mapRow.KeyId, &mapRow.ValExt, &mapRow.ValInt, &mapRow.Payload)
		if err != nil {
			return nil, err
		}
		*m = append(*m, *mapRow)
	}
	return m, nil
}

func (m *Mappings) GetByIntValue(nsId int64, keyId int64, intVal string) (*Mappings, error) {
	sql := "SELECT _id, ns_id, key_id, val_ext, val_int, payload FROM mappings WHERE ns_id = $1 AND key_id = $2 AND val_int = $3"
	rows, err := app.Db.Query(sql, nsId, keyId, intVal)
	if err != nil {
		return nil, err
	}

	*m = make(Mappings, 0, 5)
	for rows.Next() {
		mapRow := new(Mapping)
		err = rows.Scan(&mapRow.Id, &mapRow.NamespaceId, &mapRow.KeyId, &mapRow.ValExt, &mapRow.ValInt, &mapRow.Payload)
		if err != nil {
			return nil, err
		}
		*m = append(*m, *mapRow)
	}
	return m, nil
}

func (m *Mappings) CountRowsByNamespace(nsId int64) (count int64, err error) {
	sql := "SELECT COUNT(_id) AS count FROM mappings WHERE ns_id = $1"
	rows, err := app.Db.Query(sql, nsId)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	if !rows.Next() {
		return 0, errors.New("No records found matching the specified conditions")
	}

	rows.Scan(&count)
	return count, nil
}

func (m *Mappings) CountRowsByKey(keyId int64) (count int64, err error) {
	sql := "SELECT COUNT(_id) AS count FROM mappings WHERE key_id = $1"
	rows, err := app.Db.Query(sql, keyId)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	if !rows.Next() {
		return 0, errors.New("No records found matching the specified conditions")
	}

	rows.Scan(&count)
	return count, nil
}
