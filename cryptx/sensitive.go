package cryptx

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
)

var _DefaultEncrytor *SenInfoEncryptor

const _DefaultPw = "86c7d13317c7e8993656b7e093ee67a789350e0545bc073ca9dde341c98b1363"

func init() {
	var err error
	_DefaultEncrytor, err = NewFromStrKey(_DefaultPw)
	if err != nil {
		panic(err)
	}
}

func EncryptSenInfo(info string) string {
	return _DefaultEncrytor.Encrypt(info)
}

func DecryptSenInfo(encInfo string) (string, error) {
	return _DefaultEncrytor.Decrypt(encInfo)
}

type SenInfoEncryptor struct {
	cipher cipher.AEAD
}

func NewFromStrKey(keyStr string) (*SenInfoEncryptor, error) {
	key, err := hex.DecodeString(keyStr)
	if err != nil {
		return nil, err
	}
	return NewSenInfoEncryptor(key)
}

func NewSenInfoEncryptor(key []byte) (*SenInfoEncryptor, error) {

	if len(key) != 32 {
		return nil, errors.New("invalid key len")
	}

	aes256, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	c, err := cipher.NewGCM(aes256)
	if err != nil {
		return nil, err
	}

	p := &SenInfoEncryptor{cipher: c}

	return p, nil

}

func (p *SenInfoEncryptor) Encrypt(info string) string {
	// 对于同一个key，nonce要求唯一，但我们这里的应用场景是加密一些配置文件里面的敏感信息，数量比较小，随机冲突的几率几乎没有
	// 即使冲突了，问题也不大
	nonce := make([]byte, 12)
	_, _ = rand.Read(nonce)
	en := p.cipher.Seal(nil, nonce, []byte(info), nil)
	buf := make([]byte, 12+len(en))
	copy(buf, nonce)
	copy(buf[12:], en)
	return base64.URLEncoding.EncodeToString(buf)
}

func (p *SenInfoEncryptor) Decrypt(encInfo string) (string, error) {
	buf, err := base64.URLEncoding.DecodeString(encInfo)
	if err != nil {
		return "", err
	}
	if len(buf) < 12 {
		return "", errors.New("invalid.format")
	}
	nonce := buf[:12]
	r, err := p.cipher.Open(nil, nonce, buf[12:], nil)
	if err != nil {
		return "", err
	}

	return string(r), nil
}
