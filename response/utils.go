package response

import "fmt"

var bwSuffixes = map[string]uint64{
	"bps":  1,
	"kbps": 1000,
	"Mbps": 1000 * 1000,
	"Gbps": 1000 * 1000 * 1000,
}

func parseBW(s string) (uint64, error) {
	var (
		val uint64
		div uint64 = 1
		i   int
		dp  bool
	)

	for i < len(s) && (s[i] >= '0' && s[i] <= '9' || s[i] == '.' || s[i] == ' ') {
		if s[i] == '.' {
			dp = true
		} else if s[i] != ' ' {
			val = val*10 + uint64(s[i]) - uint64('0')
			if dp {
				div *= 10
			}
		}
		i++
	}

	var mul uint64 = 1
	if i < len(s) {
		if m, ok := bwSuffixes[s[i:]]; ok {
			mul = m
		} else {
			return 0, fmt.Errorf("Unknown suffix `%s'", s[i:])
		}
	}

	return val * mul / div, nil
}
