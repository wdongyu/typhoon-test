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
	MsGatewayUrl = "http://localhost:31380/typhoon-backend?index=1"

	QuieGatewayUrl = "http://localhost:31380/typhoon-backend-quie?index=1"

	MsUpdateUrl = "http://localhost:32088/apis/managedservice/typhoon-microservices-typhoon?namespace=typhoon"

	QuieUpdateUrl = "http://localhost:32088/apis/quie/update"

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

	turns := 80
	var stop = make(chan bool)
	count := make(chan int64, turns)
	for i := 1; i <= turns; i++ {
		go sendReq(i, count)
		if i == 20 {
			go msUpdateReq()
		}
		time.Sleep(time.Duration(interval)* time.Millisecond)
	}

	//timeSum := int64(0)
	//for t := range count {
	//	timeSum += t
	//}
	//log.Println(timeSum)

	<-stop
}

func sendReq(index int, count chan int64) {
	log.Printf("%d. Http request begin ...\n", index)
	begin := time.Now()
	req, err := http.NewRequest("GET", MsGatewayUrl, nil)
	if err != nil {
		log.Printf("Fail to create http request : %v\n", err)
		return
	}
	client := &http.Client{}
	for {
		res, err := client.Do(req)
		if err != nil {
			log.Printf("Fail to send http request : %v\n", err)
			return
		}
		if res.StatusCode == 200 {
			t := time.Now().Sub(begin).Milliseconds()
			//count<-t
			//if index == 80 {
			//	close(count)
			//}
			log.Printf("%d. %d\n", index, t)
			log.Printf("%d. %v\n", index, res.Status)
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
			return
		} else {
			time.Sleep(time.Duration(100)*time.Millisecond)
			log.Printf("%d. %s Retry ... \n", index, res.Status)
		}
	}
}

func quieUpdateReq() {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	time.Sleep(time.Duration(r.Intn(DefaultInterval))*time.Millisecond)

	log.Println("Update request begin ...")
	body := []byte(`{"namespace": "typhoon-quie", "rootService": "typhoon-backend", 
					"targetService": "typhoon-microservices-typhoon", 
					"revokeSubset": "v1", "deploySubset": "v2"}`)
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
	log.Println("Update request : " + res.Status + " !!!")
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
