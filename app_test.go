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

	exists1, err := agg.ExistsApp("app1")
	assert.Nil(t, err)
	assert.False(t, exists1)

	event, err := agg.CreateApp("app1")
	assert.Equal(t, app.CreatedAppEvent{AppId: "app1"}, event)
	assert.Nil(t, err)

	exists1, err = agg.ExistsApp("app1")
	assert.Nil(t, err)
	assert.True(t, exists1)

	event, err = agg.CreateApp("app1")
	assert.NotNil(t, err)

	event, err = agg.CreateApp("app2")
	assert.Equal(t, app.CreatedAppEvent{AppId: "app2"}, event)
	assert.Nil(t, err)

	exists2, err := agg.ExistsApp("app2")
	assert.Nil(t, err)
	assert.True(t, exists2)
}

func TestRemove(t *testing.T) {
	db := openBoltDB()
	defer closeBoltDB(db)

	agg := app.NewAggregate(db)

	err := agg.RemoveApp("app1")
	assert.NotNil(t, err)

	exists1, err := agg.ExistsApp("app1")
	assert.Nil(t, err)
	assert.False(t, exists1)

	event, err := agg.CreateApp("app1")
	assert.Equal(t, app.CreatedAppEvent{AppId: "app1"}, event)
	assert.Nil(t, err)

	err = agg.RemoveApp("app1")
	assert.Nil(t, err)
}

func TestProcessEvents(t *testing.T) {
	db := openBoltDB()
	defer closeBoltDB(db)

	agg := app.NewAggregate(db)

	exists1, err := agg.ExistsApp("app1")
	assert.Nil(t, err)
	assert.False(t, exists1)

	err = agg.OnCreatedApp(app.CreatedAppEvent{AppId: "app1"})
	assert.Nil(t, err)

	exists1, err = agg.ExistsApp("app1")
	assert.Nil(t, err)
	assert.True(t, exists1)

	err = agg.OnCreatedApp(app.CreatedAppEvent{AppId: "app1"})
	assert.Nil(t, err)

	err = agg.OnCreatedApp(app.CreatedAppEvent{AppId: "app2"})
	assert.Nil(t, err)

	exists2, err := agg.ExistsApp("app2")
	assert.Nil(t, err)
	assert.True(t, exists2)

	err = agg.OnRemovedApp(app.RemovedAppEvent{AppId: "app1"})
	assert.Nil(t, err)

	exists1, err = agg.ExistsApp("app1")
	assert.Nil(t, err)
	assert.False(t, exists1)
}

func TestMix(t *testing.T) {
	db := openBoltDB()
	defer closeBoltDB(db)

	agg := app.NewAggregate(db)

	exists1, err := agg.ExistsApp("app1")
	assert.Nil(t, err)
	assert.False(t, exists1)

	exists2, err := agg.ExistsApp("app2")
	assert.Nil(t, err)
	assert.False(t, exists2)

	event, err := agg.CreateApp("app1")
	assert.Equal(t, app.CreatedAppEvent{AppId: "app1"}, event)
	assert.Nil(t, err)

	exists1, err = agg.ExistsApp("app1")
	assert.Nil(t, err)
	assert.True(t, exists1)

	err = agg.OnCreatedApp(app.CreatedAppEvent{AppId: "app2"})
	assert.Nil(t, err)

	exists2, err = agg.ExistsApp("app2")
	assert.Nil(t, err)
	assert.True(t, exists2)

	err = agg.RemoveApp("app2")
	assert.Nil(t, err)

	exists2, err = agg.ExistsApp("app2")
	assert.Nil(t, err)
	assert.False(t, exists2)
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
