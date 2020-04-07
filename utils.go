package alipay

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"strings"
)

const (
	pemDelim              = "-----"
	pemTagBegin           = pemDelim + "BEGIN"
	pemTagEnd             = pemDelim + "END"
	pemLabelRSAPrivateKey = "RSA PRIVATE KEY"
	pemLabelPublicKey     = "PUBLIC KEY"
)

func parseRSAPublicKey(pemData string) (key *rsa.PublicKey, err error) {
	block, err := decodePEMData(pemData, pemLabelPublicKey)
	if err != nil {
		return
	}
	itf, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return
	}
	key, ok := itf.(*rsa.PublicKey)
	if !ok {
		err = errors.New("alipay: unknown type of public key")
	}
	return
}

func parseRSAPrivateKey(pemData string) (key *rsa.PrivateKey, err error) {
	block, err := decodePEMData(pemData, pemLabelRSAPrivateKey)
	if err != nil {
		return
	}
	key, err = x509.ParsePKCS1PrivateKey(block.Bytes)
	return
}

func decodePEMData(data string, label string) (block *pem.Block, err error) {
	standardizedData := standardizePEMData(data, label)
	block, _ = pem.Decode([]byte(standardizedData))
	if block == nil {
		err = errors.New("alipay: failed to decode PEM block")
	}
	return
}

func standardizePEMData(data string, label string) string {
	standardizedData := strings.TrimSpace(data)
	if !strings.Contains(standardizedData, pemTagBegin) {
		standardizedData = pemTagBegin + " " + label + " " + pemDelim + "\n" + standardizedData
	}
	if !strings.Contains(standardizedData, pemTagEnd) {
		standardizedData += "\n" + pemTagEnd + " " + label + " " + pemDelim
	}
	return standardizedData
}
