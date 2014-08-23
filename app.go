package app

import (
    "fmt"
)

type Domain struct {
    ids map[string]bool
}

func NewDomain() *Domain {
    return &Domain{ids: make(map[string]bool)}
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