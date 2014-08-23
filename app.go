package app

import (
	"fmt"
    "os"
	"github.com/boltdb/bolt"
)

type Domain struct {
    db *bolt.DB
	ids map[string]bool
}

func NewDomain() *Domain {
	db, err := bolt.Open("app.db", 0600, nil)
	if err != nil {
		panic(err)
	}
    return &Domain{ db: db, ids: make(map[string]bool)}
}

func (self *Domain) Close() {
    self.db.Close()
    os.Remove("app.db") // TODO
}

type CreatedAppEvent struct {
	Id string
}

func (self *Domain) Create(id string) error {
	if self.ids[id] {
		return fmt.Errorf("ID already exists")
	}
	event := CreatedAppEvent{Id: id}
	self.OnCreatedAppEvent(event)
	return nil
}

func (self *Domain) OnCreatedAppEvent(event CreatedAppEvent) {
	self.ids[event.Id] = true
}

type RemovedAppEvent struct {
	Id string
}

func (self *Domain) Remove(id string) error {
	if self.ids[id] != true {
		return fmt.Errorf("ID does not exist")
	}
	event := RemovedAppEvent{Id: id}
	self.OnRemovedAppEvent(event)
	return nil
}

func (self *Domain) OnRemovedAppEvent(event RemovedAppEvent) {
	delete(self.ids, event.Id)
}
