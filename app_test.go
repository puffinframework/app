package app_test

import (
	"github.com/puffinframework/app"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreate(t *testing.T) {
	domain := app.NewDomain()
    defer domain.Close()

	err := domain.Create("app1")
	assert.Nil(t, err)

	err = domain.Create("app1")
	assert.NotNil(t, err)

	err = domain.Create("app2")
	assert.Nil(t, err)
}

func TestRemove(t *testing.T) {
	domain := app.NewDomain()
    defer domain.Close()

	err := domain.Remove("app1")
	assert.NotNil(t, err)

	err = domain.Create("app1")
	assert.Nil(t, err)

	err = domain.Remove("app1")
	assert.Nil(t, err)
}
