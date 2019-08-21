package signal

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// getHost tries its best to return the request host.
func getHost(r *http.Request) string {
	if r.URL.IsAbs() {
		host := r.Host
		// Slice off any port information.
		if i := strings.Index(host, ":"); i != -1 {
			host = host[:i]
		}
		return host
	}
	return r.URL.Host
}

// HTTPSDPServer starts a HTTP Server that consumes SDPs
func HTTPSDPServer() (chan string, chan string) {
	port := flag.Int("port", 8080, "port of http server")
	dir := flag.String("dir", "static", "base directory of file server")
	flag.Parse()

	sdpInChan := make(chan string)
	sdpOutChan := make(chan string)

	http.HandleFunc("/sdp", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("/sdp connected from %s", r.Host)
		body, _ := ioutil.ReadAll(r.Body)
		// process request of sdp
		sdpInChan <- string(body)
		// send response of sdp
		fmt.Fprintf(w, <-sdpOutChan)
	})

	// http server for static files
	fs := http.FileServer(http.Dir(*dir))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	go func() {
		err := http.ListenAndServe(":"+strconv.Itoa(*port), nil)
		if err != nil {
			log.Println(err)
			panic(err)
		}
	}()

	log.Println("\nWebRTC SFU example server is started\n")
	log.Printf("started http and file server on :%d and %s", *port, *dir)
	return sdpInChan, sdpOutChan
}
