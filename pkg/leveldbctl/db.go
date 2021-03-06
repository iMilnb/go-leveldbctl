package leveldbctl

import (
	"fmt"
	"os"
	"path"

	leveldb "github.com/syndtr/goleveldb/leveldb"
)

func dbexists(dbpath string) bool {
	_, err := os.Stat(path.Join(dbpath, "LOG"))
	return err == nil
}

func Init(dbpath string) error {
	if dbexists(dbpath) {
		return fmt.Errorf("%s was initialized as leveldb", dbpath)
	}

	db, err := leveldb.OpenFile(dbpath, nil)
	if err != nil {
		return fmt.Errorf("cannot open leveldb")
	}
	defer db.Close()

	return nil
}

func Put(dbpath, key, value string) error {
	if !dbexists(dbpath) {
		return fmt.Errorf("%s is not leveldb", dbpath)
	}

	db, err := leveldb.OpenFile(dbpath, nil)
	if err != nil {
		return fmt.Errorf("cannot open leveldb")
	}
	defer db.Close()

	err = db.Put([]byte(key), []byte(value), nil)
	return err
}

func Get(dbpath, key string) (string, bool, error) {
	if !dbexists(dbpath) {
		return "", false, fmt.Errorf("%s is not leveldb", dbpath)
	}

	db, err := leveldb.OpenFile(dbpath, nil)
	if err != nil {
		return "", false, fmt.Errorf("cannot open leveldb")
	}
	defer db.Close()

	has, err := db.Has([]byte(key), nil)
	if err != nil {
		return "", false, fmt.Errorf("cannot open leveldb")
	}
	if !has {
		return "", false, nil
	}

	value, err := db.Get([]byte(key), nil)
	if err != nil {
		return "", true, fmt.Errorf("cannot get value")
	}
	return string(value), true, nil
}

func Delete(dbpath, key string) error {
	if !dbexists(dbpath) {
		return fmt.Errorf("%s is not leveldb", dbpath)
	}

	db, err := leveldb.OpenFile(dbpath, nil)
	if err != nil {
		return fmt.Errorf("cannot open leveldb")
	}
	defer db.Close()

	err = db.Delete([]byte(key), nil)
	return err
}

func Walk(dbpath string, f func(string, string)) error {
	if !dbexists(dbpath) {
		return fmt.Errorf("%s is not leveldb", dbpath)
	}

	db, err := leveldb.OpenFile(dbpath, nil)
	if err != nil {
		return fmt.Errorf("cannot open leveldb")
	}
	defer db.Close()

	s, err := db.GetSnapshot()
	if err != nil {
		return fmt.Errorf("cannot make snapshot leveldb")
	}
	defer s.Release()

	i := s.NewIterator(nil, nil)
	for i.Next() {
		key := string(i.Key())
		value := string(i.Value())
		f(key, value)
	}

	return nil
}
