package algo

import (
	"github.com/lithammer/shortuuid"
)

// UniqueID ...
func UniqueID(val string) string {
	// URL safe uid
	uid := shortuuid.NewWithNamespace(val)
	return uid
}
