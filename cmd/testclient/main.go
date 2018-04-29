package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

func main() {
	if len(os.Args) != 2 || os.Getenv("HTTP_PROXY") == "" {
		fmt.Println("usage example: HTTP_PROXY=socks5://user:pass@someIP testclient https://www.google.com")
		os.Exit(1)
	}

	c := http.Client{}
	for i := 0; i < 1; i++ {
		resp, err := c.Get(os.Args[1])
		if err != nil {
			fmt.Println(err)
			return
		}

		b, _ := ioutil.ReadAll(resp.Body)
		fmt.Println(string(b))

		time.Sleep(1 * time.Second)
	}

}
