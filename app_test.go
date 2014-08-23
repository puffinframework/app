package app_test

import (
	"github.com/boltdb/bolt"
	"github.com/puffinframework/app"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestCreate(t *testing.T) {
	db := openBoltDB()
	defer closeBoltDB(db)

	agg := app.NewAggregate(db)

	err := agg.CreateApp("app1")
	assert.Nil(t, err)

	err = agg.CreateApp("app1")
	assert.NotNil(t, err)

	err = agg.CreateApp("app2")
	assert.Nil(t, err)
}

func TestRemove(t *testing.T) {
	db := openBoltDB()
	defer closeBoltDB(db)

	agg := app.NewAggregate(db)

	err := agg.RemoveApp("app1")
	assert.Nil(t, err)

	err = agg.CreateApp("app1")
	assert.Nil(t, err)

	err = agg.RemoveApp("app1")
	assert.Nil(t, err)
}

func openBoltDB() *bolt.DB {
	db, err := bolt.Open("test.db", 0600, nil)
	if err != nil {
		panic(err)
	}
	return db
}

func closeBoltDB(db *bolt.DB) {
	db.Close()
	os.Remove("test.db")
}
