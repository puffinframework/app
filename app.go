package app

import (
	"fmt"
	"github.com/boltdb/bolt"
)

const (
	appsBucketName string = "PuffinApps"
)

type Aggregate struct {
	db *bolt.DB
}

type CreatedAppEvent struct {
	AppId string
}

type RemovedAppEvent struct {
	AppId string
}

func NewAggregate(db *bolt.DB) *Aggregate {
	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(appsBucketName))
		return err
	})
	return &Aggregate{db: db}
}

func (self *Aggregate) CreateApp(appId string) (event CreatedAppEvent, err error) {
	self.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(appsBucketName))
		event, err = createApp(b, appId)
		if err != nil {
			return err
		}
		return onCreatedAppEvent(b, event)
	})
	return
}

func (self *Aggregate) OnCreatedApp(event CreatedAppEvent) error {
	return self.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(appsBucketName))
		return onCreatedAppEvent(b, event)
	})
}

func createApp(b *bolt.Bucket, appId string) (event CreatedAppEvent, err error) {
	if existsApp(b, appId) {
		err = fmt.Errorf("ID already exists")
		return
	}
	event = CreatedAppEvent{AppId: appId}
	return
}

func onCreatedAppEvent(b *bolt.Bucket, event CreatedAppEvent) error {
	return b.Put([]byte(event.AppId), []byte{1})
}

func (self *Aggregate) RemoveApp(appId string) (event RemovedAppEvent, err error) {
	self.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(appsBucketName))
		event, err = removeApp(b, appId)
		if err != nil {
			return err
		}
		return onRemovedAppEvent(b, event)
	})
	return
}

func (self *Aggregate) OnRemovedApp(event RemovedAppEvent) error {
	return self.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(appsBucketName))
		return onRemovedAppEvent(b, event)
	})
}

func removeApp(b *bolt.Bucket, appId string) (event RemovedAppEvent, err error) {
	if !existsApp(b, appId) {
		err = fmt.Errorf("ID does not exist")
		return
	}
	event = RemovedAppEvent{AppId: appId}
	return
}

func onRemovedAppEvent(b *bolt.Bucket, event RemovedAppEvent) error {
	return b.Delete([]byte(event.AppId))
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
