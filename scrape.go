package main

import (
	"log"
	"time"

	"github.com/asggo/store"
)

var conf Config

func getStoreConn() *store.Store {
	var ds *store.Store
	var err error

	for tries := 1; tries < 20; tries += 2 {
		// Create our connection to the key value store.
		ds, err = store.NewStore(conf.dbFile)
		if err != nil {
			log.Printf("[-] Cannot open database: %s\n", err)
			time.Sleep(1 << uint(tries) * time.Millisecond)
		} else {
			break
		}
	}

	return ds

}

func initStore(s *store.Store) {
	s.CreateBucket("awskeys")
	s.CreateBucket("emails")
	s.CreateBucket("creds")
	s.CreateBucket("privkeys")
	s.CreateBucket("keywords")
	s.CreateBucket("pastes")
}

func cleanKeys() {
	now := time.Now()

	for key, _ := range conf.keys {
		if now.Sub(conf.keys[key]) > conf.maxTime {
			delete(conf.keys, key)
		}
	}
}

func scrape() {
	scrapePastes()
	scrapeGists()
}

func main() {
	conf = newConfig()

	// Ensure our database is initialized.
	ds := getStoreConn()
	initStore(ds)
	ds.Close()

	for {
		conf.ds = getStoreConn()
		scrape()
		conf.ds.Close()
		time.Sleep(conf.sleep)
		cleanKeys()
	}
}
