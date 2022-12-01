package main

import (
	"encoding/json"
	"log"
	"net"
	"net/http"
    "bytes"
    "sync"
    "context"
)

type Addrs struct {
    Control string `json:"control"`
    Data string `json:"data"`
}

func main() {
    wg := sync.WaitGroup{}
    wg.Add(2)

    controlIface, err := net.InterfaceByName("eth0")
    if err != nil {
        log.Fatal(err)
    }
    controlAddrs, err := controlIface.Addrs()
    if err != nil {
        log.Fatal(err)
    }
    controlAddr := controlAddrs[0].String()
    controlAddr = controlAddr[:len(controlAddr)-3]

    dataIface, err := net.InterfaceByName("eth1")
    if err != nil {
        log.Fatal(err)
    }
    dataAddrs, err := dataIface.Addrs()
    if err != nil {
        log.Fatal(err)
    }
    dataAddr := dataAddrs[0].String()
    dataAddr = dataAddr[:len(dataAddr)-3]

    m := http.NewServeMux()
    s := http.Server{Addr: ":80", Handler: m}
    m.HandleFunc("/shutdown", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("OK"))
        go func() {
            s.Shutdown(context.Background())
        }()
    })
    m.HandleFunc("/peers", func(w http.ResponseWriter, r *http.Request) {
        log.Print("Received peers request")
        decoder := json.NewDecoder(r.Body)
        var addrs []Addrs
        err := decoder.Decode(&addrs)
        if err != nil {
            log.Fatal(err)
        }

        // check if addrs has 2 entries
        if len(addrs) == 2 {
            log.Println("Received 2 peers")
            // find the addr which is not this one
            var other Addrs
            for _, addr := range addrs {
                if addr.Control != controlAddr {
                    other = addr
                }
            }
            log.Println("Other peer is", other)

            go func() {
                // send ping to the data address of the other node
                resp, err := http.Post("http://" + other.Data + "/ping", "application/json", bytes.NewBuffer([]byte("ping")))
                if err != nil {
                    log.Fatal(err)
                }
                log.Println(resp)
                // send done to the coordinator
                resp, err = http.Post("http://coordinator:80/done", "application/json", bytes.NewBuffer([]byte("done")))
                if err != nil {
                    log.Fatal(err)
                }
            }()
        }

		w.Write([]byte("OK"))
    })
    // handle ping and respond with pong
    m.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
        log.Print("Received ping")
        w.Write([]byte("pong"))
    })
    go func() {
        if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatal(err)
        }
        wg.Done()
    }()

    go func() {
        addrs := Addrs{controlAddr, dataAddr}
        data, err := json.Marshal(addrs)
        if err != nil {
            log.Fatal(err)
        }
        log.Println("Registering...")
        resp, err := http.Post("http://coordinator:80/register", "text/plain", bytes.NewBuffer(data))
        if err != nil {
            log.Fatal(err)
        }
        log.Println(resp)
        wg.Done()
    }()

    wg.Wait()
    log.Println("Done")
}
