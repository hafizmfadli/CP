package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"
)

// WaitGroup is used to wait for the program to finish goroutines.
var wg sync.WaitGroup

func main() {
	file, err := ioutil.ReadFile("./cmd/test3/bodysatu.json")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	filedua, err := ioutil.ReadFile("./cmd/test3/bodydua.json")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	

	reqbodies := [][]byte{ file, filedua  }
	targetUrl := "http://localhost:4000/v1/checkout"
	success := 0
	failed := 0
	for i := 0; i < 72; i++ {
		for _, v := range reqbodies {
			wg.Add(1)
			go func(reqbody []byte){
				resp, err := http.Post(targetUrl, "application/json", bytes.NewBuffer(reqbody))
				if err != nil {
					// biasanya error kalo udah lebih dari 10s
					failed += 1
					fmt.Println("SUCCESS : ", success)
					fmt.Println("FAILED : ", failed)
					log.Fatal(err)
				}
				defer resp.Body.Close()
				var res map[string]interface{}
				json.NewDecoder(resp.Body).Decode(&res)
				
				if resp.StatusCode >= 500 {
					failed += 1
				}else {
					success += 1
				}
				fmt.Println(res)
				wg.Done()
			}(v)
		}
	}

	wg.Wait()
	fmt.Println("SUCCESS : ", success)
	fmt.Println("FAILED : ", failed)
}