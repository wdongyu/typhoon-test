package disruption

import (
	"fmt"
	"github.com/wdongyu/typhoon-test/util"
	"log"
	"net/http"
	"strings"
	"time"
)

const (
	GatewayUrlPrefix = "http://localhost:31380/typhoon-backend?index="

	VersionHeader = "X-Version"

	TyphoonHeaderPrefix = "typhoon-microservices-typhoon-"

	TURNS = 60
)

func init() {
	log.SetFlags(log.Lmicroseconds)
}

func Process(alg string, interval int)  {
	var stop = make(chan bool)
	count := make(chan int64, TURNS)
	for i := 1; i <= TURNS; i++ {
		go sendReq(i, stop, count)
		if i == 20 {
			if alg == "Quiescence" {
				go util.QuieUpdateReq(interval)
			} else if alg == "ManagedService" {
				go util.MsUpdateReq(interval)
			} else {}
		}
		time.Sleep(time.Duration(interval)* time.Millisecond)
	}

	timeSum := int64(0)
	for t := range count {
		timeSum += t
	}
	log.Printf("Total time : %d", timeSum)

	<-stop
}

func sendReq(index int, stop chan bool, count chan int64) {
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
			t := time.Now().Sub(begin).Milliseconds()
			count<-t
			if index == TURNS {
				close(count)
			}
			log.Printf("%d. response time : %d\n", index, t)
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
			if index == TURNS {
				stop<-true
			}
			return
		} else {
			time.Sleep(time.Duration(100)*time.Millisecond)
			log.Printf("%d. %s Retry ... \n", index, res.Status)
		}
	}
}
