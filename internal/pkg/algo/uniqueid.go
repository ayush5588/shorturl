package algo

import (
	"github.com/lithammer/shortuuid"
)

// UniqueID ...
func UniqueID(val string) string {
	uid := shortuuid.NewWithNamespace(val)
	return uid
}
