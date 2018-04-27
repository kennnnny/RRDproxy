package main

import (
	"encoding/json"
	"fmt"
	"time"
)

var data = `{"ato":40000000,"cong_ctl":{"snd_ssthresh":2147483647,"rcv_ssthresh":43690,"snd_cwnd_bytes":0,"snd_cwnd_segs":10},"flow_ctl":{"rcv_wnd":43690},"last_ack_rcvd":260000000,"last_data_rcvd":1536000000,"last_data_sent":1260000000,"opts":{"sack":true,"tmstamps":true,"wscale":7},"peer_opts":{"sack":true,"tmstamps":true,"wscale":7},"rcv_mss":536,"rto":204000000,"rtt":1692000,"rttvar":3334000,"snd_mss":27264,"state":"close-wait","sys":{"path_mtu":65535,"adv_mss":65483,"ca_state":0,"rexmits":0,"backoffs":0,"wnd_ka_probes":0,"unacked_segs":0,"sacked_segs":0,"lost_segs":0,"retrans_segs":0,"fack_segs":0,"reord_segs":3,"rcv_rtt":0,"total_retrans_segs":0,"pacing_rate":322150505,"thru_bytes_acked":3202,"thru_bytes_rcvd":962,"segs_out":25,"segs_in":26,"not_sent_bytes":0,"min_rtt":23000,"data_segs_out":12,"data_segs_in":12}}`

func main() {
	// info := &tcpinfo.Info{}
	info := &MyData{}
	fmt.Println(json.Unmarshal([]byte(data), &info))
	fmt.Printf("%+v", info)
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
