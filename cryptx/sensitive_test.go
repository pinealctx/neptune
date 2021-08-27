package cryptx

import (
	"fmt"
	"testing"
)

func TestSenInfoEncryptor(t *testing.T) {

	passwordBeforeEncrypt := "myIQvhRYZVEha3vy"

	enP := EncryptSenInfo(passwordBeforeEncrypt)
	fmt.Printf("encryt.password:\n%s\n", enP)

}

func TestSenInfoDecryptor(t *testing.T) {

	passwordBeforeDecrypt := "Ue86TDr7Hrbs_c1XWzijFvZzEUNLq_Y_Ya0VRw=="

	deP, err := DecryptSenInfo(passwordBeforeDecrypt)
	if err != nil {
		panic(err)
	}
	fmt.Printf("decrypt.password:\n%s\n", deP)

}

// Ue86TDr7Hrbs_c1XWzijFvZzEUNLq_Y_Ya0VRw==
