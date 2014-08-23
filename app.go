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

		v := b.Get([]byte(appId))
		if v != nil {
			return fmt.Errorf("ID already exists")
		}

		event := CreatedAppEvent{AppId: appId}
		return self.onCreatedAppEvent(b, event)
	})
}

func (self *Aggregate) onCreatedAppEvent(b *bolt.Bucket, event CreatedAppEvent) error {
	return b.Put([]byte(event.AppId), []byte{1})
}

func (self *Aggregate) RemoveApp(appId string) error {
	return self.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(appsBucketName))
		if err != nil {
			return err
		}

		event := RemovedAppEvent{AppId: appId}
		return self.onRemovedAppEvent(b, event)
	})
}

func (self *Aggregate) onRemovedAppEvent(b *bolt.Bucket, event RemovedAppEvent) error {
	return b.Delete([]byte(event.AppId))
}
