package disruption

import (
	"encoding/json"
	"fmt"
	"github.com/wdongyu/typhoon-test/util"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	GatewayUrlPrefix = "http://localhost:31380/typhoon-backend?index="

	TyphoonManagedServiceUrl = "http://localhost:32088/apis/managedservice?" +
		"projectName=typhoon-pm&&app=typhoon-microservices-typhoon&&namespace=typhoon"
	WindManagedServiceUrl = "http://localhost:32088/apis/managedservice?" +
		"projectName=typhoon-pm&&app=typhoon-microservices-windcontroller&&namespace=typhoon"

	VersionHeader = "X-Version"

	TyphoonHeaderPrefix = "typhoon-microservices-typhoon-"
	WindHeaderPrefix = "typhoon-microservices-windcontroller-"

	TURNS = 40
)

func init() {
	log.SetFlags(log.Lmicroseconds)
}

func Process(alg string, interval int)  {
	var stop = make(chan bool)
	respTime := make(chan int64, TURNS)
	ts := make(chan string, TURNS)
	for i := 1; i <= TURNS; i++ {
		go sendReq(alg, i, stop, respTime, ts)
		if i == 20 {
			if alg == "Quiescence" {
				go util.QuieUpdateReq(interval)
			} else if alg == "CompEvo" {
				go util.MsUpdateReq(interval)
			} else {}
		}
		time.Sleep(time.Duration(interval)* time.Millisecond)
	}

	timeSum := int64(0)
	disSum := 0
	errCount := 0
	for i := 0; i < TURNS; i++ {
		t := <-respTime
		timeSum += t
	}
	log.Printf("Total time : %d, average: %d\n", timeSum, timeSum/(TURNS-1))

	var begin, end string
	var beginTs, endTs int64
	if alg == "Quiescence" {
		begin = util.QuieUpdateBegin
		end = util.QuieUpdateEnd
		log.Printf("Update elapsed time : %d\n", util.QuieUpdateElapse)
	} else if alg == "CompEvo" {
		for _, url := range []string{TyphoonManagedServiceUrl} {
			req, err := http.NewRequest("GET", url, nil)
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
			if beginTs == 0 || beginTs > msList.ManagedServices[0].Status.LastUpdateBegin {
				beginTs = msList.ManagedServices[0].Status.LastUpdateBegin
			}

			if endTs == 0 || endTs < msList.ManagedServices[0].Status.LastUpdateEnd {
				endTs = msList.ManagedServices[0].Status.LastUpdateEnd
			}
		}

		begin = strconv.FormatInt(beginTs, 10)
		end = strconv.FormatInt(endTs, 10)
		log.Printf("Update elapsed time : %d\n", (endTs-beginTs)/1000000)
	}
	for i := 1; i <= TURNS-1; i++ {
		t := <-ts
		if util.Between(begin, end, t) {
			tokens := strings.Split(t, "-")
			//log.Println(tokens[2])
			res, err := strconv.Atoi(tokens[2])
			if err != nil {
				log.Printf("Fail to convert %s : %v\n", tokens[2], err)
			} else {
				errCount += 1
				disSum += res
			}
		}
	}
	log.Printf("Total disrupted time : %d, disrupted count: %d\n", disSum + errCount*50, errCount)

	<-stop
}

func sendReq(alg string, index int, stop chan bool, respTime chan int64, ts chan string) {
	log.Printf("%d. Http request begin ...\n", index)
	begin := time.Now()
	req, err := http.NewRequest("GET", fmt.Sprintf("%s%d", GatewayUrlPrefix, index), nil)
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
			end := time.Now()
			t := end.Sub(begin).Milliseconds()
			if index != 1 {
				respTime<-t
				ts<-fmt.Sprintf("%v-%v-%d", begin.UnixNano(), end.UnixNano(), t)
			}
			if index == TURNS {
				close(respTime)
				close(ts)
			}

			log.Printf("%d. response time : %d\n", index, t)
			log.Printf("%d. %v\n", index, res.Status)
			typhoonHeader := ""
			//windHeader := ""
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
				//if strings.HasPrefix(value, WindHeaderPrefix) {
				//	if windHeader == "" {
				//		windHeader = value
				//	} else {
				//		if windHeader != value {
				//			log.Printf("Version header %s : %s conflict ...\n", windHeader, value)
				//		}
				//	}
				//}
			}
			log.Printf("%d. Response : %s\n", index, typhoonHeader)
			//if index == TURNS {
			//	stop<-true
			//}
			return
		} else {
			time.Sleep(time.Duration(100)*time.Millisecond)
			log.Printf("%d. %s Retry ... \n", index, res.Status)
		}
	}
}
