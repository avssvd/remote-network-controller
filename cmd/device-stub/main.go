package main

import (
	"fmt"
	"github.com/reiver/go-oi"
	"github.com/reiver/go-telnet"
	"log"
	"time"
)

func main() {

	//var handler telnet.Handler = telnet.EchoHandler
	var handler telnet.Handler = StubDeviceTelnetHandler{}

	err := telnet.ListenAndServe(":5555", handler)
	if nil != err {
		//@TODO: Handle this error better.
		panic(err)
	}
}

type StubDeviceTelnetHandler struct {}

func (handler StubDeviceTelnetHandler) ServeTELNET(ctx telnet.Context, w telnet.Writer, r telnet.Reader) {
	clientChan := make(chan []byte)
	go GetClientChan(r, clientChan)
	spamChan := make(chan []byte)
	go GetSpamChan(spamChan)


	for {
		select {
		case msg := <- clientChan:
			_, err := oi.LongWrite(w, []byte(fmt.Sprintf("Answer: %s\n", string(msg))))
			if err != nil {
				log.Printf("ServeTELNET clientChan: %s/n", err.Error())
			}

		case spam := <- spamChan:
			_, err := oi.LongWrite(w, spam)
			if err != nil {
				log.Printf("ServeTELNET spamChan: %s/n", err.Error())
			}
		}
	}
}

func GetClientChan(r telnet.Reader, ch chan []byte) {
	var buffer [1]byte // Seems like the length of the buffer needs to be small, otherwise will have to wait for buffer to fill up.
	p := buffer[:]
	for {
		n, err := r.Read(p)
		if err != nil {
			log.Printf("GetClientChan: %s\n", err.Error())
			continue
		}
		if n > 0 {
			ch <- p
		}
	}
}

func GetSpamChan(ch chan []byte) {
	for {
		ch <- []byte(time.Now().String())
		time.Sleep(3*time.Second)
	}
}
