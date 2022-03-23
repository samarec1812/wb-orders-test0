package main

import (
	"bytes"
	"encoding/json"
	"github.com/nats-io/stan.go"
	orders "github.com/samarec1812/wb-orders-test0"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
	"sync"
	"time"
)

func main() {

	sc, _ := stan.Connect("prod", "sub")
	defer sc.Close()
	//var order orders.Orders
	_, err := sc.Subscribe("samarec", func(m *stan.Msg) {
		var order orders.Orders
		err := json.Unmarshal(m.Data, &order)
		if err != nil {
			logrus.Errorf("error model")
		} else {
			url := "http://127.0.0.1:8089/api/send"
			_, err = http.Post(url, "application/json", bytes.NewBuffer(m.Data))
			if err != nil {
				logrus.Errorf("error create %s", err.Error())
			}
		}
	}, stan.AckWait(5* time.Second), stan.DurableName("sub-1"))

	if err != nil {
		log.Fatalf("error with subscribe %s", err.Error())
		return
	}

	w := sync.WaitGroup{}
	w.Add(1)
	w.Wait()
}
