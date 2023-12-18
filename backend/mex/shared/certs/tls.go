package certs

import (
	"context"
	//nolint: gosec
	"crypto/sha1"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"strings"

	L "github.com/d4l-data4life/mex/mex/shared/log"
)

const certTemplate = `
- certificate s/n: 0x%x, version: 0x%02x
  * SHA-1   fingerprint: %s
  * SHA-256 fingerprint: %s
  * Subject: %s
  * Issuer : %s
  * Is CA  : %v`

func GetTLSClientConfig(ctx context.Context, log L.Logger, certReaders []io.Reader) (*x509.CertPool, error) {
	rootCAs, err := x509.SystemCertPool()
	if err != nil {
		return nil, err
	}

	if certReaders == nil {
		return rootCAs, nil
	}

	if len(certReaders) == 0 {
		return rootCAs, nil
	}

	additionalCerts, err := readCertificates(certReaders)
	if err != nil {
		return nil, err
	}

	for _, cert := range additionalCerts {
		rootCAs.AddCert(cert)

		//nolint: gosec
		s1 := sha1.Sum(cert.Raw)
		s256 := sha256.Sum256(cert.Raw)
		log.Info(ctx, L.Messagef(certTemplate,
			cert.SerialNumber.Int64(), cert.Version,
			hex(s1[:], ":"), hex(s256[:], ":"),
			cert.Subject.String(), cert.Issuer.String(),
			cert.IsCA,
		))
	}

	return rootCAs, nil
}

func readCertificates(readers []io.Reader) ([]*x509.Certificate, error) {
	totalBlocks := []*pem.Block{}
	for _, r := range readers {
		blocks, err := readPEMBlocks(r)
		if err != nil {
			return nil, err
		}
		totalBlocks = append(totalBlocks, blocks...)
	}

	certs := []*x509.Certificate{}
	for _, block := range totalBlocks {
		if block.Type != "CERTIFICATE" {
			continue
		}

		cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return nil, err
		}

		certs = append(certs, cert)
	}

	return certs, nil
}

func readPEMBlocks(r io.Reader) ([]*pem.Block, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	var block *pem.Block
	blocks := []*pem.Block{}

	for {
		block, data = pem.Decode(data)
		if block == nil {
			break
		}

		blocks = append(blocks, block)
	}

	return blocks, nil
}

func hex(hash []byte, sep string) string {
	s := make([]string, len(hash))

	for i, b := range hash {
		s[i] = fmt.Sprintf("%02X", b)
	}

	return strings.Join(s, sep)
}
