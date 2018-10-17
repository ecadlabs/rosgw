package response

import (
	"bufio"
	"io"
	"strings"
)

type TrafficMonitorResponse struct {
	Name                  string `json:"name"`
	RxPacketsPerSecond    uint64 `json:"rx_packets_per_second"`
	RxBitsPerSecond       uint64 `json:"rx_bits_per_second"`
	FpRxPacketsPerSecond  uint64 `json:"fp_rx_packets_per_second"`
	FpRxBitsPerSecond     uint64 `json:"fp_rx_bits_per_second"`
	RxDropsPerSecond      uint64 `json:"rx_drops_per_second"`
	RxErrorsPerSecond     uint64 `json:"rx_errors_per_second"`
	TxPacketsPerSecond    uint64 `json:"tx_packets_per_second"`
	TxBitsPerSecond       uint64 `json:"tx_bits_per_second"`
	FpTxPacketsPerSecond  uint64 `json:"fp_tx_packets_per_second"`
	FpTxBitsPerSecond     uint64 `json:"fp_tx_bits_per_second"`
	TxDropsPerSecond      uint64 `json:"tx_drops_per_second"`
	TxQueueDropsPerSecond uint64 `json:"tx_queue_drops_per_second"`
	TxErrorsPerSecond     uint64 `json:"tx_errors_per_second"`
}

type TrafficMonitor struct {
	C        chan<- *TrafficMonitorResponse
	NonBlock bool
}

func (t *TrafficMonitor) ParseResponse(r io.Reader) error {
	s := bufio.NewScanner(r)

	var hasData bool
	res := new(TrafficMonitorResponse)
	for s.Scan() {
		line := s.Text()
		if line == "" {
			if t.NonBlock {
				select {
				case t.C <- res:
				default:
				}
			} else {
				t.C <- res
			}
			res = new(TrafficMonitorResponse)
			hasData = false
			continue
		}

		f := strings.SplitN(line, ":", 2)
		if len(f) != 2 {
			continue
		}

		hasData = true

		key := strings.TrimSpace(f[0])
		val := strings.TrimSpace(f[1])

		switch key {
		case "name":
			res.Name = val
		default:
			v, _ := parseBW(val)
			switch key {
			case "rx-packets-per-second":
				res.RxPacketsPerSecond = v
			case "rx-bits-per-second":
				res.RxBitsPerSecond = v
			case "fp-rx-packets-per-second":
				res.FpRxPacketsPerSecond = v
			case "fp-rx-bits-per-second":
				res.FpRxBitsPerSecond = v
			case "rx-drops-per-second":
				res.RxDropsPerSecond = v
			case "rx-errors-per-second":
				res.RxErrorsPerSecond = v
			case "tx-packets-per-second":
				res.TxPacketsPerSecond = v
			case "tx-bits-per-second":
				res.TxBitsPerSecond = v
			case "fp-tx-packets-per-second":
				res.FpTxPacketsPerSecond = v
			case "fp-tx-bits-per-second":
				res.FpTxBitsPerSecond = v
			case "tx-drops-per-second":
				res.TxDropsPerSecond = v
			case "tx-queue-drops-per-second":
				res.TxQueueDropsPerSecond = v
			case "tx-errors-per-second":
				res.TxErrorsPerSecond = v
			}
		}
	}

	if hasData {
		t.C <- res
	}

	return s.Err()
}
