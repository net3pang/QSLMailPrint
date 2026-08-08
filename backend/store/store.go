package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type Store struct {
	Tasks     []map[string]any `json:"tasks"`
	Contacts  []map[string]any `json:"contacts"`
	UpdatedAt string           `json:"updatedAt"`
}

type Database struct {
	mu   sync.Mutex
	path string
	data Store
}

func Open(configPath string) (*Database, error) {
	path := configPath
	if path == "" {
		path = os.Getenv("QSL_MAIL_DB")
	}
	if path == "" {
		base, err := os.UserConfigDir()
		if err != nil {
			return nil, err
		}
		path = filepath.Join(base, "QSLMail", "qsl-mail-data.json")
	}
	database := &Database{path: path, data: Store{Tasks: []map[string]any{}, Contacts: []map[string]any{}}}
	content, err := os.ReadFile(path)
	if err == nil && len(content) > 0 {
		if err := json.Unmarshal(content, &database.data); err != nil {
			return nil, fmt.Errorf("读取本地数据失败: %w", err)
		}
	} else if !errors.Is(err, os.ErrNotExist) {
		return nil, err
	}
	return database, nil
}

func (database *Database) collection(name string) (*[]map[string]any, error) {
	if name == "tasks" {
		return &database.data.Tasks, nil
	}
	if name == "contacts" {
		return &database.data.Contacts, nil
	}
	return nil, errors.New("未知数据集合")
}

func (database *Database) saveLocked() error {
	database.data.UpdatedAt = time.Now().Format(time.RFC3339)
	if err := os.MkdirAll(filepath.Dir(database.path), 0755); err != nil {
		return err
	}
	content, err := json.MarshalIndent(database.data, "", "  ")
	if err != nil {
		return err
	}
	temporary := database.path + ".tmp"
	if err := os.WriteFile(temporary, content, 0600); err != nil {
		return err
	}
	return os.Rename(temporary, database.path)
}

func (database *Database) Collection(name string) ([]map[string]any, error) {
	database.mu.Lock()
	defer database.mu.Unlock()
	collection, err := database.collection(name)
	if err != nil {
		return nil, err
	}
	return append([]map[string]any(nil), (*collection)...), nil
}

func (database *Database) SaveRecord(name string, item map[string]any) (map[string]any, error) {
	database.mu.Lock()
	defer database.mu.Unlock()
	collection, err := database.collection(name)
	if err != nil {
		return nil, err
	}
	item["updatedAt"] = time.Now().Format(time.RFC3339)
	updated := false
	if id, ok := item["id"].(string); ok && id != "" {
		for index, existing := range *collection {
			if existing["id"] == id {
				(*collection)[index] = item
				updated = true
				break
			}
		}
	}
	if !updated {
		*collection = append([]map[string]any{item}, (*collection)...)
	}
	if err := database.saveLocked(); err != nil {
		return nil, err
	}
	return item, nil
}

func (database *Database) DeleteRecord(name, id string) error {
	database.mu.Lock()
	defer database.mu.Unlock()
	collection, err := database.collection(name)
	if err != nil {
		return err
	}
	filtered := (*collection)[:0]
	for _, item := range *collection {
		if item["id"] != id {
			filtered = append(filtered, item)
		}
	}
	*collection = filtered
	return database.saveLocked()
}
