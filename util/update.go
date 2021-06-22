package util

import (
	"bytes"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	TyphoonMsUpdateUrl = "http://localhost:32088/apis/managedservice/typhoon-microservices-typhoon?namespace=typhoon"
	WindMsUpdateUrl = "http://localhost:32088/apis/managedservice/typhoon-microservices-windcontroller?namespace=typhoon"

	QuieUpdateUrl = "http://localhost:32088/apis/quie/update"

	CanaryUpdateUrl = "http://localhost:32088/apis/virtualservice/typhoon-microservices-typhoon?namespace=typhoon"
	//CanaryUpdateUrl = "http://localhost:32088/apis/virtualservice/typhoon-microservices-windcontroller?namespace=typhoon"
)

var (
	QuieUpdateBegin string
	QuieUpdateEnd string
	QuieUpdateElapse int64
)

func QuieUpdateReq(interval int) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	time.Sleep(time.Duration(r.Intn(interval))*time.Millisecond)

	log.Println("Quiescence Update request begin ...")
	begin := time.Now()
	body := []byte(`{"namespace": "typhoon", "rootService": "typhoon-backend",
					"targetService": "typhoon-microservices-typhoon", 
					"revokeSubset": "1bf67d52", "deploySubset": "550c1013"}`)
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

	log.Println("Quiescence Update response : " + res.Status)
	end := time.Now()
	log.Printf("Elapsed time : %d", end.Sub(begin).Milliseconds())

	QuieUpdateBegin = strconv.FormatInt(begin.UnixNano(), 10)
	QuieUpdateEnd = strconv.FormatInt(end.UnixNano(), 10)
	QuieUpdateElapse = end.Sub(begin).Milliseconds()
}

func MsUpdateReq(interval int) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	time.Sleep(time.Duration(r.Intn(interval))*time.Millisecond)

	log.Println("ManagedService Update request begin ...")
	// begin := time.Now()
	body := []byte(`{"routeSubset": "550c1013"}`)
	req1, err := http.NewRequest("PATCH", TyphoonMsUpdateUrl, bytes.NewBuffer(body))
	if err != nil {
		log.Printf("Fail to create update request : %v\n", err)
		return
	}
	//req2, err := http.NewRequest("PATCH", WindMsUpdateUrl, bytes.NewBuffer(body))
	//if err != nil {
	//	log.Printf("Fail to create update request : %v\n", err)
	//	return
	//}
	client := &http.Client{}
	res1, err := client.Do(req1)
	if err != nil {
		log.Printf("Fail to send update request : %v\n", err)
		return
	}
	//res2, err := client.Do(req2)
	//if err != nil {
	//	log.Printf("Fail to send update request : %v\n", err)
	//	return
	//}
	log.Println("Typhoon ManagedService response : " + res1.Status)
	//log.Println("WindController ManagedService response : " + res2.Status)
	// log.Printf("Elapsed time : %d", time.Now().Sub(begin).Milliseconds())
}

func CanarypdateReq(interval int) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	time.Sleep(time.Duration(r.Intn(interval))*time.Millisecond)

	log.Println("Canary Update request begin ...")
	begin := time.Now()
	body := []byte(`{"routeSubset": "550c1013"}`)
	req, err := http.NewRequest("PATCH", CanaryUpdateUrl, bytes.NewBuffer(body))
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
	log.Println("Canary Update response : " + res.Status)
	log.Printf("Elapsed time : %d", time.Now().Sub(begin).Milliseconds())
}

func Between(begin string, end string, ts string) bool {
	tokens := strings.Split(ts, "-")
	if  (begin < tokens[0] && tokens[0] < end) ||
			(begin < tokens[1] && tokens[1] < end) {
		return true
	}
	return false
}

