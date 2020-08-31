package DBConnection

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/dgraph-io/badger"
	"log"
	"time"
)

var DB *badger.DB
func init() {
	var err error
	DB, err = badger.Open(badger.DefaultOptions("DB/tables").WithTruncate(true))
	if err != nil {
		log.Fatal("error while opening DB", err)
	}
}
func AddToDB(key, message []byte)error{

	check := getKey(key) // check if key already present,
	if check{
		fmt.Printf("\n%s\n","****calling merge****")
		err := doMerge(key,message) // if it does using merge instead of update
		fmt.Printf("\n%s\n","****merge completed****")
		if err != nil {
			return err
		}
		return nil
	} else{ //if key is not already present doing add operation
		fmt.Printf("\n%s\n","****calling update****")
		err := DB.Update(func(txn1 *badger.Txn) error {
			err := txn1.Set(key,message)
			if err != nil{
				return err
			}
			return nil
		})
		fmt.Printf("\n%s\n","****update completed****")
		if err != nil {
			log.Println(errors.New("Error Updating DB"))
			return err
		}
	}

	return nil
}

func GetFromDB(key []byte) ([]byte, error){
	var dst []byte
	err := DB.View(func(txn2 *badger.Txn)error{
		defer txn2.Discard()

		item,err := txn2.Get(key)
		if err != nil {
			return err
		}
		err = item.Value(func(val []byte) error {
			dst = append(dst,val...)
			return nil
		})
		if err != nil{
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return dst,nil
}

func getKey(key []byte)bool{
	var dbKey []byte
	txn3 := DB.NewTransaction(false)
	defer txn3.Discard()
	item, err := txn3.Get(key)
	if err != nil{
		return false
	}
	if bytes.Equal(item.KeyCopy(dbKey), key){
		return true
	}
	return false
}

func doMerge(key,message []byte) error{
	add := badger.MergeFunc(func(oldValue,newValue []byte)[]byte{
		return append(oldValue, newValue...)
	})
	mrgOpr := DB.GetMergeOperator(key,add,200*time.Millisecond)
	err := mrgOpr.Add(message)
	if err != nil{
		return err
	}
	return nil
}

func CloseDB() error {
	return DB.Close()
}

func DropAll() error{
	err := DB.DropAll()
	if err != nil {
		return err
	}
	fmt.Println("***DROPPED All DATA*** ")
	return nil
}