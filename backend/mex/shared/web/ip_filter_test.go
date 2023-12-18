package web

import (
	"net/http"
	"net/http/httptest"
	"testing"

	L "github.com/d4l-data4life/mex/mex/shared/log"
)

type InOut struct {
	addr  string
	block bool
}

func TestIPFilter(t *testing.T) {
	tests := []struct {
		name      string
		whitelist []string
		cases     []InOut
	}{
		{
			name:      "A",
			whitelist: []string{"1.2.3.4"},
			cases: []InOut{
				{addr: "1.2.3.4", block: false},
				{addr: "1.2.3.5", block: true},
				{addr: "0.0.0.0", block: true},
			},
		},
		{
			name:      "B",
			whitelist: []string{"100.100.100.100/24"},
			cases: []InOut{
				{addr: "1.2.3.4", block: true},
				{addr: "0.0.0.0", block: true},
				{addr: "100.100.100.0", block: false},
				{addr: "100.100.100.128", block: false},
				{addr: "100.100.100.255", block: false},
				{addr: "100.100.101.0", block: true},
			},
		},
		{
			name:      "C",
			whitelist: []string{"100.100.100.100/24", "1.2.3.4/16", "192.168.0.1"},
			cases: []InOut{
				{addr: "0.0.0.0", block: true},
				{addr: "100.100.100.0", block: false},
				{addr: "100.100.100.128", block: false},
				{addr: "100.100.101.0", block: true},
				{addr: "1.2.3.4", block: false},
				{addr: "1.2.5.123", block: false},
				{addr: "1.2.125.123", block: false},
				{addr: "192.168.0.1", block: false},
				{addr: "192.168.0.0", block: true},
			},
		},
	}

	log := &L.NullLogger{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := NewIPFilter(log, tt.whitelist)

			for _, c := range tt.cases {
				rr := httptest.NewRecorder()
				req, _ := http.NewRequest("GET", "/", nil)
				req.Header.Set("X-Real-IP", c.addr)

				f(MyHandler{}).ServeHTTP(rr, req)

				if (rr.Code == http.StatusForbidden) != c.block {
					t.Errorf("incorrect IP filter behavior: block wanted: %v, blocked: %v, address: %s", c.block, rr.Code == http.StatusForbidden, c.addr)
				}
			}
		})
	}
}

type MyHandler struct{}

func (x MyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
}
