package kvdb

import (
	"os"

	"github.com/boltdb/bolt"
)

func InitDB() (*bolt.DB, error) {

	//ex, _ := os.Executable()

	//exPath := filepath.Dir(ex)
	//fmt.Println(exPath)

	if _, err := os.Stat("./db"); os.IsNotExist(err) {
		os.Mkdir("./db", 0700)
	}

	db, err := bolt.Open("./db/go-kb.db", 0644, nil)
	if err != nil {
		panic(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("stacks"))
		if err != nil {
			return err
		}
		return nil
	})

	return db, err

}
