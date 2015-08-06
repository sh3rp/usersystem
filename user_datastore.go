package usersystem

import (
	"encoding/json"
	"fmt"
	"github.com/boltdb/bolt"
	"log"
	"time"
)

const (
	USER_DB_FILE = "user.db"
	USER_TABLE   = "user"
	CONFIG_TABLE = "config"
)

type UserDS struct {
	DB *bolt.DB
}

func New() UserDS {
	return UserDS{}
}

func (userDS *UserDS) Open() {
	db, err := bolt.Open(USER_DB_FILE, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	userDS.DB = db
}

func (userDS *UserDS) PutConfig(config Config) {
	db := userDS.DB
	err := db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(CONFIG_TABLE))
		if err != nil {
			log.Fatal(err)
			return err
		}
		encoded, err := json.Marshal(config)
		return bucket.Put([]byte("default"), encoded)
	})
	if err != nil {
		log.Fatal(err)
	}
}

func (userDS *UserDS) GetConfig() (Config, error) {
	var config Config
	db := userDS.DB
	err := db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(CONFIG_TABLE))
		conf := bucket.Get([]byte("default"))
		if conf != nil {
			err = json.Unmarshal(conf, &config)
			if err != nil {
				log.Printf("Error unmarshalling config")
				log.Fatal(err)
			}
			return err
		} else {
			config = Config{CurrentId: 0}
			encoded, jsonErr := json.Marshal(config)
			if jsonErr != nil {
				return jsonErr
			}
			return bucket.Put([]byte("default"), encoded)
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	return config, err
}

func (userDS *UserDS) NextID() int32 {
	conf, err := userDS.GetConfig()
	if err != nil {
		conf = Config{CurrentId: 0}
		userDS.PutConfig(conf)
	}
	conf.CurrentId = conf.CurrentId + 1
	userDS.PutConfig(conf)
	return conf.CurrentId
}

func (userDS *UserDS) NewUser(name string) User {
	var user User
	user, err := userDS.GetUser(name)
	if err != nil {
		user = User{
			Id:      userDS.NextID(),
			Name:    name,
			Created: time.Now(),
		}
		userDS.UpdateUser(user)
	}
	return user
}

func (userDS *UserDS) UpdateUser(user User) {
	db := userDS.DB
	err := db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(USER_TABLE))
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

func (userDS *UserDS) GetUser(name string) (User, error) {
	var user User
	db := userDS.DB
	err := db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(USER_TABLE))
		v := bucket.Get([]byte(name))
		err := json.Unmarshal(v, &user)
		if err != nil {
			log.Fatal(err)
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
		return user, err
	} else {
		return user, nil
	}
}

type User struct {
	Id      int32
	Name    string
	Email   string
	Created time.Time
}

type Config struct {
	CurrentId int32
}

func (user User) String() string {
	str := "[User] "
	str += "id=" + fmt.Sprint(user.Id)
	str += " name=" + user.Name + " email=" + user.Email
	str += " time=" + user.Created.String()
	return str
}
