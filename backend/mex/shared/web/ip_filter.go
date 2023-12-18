package web

import (
	"context"
	"net"
	"net/http"

	"github.com/d4l-data4life/mex/mex/shared/constants"
	L "github.com/d4l-data4life/mex/mex/shared/log"
)

type IPFilter struct {
	log L.Logger

	// List of IP addresses or CIDR strings
	allowedIPNets []*net.IPNet

	next http.Handler
}

func NewIPFilter(log L.Logger, whitelistedIPs []string) func(http.Handler) http.Handler {
	ctx := context.Background()

	allowedIPNets := []*net.IPNet{}
	for _, allowedIP := range whitelistedIPs {
		ip := net.ParseIP(allowedIP)
		if ip != nil {
			allowedIPNets = append(allowedIPNets, &net.IPNet{
				IP:   ip,
				Mask: net.CIDRMask(32, 32),
			})
			continue
		}

		_, ipNet, err := net.ParseCIDR(allowedIP)
		if err == nil {
			allowedIPNets = append(allowedIPNets, ipNet)
		}
	}

	log.Info(ctx, L.Messagef("%d allowed IPs: %v", len(whitelistedIPs), whitelistedIPs))
	if len(whitelistedIPs) == 0 {
		log.Warn(ctx, L.Message("IP whitelist is empty; no request will pass"))
	}

	return func(next http.Handler) http.Handler {
		return &IPFilter{log: log, allowedIPNets: allowedIPNets, next: next}
	}
}

func (f *IPFilter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	realRemoteIP := r.Header.Get(constants.HTTPHeaderRealIP)
	if realRemoteIP == "" {
		realRemoteIP = r.RemoteAddr
	}

	ip := net.ParseIP(realRemoteIP)
	if ip == nil {
		f.log.Warn(r.Context(), L.Messagef("IP denied, because it could not be parsed: [%s]", r.RemoteAddr))
		w.WriteHeader(http.StatusForbidden)
	}

	if f.isAllowed(ip) {
		f.next.ServeHTTP(w, r)
	} else {
		f.log.Warn(r.Context(), L.Messagef("IP denied: [%s]", r.RemoteAddr))
		w.WriteHeader(http.StatusForbidden)
	}
}

func (f *IPFilter) isAllowed(ip net.IP) bool {
	for _, allowedIPNet := range f.allowedIPNets {
		if allowedIPNet.Contains(ip) {
			return true
		}
	}
	return false
}
