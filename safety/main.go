package main

import (
	"bytes"
	"flag"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

const (
	GatewayUrl = "http://localhost:31380/typhoon-backend?index=1"

	// UpdateUrl = "http://localhost:32088/apis/managedservice/typhoon-microservices-typhoon?namespace=typhoon"
	UpdateUrl = "http://localhost:32088/apis/virtualservice/typhoon-microservices-typhoon?namespace=typhoon"


	VersionHeader = "X-Version"

	DefaultInterval = 150

	TyphoonHeaderPrefix = "typhoon-microservices-typhoon-"
)

func init() {
	log.SetFlags(log.Lmicroseconds)
}

func main()  {
	var interval int
	flag.IntVar(&interval, "interval", DefaultInterval, `interval to send http request`)
	flag.Parse()

	var stop = make(chan bool)
	turns := 40
	for i := 1; i <= turns; i++ {
		log.Printf("#%d.\n", i)
		go sendReq()
		if i == 20 {
			go updateReq()
		}
		time.Sleep(time.Duration(interval)* time.Millisecond)
	}

	<-stop
}

func sendReq() {
	req, err := http.NewRequest("GET", GatewayUrl, nil)
	if err != nil {
		log.Printf("Fail to create http request : %v\n", err)
		return
	}
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Printf("Fail to send http request : %v\n", err)
		return
	}

	log.Println(res.Status)
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
	log.Printf("Response : %s\n", typhoonHeader)
}

func updateReq() {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	time.Sleep(time.Duration(r.Intn(DefaultInterval))*time.Millisecond)

	body := []byte(`{"routeSubset": "822d65df"}`)
	req, err := http.NewRequest("PATCH", UpdateUrl, bytes.NewBuffer(body))
	if err != nil {
		log.Printf("Fail to create update request : %v\n", err)
		return
	}
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Printf("Fail to send update request : %v\n", err)
		return
	}
	log.Println("Update request : " + res.Status)
}
