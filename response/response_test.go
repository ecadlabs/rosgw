package response

import (
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
	rd := strings.NewReader(monData)
	ch := make(chan *TrafficMonitorResponse, 10)
	rp := TrafficMonitor{
		C: ch,
	}

	err := rp.ParseResponse(rd)
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
