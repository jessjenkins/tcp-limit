package main

import (
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	"golang.org/x/net/netutil"
)

func main() {
	log.Println("Starting")

	go func() {
		connectionCount := 3

		l, err := net.Listen("tcp", "localhost:56700")

		if err != nil {
			log.Fatalf("Listen: %v", err)
		}

		defer l.Close()

		l = netutil.LimitListener(l, connectionCount)

		http.HandleFunc("/slow", slow)

		log.Fatal(http.Serve(l, nil))
	}()

	var wg sync.WaitGroup

	for i := 0; i < 12; i++ {
		wg.Add(1)
		time.Sleep(100 * time.Millisecond)
		go func() {
			defer wg.Done()
			getter(i)
		}()
	}

	wg.Wait()
}

func getter(i int) {
	log.Printf("Getting %d\n", i)
	t := time.Now()

	hc := http.DefaultClient
	hc.Timeout = 1 * time.Second

	resp, err := hc.Get("http://localhost:56700/slow")
	if err != nil {
		log.Printf("Error from http get %d: %v", i, err)
		return
	}
	defer resp.Body.Close()

	d := time.Since(t).Milliseconds()
	log.Printf("Response %d (in %dms): %d %s\n", i, d, resp.StatusCode, resp.Status)
}

func slow(w http.ResponseWriter, req *http.Request) {
	time.Sleep(500 * time.Millisecond)
	w.Write([]byte("Hello world!"))
}
