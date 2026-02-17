package repositories

import (
	"encoding/json"
	"os"
	"sync"
	"time"

	"streaming-system/internal/models"
)

type ProgressRepoJSON struct {
	filePath string
	mu       sync.Mutex
}

func NewProgressRepoJSON(filePath string) *ProgressRepoJSON {
	return &ProgressRepoJSON{filePath: filePath}
}

func (r *ProgressRepoJSON) GetAll() ([]models.PlaybackProgress, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	data, err := os.ReadFile(r.filePath)
	if err != nil {
		return nil, err
	}

	if len(data) == 0 {
		return []models.PlaybackProgress{}, nil
	}

	var items []models.PlaybackProgress
	if err := json.Unmarshal(data, &items); err != nil {
		return nil, err
	}
	return items, nil
}

func (r *ProgressRepoJSON) SaveAll(items []models.PlaybackProgress) error {
	out, err := json.MarshalIndent(items, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(r.filePath, out, 0644)
}

func (r *ProgressRepoJSON) GetByUserAndContent(userID, contentID string) (*models.PlaybackProgress, error) {
	items, err := r.GetAll()
	if err != nil {
		return nil, err
	}
	for _, it := range items {
		if it.UserID == userID && it.ContentID == contentID {
			tmp := it
			return &tmp, nil
		}
	}
	return nil, nil
}

func (r *ProgressRepoJSON) Upsert(p models.PlaybackProgress) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	data, err := os.ReadFile(r.filePath)
	if err != nil {
		return err
	}

	var items []models.PlaybackProgress
	if len(data) > 0 {
		if err := json.Unmarshal(data, &items); err != nil {
			return err
		}
	}

	found := false
	for i := range items {
		if items[i].UserID == p.UserID && items[i].ContentID == p.ContentID {
			p.CreatedAt = items[i].CreatedAt
			p.UpdatedAt = time.Now()
			items[i] = p
			found = true
			break
		}
	}

	if !found {
		p.CreatedAt = time.Now()
		p.UpdatedAt = time.Now()
		items = append(items, p)
	}

	out, err := json.MarshalIndent(items, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(r.filePath, out, 0644)
}
