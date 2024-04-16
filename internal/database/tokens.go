package database

import "time"

func (db *DB) RevokeToken(token string) error {
	dbs, err := db.loadDB()
	if err != nil {
		return err
	}

	dbs.RevokedTokens[token] = time.Now().UTC()

	err = db.writeDB(dbs)
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) IsRevoked(token string) (bool, error) {
	dbs, err := db.loadDB()
	if err != nil {
		return false, err
	}

	_, ok := dbs.RevokedTokens[token]
	return ok, nil
}
