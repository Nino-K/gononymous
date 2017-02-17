package cert

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"fmt"
	"log"
	"math/big"
	"net"
	"os"
	"time"
)

type Generator struct {
	Addrs   []net.IP
	OutPath string
}

func (g *Generator) GenerateSrvCertKey() ([]byte, []byte, error) {
	srvKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatalf("generating server key: %v", err)
	}
	serverCertTemplate, err := certTemplate()
	if err != nil {
		log.Fatalf("creating server template: %v", err)
	}
	//serverCertTemplate.IsCA = true
	serverCertTemplate.KeyUsage = x509.KeyUsageCertSign | x509.KeyUsageDigitalSignature
	serverCertTemplate.ExtKeyUsage = []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth}
	serverCertTemplate.IPAddresses = g.Addrs

	certPEM, err := createCert(serverCertTemplate, serverCertTemplate, &srvKey.PublicKey, srvKey)
	if err != nil {
		return nil, nil, err
	}
	key := pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(srvKey)}
	keyPEM := pem.EncodeToMemory(&key)

	err = g.outWriter(certPEM, "cert.pem")
	if err != nil {
		return nil, nil, err
	}
	err = g.outWriter(keyPEM, "key.pem")
	if err != nil {
		return nil, nil, err
	}

	return certPEM, keyPEM, nil
}

func (g *Generator) outWriter(encodedPEM []byte, name string) error {
	path := g.OutPath + "/" + name
	f, err := os.Create(path)
	if err != nil {
		fmt.Println(err)
		return errors.New(fmt.Sprintf("OutWriter creating file %s: %v", path, err))
	}
	defer f.Close()
	_, err = f.Write(encodedPEM)
	if err != nil {
		return err
	}
	return nil
}

func certTemplate() (*x509.Certificate, error) {
	// generate a random serial number (a real cert authority would have some logic behind this)
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return nil, errors.New("failed to generate serial number: " + err.Error())
	}

	tmpl := x509.Certificate{
		SerialNumber:          serialNumber,
		Subject:               pkix.Name{Organization: []string{"gononymous"}},
		SignatureAlgorithm:    x509.SHA256WithRSA,
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(time.Hour), //TODO we need to fix this, this is only valid for an hour
		BasicConstraintsValid: true,
	}
	return &tmpl, nil
}

func createCert(template, parent *x509.Certificate, pub interface{}, parentPriv interface{}) (certPEM []byte, err error) {
	certDER, err := x509.CreateCertificate(rand.Reader, template, parent, pub, parentPriv)
	if err != nil {
		return
	}
	// PEM encode the certificate (this is a standard TLS encoding)
	b := pem.Block{Type: "CERTIFICATE", Bytes: certDER}
	certPEM = pem.EncodeToMemory(&b)
	return
}
