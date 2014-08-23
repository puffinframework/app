package app_test

import (
    "github.com/puffinframework/app"
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {
    domain := app.NewDomain()

    err := domain.Create("app1")
    assert.Nil(t, err)

    err = domain.Create("app1")
    assert.NotNil(t, err)

    err = domain.Create("app2")
    assert.Nil(t, err)
}

func TestRemove(t *testing.T) {
    domain := app.NewDomain()

    err := domain.Remove("app1")
    assert.NotNil(t, err)

    err = domain.Create("app1")
    assert.Nil(t, err)

    err = domain.Remove("app1")
    assert.Nil(t, err)
}
