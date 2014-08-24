package app

import (
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/puffinframework/event"
)

const (
	appsBucketName  string     = "PuffinApps"
	CreatedAppEvent event.Type = "CreatedAppEvent"
	RemovedAppEvent event.Type = "RemovedAppEvent"
)

type Aggregate struct {
	db *bolt.DB
}

func NewCreatedAppEvent(appId string) event.Event {
	return event.NewEvent(CreatedAppEvent, 1, appId)
}

func NewRemovedAppEvent(appId string) event.Event {
	return event.NewEvent(RemovedAppEvent, 1, appId)
}

func NewAggregate(db *bolt.DB) *Aggregate {
	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(appsBucketName))
		return err
	})
	return &Aggregate{db: db}
}

func (self *Aggregate) CreateApp(appId string) (evt event.Event, err error) {
	self.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(appsBucketName))
		evt, err = createApp(b, appId)
		if err != nil {
			return err
		}
		return onCreatedAppEvent(b, evt)
	})
	return
}

func (self *Aggregate) OnCreatedApp(evt event.Event) error {
	return self.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(appsBucketName))
		return onCreatedAppEvent(b, evt)
	})
}

func createApp(b *bolt.Bucket, appId string) (evt event.Event, err error) {
	if existsApp(b, appId) {
		err = fmt.Errorf("ID already exists")
		return
	}
	evt = NewCreatedAppEvent(appId)
	return
}

func onCreatedAppEvent(b *bolt.Bucket, evt event.Event) error {
	appId := evt.Data().(string)
	return b.Put([]byte(appId), []byte{1})
}

func (self *Aggregate) RemoveApp(appId string) (evt event.Event, err error) {
	self.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(appsBucketName))
		evt, err = removeApp(b, appId)
		if err != nil {
			return err
		}
		return onRemovedAppEvent(b, evt)
	})
	return
}

func (self *Aggregate) OnRemovedApp(evt event.Event) error {
	return self.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(appsBucketName))
		return onRemovedAppEvent(b, evt)
	})
}

func removeApp(b *bolt.Bucket, appId string) (evt event.Event, err error) {
	if !existsApp(b, appId) {
		err = fmt.Errorf("ID does not exist")
		return
	}
	evt = NewRemovedAppEvent(appId)
	return
}

func onRemovedAppEvent(b *bolt.Bucket, evt event.Event) error {
	appId := evt.Data().(string)
	return b.Delete([]byte(appId))
}
func (self *Aggregate) ExistsApp(appId string) (exists bool, err error) {
	err = self.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(appsBucketName))
		exists = existsApp(b, appId)
		return nil
	})
	return
}

func existsApp(b *bolt.Bucket, appId string) bool {
	v := b.Get([]byte(appId))
	if v != nil {
		return true
	}
	return false
}
