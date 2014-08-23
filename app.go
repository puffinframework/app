package app

import (
	"fmt"
	"github.com/boltdb/bolt"
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
	return &Aggregate{db: db}
}

const (
	appsBucketName = "PuffinApps"
)

func (self *Aggregate) CreateApp(appId string) error {
	return self.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(appsBucketName))
		if err != nil {
			return err
		}
		event, err := createApp(b, appId)
		if err != nil {
			return err
		}
		return onCreatedAppEvent(b, event)
	})
}

func (self *Aggregate) OnCreatedApp(event CreatedAppEvent) error {
	return self.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(appsBucketName))
		if err != nil {
			return err
		}
		return onCreatedAppEvent(b, event)
	})
}

func createApp(b *bolt.Bucket, appId string) (event CreatedAppEvent, err error) {
	v := b.Get([]byte(appId))
	if v != nil {
		err = fmt.Errorf("ID already exists")
		return
	}
	event = CreatedAppEvent{AppId: appId}
	return
}

func onCreatedAppEvent(b *bolt.Bucket, event CreatedAppEvent) error {
	return b.Put([]byte(event.AppId), []byte{1})
}

func (self *Aggregate) RemoveApp(appId string) error {
	return self.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(appsBucketName))
		if err != nil {
			return err
		}
		event, err := removeApp(b, appId)
		if err != nil {
			return err
		}
		return onRemovedAppEvent(b, event)
	})
}

func (self *Aggregate) OnRemovedApp(event RemovedAppEvent) error {
	return self.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(appsBucketName))
		if err != nil {
			return err
		}
		return onRemovedAppEvent(b, event)
	})
}

func removeApp(b *bolt.Bucket, appId string) (event RemovedAppEvent, err error) {
	event = RemovedAppEvent{AppId: appId}
	return
}

func onRemovedAppEvent(b *bolt.Bucket, event RemovedAppEvent) error {
	return b.Delete([]byte(event.AppId))
}
