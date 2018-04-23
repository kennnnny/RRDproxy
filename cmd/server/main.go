package main

import (
	"flag"
	"log"
	"os"

	"github.com/dt665m/gosocks5"
)

var (
	port     string
	login    string
	password string
)

func init() {
	flag.StringVar(&port, "p", "8080", "socks5 listen port")
	flag.StringVar(&login, "user", "", "socks5 username")
	flag.StringVar(&password, "password", "", "socks5 password")
	flag.Parse()
}

func main() {
	conf := &socks5.Config{}
	conf.Logger = log.New(os.Stderr, "", log.LstdFlags|log.Lshortfile)
	if login != "" && password != "" {
		conf.Credentials = socks5.StaticCredentials(map[string]string{login: password})
	} else {
		conf.Logger.Println("WARN: no authentication provided (free use mode)")
	}
	server, err := socks5.New(conf)
	if err != nil {
		panic(err)
	}

	conf.Logger.Printf("SOCKS5 server started on port %s", port)
	if err := server.ListenAndServe("tcp", ":"+port); err != nil {
		panic(err)
	}
}
