package arrowheadfunctions

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
)

func Hello() {
	fmt.Println("Daniel-sama")
}

func GetClient(certFile string, keyFile string) *http.Client {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		log.Panic("Certficate load error. ", err)
	}
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			Certificates: []tls.Certificate{cert},
		},
	}
	client := &http.Client{Transport: tr}
	return client
}
