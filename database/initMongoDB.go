package database

import (
	"log"
	"main/env"
	"time"

	"gopkg.in/mgo.v2"
)

var session *mgo.Session

// InitMongoDB 初始化一个MongoDB会话，并持有该链接
func init() {

	dialInfo := &mgo.DialInfo{
		Addrs:     []string{env.MongoDBHost},
		Direct:    false,
		Timeout:   time.Second * 60,
		Username:  env.MongoDBUser,
		Password:  env.MongoDBPassword,
		PoolLimit: env.MongoDBPoolLimit,
	}

	var err error
	session, err = mgo.DialWithInfo(dialInfo)
	if err != nil {
		log.Panic(err)
	}
	session.SetMode(mgo.Monotonic, true)
}

type SessionStore struct {
	session *mgo.Session
}

func (d *SessionStore) C(name string) *mgo.Collection {
	return d.session.DB(env.MongoDB).C(name)
}

//为每一HTTP请求创建新的DataStore对象
func NewSessionStore() *SessionStore {

	ds := &SessionStore{
		session: session.Copy(),
	}
	return ds
}

func (d *SessionStore) Close() {
	d.session.Close()
}

func GetErrNotFound() error {
	return mgo.ErrNotFound
}

func Over() {
	session.Close()
}
