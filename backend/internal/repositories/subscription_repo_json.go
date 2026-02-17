package repositories

import (
	"encoding/json"
	"os"
	"sync"
	"time"

	"streaming-system/internal/models"
)

type SubscriptionRepoJSON struct {
	filePath string
	mu       sync.Mutex
}

func NewSubscriptionRepoJSON(filePath string) *SubscriptionRepoJSON {
	return &SubscriptionRepoJSON{filePath: filePath}
}

func (r *SubscriptionRepoJSON) GetAll() ([]models.Subscription, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	data, err := os.ReadFile(r.filePath)
	if err != nil {
		return nil, err
	}

	if len(data) == 0 {
		return []models.Subscription{}, nil
	}

	var subs []models.Subscription
	if err := json.Unmarshal(data, &subs); err != nil {
		return nil, err
	}
	return subs, nil
}

func (r *SubscriptionRepoJSON) GetByUser(userID string) (*models.Subscription, error) {
	subs, err := r.GetAll()
	if err != nil {
		return nil, err
	}
	for _, s := range subs {
		if s.UserID == userID {
			ss := s
			return &ss, nil
		}
	}
	return nil, nil
}

func (r *SubscriptionRepoJSON) IsActive(userID string) (bool, error) {
	s, err := r.GetByUser(userID)
	if err != nil {
		return false, err
	}
	if s == nil {
		return false, nil
	}
	if s.Status != models.SubActive {
		return false, nil
	}
	return time.Now().Before(s.ExpiresAt), nil
}
