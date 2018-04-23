
# Golang SOCKS5

- Tested on Go 1.10
- A naive server is included in cmd/server
- net/http's client can do SOCKS5 when HTTP_PROXY environment variable is set.  An example is given in cmd/testclient

# Example
```go
// Create a SOCKS5 server
conf := &socks5.Config{}
server, err := socks5.New(conf)
if err != nil {
  panic(err)
}

// Create SOCKS5 proxy on localhost port 8000
if err := server.ListenAndServe("tcp", "127.0.0.1:8000"); err != nil {
  panic(err)
}
```

## Credits
- Modified from github.com/armon/go-socks5