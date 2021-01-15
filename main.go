package main

import (
	"flag"
	"github.com/wdongyu/typhoon-test/disruption"
	"github.com/wdongyu/typhoon-test/safety"
	"github.com/wdongyu/typhoon-test/timeliness"
	"log"
)

const (
	DefaultInterval = 750
	DefaultObject = "Safety"
	DefaultAlg = "ManagedService"
)

func init() {
	log.SetFlags(log.Lmicroseconds)
}

func main() {
	var expType, alg string
	var interval int
	flag.IntVar(&interval, "interval", DefaultInterval, `interval to send http request`)
	flag.StringVar(&expType, "object", DefaultObject, `object of the test`)
	flag.StringVar(&alg, "alg", DefaultAlg, `algorithm will be used`)
	flag.Parse()

	switch expType {
		case "Safety" : safety.Process(alg, interval)
		case "Timeliness": timeliness.Process(alg, interval)
		case "Disruption": disruption.Process(alg, interval)
	}
}
