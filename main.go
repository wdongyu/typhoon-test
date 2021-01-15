package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

const (
	GatewayUrlPrefix = "http://localhost:31380/typhoon-backend?index="

	VersionHeader = "X-Version"

	TyphoonHeaderPrefix = "typhoon-microservices-typhoon-"

	DefaultInterval = 50
	DefaultObject = "safety"
	DefaultAlg = ""

)

func init() {
	log.SetFlags(log.Lmicroseconds)
}

func main() {
	var object, alg string
	var interval int
	flag.IntVar(&interval, "interval", DefaultInterval, `interval to send http request`)
	flag.StringVar(&object, "object", DefaultObject, `object of the test`)
	flag.StringVar(&alg, "alg", DefaultAlg, `algorithm will be used`)
	flag.Parse()

	var stop = make(chan bool)
	turns := 80
	for i := 1; i <= turns; i++ {
		log.Printf("#%d.\n", i)
		go sendReq(i)
		if i == 20 {
			if object == "safety" && alg == "" {
				go CanarypdateReq()
			} else if alg == "Quiescence" {
				go QuieUpdateReq()
			} else if alg == "ManagedService" {
				go MsUpdateReq()
			}
		}
		time.Sleep(time.Duration(interval)* time.Millisecond)
	}

	<-stop
}

func sendReq(index int) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s%d", GatewayUrlPrefix, index), nil)
	if err != nil {
		log.Printf("%d. Fail to create http request : %v\n", index, err)
		return
	}
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Printf("%d. Fail to send http request : %v\n", index, err)
		return
	}

	log.Printf("%d. %s ", index, res.Status)
	typhoonHeader := ""
	for _, value := range res.Header[VersionHeader] {
		if strings.HasPrefix(value, TyphoonHeaderPrefix) {
			if typhoonHeader == "" {
				typhoonHeader = value
			} else {
				if typhoonHeader != value {
					log.Printf("Version header %s : %s conflict ...\n", typhoonHeader, value)
				}
			}
		}
	}
	log.Printf("%d. Response : %s\n", index, typhoonHeader)
}
