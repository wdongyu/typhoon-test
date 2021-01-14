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

	MsUpdateUrl = "http://localhost:32088/apis/managedservice/typhoon-microservices-typhoon?namespace=typhoon"

	QuieUpdateUrl = "http://localhost:32088/apis/quie/update"

	VersionHeader = "X-Version"

	DefaultInterval = 1500

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
	turns := 10
	for i := 1; i <= turns; i++ {
		log.Printf("#%d.\n", i)
		go sendReq(i)
		//if i == 20 {
		//	go quieUpdateReq()
		//}
		time.Sleep(time.Duration(interval)* time.Millisecond)
	}

	<-stop
}

func sendReq(index int) {
	req, err := http.NewRequest("GET", GatewayUrl, nil)
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

func quieUpdateReq() {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	time.Sleep(time.Duration(r.Intn(DefaultInterval))*time.Millisecond)

	log.Println("Update request begin ...")
	begin := time.Now()
	body := []byte(`{"namespace": "typhoon", "rootService": "typhoon-backend",
					"targetService": "typhoon-microservices-typhoon", 
					"revokeSubset": "f7400817", "deploySubset": "822d65df"}`)
	req, err := http.NewRequest("POST", QuieUpdateUrl, bytes.NewBuffer(body))
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
	log.Println(time.Now().Sub(begin).Milliseconds())
}

func msUpdateReq() {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	time.Sleep(time.Duration(r.Intn(DefaultInterval))*time.Millisecond)

	log.Println("Update request begin ...")
	begin := time.Now()
	body := []byte(`{"routeSubset": "822d65df"}`)
	req, err := http.NewRequest("PATCH", MsUpdateUrl, bytes.NewBuffer(body))
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
	log.Println(time.Now().Sub(begin).Milliseconds())
}
