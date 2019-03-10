package main

import (
	"bufio"
	"io/ioutil"
	"log"
	"os"

	"github.com/gorilla/websocket"
)

var uri = "wss://eqsjdaluec.execute-api.us-west-1.amazonaws.com/dev"

func main() {
	con, res, err := websocket.DefaultDialer.Dial(uri, nil)
	if err != nil {
		b, _ := ioutil.ReadAll(res.Body)
		log.Print("resp body: ", string(b))
		log.Print("res: ", res)
		log.Fatal("err: ", err)
	}

	go func() {
		for {
			_, b, err := con.ReadMessage()
			if err != nil {
				log.Print("read err")
				log.Fatal(err)
			}
			log.Print(string(b))
		}
	}()

	scan := bufio.NewScanner(os.Stdin)
	for scan.Scan() {
		log.Print("sending ")
		err = con.WriteMessage(1, scan.Bytes())
		if err != nil {
			log.Fatal(err)
		}
	}

}
