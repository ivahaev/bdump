package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/boltdb/bolt"
)

var fileName string

func main() {
	args := os.Args
	if len(args) < 2 {
		fmt.Println("No db file provided")
		os.Exit(1)
	}
	fileName := args[1]
	db, err := bolt.Open(fileName, 0600, nil)
	if err != nil {
		fmt.Println("Can't open db", err)
		os.Exit(1)
	}
	dump(db)
}

func dump(db *bolt.DB) {
	data := map[string]interface{}{}
	db.View(func(tx *bolt.Tx) error {
		tx.ForEach(func(name []byte, b *bolt.Bucket) error {
			data[string(name)] = readBucket(b)
			return nil
		})
		return nil
	})
	encoded, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Println("Can't marshal data to JSON", err)
		os.Exit(1)
	}
	dumpFileName := path.Base(fileName)
	dumpFileName = strings.TrimRight(dumpFileName, path.Ext(dumpFileName)) + "_dump.json"
	err = ioutil.WriteFile(dumpFileName, encoded, 0644)
	if err != nil {
		fmt.Println("Can't write dump file", err)
		os.Exit(1)
	}
	fmt.Println("Database dumped to file " + dumpFileName)
}

func readBucket(b *bolt.Bucket) map[string]interface{} {
	data := map[string]interface{}{}
	b.ForEach(func(k, v []byte) error {
		if subB := b.Bucket(k); subB != nil {
			data[string(k)] = readBucket(subB)
			return nil
		}
		var _data interface{}
		err := json.Unmarshal(v, &_data)
		if err != nil {
			fmt.Println("Can't unmarshal data", string(v), err)
			return nil
		}
		data[string(k)] = _data
		return nil
	})
	return data
}
