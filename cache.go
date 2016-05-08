package main

import (
	"fmt"
	"github.com/boltdb/bolt"
	"log"
	"strings"
)

var stickerBucket = []byte("stickers")

// Always follow with `defer db.Close()`.
func initDatabase(dbFilename string) *bolt.DB {
	log.Print("Opening database for caching: ", dbFilename)
	db, err := bolt.Open(dbFilename, 0600, nil)
	check(err)

	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(stickerBucket)
		if err != nil {
			return fmt.Errorf("Create bucket: %s", err)
		}
		return nil
	})

	return db
}

func keyFromGlyphsList (glyphsList []string) []byte {
	key := "ver:" + programVersion + "|glyphs:" + strings.Join(glyphsList, ",")
	return []byte(key)
}

// Check the cache for an image already generated.
func checkCache(glyphsList []string, db *bolt.DB) []byte {
	key := keyFromGlyphsList(glyphsList)
	var v []byte

	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(stickerBucket)
		v = b.Get(key)
		return nil
	})
	check(err)

	return v
}

func addToCache(glyphsList []string, fileID string, db *bolt.DB) {
	key := keyFromGlyphsList(glyphsList)
	id := []byte(fileID)

	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(stickerBucket)
		e := b.Put(key, id)
		return e
	})
	check (err)
}

func dumpCache(db *bolt.DB) {
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(stickerBucket)
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			fmt.Printf("%80s -> %80s", k, v)
		}
		return nil
	})
}