package main

import (
	"bufio"
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/gorilla/websocket"
)

var sessionId = flag.String("sessionid", "sessionID", "The unique id that is used for each session")
var clientId = flag.String("clientId", "clientId", "The unique id that is used to designate client")
var srvAddr = flag.String("addr", "localhost:9797", "The addrees:port of the server")

func main() {
	flag.Parse()

	u := url.URL{Scheme: "wss", Host: *srvAddr, Path: *sessionId}

	clientIDSuffix := time.Now().UnixNano()
	header := http.Header{}
	header.Add("CLIENT_ID", fmt.Sprintf("%s-%d", *clientId, clientIDSuffix))

	tlsConfig := &tls.Config{InsecureSkipVerify: true}
	dialer := websocket.Dialer{TLSClientConfig: tlsConfig}
	c, _, err := dialer.Dial(u.String(), header)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()
	go func() {
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				fmt.Println("read:", err)
				return
			}
			fmt.Println("\n" + string(message))
		}
	}()

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf(*clientId + ": ")
		text, _ := reader.ReadString('\n')
		err := c.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("[%s]: %s", *clientId, text)))
		if err != nil {
			fmt.Println("write:", err)
			return
		}
	}

}
