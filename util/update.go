package util

import (
	"bytes"
	"log"
	"math/rand"
	"net/http"
	"time"
)

const (
	MsUpdateUrl = "http://localhost:32088/apis/managedservice/typhoon-microservices-typhoon?namespace=typhoon"

	QuieUpdateUrl = "http://localhost:32088/apis/quie/update"

	CanaryUpdateUrl = "http://localhost:32088/apis/virtualservice/typhoon-microservices-typhoon?namespace=typhoon"
)

func QuieUpdateReq(interval int) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	time.Sleep(time.Duration(r.Intn(interval))*time.Millisecond)

	log.Println("Quiescence Update request begin ...")
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
	log.Println("Quiescence Update response : " + res.Status)
	log.Printf("Elapsed time : %d", time.Now().Sub(begin).Milliseconds())
}

func MsUpdateReq(interval int) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	time.Sleep(time.Duration(r.Intn(interval))*time.Millisecond)

	log.Println("ManagedService Update request begin ...")
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
	log.Println("ManagedService response : " + res.Status)
	log.Printf("Elapsed time : %d", time.Now().Sub(begin).Milliseconds())
}

func CanarypdateReq(interval int) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	time.Sleep(time.Duration(r.Intn(interval))*time.Millisecond)

	log.Println("Canary Update request begin ...")
	begin := time.Now()
	body := []byte(`{"routeSubset": "822d65df"}`)
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

