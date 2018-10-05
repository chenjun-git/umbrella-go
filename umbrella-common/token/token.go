package token

import (
	"crypto/ecdsa"
	"crypto/md5"
	"crypto/rand"
	"crypto/x509"
	"encoding/base64"
	"encoding/binary"
	"encoding/pem"
	"errors"
	"math/big"
	"sync"
)

type counter struct {
	sync.Mutex
	num uint32
}

var count counter

type Token struct {
	IssueTime uint32
	TTL       uint16
	UserID    string

	Mask1 int64 // not used in v1, v2
	Mask2 int64 // not userd in v1, v2
}

var privateKey *ecdsa.PrivateKey
var publicKeys []*ecdsa.PublicKey

func packLeadingZero32(bs []byte) []byte {
	if len(bs) >= 32 {
		return bs
	}

	packBytes := make([]byte, 32-len(bs))

	return append(packBytes, bs...)
}

func unpackLeadingZero(bs []byte) []byte {
	i := 0
	for i < len(bs) && bs[i] == 0 {
		i++
	}

	return bs[i:]
}

func InitPrivateKey(key []byte) {
	block, _ := pem.Decode(key)
	if block == nil {
		panic("private key invalid")
	}

	prvk, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		panic(err)
	}
	privateKey = prvk
}

func InitPublicKeys(keys ...[]byte) {
	if len(keys) == 0 {
		panic("no public key")
	}

	publicKeys = make([]*ecdsa.PublicKey, 0)
	for _, key := range keys {
		block, _ := pem.Decode(key)
		if block == nil {
			panic("public key invalid")
		}

		pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			panic(err)
		}

		pubKey := pubInterface.(*ecdsa.PublicKey)
		publicKeys = append(publicKeys, pubKey)
	}
}

func EncryptAccessToken(version int, tk *Token) (string, error) {
	count.Lock()
	seq := count.num
	count.num += 1
	count.Unlock()

	switch version {
	case 1:
		data, err := tk.encryptV1(seq)
		if err != nil {
			return "", err
		}

		token := base64.URLEncoding.EncodeToString(data)
		return token, nil
	default:
		return "", errors.New("invalid version")
	}
}

func DecryptAccessToken(token string) (*Token, error) {
	if len(token) == 0 {
		return nil, errors.New("empty token")
	}

	data, err := base64.URLEncoding.DecodeString(token)
	if err != nil {
		return nil, errors.New("invalid base64 token")
	}

	switch int(data[0]) {
	case 1:
		token := &Token{}
		err := token.decryptV1(data)
		if err != nil {
			return nil, err
		}
		return token, nil
	default:
		return nil, errors.New("invalid version")
	}
}

func GetTokenVersion(token string) (int, error) {
	data, err := base64.URLEncoding.DecodeString(token)
	if err != nil {
		return 0, errors.New("invalid base64 token")
	}

	return int(data[0]), nil
}

func (t *Token) encryptV1(seq uint32) ([]byte, error) {
	var datas = make([]byte, 10, 106)

	datas[0] = 0x01                                // version
	binary.LittleEndian.PutUint32(datas[1:5], seq) // 只要3个字节;下面覆盖了1个字节
	binary.LittleEndian.PutUint32(datas[4:8], t.IssueTime)
	binary.LittleEndian.PutUint16(datas[8:10], uint16(t.TTL))

	datas = append(datas, []byte(t.UserID)...)

	h := md5.New()
	h.Write(datas)
	hashed := h.Sum(nil)
	r, s, err := ecdsa.Sign(rand.Reader, privateKey, hashed)
	if err != nil {
		return nil, err
	}

	datas = append(datas, packLeadingZero32(r.Bytes())...)
	datas = append(datas, packLeadingZero32(s.Bytes())...)

	return datas, nil
}

func (t *Token) decryptV1(data []byte) error {
	if len(data) < 106 {
		return errors.New("invalid token length")
	}
	h := md5.New()
	h.Write(data[:42])
	hashed := h.Sum(nil)
	r := big.NewInt(0)
	s := big.NewInt(0)
	r = r.SetBytes(unpackLeadingZero(data[42:74]))
	s = s.SetBytes(unpackLeadingZero(data[74:]))

	valid := false
	for _, pubk := range publicKeys {
		ok := ecdsa.Verify(pubk, hashed, r, s)
		if ok {
			valid = true
			break
		}
	}

	if !valid {
		return errors.New("sign verify failed")
	}

	t.IssueTime = uint32(binary.LittleEndian.Uint32(data[4:8]))
	t.TTL = uint16(binary.LittleEndian.Uint16(data[8:10]))
	t.UserID = string(data[10:42])

	return nil
}
