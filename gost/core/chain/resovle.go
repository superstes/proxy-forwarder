package chain

import (
	"context"
	"fmt"
	"net"

	"proxy_forwarder/gost/core/hosts"
	"proxy_forwarder/gost/core/resolver"
	"proxy_forwarder/log"
)

func Resolve(ctx context.Context, network, addr string, r resolver.Resolver, hosts hosts.HostMapper) (string, error) {
	if addr == "" {
		return addr, nil
	}

	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		return "", err
	}
	if host == "" {
		return addr, nil
	}

	if hosts != nil {
		if ips, _ := hosts.Lookup(ctx, network, host); len(ips) > 0 {
			log.Debug("resolve", fmt.Sprintf("hit host mapper: %s -> %s", host, ips))
			return net.JoinHostPort(ips[0].String(), port), nil
		}
	}

	if r != nil {
		ips, err := r.Resolve(ctx, network, host)
		if err != nil {
			if err == resolver.ErrInvalid {
				return addr, nil
			}
			log.Error("resolve", err)
		}
		if len(ips) == 0 {
			return "", fmt.Errorf("resolver: domain %s does not exist", host)
		}
		return net.JoinHostPort(ips[0].String(), port), nil
	}
	return addr, nil
}
