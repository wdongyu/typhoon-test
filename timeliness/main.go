package timeliness

import (
	"encoding/json"
	"fmt"
	"github.com/wdongyu/typhoon-test/util"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

const (
	GatewayUrlPrefix = "http://localhost:31380/typhoon-backend?index="

	VersionHeader = "X-Version"

	TyphoonHeaderPrefix = "typhoon-microservices-typhoon-"
	WindHeaderPrefix = "typhoon-microservices-windcontroller-"

	TURNS = 60

	TyphoonManagedServiceUrl = "http://localhost:32088/apis/managedservice?" +
		"projectName=typhoon-pm&&app=typhoon-microservices-typhoon&&namespace=typhoon"
	WindManagedServiceUrl = "http://localhost:32088/apis/managedservice?" +
		"projectName=typhoon-pm&&app=typhoon-microservices-windcontroller&&namespace=typhoon"
)

func init() {
	log.SetFlags(log.Lmicroseconds)
}

func Process(alg string, interval int)  {
	var stop = make(chan bool)
	for i := 1; i <= TURNS; i++ {
		log.Printf("#%d.\n", i)
		go sendReq(i, stop)
		if i == 20 {
			if alg == "Quiescence" {
				go util.QuieUpdateReq(interval)
			} else if alg == "CompEvo" {
				go util.MsUpdateReq(interval)
			} else {}
		}
		time.Sleep(time.Duration(interval)* time.Millisecond)
	}

	<-stop

	if alg == "Quiescence" {
		log.Printf("Update elapsed time : %d\n", util.QuieUpdateElapse)
	} else if alg == "CompEvo" {
		req, err := http.NewRequest("GET", TyphoonManagedServiceUrl, nil)
		if err != nil {
			log.Printf("Fail to create get request : %v\n", err)
			return
		}
		client := &http.Client{}
		res, err := client.Do(req)
		if err != nil {
			log.Printf("Fail to send get request : %v\n", err)
			return
		}
		resBody, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Printf("Fail to read push status response body : %v\n", err)
			return
		}
		var msList util.ManagedServiceList
		if err = json.Unmarshal(resBody, &msList); err != nil {
			log.Printf("Fail to unmarshal push status response body : %v\n", err)
			return
		}

		log.Printf("Update elapsed time : %d\n", msList.ManagedServices[0].Status.ElapseTime)
	}
}

func sendReq(index int, stop chan bool) {
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
	windHeader := ""
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
		if strings.HasPrefix(value, WindHeaderPrefix) {
			if windHeader == "" {
				windHeader = value
			} else {
				if windHeader != value {
					log.Printf("Version header %s : %s conflict ...\n", windHeader, value)
				}
			}
		}
	}
	log.Printf("%d. Response : %s, %s\n", index, windHeader,  typhoonHeader)
	if index == TURNS {
		stop<-true
	}
}
