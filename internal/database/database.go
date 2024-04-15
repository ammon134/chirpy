package database

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync"
)

type DB struct {
	mux  *sync.RWMutex
	path string
}

type DBStructure struct {
	ChirpTable ChirpTable
	UserTable  UserTable
}

func NewDB(path string) (*DB, error) {
	// ensure db exists
	db := &DB{
		path: path,
		mux:  &sync.RWMutex{},
	}
	err := db.ensureDB()
	return db, err
}

func (db *DB) ensureDB() error {
	_, err := os.Stat(db.path)
	if errors.Is(err, os.ErrNotExist) {
		return db.createDB()
	}
	return err
}

func (db *DB) createDB() error {
	dbs := DBStructure{
		ChirpTable: ChirpTable{
			Chirps:    map[int]Chirp{},
			NextIndex: 1,
		},
		UserTable: UserTable{
			Users:     map[string]User{},
			NextIndex: 1,
		},
	}
	return db.writeDB(dbs)
}

func (db *DB) loadDB() (DBStructure, error) {
	db.mux.RLock()
	defer db.mux.RUnlock()

	dbs := DBStructure{}
	file, err := os.ReadFile(db.path)
	if err != nil {
		return dbs, err
	}

	err = json.Unmarshal(file, &dbs)
	if err != nil {
		fmt.Println("error here")
		return dbs, err
	}

	return dbs, nil
}

func (db *DB) writeDB(dbs DBStructure) error {
	db.mux.Lock()
	defer db.mux.Unlock()

	file, err := json.Marshal(dbs)
	if err != nil {
		return err
	}
	err = os.WriteFile(db.path, file, 0666)
	if err != nil {
		return err
	}
	return nil
}
