package main

import (
	"fmt"
	"log"
	"time"

	"github.com/asdine/storm"
	"github.com/asdine/storm/codec/gob"
	"github.com/satori/go.uuid"
	"github.com/tsingson/experiments/bbolt/vod"
	"github.com/tsingson/fastweb/fasthttputils"

	"v/github.com/sanity-io/litter@v1.1.0"
)

func main() {
	path, _ := fasthttputils.GetCurrentExecDir()
	db, err := storm.Open(path+"/test2.db", storm.Codec(gob.Codec))
	db2, err := storm.Open(path+"/test3.db", storm.Codec(gob.Codec))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	/**
	db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("MyBucket"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})
	*/
	// Entries
	err = addEntry(db, "apple", 100, time.Now())
	if err != nil {
		log.Fatal(err)
	}
	err = addEntry(db, "bread", 300, time.Now().AddDate(0, 0, -2))
	if err != nil {
		log.Fatal(err)
	}

	var one vod.Entry
	err = db.One("Food", "apple", &one)
	if err != nil {
		fmt.Println("--------------------- not found ------------------")
	}
	fmt.Println("--------------------- found ------------------")
	litter.Dump(one)

	var data3 []vod.Entry

	db.All(&data3)
	fmt.Println("---------------------------------------------")
	litter.Dump(data3)

}

func addEntry(db *storm.DB, food string, calories int, date time.Time) error {
	entry := vod.Entry{Food: food, Calories: calories, Date: date}
	id, _ := uuid.NewV4()

	config := new(vod.Config)
	config.Height = 234234
	config2 := config
	entry.Config = []*vod.Config{config, config2}

	entry.EntryID = id.String()
	err := db.Save(&entry)

	return err
}
