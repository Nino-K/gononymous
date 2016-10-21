package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

func main() {
	addr := "localhost:9988"
	u := url.URL{Scheme: "ws", Host: addr, Path: "/sessionHandle"}

	header := http.Header{}
	clientID := strconv.FormatInt(time.Now().UnixNano(), 10)
	fmt.Println("clientID", clientID)
	header.Add("CLIENT_ID", clientID)
	c, _, err := websocket.DefaultDialer.Dial(u.String(), header)
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
			fmt.Printf("recv: %s\n", message)
		}
	}()

	reader := bufio.NewReader(os.Stdin)
	fmt.Print(clientID + ": ")
	for {
		text, _ := reader.ReadString('\n')
		err := c.WriteMessage(websocket.TextMessage, []byte(text))
		if err != nil {
			fmt.Println("write:", err)
			return
		}
		time.Sleep(time.Second * 5)
	}

}
