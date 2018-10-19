package response

import (
	"fmt"
	"unicode"
)

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

	for i < len(s) && (s[i] >= '0' && s[i] <= '9' || s[i] == '.' || unicode.IsSpace(rune(s[i]))) {
		if s[i] == '.' {
			dp = true
		} else if !unicode.IsSpace(rune(s[i])) {
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

func parseKeyValuePairs(s string, isFixed func(string) int) (values map[string]string, args []string) {
	values = make(map[string]string)

	var i int
	for i < len(s) {
		for i < len(s) && unicode.IsSpace(rune(s[i])) {
			i++
		}

		start := i
		for i < len(s) && !unicode.IsSpace(rune(s[i])) && s[i] != '=' {
			i++
		}

		token := s[start:i]

		if i == len(s) || unicode.IsSpace(rune(s[i])) {
			args = append(args, token)

			if i == len(s) {
				break
			}

			i++
			continue
		}

		i++
		if i == len(s) {
			break
		}

		if isFixed != nil {
			if w := isFixed(token); w > 0 {
				start = i
				for i < len(s) && w != 0 {
					i++
					w--
				}

				values[token] = s[start:i]
				continue
			}
		}

		// parameter
		if s[i] == '"' {
			i++
			str := make([]byte, 0, len(s)-i)

		Loop:
			for i < len(s) {
				switch s[i] {
				case '"':
					i++
					break Loop
				case '\\':
					if i < len(s)-1 && s[i+1] == '"' {
						str = append(str, '"')
						i++
					} else {
						str = append(str, s[i])
					}
				default:
					str = append(str, s[i])
				}
				i++
			}

			values[token] = string(str)
		} else {
			start = i
			for i < len(s) && !unicode.IsSpace(rune(s[i])) && s[i] != '=' {
				i++
			}

			values[token] = s[start:i]
		}
	}

	return
}
