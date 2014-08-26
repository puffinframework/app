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

	createdAppEvent, err := agg.CreateApp("app1")
	assert.Equal(t, app.CreatedAppEvent, createdAppEvent.Header().Type)
	assert.Equal(t, "app1", createdAppEvent.Data())
	assert.Nil(t, err)

	exists1, err = agg.ExistsApp("app1")
	assert.Nil(t, err)
	assert.True(t, exists1)

	_, err = agg.CreateApp("app1")
	assert.NotNil(t, err)

	createdAppEvent, err = agg.CreateApp("app2")
	assert.Equal(t, app.CreatedAppEvent, createdAppEvent.Header().Type)
	assert.Equal(t, "app2", createdAppEvent.Data())
	assert.Nil(t, err)

	exists2, err := agg.ExistsApp("app2")
	assert.Nil(t, err)
	assert.True(t, exists2)
}

func TestRemove(t *testing.T) {
	db := openBoltDB()
	defer closeBoltDB(db)

	agg := app.NewAggregate(db)

	_, err := agg.RemoveApp("app1")
	assert.NotNil(t, err)

	exists1, err := agg.ExistsApp("app1")
	assert.Nil(t, err)
	assert.False(t, exists1)

	createdAppEvent, err := agg.CreateApp("app1")
	assert.Equal(t, app.CreatedAppEvent, createdAppEvent.Header().Type)
	assert.Equal(t, "app1", createdAppEvent.Data())
	assert.Nil(t, err)

	removedAppEvent, err := agg.RemoveApp("app1")
	assert.Equal(t, app.RemovedAppEvent, removedAppEvent.Header().Type)
	assert.Equal(t, "app1", removedAppEvent.Data())
	assert.Nil(t, err)
}

func TestProcessEvents(t *testing.T) {
	db := openBoltDB()
	defer closeBoltDB(db)

	agg := app.NewAggregate(db)

	exists1, err := agg.ExistsApp("app1")
	assert.Nil(t, err)
	assert.False(t, exists1)

	err = agg.OnCreatedApp(app.NewCreatedAppEvent("app1"))
	assert.Nil(t, err)

	exists1, err = agg.ExistsApp("app1")
	assert.Nil(t, err)
	assert.True(t, exists1)

	err = agg.OnCreatedApp(app.NewCreatedAppEvent("app1"))
	assert.Nil(t, err)

	err = agg.OnCreatedApp(app.NewCreatedAppEvent("app2"))
	assert.Nil(t, err)

	exists2, err := agg.ExistsApp("app2")
	assert.Nil(t, err)
	assert.True(t, exists2)

	err = agg.OnRemovedApp(app.NewCreatedAppEvent("app1"))
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

	createdAppEvent, err := agg.CreateApp("app1")
	assert.Equal(t, app.CreatedAppEvent, createdAppEvent.Header().Type)
	assert.Equal(t, "app1", createdAppEvent.Data())
	assert.Nil(t, err)

	exists1, err = agg.ExistsApp("app1")
	assert.Nil(t, err)
	assert.True(t, exists1)

	err = agg.OnCreatedApp(app.NewCreatedAppEvent("app2"))
	assert.Nil(t, err)

	exists2, err = agg.ExistsApp("app2")
	assert.Nil(t, err)
	assert.True(t, exists2)

	removedAppEvent, err := agg.RemoveApp("app2")
	assert.Equal(t, app.RemovedAppEvent, removedAppEvent.Header().Type)
	assert.Equal(t, "app2", removedAppEvent.Data())
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
	defer os.Remove(db.Path())
	db.Close()
}
