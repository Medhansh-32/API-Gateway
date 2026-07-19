package service

import (
	"log"
	"sync"

	"github.com/medhansh-32/api-gateway/internal/models"
)

type wrrState struct {
	mu      sync.Mutex
	weights map[string]int
}

func newWRRState() *wrrState {
	return &wrrState{weights: make(map[string]int)}
}

func (s *wrrState) getTarget(serviceID string, targets []models.TargetConfig) (models.TargetConfig, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(targets) == 0 {
		log.Printf("[WRR] service=%s no targets provided", serviceID)
		return models.TargetConfig{}, false
	}

	var total int
	var best *models.TargetConfig
	var bestKey string

	for i := range targets {
		t := targets[i]
		w := t.Weight
		if w <= 0 {
			log.Printf("[WRR] service=%s target=%s has weight<=0 (%d), defaulting to 1", serviceID, t.URL, t.Weight)
			w = 1
		}

		s.weights[t.URL] += w
		total += w

		log.Printf("[WRR] service=%s target=%s weight=%d currentScore=%d", serviceID, t.URL, w, s.weights[t.URL])

		if best == nil || s.weights[t.URL] > s.weights[bestKey] {
			best = &targets[i]
			bestKey = t.URL
		}
	}

	log.Printf("[WRR] service=%s total=%d picked=%s scoreBeforeDeduction=%d", serviceID, total, bestKey, s.weights[bestKey])

	s.weights[bestKey] -= total

	log.Printf("[WRR] service=%s picked=%s scoreAfterDeduction=%d allScores=%v", serviceID, bestKey, s.weights[bestKey], s.weights)

	return *best, true
}

type WRRManager struct {
	mu       sync.Mutex
	services map[string]*wrrState
}

func NewWRRManager() *WRRManager {
	return &WRRManager{services: make(map[string]*wrrState)}
}

func (m *WRRManager) GetTarget(serviceID string, targets []models.TargetConfig) (models.TargetConfig, bool) {
	m.mu.Lock()
	state, ok := m.services[serviceID]
	if !ok {
		log.Printf("[WRR] service=%s creating new wrrState (first request for this service)", serviceID)
		state = newWRRState()
		m.services[serviceID] = state
	}
	m.mu.Unlock()

	target, ok := state.getTarget(serviceID, targets)
	if !ok {
		log.Printf("[WRR] service=%s GetTarget failed, no target returned", serviceID)
	} else {
		log.Printf("[WRR] service=%s ==> returning target=%s", serviceID, target.URL)
	}
	return target, ok
}