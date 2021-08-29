package random

import (
	/* #nosec */
	"crypto/md5"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/sha3"
)

func MD5UUID() string {
	id := uuid.NewV1()
	/* #nosec */
	h := md5.New()
	return writeHex(h, id[:])
}

func SHA256UUID() string {
	id := uuid.NewV1()
	h := sha3.New256()
	return writeHex(h, id[:])
}
