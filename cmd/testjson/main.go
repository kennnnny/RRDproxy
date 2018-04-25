package main

import (
	"encoding/json"
	"fmt"

	"github.com/mikioh/tcpinfo"
)

var data = `{"ato":40000000,"cong_ctl":{"snd_ssthresh":2147483647,"rcv_ssthresh":43690,"snd_cwnd_bytes":0,"snd_cwnd_segs":10},"flow_ctl":{"rcv_wnd":43690},"last_ack_rcvd":260000000,"last_data_rcvd":1536000000,"last_data_sent":1260000000,"opts":{"sack":true,"tmstamps":true,"wscale":7},"peer_opts":{"sack":true,"tmstamps":true,"wscale":7},"rcv_mss":536,"rto":204000000,"rtt":1692000,"rttvar":3334000,"snd_mss":27264,"state":"close-wait","sys":{"path_mtu":65535,"adv_mss":65483,"ca_state":0,"rexmits":0,"backoffs":0,"wnd_ka_probes":0,"unacked_segs":0,"sacked_segs":0,"lost_segs":0,"retrans_segs":0,"fack_segs":0,"reord_segs":3,"rcv_rtt":0,"total_retrans_segs":0,"pacing_rate":322150505,"thru_bytes_acked":3202,"thru_bytes_rcvd":962,"segs_out":25,"segs_in":26,"not_sent_bytes":0,"min_rtt":23000,"data_segs_out":12,"data_segs_in":12}}`

func main() {
	info := &tcpinfo.Info{}
	fmt.Println(json.Unmarshal([]byte(data), &info))
	fmt.Printf("%+v", info)
}
