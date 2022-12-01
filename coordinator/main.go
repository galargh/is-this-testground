package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
    "context"
    "sync"
)

type Addrs struct {
	Control string `json:"control"`
	Data    string `json:"data"`
}

func main() {
    wg1 := sync.WaitGroup{}
    wg1.Add(2)

    wg2 := sync.WaitGroup{}
    wg2.Add(2)

    m := http.NewServeMux()
    s := http.Server{Addr: ":80", Handler: m}
    // create a counter
	// list of Addrs
	all := []Addrs{}

    go func() {
        wg1.Wait()
        // send all to all addrs specified in the list
        for _, addr := range all {
            log.Println("Sending to", addr)
            data, err := json.Marshal(all)
            if err != nil {
                log.Fatal(err)
            }
            // send the data to the control address
            resp, err := http.Post("http://" + addr.Control + "/peers", "application/json", bytes.NewBuffer(data))
            if err != nil {
                log.Fatal(err)
            }
            log.Println(resp)
        }
    }()

    go func() {
        wg2.Wait()
        log.Println("Shutting down")
        // tell all runners to shutdown
        for _, addr := range all {
            log.Println("Sending to", addr)
            // send the data to the control address
            resp, err := http.Post("http://" + addr.Control + "/shutdown", "application/json", bytes.NewBuffer([]byte("shutdown")))
            if err != nil {
                log.Fatal(err)
            }
            log.Println(resp)
        }
        s.Shutdown(context.Background())
    }()

	m.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
        log.Println("Received register request")
		// extract JSON data from the body of the request
        decoder := json.NewDecoder(r.Body)
        var addrs Addrs
        err := decoder.Decode(&addrs)
        if err != nil {
            log.Fatal(err)
        }

		// add the new Addrs to the list
		all = append(all, addrs)

		w.Write([]byte("OK"))

        wg1.Done()
	})
    // handle done and respond with OK
    m.HandleFunc("/done", func(w http.ResponseWriter, r *http.Request) {
        log.Print("Received done")
        w.Write([]byte("OK"))
        wg2.Done()
    })
	if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
        log.Fatal(err)
    }
}
