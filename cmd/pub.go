package main

import (
	"fmt"
	stan "github.com/nats-io/stan.go"
	"io/ioutil"
	"log"
)

func main() {
	sc, _ := stan.Connect("prod", "simple-pub")
	defer sc.Close()

	for i := 1; i < 3; i++ {
		msg, err := ioutil.ReadFile(fmt.Sprintf("example/model%d.json", i))
		if err != nil {
			log.Fatalf("error reading file %s", err.Error())
		}
		sc.Publish("samarec", msg)
		fmt.Println(string(msg))
	}

	/*for i := 0;; i++ {
		if i == 5 {
			break
		}
		sc.Publish("samarec", msg)
		time.Sleep(2 * time.Second)
	}*/
}
