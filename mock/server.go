package mock

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
)

// FirewallRule mirrors the pfSense REST API v2 firewall rule object.
type FirewallRule struct {
	ID               int      `json:"id"`
	Type             string   `json:"type"`
	Interface        []string `json:"interface"`
	Ipprotocol       string   `json:"ipprotocol"`
	Source           string   `json:"source"`
	Destination      string   `json:"destination"`
	Descr            string   `json:"descr,omitempty"`
	Disabled         bool     `json:"disabled,omitempty"`
	Floating         bool     `json:"floating,omitempty"`
	Direction        string   `json:"direction,omitempty"`
	Protocol         string   `json:"protocol,omitempty"`
	Log              bool     `json:"log,omitempty"`
	Gateway          string   `json:"gateway,omitempty"`
	Statetype        string   `json:"statetype,omitempty"`
	Quick            bool     `json:"quick,omitempty"`
	Sched            string   `json:"sched,omitempty"`
	Tracker          int64    `json:"tracker,omitempty"`
	Tag              string   `json:"tag,omitempty"`
	CreatedBy        string   `json:"created_by,omitempty"`
	CreatedTime      int64    `json:"created_time,omitempty"`
	UpdatedBy        string   `json:"updated_by,omitempty"`
	UpdatedTime      int64    `json:"updated_time,omitempty"`
	AssociatedRuleId string   `json:"associated_rule_id,omitempty"`
	Ackqueue         string   `json:"ackqueue,omitempty"`
	Defaultqueue     string   `json:"defaultqueue,omitempty"`
	Dnpipe           string   `json:"dnpipe,omitempty"`
	Pdnpipe          string   `json:"pdnpipe,omitempty"`
	SourcePort       string   `json:"source_port,omitempty"`
	DestinationPort  string   `json:"destination_port,omitempty"`
}

// Server is a mock pfSense REST API server backed by in-memory state.
type Server struct {
	*httptest.Server
	rules  map[int]FirewallRule
	nextID atomic.Int32
	mu     sync.RWMutex
}

// NewServer starts a mock pfSense REST API server.
func NewServer() *Server {
	s := &Server{rules: make(map[int]FirewallRule)}
	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/v2/firewall/rule", s.createRule)
	mux.HandleFunc("GET /api/v2/firewall/rule", s.getRule)
	mux.HandleFunc("PUT /api/v2/firewall/rule", s.updateRule)
	mux.HandleFunc("DELETE /api/v2/firewall/rule", s.deleteRule)
	mux.HandleFunc("POST /api/v2/auth/jwt", s.authJWT)
	s.Server = httptest.NewServer(mux)
	return s
}
