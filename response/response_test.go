package response

import (
	"encoding/json"
	"strings"
	"testing"
)

var monData = `name:    bridge
rx-packets-per-second:         8
rx-bits-per-second:   9.5kbps
fp-rx-packets-per-second:         8
fp-rx-bits-per-second:   9.5kbps
rx-drops-per-second:         0
rx-errors-per-second:         0
tx-packets-per-second:        12
tx-bits-per-second:  84.8kbps
fp-tx-packets-per-second:        10
fp-tx-bits-per-second:  82.1kbps
tx-drops-per-second:         0
tx-queue-drops-per-second:         0
tx-errors-per-second:         0

name:   bridge
rx-packets-per-second:        4
rx-bits-per-second:  4.8kbps
fp-rx-packets-per-second:        4
fp-rx-bits-per-second:  4.4kbps
rx-drops-per-second:        0
rx-errors-per-second:        0
tx-packets-per-second:        3
tx-bits-per-second:  4.3kbps
fp-tx-packets-per-second:        0
fp-tx-bits-per-second:     0bps
tx-drops-per-second:        0
tx-queue-drops-per-second:        0
tx-errors-per-second:        0

name:    bridge
rx-packets-per-second:         8
rx-bits-per-second:   4.7kbps
fp-rx-packets-per-second:         8
fp-rx-bits-per-second:   4.7kbps
rx-drops-per-second:         0
rx-errors-per-second:         0
tx-packets-per-second:         7
tx-bits-per-second:  11.2kbps
fp-tx-packets-per-second:         0
fp-tx-bits-per-second:      0bps
tx-drops-per-second:         0
tx-queue-drops-per-second:         0
tx-errors-per-second:         0`

var monRes = []TrafficMonitorResponse{
	{Name: "bridge", RxPacketsPerSecond: 0x8, RxBitsPerSecond: 0x251c, FpRxPacketsPerSecond: 0x8, FpRxBitsPerSecond: 0x251c, RxDropsPerSecond: 0x0, RxErrorsPerSecond: 0x0, TxPacketsPerSecond: 0xc, TxBitsPerSecond: 0x14b40, FpTxPacketsPerSecond: 0xa, FpTxBitsPerSecond: 0x140b4, TxDropsPerSecond: 0x0, TxQueueDropsPerSecond: 0x0, TxErrorsPerSecond: 0x0},
	{Name: "bridge", RxPacketsPerSecond: 0x4, RxBitsPerSecond: 0x12c0, FpRxPacketsPerSecond: 0x4, FpRxBitsPerSecond: 0x1130, RxDropsPerSecond: 0x0, RxErrorsPerSecond: 0x0, TxPacketsPerSecond: 0x3, TxBitsPerSecond: 0x10cc, FpTxPacketsPerSecond: 0x0, FpTxBitsPerSecond: 0x0, TxDropsPerSecond: 0x0, TxQueueDropsPerSecond: 0x0, TxErrorsPerSecond: 0x0},
	{Name: "bridge", RxPacketsPerSecond: 0x8, RxBitsPerSecond: 0x125c, FpRxPacketsPerSecond: 0x8, FpRxBitsPerSecond: 0x125c, RxDropsPerSecond: 0x0, RxErrorsPerSecond: 0x0, TxPacketsPerSecond: 0x7, TxBitsPerSecond: 0x2bc0, FpTxPacketsPerSecond: 0x0, FpTxBitsPerSecond: 0x0, TxDropsPerSecond: 0x0, TxQueueDropsPerSecond: 0x0, TxErrorsPerSecond: 0x0},
}

func TestMonitor(t *testing.T) {
	ch := make(chan *TrafficMonitorResponse, 10)
	rp := TrafficMonitor{
		C: ch,
	}

	err := rp.ParseResponse(strings.NewReader(monData))
	if err != nil {
		t.Fatal(err)
	}

	i := 0
Loop:
	for {
		select {
		case v := <-ch:
			if *v != monRes[i] {
				t.Fatalf("got %v, want %v\n", v, monRes[i])
			}
			i++
		default:
			break Loop
		}
	}
}

var ifData = `Flags: D - dynamic, X - disabled, R - running, S - slave
 0  R  name="ether1" default-name="ether1" type="ether" mtu=1500
       actual-mtu=1500 l2mtu=1598 max-l2mtu=4074
       mac-address=CC:2D:E0:79:1C:85 last-link-down-time=oct/17/2018 12:19:57
       last-link-up-time=oct/17/2018 12:20:12 link-downs=1

 1   S name="ether2-master" default-name="ether2" type="ether" mtu=1500
       actual-mtu=1500 l2mtu=1598 max-l2mtu=4074
       mac-address=CC:2D:E0:79:1C:86 link-downs=0

 2   S name="ether3" default-name="ether3" type="ether" mtu=1500
       actual-mtu=1500 l2mtu=1598 max-l2mtu=4074
       mac-address=CC:2D:E0:79:1C:87 link-downs=0

 3   S name="ether4" default-name="ether4" type="ether" mtu=1500
       actual-mtu=1500 l2mtu=1598 max-l2mtu=4074
       mac-address=CC:2D:E0:79:1C:88 link-downs=0

 4   S name="ether5" default-name="ether5" type="ether" mtu=1500
       actual-mtu=1500 l2mtu=1598 max-l2mtu=4074
       mac-address=CC:2D:E0:79:1C:89 link-downs=0

 5   S name="sfp1" default-name="sfp1" type="ether" mtu=1500 actual-mtu=1500
       l2mtu=1600 max-l2mtu=4076 mac-address=CC:2D:E0:79:1C:8A link-downs=0

 6  RS name="wlan1" default-name="wlan1" type="wlan" mtu=1500 actual-mtu=1500
       l2mtu=1600 max-l2mtu=2290 mac-address=CC:2D:E0:79:1C:8C
       last-link-down-time=oct/18/2018 15:34:50
       last-link-up-time=oct/18/2018 17:59:02 link-downs=7

 7  RS name="wlan2" default-name="wlan2" type="wlan" mtu=1500 actual-mtu=1500
       l2mtu=1600 max-l2mtu=2290 mac-address=CC:2D:E0:79:1C:8B
       last-link-down-time=oct/19/2018 11:02:47
       last-link-up-time=oct/19/2018 11:08:45 link-downs=251

 8  R  ;;; defconf
       name="bridge" type="bridge" mtu=auto actual-mtu=1500 l2mtu=1598
       mac-address=CC:2D:E0:79:1C:86 last-link-up-time=oct/16/2018 21:17:30
       link-downs=0

 9  R  name="pppoe-out1" type="pppoe-out" mtu=1480 actual-mtu=1480
       last-link-down-time=oct/17/2018 08:01:04
       last-link-up-time=oct/17/2018 15:03:47 link-downs=1

`

var ifJSON = `[{"index":0,"flags":"R","name":"ether1","default_name":"ether1","type":"ether","mtu":1500,"actual_mtu":1500,"l2mtu":1598,"max_l2mtu":4074,"mac_address":"CC:2D:E0:79:1C:85","last_link_down_time":"2018-10-17T12:19:57Z","last_link_up_time":"2018-10-17T12:20:12Z","link_downs":1},{"index":1,"flags":"S","name":"ether2-master","default_name":"ether2","type":"ether","mtu":1500,"actual_mtu":1500,"l2mtu":1598,"max_l2mtu":4074,"mac_address":"CC:2D:E0:79:1C:86","link_downs":0},{"index":2,"flags":"S","name":"ether3","default_name":"ether3","type":"ether","mtu":1500,"actual_mtu":1500,"l2mtu":1598,"max_l2mtu":4074,"mac_address":"CC:2D:E0:79:1C:87","link_downs":0},{"index":3,"flags":"S","name":"ether4","default_name":"ether4","type":"ether","mtu":1500,"actual_mtu":1500,"l2mtu":1598,"max_l2mtu":4074,"mac_address":"CC:2D:E0:79:1C:88","link_downs":0},{"index":4,"flags":"S","name":"ether5","default_name":"ether5","type":"ether","mtu":1500,"actual_mtu":1500,"l2mtu":1598,"max_l2mtu":4074,"mac_address":"CC:2D:E0:79:1C:89","link_downs":0},{"index":5,"flags":"S","name":"sfp1","default_name":"sfp1","type":"ether","mtu":1500,"actual_mtu":1500,"l2mtu":1600,"max_l2mtu":4076,"mac_address":"CC:2D:E0:79:1C:8A","link_downs":0},{"index":6,"flags":"RS","name":"wlan1","default_name":"wlan1","type":"wlan","mtu":1500,"actual_mtu":1500,"l2mtu":1600,"max_l2mtu":2290,"mac_address":"CC:2D:E0:79:1C:8C","last_link_down_time":"2018-10-18T15:34:50Z","last_link_up_time":"2018-10-18T17:59:02Z","link_downs":7},{"index":7,"flags":"RS","name":"wlan2","default_name":"wlan2","type":"wlan","mtu":1500,"actual_mtu":1500,"l2mtu":1600,"max_l2mtu":2290,"mac_address":"CC:2D:E0:79:1C:8B","last_link_down_time":"2018-10-19T11:02:47Z","last_link_up_time":"2018-10-19T11:08:45Z","link_downs":251},{"index":8,"flags":"R","name":"bridge","type":"bridge","actual_mtu":1500,"l2mtu":1598,"mac_address":"CC:2D:E0:79:1C:86","last_link_up_time":"2018-10-16T21:17:30Z","link_downs":0},{"index":9,"flags":"R","name":"pppoe-out1","type":"pppoe-out","mtu":1480,"actual_mtu":1480,"last_link_down_time":"2018-10-17T08:01:04Z","last_link_up_time":"2018-10-17T15:03:47Z","link_downs":1}]`

func TestInterfaces(t *testing.T) {
	var ifs, ex Interfaces

	if err := ifs.ParseResponse(strings.NewReader(ifData)); err != nil {
		t.Fatal(err)
	}

	if err := json.Unmarshal([]byte(ifJSON), &ex); err != nil {
		t.Fatal(err)
	}

	if len(ifs) != len(ex) {
		t.Fatalf("got %v, want %v\n", len(ifs), len(ex))
	}

	for i := range ifs {
		v1 := *ifs[i]
		v1.LastLinkDownTime = nil
		v1.LastLinkUpTime = nil

		v2 := *ex[i]
		v2.LastLinkDownTime = nil
		v2.LastLinkUpTime = nil

		if v1 != v2 {
			t.Fatalf("got %v, want %v\n", v1, v2)
		}

		t1 := ifs[i].LastLinkUpTime
		t2 := ex[i].LastLinkUpTime

		if t1 == nil && t2 != nil || t1 != nil && t2 == nil || t1 != nil && t2 != nil && !t1.Equal(*t2) {
			t.Fatalf("got %v, want %v\n", *t1, *t2)
		}

		t1 = ifs[i].LastLinkDownTime
		t2 = ex[i].LastLinkDownTime

		if t1 == nil && t2 != nil || t1 != nil && t2 == nil || t1 != nil && t2 != nil && !t1.Equal(*t2) {
			t.Fatalf("got %v, want %v\n", t1, t2)
		}
	}
}
