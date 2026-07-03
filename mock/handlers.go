package mock

import (
	"encoding/json"
	"net/http"
	"strconv"
)

func (s *Server) createRule(w http.ResponseWriter, r *http.Request) {
	var rule FirewallRule
	if err := json.NewDecoder(r.Body).Decode(&rule); err != nil {
		writeResponse(w, http.StatusBadRequest, nil)
		return
	}
	rule.ID = int(s.nextID.Add(1))
	s.mu.Lock()
	s.rules[rule.ID] = rule
	s.mu.Unlock()
	writeResponse(w, http.StatusCreated, rule)
}

func (s *Server) getRule(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		writeResponse(w, http.StatusBadRequest, nil)
		return
	}
	s.mu.RLock()
	rule, ok := s.rules[id]
	s.mu.RUnlock()
	if !ok {
		writeResponse(w, http.StatusNotFound, nil)
		return
	}
	writeResponse(w, http.StatusOK, rule)
}

func (s *Server) updateRule(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		writeResponse(w, http.StatusBadRequest, nil)
		return
	}
	var update FirewallRule
	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		writeResponse(w, http.StatusBadRequest, nil)
		return
	}
	s.mu.Lock()
	existing, ok := s.rules[id]
	if !ok {
		s.mu.Unlock()
		writeResponse(w, http.StatusNotFound, nil)
		return
	}
	update.ID = existing.ID
	s.rules[id] = update
	s.mu.Unlock()
	writeResponse(w, http.StatusOK, update)
}

func (s *Server) deleteRule(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		writeResponse(w, http.StatusBadRequest, nil)
		return
	}
	s.mu.Lock()
	_, ok := s.rules[id]
	delete(s.rules, id)
	s.mu.Unlock()
	if !ok {
		writeResponse(w, http.StatusNotFound, nil)
		return
	}
	writeResponse(w, http.StatusOK, nil)
}

func (s *Server) authJWT(w http.ResponseWriter, _ *http.Request) {
	writeResponse(w, http.StatusOK, map[string]string{"token": "mock-jwt-token"})
}

func parseID(r *http.Request) (int, error) {
	return strconv.Atoi(r.URL.Query().Get("id"))
}
