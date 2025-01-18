package id

import (
	"crypto"
	"crypto/ecdh"
	"crypto/rand"
	"crypto/rsa"
	"fmt"
)

type GenerateKeyOpts struct {
	Type string
	Size *int
}

func GenerateKey(opt GenerateKeyOpts) (identity KeyPackageIdentity, err error) {
	var pk crypto.PrivateKey
	switch opt.Type {
	case "rsa":
		keySize := 2048
		if opt.Size != nil {
			keySize = *opt.Size
		}
		pk, err = rsa.GenerateKey(rand.Reader, keySize)
		if err != nil {
			return nil, err
		}
		identity = createRSAIdKey(pk.(*rsa.PrivateKey))
	case "ec25519":
		pk, err = ecdh.X25519().GenerateKey(rand.Reader)
		if err != nil {
			return nil, err
		}
		identity = createECDHIdKey(pk.(*ecdh.PrivateKey))
	case "p256":
		pk, err = ecdh.P256().GenerateKey(rand.Reader)
		if err != nil {
			return nil, err
		}
		identity = createECDHIdKey(pk.(*ecdh.PrivateKey))
	case "p384":
		pk, err = ecdh.P384().GenerateKey(rand.Reader)
		if err != nil {
			return nil, err
		}
		identity = createECDHIdKey(pk.(*ecdh.PrivateKey))
	case "p521":
		pk, err = ecdh.P521().GenerateKey(rand.Reader)
		if err != nil {
			return nil, err
		}
		identity = createECDHIdKey(pk.(*ecdh.PrivateKey))
	default:
		return nil, fmt.Errorf("unsupported key type: %s", opt.Type)
	}

	return identity, nil
}
