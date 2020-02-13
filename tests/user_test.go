package tests

import (
	"github.com/8luebottle/go-blog/api/models"
	"testing"
)

func TestUserValidate(t *testing.T)  {
	user := models.User{}
	if err := user.Validate("aa"); err != nil {
		t.Log(err)
	}
	t.Log("asdasd")
}