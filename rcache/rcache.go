package rcache

import (
	"bytes"
	"github.com/astaxie/beego/context"
	"github.com/boltdb/bolt"
	"log"
)

type Resource struct {
	Key   []byte
	Value []byte
}

var (
	ResourceChan chan Resource
)

var (
	DBResource = []byte("resource")
)

var db *bolt.DB

func init() {

	ResourceChan = make(chan Resource, 32)

	var err error
	db, err = bolt.Open("./boltdbs/book_server.db", 0600, nil)
	if err != nil {
		log.Panic(err)
	}

	var tables = [][]byte{DBResource}
	err = db.Update(func(tx *bolt.Tx) error {
		for _, l := range tables {
			_, err := tx.CreateBucketIfNotExists(l)
			if err != nil {
				panic(err)
			}
		}
		return nil
	})

	if err != nil {
		panic(err)
	}
}

func GetDb() *bolt.DB {
	return db
}

func Saver() {
	for {
		select {
		case res := <-ResourceChan:
			{
				_ = db.Update(func(tx *bolt.Tx) error {
					b := tx.Bucket(DBResource)
					if b != nil {
						err := b.Put(res.Key, res.Value)
						if err != nil {
							return err
						}
					}
					return nil
				})

			}
		}
	}
}

func FlushAll() {
	_ = db.Update(func(tx *bolt.Tx) error {
		// 结束以后将桶创建回来
		// 按照要求删除桶
		tx.DeleteBucket(DBResource)
		tx.CreateBucketIfNotExists(DBResource)
		return nil
	})
}

func TestCache(ctx *context.Context) {
	var buf bytes.Buffer
	_ = db.View(func(tx *bolt.Tx) error {

		b := tx.Bucket(DBResource)
		if b != nil {
			_ = b.ForEach(func(k, v []byte) error {
				buf.Write(bytes.Join([][]byte{k, v}, []byte(" : ")))
				return nil
			})
		}
		return nil
	})
	ctx.ResponseWriter.Write(buf.Bytes())
}
