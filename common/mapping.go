package common

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/json"
	"github.com/Supernomad/quantum/ecdh"
	"net"
)

type Mapping struct {
	StringAddress string

	PublicKey []byte
	Address   *net.UDPAddr `json:"-"`
	SecretKey []byte       `json:"-"`
	Cipher    cipher.AEAD  `json:"-"`
}

func (m *Mapping) String() string {
	buf, _ := json.Marshal(m)
	return string(buf)
}

func ParseMapping(data string, privkey []byte) (*Mapping, error) {
	var mapping Mapping
	json.Unmarshal([]byte(data), &mapping)

	addr, err := net.ResolveUDPAddr("udp", mapping.StringAddress)
	if err != nil {
		return nil, err
	}

	mapping.Address = addr
	mapping.SecretKey = ecdh.GenerateSharedSecret(mapping.PublicKey, privkey)

	block, err := aes.NewCipher(mapping.SecretKey)
	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	mapping.Cipher = aesgcm

	return &mapping, nil
}

func NewMapping(address string, pubkey []byte) *Mapping {
	return &Mapping{
		StringAddress: address,
		PublicKey:     pubkey,
	}
}
