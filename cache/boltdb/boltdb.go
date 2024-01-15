package boltdb

import (
	"SQLGuardian/consts"
	"github.com/boltdb/bolt"
)

type BoltDB struct {
	dbPath     string
	db         *bolt.DB
	bucketName string
}

// NewBoltDB 创建一个 BoltDB 实例
func NewBoltDB() *BoltDB {
	dbPath := consts.DBPath
	bucketName := consts.BucketName
	db, err := bolt.Open(dbPath, 0600, nil)
	if err != nil {
		panic(err)
	}
	return &BoltDB{dbPath: dbPath, db: db, bucketName: bucketName}
}

// Close 关闭数据库连接
func (bu *BoltDB) Close() {
	if bu.db != nil {
		bu.db.Close()
	}
}

// CreateBucketIfNotExists 创建指定的 bucket，如果不存在的话
func (bu *BoltDB) CreateBucketIfNotExists() error {
	return bu.db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bu.bucketName))
		return err
	})
}

// InsertData 向指定 bucket 中插入数据
func (bu *BoltDB) InsertData(key string, value []byte) error {
	return bu.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bu.bucketName))
		if bucket != nil {
			return bucket.Put([]byte(key), value)
		}
		return bolt.ErrBucketNotFound
	})
}

// QueryData 查询指定 bucket 中的数据
func (bu *BoltDB) QueryData(key string) ([]byte, error) {
	var result []byte
	err := bu.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bu.bucketName))
		if bucket != nil {
			result = bucket.Get([]byte(key))
			return nil
		}
		return bolt.ErrBucketNotFound
	})
	return result, err
}

// DeleteData 删除指定 bucket 中的数据
func (bu *BoltDB) DeleteData(key string) error {
	return bu.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bu.bucketName))
		if bucket != nil {
			return bucket.Delete([]byte(key))
		}
		return bolt.ErrBucketNotFound
	})
}

//type Test struct {
//	Id   int    `json:"id"`
//	Name string `json:"name"`
//}
//
//func test() {
//	// 示例用法
//	dbUtils := NewBoltDB()
//	defer dbUtils.Close()
//
//	// 创建 bucket
//	err := dbUtils.CreateBucketIfNotExists()
//	if err != nil {
//		return
//	}
//
//	key := []byte("test")
//
//	list := []Test{
//		{Id: 1, Name: "test1"},
//		{Id: 2, Name: "test2"},
//	}
//	fmt.Println(list)
//	jsonStr, err := json.Marshal(list)
//	fmt.Println(jsonStr)
//	// 插入数据
//	err = dbUtils.InsertData(key, jsonStr)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// 查询数据
//	result, err := dbUtils.QueryData(key)
//	if err != nil {
//		log.Fatal(err)
//	}
//	var list2 []Test
//	err = json.Unmarshal(result, &list2)
//	log.Println(list2)
//
//	//// 删除数据
//	err = dbUtils.DeleteData(key)
//	if err != nil {
//		log.Fatal(err)
//	}
//}
