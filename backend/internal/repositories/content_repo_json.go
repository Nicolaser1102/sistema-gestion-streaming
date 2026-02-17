package repositories

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"

	"time"

	"streaming-system/internal/models"
)

type ContentRepoJSON struct {
	filePath string
	mu       sync.Mutex
}

func NewContentRepoJSON(filePath string) *ContentRepoJSON {
	return &ContentRepoJSON{filePath: filePath}
}

func (r *ContentRepoJSON) GetAll() ([]models.Content, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	data, err := os.ReadFile(r.filePath)
	if err != nil {
		return nil, err
	}

	var contents []models.Content
	if len(data) == 0 {
		return []models.Content{}, nil
	}

	if err := json.Unmarshal(data, &contents); err != nil {
		return nil, err
	}

	return contents, nil
}

func (r *ContentRepoJSON) FindByID(id string) (*models.Content, error) {
	contents, err := r.GetAll()
	if err != nil {
		return nil, err
	}

	for _, c := range contents {
		if c.ID == id && c.Active {
			cc := c
			return &cc, nil
		}
	}
	return nil, nil
}

func (r *ContentRepoJSON) SearchAndFilter(query, genre, ctype string, year int) ([]models.Content, error) {
	contents, err := r.GetAll()
	if err != nil {
		return nil, err
	}

	var result []models.Content

	for _, c := range contents {

		if !c.Active {
			continue
		}

		if query != "" && !strings.Contains(strings.ToLower(c.Title), strings.ToLower(query)) {
			continue
		}

		if genre != "" && !strings.EqualFold(c.Genre, genre) {
			continue
		}

		if ctype != "" && string(c.Type) != strings.ToUpper(ctype) {
			continue
		}

		if year != 0 && c.Year != year {
			continue
		}

		result = append(result, c)
	}

	return result, nil
}

func (r *ContentRepoJSON) Create(content models.Content) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	data, err := os.ReadFile(r.filePath)
	if err != nil {
		return err
	}

	var contents []models.Content
	if len(data) > 0 {
		if err := json.Unmarshal(data, &contents); err != nil {
			return err
		}
	}

	contents = append(contents, content)

	out, err := json.MarshalIndent(contents, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(r.filePath, out, 0644)
}

func (r *ContentRepoJSON) Update(id string, updated models.Content) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	data, err := os.ReadFile(r.filePath)
	if err != nil {
		return err
	}

	var contents []models.Content
	if len(data) > 0 {
		if err := json.Unmarshal(data, &contents); err != nil {
			return err
		}
	}

	found := false

	for i, c := range contents {
		if c.ID == id {
			updated.ID = id
			updated.CreatedAt = c.CreatedAt
			updated.UpdatedAt = time.Now()
			contents[i] = updated
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("contenido no encontrado")
	}

	out, err := json.MarshalIndent(contents, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(r.filePath, out, 0644)
}

func (r *ContentRepoJSON) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	data, err := os.ReadFile(r.filePath)
	if err != nil {
		return err
	}

	var contents []models.Content
	if len(data) > 0 {
		if err := json.Unmarshal(data, &contents); err != nil {
			return err
		}
	}

	found := false

	for i, c := range contents {
		if c.ID == id {
			contents[i].Active = false
			contents[i].UpdatedAt = time.Now()
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("contenido no encontrado")
	}

	out, err := json.MarshalIndent(contents, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(r.filePath, out, 0644)
}

func (r *ContentRepoJSON) GetAllIncludingInactive() ([]models.Content, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	data, err := os.ReadFile(r.filePath)
	if err != nil {
		return nil, err
	}

	var contents []models.Content
	if len(data) == 0 {
		return []models.Content{}, nil
	}

	if err := json.Unmarshal(data, &contents); err != nil {
		return nil, err
	}

	return contents, nil
}
