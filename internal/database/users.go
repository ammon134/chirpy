package database

import "errors"

type UserTable struct {
	Users     map[int]User `json:"users"`
	NextIndex int          `json:"next_index"`
}

type User struct {
	Email          string `json:"email"`
	HashedPassword []byte `json:"hashed_password"`
	IsChirpyRed    bool   `json:"is_chirpy_red"`
	ID             int    `json:"id"`
}

var ErrAlreadyExist = errors.New("already exitst")

func (db *DB) CreateUser(email string, hash []byte) (User, error) {
	dbs, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	user := User{
		ID:             dbs.UserTable.NextIndex,
		Email:          email,
		IsChirpyRed:    false,
		HashedPassword: hash,
	}
	_, err = db.GetUserByEmail(user.Email)
	if err != nil && !errors.Is(err, ErrNotExist) {
		return User{}, ErrAlreadyExist
	}

	dbs.UserTable.Users[dbs.ChirpTable.NextIndex] = user
	dbs.UserTable.NextIndex++

	err = db.writeDB(dbs)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (db *DB) GetUserByID(id int) (User, error) {
	dbs, err := db.loadDB()
	if err != nil {
		return User{}, err
	}
	user, ok := dbs.UserTable.Users[id]
	if !ok {
		return User{}, ErrNotExist
	}
	return user, nil
}

func (db *DB) GetUserByEmail(email string) (User, error) {
	dbs, err := db.loadDB()
	if err != nil {
		return User{}, err
	}
	for _, user := range dbs.UserTable.Users {
		if user.Email == email {
			return user, nil
		}
	}
	return User{}, ErrNotExist
}

func (db *DB) UpdateUser(id int, email string, hashedPassword []byte) (User, error) {
	dbs, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	user, err := db.GetUserByID(id)
	if err != nil {
		return User{}, err
	}

	user.Email = email
	user.HashedPassword = hashedPassword

	dbs.UserTable.Users[id] = user

	err = db.writeDB(dbs)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (db *DB) UpgradeUser(id int) error {
	dbs, err := db.loadDB()
	if err != nil {
		return err
	}

	user, err := db.GetUserByID(id)
	if err != nil {
		return err
	}

	user.IsChirpyRed = true

	dbs.UserTable.Users[id] = user

	err = db.writeDB(dbs)
	if err != nil {
		return err
	}

	return nil
}
