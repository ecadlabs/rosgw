package response

import (
	"bufio"
	"io"
	"strconv"
	"strings"
	"time"
)

// Command: interface print detail
type Interface struct {
	Index            int        `json:"index"`
	Flags            string     `json:"flags"`
	Name             string     `json:"name"`
	DefaultName      string     `json:"default_name,omitempty"`
	Type             string     `json:"type,omitempty"`
	MTU              int        `json:"mtu,omitempty"`
	ActualMTU        int        `json:"actual_mtu,omitempty"`
	L2MTU            int        `json:"l2mtu,omitempty"`
	MaxL2MTU         int        `json:"max_l2mtu,omitempty"`
	MACAddress       string     `json:"mac_address,omitempty"`
	LastLinkDownTime *time.Time `json:"last_link_down_time,omitempty"`
	LastLinkUpTime   *time.Time `json:"last_link_up_time,omitempty"`
	LinkDowns        int        `json:"link_downs"`
}

type Interfaces []*Interface

const timeFmt = "Jan/02/2006 15:04:05"

func (ifs *Interfaces) ParseResponse(r io.Reader) error {
	iface := new(Interface)
	var hasData bool

	s := bufio.NewScanner(r)
	for s.Scan() {
		line := s.Text()

		if strings.HasPrefix(line, "Flags:") {
			continue
		}

		if line == "" && hasData {
			*ifs = append(*ifs, iface)
			iface = new(Interface)
			hasData = false
			continue
		}

		hasData = true

		values, args := parseKeyValuePairs(line, func(s string) int {
			// Time values aren't enquoted
			if strings.HasSuffix(s, "-time") {
				return 20
			}
			return -1
		})

		if len(args) >= 2 {
			// New entry header
			i, _ := strconv.ParseInt(args[0], 10, 32)
			iface.Index = int(i)
			iface.Flags = args[1]
		}

		for key, val := range values {
			switch key {
			case "name":
				iface.Name = val
			case "default-name":
				iface.DefaultName = val
			case "type":
				iface.Type = val
			case "mac-address":
				iface.MACAddress = val
			case "mtu", "actual-mtu", "l2mtu", "max-l2mtu", "link-downs":
				i, _ := strconv.ParseInt(val, 10, 32)
				switch key {
				case "mtu":
					iface.MTU = int(i)
				case "actual-mtu":
					iface.ActualMTU = int(i)
				case "l2mtu":
					iface.L2MTU = int(i)
				case "max-l2mtu":
					iface.MaxL2MTU = int(i)
				case "link-downs":
					iface.LinkDowns = int(i)
				}
			case "last-link-down-time", "last-link-up-time":
				if tv, err := time.Parse(timeFmt, val); err == nil {
					switch key {
					case "last-link-down-time":
						iface.LastLinkDownTime = &tv
					case "last-link-up-time":
						iface.LastLinkUpTime = &tv
					}
				}
			}
		}
	}

	if hasData {
		*ifs = append(*ifs, iface)
	}

	return s.Err()
}

var _ ResponseParser = &Interfaces{}
