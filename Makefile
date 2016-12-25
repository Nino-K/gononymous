default: test 

test:            
	go get -u github.com/onsi/ginkgo/ginkgo
	go get -u github.com/onsi/gomega
	go test -v -race ./...
