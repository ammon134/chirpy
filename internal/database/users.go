package database

import "errors"

type UserTable struct {
	Users     map[string]User `json:"users"`
	NextIndex int             `json:"next_index"`
}

type User struct {
	Email string `json:"email"`
	Hash  []byte `json:"hash,omitempty"`
	ID    int    `json:"id"`
}

func (db *DB) CreateUser(email string, hash []byte) (User, error) {
	dbs, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	user := User{
		ID:    dbs.UserTable.NextIndex,
		Email: email,
		Hash:  hash,
	}
	if _, ok := dbs.UserTable.Users[user.Email]; ok {
		return User{}, errors.New("email is already used")
	} else {
		dbs.UserTable.Users[user.Email] = user
	}
	dbs.UserTable.NextIndex++

	err = db.writeDB(dbs)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (db *DB) GetUser(email string) (User, error) {
	dbs, err := db.loadDB()
	if err != nil {
		return User{}, err
	}
	if _, ok := dbs.UserTable.Users[email]; !ok {
		return User{}, errors.New("user does not exist")
	} else {
		return dbs.UserTable.Users[email], nil
	}
}
