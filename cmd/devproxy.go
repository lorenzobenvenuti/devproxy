package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/lorenzobenvenuti/devproxy"
)

func main() {
	localPort := flag.Int("lp", 9090, "Local port")
	remoteHost := flag.String("rh", "", "Remote host")
	flag.Parse()
	url, urlErr := url.Parse(*remoteHost)
	if urlErr != nil {
		panic(urlErr)
	}
	//eventBus := devproxylib.NewEventBus()
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)
	logger.Printf("Starting devproxylib  with %d, %s\n", *localPort, *url)
	handlers := []devproxy.Handler{devproxy.NewLoggerHandler(logger)}
	ph := devproxy.NewProxyHandler(devproxy.NewStringIdGenerator(), devproxy.NewHandlerChainFactory(logger, &http.Transport{}, handlers), url)
	err := http.ListenAndServe(fmt.Sprintf(":%d", *localPort), ph)
	if err != nil {
		panic(err)
	}
}
