package repositories

import (
	"encoding/json"
	"os"
	"sync"
	"time"

	"streaming-system/internal/models"
)

type MyListRepoJSON struct {
	filePath string
	mu       sync.Mutex
}

func NewMyListRepoJSON(filePath string) *MyListRepoJSON {
	return &MyListRepoJSON{filePath: filePath}
}

func (r *MyListRepoJSON) GetAll() ([]models.MyListItem, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	data, err := os.ReadFile(r.filePath)
	if err != nil {
		return nil, err
	}

	var items []models.MyListItem
	if len(data) == 0 {
		return []models.MyListItem{}, nil
	}

	if err := json.Unmarshal(data, &items); err != nil {
		return nil, err
	}
	return items, nil
}

func (r *MyListRepoJSON) SaveAll(items []models.MyListItem) error {
	out, err := json.MarshalIndent(items, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(r.filePath, out, 0644)
}

func (r *MyListRepoJSON) GetByUser(userID string) ([]models.MyListItem, error) {
	items, err := r.GetAll()
	if err != nil {
		return nil, err
	}

	var result []models.MyListItem
	for _, it := range items {
		if it.UserID == userID {
			result = append(result, it)
		}
	}
	return result, nil
}

func (r *MyListRepoJSON) Add(userID, contentID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	data, err := os.ReadFile(r.filePath)
	if err != nil {
		return err
	}

	var items []models.MyListItem
	if len(data) > 0 {
		if err := json.Unmarshal(data, &items); err != nil {
			return err
		}
	}

	// evitar duplicados
	for _, it := range items {
		if it.UserID == userID && it.ContentID == contentID {
			return nil
		}
	}

	items = append(items, models.MyListItem{
		UserID:    userID,
		ContentID: contentID,
		CreatedAt: time.Now(),
	})

	out, err := json.MarshalIndent(items, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(r.filePath, out, 0644)
}

func (r *MyListRepoJSON) Remove(userID, contentID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	data, err := os.ReadFile(r.filePath)
	if err != nil {
		return err
	}

	var items []models.MyListItem
	if len(data) > 0 {
		if err := json.Unmarshal(data, &items); err != nil {
			return err
		}
	}

	var filtered []models.MyListItem
	for _, it := range items {
		if it.UserID == userID && it.ContentID == contentID {
			continue
		}
		filtered = append(filtered, it)
	}

	out, err := json.MarshalIndent(filtered, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(r.filePath, out, 0644)
}
