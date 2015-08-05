package usersystem

import (
	"encoding/json"
	"fmt"
	"github.com/boltdb/bolt"
	"log"
	"time"
)

type UserDS struct {
	DB *bolt.DB
}

func (userDS UserDS) Open() {
	db, err := bolt.Open("user.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	userDS.DB = db
}

func (userDS UserDS) UpdateUser(user User) {
	db := userDS.DB
	err := db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte("user"))
		if err != nil {
			return err
		}
		encoded, err := json.Marshal(user)
		return bucket.Put([]byte(user.Name), encoded)
	})
	if err != nil {
		log.Fatal(err)
	}
}

func (userDS UserDS) GetUser(name string) User {
	var user User
	db := userDS.DB
	db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("user"))
		v := bucket.Get([]byte(name))
		err := json.Unmarshal(v, &user)
		if err != nil {
			log.Fatal(err)
		}
		return nil
	})
	return user
}

type User struct {
	Id      int32
	Name    string
	Email   string
	Created time.Time
}

func (user User) String() string {
	str := "[User] "
	str += "id=" + fmt.Sprint(user.Id)
	str += " name=" + user.Name + " email=" + user.Email
	str += " time=" + user.Created.String()
	return str
}
