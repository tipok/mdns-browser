package discovery

import (
	"context"
	"fmt"
	"mdns-browser/internal/data"
	"strconv"
	"strings"

	"github.com/hashicorp/mdns"
)

func unescapeDNSName(s string) string {
	var b strings.Builder
	for i := 0; i < len(s); i++ {
		if s[i] == '\\' && i+1 < len(s) {
			// Numeric escape \DDD (three decimal digits)
			if i+3 < len(s) &&
				s[i+1] >= '0' && s[i+1] <= '9' &&
				s[i+2] >= '0' && s[i+2] <= '9' &&
				s[i+3] >= '0' && s[i+3] <= '9' {
				val, err := strconv.Atoi(s[i+1 : i+4])
				if err == nil && val >= 0 && val <= 255 {
					b.WriteByte(byte(val))
					i += 3
					continue
				}
			}
			// Single char escape: skip the backslash, keep next char
			i++
			b.WriteByte(s[i])
			continue
		}
		b.WriteByte(s[i])
	}
	return b.String()
}

func ListAllServices(ctx context.Context, addCh chan data.ListItem) error {
	entriesCh := make(chan *mdns.ServiceEntry, 100)
	go func() {
		defer close(addCh)
		for {
			select {
			case <-ctx.Done():
				return
			case entry, ok := <-entriesCh:
				if !ok {
					return
				}
				it := data.ListItem{
					Name:       unescapeDNSName(entry.Name),
					Host:       entry.Host,
					AddrV4:     entry.AddrV4.String(),
					AddrV6:     entry.AddrV6IPAddr.String(),
					Port:       entry.Port,
					Info:       entry.Info,
					InfoFields: entry.InfoFields,
				}
				select {
				case <-ctx.Done():
					return
				case addCh <- it:
				}
			}
		}
	}()

	for _, svc := range Services {
		select {
		case <-ctx.Done():
			close(entriesCh)
			return ctx.Err()
		default:
		}

		params := mdns.DefaultParams(svc)
		params.Entries = entriesCh
		params.Logger = NoopLogLogger
		err := mdns.QueryContext(ctx, params)
		if err != nil {
			close(entriesCh)
			return fmt.Errorf("error querying for %s: %s", svc, err)
		}
	}

	close(entriesCh)

	return nil
}
