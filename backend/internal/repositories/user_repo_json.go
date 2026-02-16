package repositories

import (
	"encoding/json"
	"errors"
	"os"
	"strings"
	"sync"

	"streaming-system/internal/models"
)

type UserRepoJSON struct {
	filePath string
	mu       sync.Mutex
}

func NewUserRepoJSON(filePath string) *UserRepoJSON {
	return &UserRepoJSON{
		filePath: filePath,
	}
}

// Obtener todos los usuarios
func (r *UserRepoJSON) GetAll() ([]models.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	file, err := os.ReadFile(r.filePath)
	if err != nil {
		return nil, err
	}

	var users []models.User
	if len(file) == 0 {
		return []models.User{}, nil
	}

	err = json.Unmarshal(file, &users)
	if err != nil {
		return nil, err
	}

	return users, nil
}

// Buscar usuario por email
func (r *UserRepoJSON) FindByEmail(email string) (*models.User, error) {
	users, err := r.GetAll()
	if err != nil {
		return nil, err
	}

	for _, u := range users {
		if strings.EqualFold(u.Email, email) {
			return &u, nil
		}
	}

	return nil, nil
}

// Crear usuario nuevo
func (r *UserRepoJSON) Create(user models.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	file, err := os.ReadFile(r.filePath)
	if err != nil {
		return err
	}

	var users []models.User
	if len(file) > 0 {
		if err := json.Unmarshal(file, &users); err != nil {
			return err
		}
	}

	// Validar duplicado
	for _, u := range users {
		if strings.EqualFold(u.Email, user.Email) {
			return errors.New("email already exists")
		}
	}

	users = append(users, user)

	updatedData, err := json.MarshalIndent(users, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(r.filePath, updatedData, 0644)
}

func (r *UserRepoJSON) FindByID(id string) (*models.User, error) {
	users, err := r.GetAll()
	if err != nil {
		return nil, err
	}

	for _, u := range users {
		if u.ID == id {
			uu := u
			return &uu, nil
		}
	}
	return nil, nil
}
