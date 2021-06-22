package safety

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
	//GatewayUrl = "http://localhost:10080/test?index=1"

	VersionHeader = "X-Version"

	TyphoonHeaderPrefix = "typhoon-microservices-typhoon-"
	WindHeaderPrefix = "typhoon-microservices-windcontroller-"
	//TyphoonHeaderPrefix = "typhoon"

	TURNS = 40
)

func init() {
	log.SetFlags(log.Lmicroseconds)
}

func Process(alg string, interval int)  {
	var stop = make(chan bool)
	var count = make(chan int, TURNS)
	for i := 1; i <= TURNS; i++ {
		log.Printf("#%d.\n", i)
		go sendReq(i, stop, count)
		if i == 20 {
			if alg == "Direct" {
				go util.CanarypdateReq(interval)
			} else if alg == "Quiescence" {
				go util.QuieUpdateReq(interval)
			} else if alg == "CompEvo" {
				go util.MsUpdateReq(interval)
			} else {}
		}
		time.Sleep(time.Duration(interval)* time.Millisecond)
	}
	incon := int(0)
	for t := range count {
		incon += t
	}
	log.Printf("Total inconsistency : %d\n", incon)

	<-stop
}

func sendReq(index int, stop chan bool, count chan int) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s%d", GatewayUrlPrefix, index), nil)
	//req, err := http.NewRequest("GET", GatewayUrl, nil)
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

	log.Printf("%d. %s ", index, res.Status)
	typhoonHeader := ""
	for _, value := range res.Header[VersionHeader] {
		if strings.HasPrefix(value, TyphoonHeaderPrefix) {
			if typhoonHeader == "" {
				typhoonHeader = value
			} else {
				if typhoonHeader != value {
					log.Printf("Version header %s : %s conflict ...\n", typhoonHeader, value)
					count<-1
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
	if index == TURNS {
		close(count)
		stop<-true
	}
}

