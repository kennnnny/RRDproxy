// Copyright 2016 Mikio Hara. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package tcpinfo

/*
#include <netinet/tcp.h>
*/
import "C"

const (
	sysTCP_CONNECTION_INFO = C.TCP_CONNECTION_INFO

	sysTCPCI_OPT_TIMESTAMPS = C.TCPCI_OPT_TIMESTAMPS
	sysTCPCI_OPT_SACK       = C.TCPCI_OPT_SACK
	sysTCPCI_OPT_WSCALE     = C.TCPCI_OPT_WSCALE
	sysTCPCI_OPT_ECN        = C.TCPCI_OPT_ECN

	SysFlagLossRecovery       SysFlags = C.TCPCI_FLAG_LOSSRECOVERY
	SysFlagReorderingDetected SysFlags = C.TCPCI_FLAG_REORDERING_DETECTED

	sizeofTCPConnectionInfo = C.sizeof_struct_tcp_connection_info
)

type tcpConnectionInfo C.struct_tcp_connection_info
