package database

func (db *DB) CreateChirp(body string) (Chirp, error) {
	// loadDB
	dbs, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}
	// create Chirp, give it ID
	chirp := Chirp{
		ID:   dbs.NextIndex,
		Body: body,
	}
	dbs.Chirps[dbs.NextIndex] = chirp
	dbs.NextIndex++
	// writeDB
	err = db.writeDB(dbs)
	if err != nil {
		return Chirp{}, err
	}

	return chirp, nil
}

func (db *DB) GetChirps() ([]Chirp, error) {
	// loadDB
	dbs, err := db.loadDB()
	if err != nil {
		return nil, err
	}

	chirps := []Chirp{}
	for _, chirp := range dbs.Chirps {
		chirps = append(chirps, chirp)
	}
	return chirps, nil
}
