default: test 

test:            
	go fmt ./...
	go get -u github.com/onsi/ginkgo/ginkgo
	go get -u github.com/onsi/gomega
	ginkgo -r -race -randomizeAllSpecs 

certs:
	go run /usr/local/go/src/crypto/tls/generate_cert.go --host localhost
	mv *.pem cert/
