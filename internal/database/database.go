package database

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"strconv"
	"sync"

	"golang.org/x/crypto/bcrypt"
)

type DB struct {
	path string
	mu   *sync.RWMutex
}

type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
    Users map[int]User `json:"users"`
}

type Chirp struct {
	ID   int    `json:"id"`
	Body string `json:"body"`
}

type User struct {
    ID int `json:"id"`
    Email string `json:"email"`
    Password string `json:"password"`
}

func NewDB(path string) (*DB, error) {
	db := &DB{
		path: path,
		mu:   &sync.RWMutex{},
	}
	err := db.ensureDB()
	return db, err
}

func (db *DB) LoginUser(email string, password string) (User, error) {
    dbStructure, err := db.loadDB()
    if err != nil {
        return User{}, err
    }

    authUser := User{}

    for _, user := range dbStructure.Users {
        if user.Email == email {
            err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
            if err != nil {
                return User{}, err
            }
            authUser = user
        }
    }

    return authUser, nil
}

func (db *DB) UpdateUser(id string, email string, password string) (User, error) {
    dbStructure, err := db.loadDB()
    if err != nil {
        return User{}, err
    }

    encrypted, err := bcrypt.GenerateFromPassword([]byte(password), 4)
    if err != nil {
        return User{}, err
    }

    intID, err := strconv.Atoi(id)
    if err != nil {
        log.Fatalf("erro id %v", err)
    }

    requestedUser := User{}
    for _, user := range dbStructure.Users {
        if user.ID == intID {
            requestedUser = user
        } else {
            continue
        }
    }

    if requestedUser.Email == "" {
       err = errors.New("Not found") 
    } 


    user := User{
        ID: intID,
        Email: email,
        Password: string(encrypted),
    }

    dbStructure.Users[intID] = user

    err = db.writeDB(dbStructure)
    if err != nil{
        return User{}, err
    }

    return user, nil
}


func (db *DB) CreateUser(email string, password string) (User, error) {
    dbStructure, err := db.loadDB()
    if err != nil {
        return User{}, err
    }

    encrypted, err := bcrypt.GenerateFromPassword([]byte(password), 4)
    if err != nil {
        return User{}, err
    }

    id := len(dbStructure.Users) + 1
    user := User{
        ID: id,
        Email: email,
        Password: string(encrypted),
    }

    dbStructure.Users[id] = user

    err = db.writeDB(dbStructure)
    if err != nil{
        return User{}, err
    }

    return user, nil
}

func (db *DB) CreateChirp(body string) (Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}

	id := len(dbStructure.Chirps) + 1
	chirp := Chirp{
		ID:   id,
		Body: body,
	}
	dbStructure.Chirps[id] = chirp

	err = db.writeDB(dbStructure)
	if err != nil {
		return Chirp{}, err
	}

	return chirp, nil
}

func (db *DB) FindChirp(id int) (Chirp, error) {
    dbStructure, err := db.loadDB()
    if err != nil {
            return Chirp{}, err
    }

    requestedChirp := Chirp{}
    for _, chirp := range dbStructure.Chirps {
        if chirp.ID == id {
            requestedChirp = chirp
        } else {
            continue
        }
    }

    if requestedChirp.Body == "" {
       err = errors.New("Not found") 
    } 
    
    return requestedChirp, err
}


func (db *DB) GetChirps() ([]Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return nil, err
	}

	chirps := make([]Chirp, 0, len(dbStructure.Chirps))
	for _, chirp := range dbStructure.Chirps {
		chirps = append(chirps, chirp)
	}

	return chirps, nil
}

func (db *DB) createDB() error {
	dbStructure := DBStructure{
		Chirps: map[int]Chirp{},
        Users: map[int]User{},
	}
	return db.writeDB(dbStructure)
}

func (db *DB) ensureDB() error {
	_, err := os.ReadFile(db.path)
	if errors.Is(err, os.ErrNotExist) {
		return db.createDB()
	}
	return err
}

func (db *DB) loadDB() (DBStructure, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	dbStructure := DBStructure{}
	dat, err := os.ReadFile(db.path)
	if errors.Is(err, os.ErrNotExist) {
		return dbStructure, err
	}
	err = json.Unmarshal(dat, &dbStructure)
	if err != nil {
		return dbStructure, err
	}

	return dbStructure, nil
}

func (db *DB) writeDB(dbStructure DBStructure) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	dat, err := json.Marshal(dbStructure)
	if err != nil {
		return err
	}

	err = os.WriteFile(db.path, dat, 0600)
	if err != nil {
		return err
	}
	return nil
}

