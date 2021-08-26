package random

import (
	/* #nosec */
	"crypto/md5"
	uuid "github.com/satori/go.uuid"
)

func MD5UUID() string {
	id := uuid.NewV1()
	/* #nosec */
	h := md5.New()
	return writeHex(h, id[:])
}
