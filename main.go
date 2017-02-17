package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/Nino-K/gononymous/cert"
	"github.com/Nino-K/gononymous/handler"
	"github.com/Nino-K/gononymous/server"
	"github.com/gorilla/websocket"
)

var (
	port = flag.Int("port", 9797, "Port that server is listening on, default is 9797")
	addr = flag.String("addr", "127.0.0.1", "Address that server is listening on, default is 127.0.0.1")
)

func main() {
	flag.Parse()

	if *port < 1024 || *port > 65535 {
		fmt.Fprintf(os.Stderr, "-port must be within range 1024-65535")
		os.Exit(1)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan,
		os.Interrupt,
		syscall.SIGINT,
		syscall.SIGKILL,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	serverAddr := fmt.Sprintf("%s:%d", *addr, *port)
	go func() {
		fmt.Printf("gononymous is listening on %s \n", serverAddr)
		err := start(serverAddr)
		if err != nil {
			log.Fatalf("start: %s\n", err)
			return
		}
	}()
	for {
		select {
		case <-sigChan:
			os.Exit(0)
		}
	}
}

func start(srvAddr string) error {
	upgrader := websocket.Upgrader{}
	sessonManager := server.NewSessionManager()
	sessionHandler := handler.NewSessionHandler(sessonManager, &upgrader)

	certPEM, keyPEM, err := cert.GenerateSrvCertKey()
	if err != nil {
		return err
	}
	serverTLSCert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		return err
	}
	tlsConfig := &tls.Config{Certificates: []tls.Certificate{serverTLSCert}}
	server := &http.Server{
		Addr:      srvAddr,
		TLSConfig: tlsConfig,
	}
	ln, err := tls.Listen("tcp", srvAddr, tlsConfig)
	if err != nil {
		return err
	}
	http.HandleFunc("/", sessionHandler.Join)
	return server.Serve(ln)
}
