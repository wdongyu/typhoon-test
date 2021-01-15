package timeliness

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

	TURNS = 40
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
			} else if alg == "ManagedService" {
				go util.MsUpdateReq(interval)
			} else {}
		}
		time.Sleep(time.Duration(interval)* time.Millisecond)
	}

	<-stop
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
}
