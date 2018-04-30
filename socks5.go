package socks5

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"time"

	"github.com/mikioh/tcp"
	"github.com/mikioh/tcpinfo"
)

const (
	socks5Version = uint8(5)
)

// Config is used to setup and configure a Server
type Config struct {
	// AuthMethods can be provided to implement custom authentication
	// By default, "auth-less" mode is enabled.
	// For password-based auth use UserPassAuthenticator.
	AuthMethods []Authenticator

	// If provided, username/password authentication is enabled,
	// by appending a UserPassAuthenticator to AuthMethods. If not provided,
	// and AUthMethods is nil, then "auth-less" mode is enabled.
	Credentials CredentialStore

	// Resolver can be provided to do custom name resolution.
	// Defaults to DNSResolver if not provided.
	Resolver NameResolver

	// Rules is provided to enable custom logic around permitting
	// various commands. If not provided, PermitAll is used.
	Rules RuleSet

	// Rewriter can be used to transparently rewrite addresses.
	// This is invoked before the RuleSet is invoked.
	// Defaults to NoRewrite.
	Rewriter AddressRewriter

	// BindIP is used for bind or udp associate
	BindIP net.IP

	// Logger can be used to provide a custom log target.
	// Defaults to stdout.
	Logger *log.Logger

	// Optional function for dialing out
	Dial func(ctx context.Context, network, addr string) (net.Conn, error)
}

//MyData is OutPut Data Structure
type MyData struct {
	Ato     time.Duration `json:"ato"`
	CongCtl struct {
		SSthresh          uint `json:"snd_ssthresh"`
		RcvThresh         uint `json:"rcv_ssthresh"`
		SenderWindowBytes uint `json:"snd_cwnd_bytes"`
		SenderWindowSegs  uint `json:"snd_cwnd_segs"`
	} `json:"cong_ctl"`
	FlowControl struct {
		ReceiverWindow uint `json:"rcv_wnd"`
	} `json:"flow_ctl"`
	LastDataReceived time.Duration `json:"last_data_rcvd"` // since last data received [FreeBSD and Linux]
	LastAckReceived  time.Duration `json:"last_ack_rcvd"`  // since last ack received [Linux only]
	LastDataSent     time.Duration `json:"last_data_sent"` // since last data sent [Linux only]
	Opts             struct {
		SACKPermitted bool `json:"sack"`
		Timestamps    bool `json:"tmstamps"`
		WindowScale   int  `json:"wscale"`
	} `json:"opts"`
	PeerOpts struct {
		SACKPermitted bool `json:"sack"`
		Timestamps    bool `json:"tmstamps"`
		WindowScale   int  `json:"wscale"`
	} `json:"peer_opts"`
	ReceiverMSS uint          `json:"rcv_mss"`
	RTO         time.Duration `json:"rto"`
	RTT         time.Duration `json:"rtt"`
	RTTVar      time.Duration `json:"rttvar"`
	SenderMSS   uint          `json:"snd_mss"`
	State       string        `json:"state"`
	System      struct {
		PathMTU                 uint          `json:"path_mtu"`           // path maximum transmission unit
		AdvertisedMSS           uint          `json:"adv_mss"`            // advertised maximum segment size
		CAState                 int           `json:"ca_state"`           // state of congestion avoidance
		Retransmissions         uint          `json:"rexmits"`            // # of retranmissions on timeout invoked
		Backoffs                uint          `json:"backoffs"`           // # of times retransmission backoff timer invoked
		WindowOrKeepAliveProbes uint          `json:"wnd_ka_probes"`      // # of window or keep alive probes sent
		UnackedSegs             uint          `json:"unacked_segs"`       // # of unack'd segments
		SackedSegs              uint          `json:"sacked_segs"`        // # of sack'd segments
		LostSegs                uint          `json:"lost_segs"`          // # of lost segments
		RetransSegs             uint          `json:"retrans_segs"`       // # of retransmitting segments in transmission queue
		ForwardAckSegs          uint          `json:"fack_segs"`          // # of forward ack segments in transmission queue
		ReorderedSegs           uint          `json:"reord_segs"`         // # of reordered segments allowed
		ReceiverRTT             time.Duration `json:"rcv_rtt"`            // current RTT for receiver
		TotalRetransSegs        uint          `json:"total_retrans_segs"` // # of retransmitted segments
		PacingRate              uint64        `json:"pacing_rate"`        // pacing rate
		ThruBytesAcked          uint64        `json:"thru_bytes_acked"`   // # of bytes for which cumulative acknowledgments have been received
		ThruBytesReceived       uint64        `json:"thru_bytes_rcvd"`    // # of bytes for which cumulative acknowledgments have been sent
		SegsOut                 uint          `json:"segs_out"`           // # of segments sent
		SegsIn                  uint          `json:"segs_in"`            // # of segments received
		NotSentBytes            uint          `json:"not_sent_bytes"`     // # of bytes not sent yet
		MinRTT                  time.Duration `json:"min_rtt"`            // current measured minimum RTT; zero means not available
		DataSegsOut             uint          `json:"data_segs_out"`      // # of segments sent containing a positive length data segment
		DataSegsIn              uint          `json:"data_segs_in"`       // # of segments received containing a positive length data segment
	} `json:"sys"`
}

// Server is reponsible for accepting connections and handling
// the details of the SOCKS5 protocol
type Server struct {
	config      *Config
	authMethods map[uint8]Authenticator
}

// New creates a new Server and potentially returns an error
func New(conf *Config) (*Server, error) {
	// Ensure we have at least one authentication method enabled
	if len(conf.AuthMethods) == 0 {
		if conf.Credentials != nil {
			conf.AuthMethods = []Authenticator{&UserPassAuthenticator{conf.Credentials}}
		} else {
			conf.AuthMethods = []Authenticator{&NoAuthAuthenticator{}}
		}
	}

	// Ensure we have a DNS resolver
	if conf.Resolver == nil {
		conf.Resolver = DNSResolver{}
	}

	// Ensure we have a rule set
	if conf.Rules == nil {
		conf.Rules = PermitAll()
	}

	// Ensure we have a log target
	if conf.Logger == nil {
		conf.Logger = log.New(os.Stdout, "", log.LstdFlags)
	}

	server := &Server{
		config: conf,
	}

	server.authMethods = make(map[uint8]Authenticator)

	for _, a := range conf.AuthMethods {
		server.authMethods[a.GetCode()] = a
	}

	return server, nil
}

// ListenAndServe is used to create a listener and serve on it
func (s *Server) ListenAndServe(network, addr string) error {
	l, err := net.Listen(network, addr)
	if err != nil {
		return err
	}
	return s.Serve(l)
}

// Serve is used to serve connections from a listener
func (s *Server) Serve(l net.Listener) error {
	for {
		conn, err := l.Accept()
		if err != nil {
			return err
		}

		tc, err := tcp.NewConn(conn)
		if err != nil {
			log.Fatal(err)
		}

		go s.Monitor(tc)
		go s.ServeConn(tc)
	}
}

// ServeConn is used to serve a single connection.
func (s *Server) ServeConn(conn net.Conn) error {
	defer conn.Close()
	bufConn := bufio.NewReader(conn)

	// Read the version byte
	version := []byte{0}
	if _, err := bufConn.Read(version); err != nil {
		s.config.Logger.Printf("[ERR] socks: Failed to get version byte: %v", err)
		return err
	}

	// Ensure we are compatible
	if version[0] != socks5Version {
		err := fmt.Errorf("Unsupported SOCKS version: %v", version)
		s.config.Logger.Printf("[ERR] socks: %v", err)
		return err
	}

	// Authenticate the connection
	authContext, err := s.authenticate(conn, bufConn)
	if err != nil {
		err = fmt.Errorf("Failed to authenticate: %v", err)
		s.config.Logger.Printf("[ERR] socks: %v", err)
		return err
	}

	request, err := NewRequest(bufConn)
	if err != nil {
		if err == unrecognizedAddrType {
			if err := sendReply(conn, addrTypeNotSupported, nil); err != nil {
				return fmt.Errorf("Failed to send reply: %v", err)
			}
		}
		return fmt.Errorf("Failed to read destination address: %v", err)
	}
	request.AuthContext = authContext
	if client, ok := conn.RemoteAddr().(*net.TCPAddr); ok {
		request.RemoteAddr = &AddrSpec{IP: client.IP, Port: client.Port}
	}

	// Process the client request
	if err := s.handleRequest(request, conn); err != nil {
		err = fmt.Errorf("Failed to handle request: %v", err)
		s.config.Logger.Printf("[ERR] socks: %v", err)
		return err
	}

	return nil
}

//Monitor monitors net.conn and shows tcp.infos
func (s *Server) Monitor(tc *tcp.Conn) {
	fmt.Println("starting monitor for", tc.RemoteAddr())
	//delete all rules first
	deletecmd := exec.Command("iptables", "-t", "mangle", "-F")
	deletecmd.Stderr = os.Stderr
	deletecmd.Stdout = os.Stdout
	if err := deletecmd.Start(); err != nil {
		fmt.Println("delete", err)
	}
	//add a normal TCPMSS rule
	addcmd := exec.Command("iptables", "-t mangle", "-I POSTROUTING", "-p tcp --tcp-flags SYN,RST SYN", "-j TCPMSS --set-mss 1492")
	addcmd.Stderr = os.Stderr
	addcmd.Stdout = os.Stdout
	if err := addcmd.Start(); err != nil {
		fmt.Println("add:", err)
	}
	for {

		//Print tcpinfo
		var o tcpinfo.Info
		var b [256]byte
		i, err := tc.Option(o.Level(), o.Name(), b[:])
		if err != nil {
			log.Println(err)
			return
		}
		data, err := json.Marshal(i)
		if err != nil {
			log.Println(err)
			return
		}

		txt := string(data)
		// info := &tcpinfo.Info{}
		info := &MyData{}
		json.Unmarshal([]byte(txt), &info)
		fmt.Printf("%+v\n", info)
		//exec.Command("iptables", "-I").Run()

		//lower MSS if retransmit happened
		switch info.System.Retransmissions {
		case 0:
			exec.Command("sudo iptables", "-t mangle -R POSTROUTING 1 -p tcp --tcp-flags SYN,RST SYN -j TCPMSS --set-mss 1492")
		default:
			exec.Command("sudo iptables", "-t mangle -R POSTROUTING 1 -p tcp --tcp-flags SYN,RST SYN -j TCPMSS --set-mss 200")
			fmt.Println("Detect a retransimission! change MSS to 200")
		}
		time.Sleep(100 * time.Millisecond)
		//command already added in linux's iptables: exec.Command("iptables", "-t mangle -I POSTROUTING -p tcp --tcp-flags SYN,RST SYN -j TCPMSS --set-mss 1492")
		//command that delete rules in iptables: exec.Command("sudo iptables", "-t mangle -F")
		//command that Replace a rule in iptables(first line):
		//exec.Command("sudo iptables", "-t mangle -R POSTROUTING 1 -p tcp --tcp-flags SYN,RST SYN -j TCPMSS --set-mss 1492")
		// iptables -t mangle -I OUTPUT -p tcp --sport 80 --tcp-flags SYN,ACK SYN,ACK -j TCPWIN --tcpwin-set 1000
	}
}
