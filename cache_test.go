package main

import (
	"testing"
	"io/ioutil"
	"path/filepath"
	"os"
	"log"
	"bytes"
	"github.com/boltdb/bolt"
	"fmt"
)

var testDatabaseFilename = "test_db.boltdb"

func tempDbSetup() (string, *bolt.DB) {
	tmpdir, err := ioutil.TempDir("","tg_glyph_bot_test")
	check (err)
	tmpfile := filepath.Join(tmpdir, testDatabaseFilename)
	db := initDatabase(tmpfile)
	return tmpfile, db
}

func tempDbTeardown(tempFile string) {
	d, _ := filepath.Split(tempFile)
	log.Printf("Deleting temp files in %s", d)
	os.RemoveAll(d)
}

func TestNonexistent(t *testing.T){
	dbName, db := tempDbSetup()
	defer db.Close()
	defer tempDbTeardown(dbName)

	result := checkCache([]string{"Doesn't Exist"}, db)
	if result != nil {
		t.Fail()
	}
}

func TestSimple(t *testing.T){
	dbName, db := tempDbSetup()
	defer db.Close()
	defer tempDbTeardown(dbName)

	g := []string{"Open All","Clear All","Discover","Truth"}
	fileID := "AgADBQADtacxG3wqnwr9PWHBQNUPuTvYvTIABF_lBI7oY8NyeRYBAAEC"

	addToCache(g, fileID, db)

	result := checkCache(g, db)
	fmt.Printf("%s -> %s", g, result)
	if bytes.Compare(result, []byte(fileID)) != 0 {
		t.Fail()
	}
}