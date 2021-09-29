package random

import (
	/* #nosec */
	"crypto/md5"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/sha3"
)

//MD5UUID new uuid first, use md5 hash
func MD5UUID() string {
	id := uuid.NewV1()
	/* #nosec */
	h := md5.New()
	return writeHex(h, id[:])
}

//SHA256UUID new uuid first, use sha256 hash
func SHA256UUID() string {
	id := uuid.NewV1()
	h := sha3.New256()
	return writeHex(h, id[:])
}
