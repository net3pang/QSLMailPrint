package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// The Go service uses one atomic JSON document so it can run without a C toolchain
// or external database driver. Electron can later replace this store with SQLite
// without changing the renderer API.
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

func openDatabase() (*Database, error) {
	path := os.Getenv("QSL_MAIL_DB")
	if path == "" {
		base, err := os.UserConfigDir()
		if err != nil {
			return nil, err
		}
		path = filepath.Join(base, "QSLMail", "qsl-mail-data.json")
	}
	database := &Database{path: path, data: Store{Tasks: []map[string]any{}, Contacts: []map[string]any{}}}
	if content, err := os.ReadFile(path); err == nil && len(content) > 0 {
		if err := json.Unmarshal(content, &database.data); err != nil {
			return nil, fmt.Errorf("读取本地数据失败: %w", err)
		}
	} else if !errors.Is(err, os.ErrNotExist) {
		return nil, err
	}
	return database, nil
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

func (database *Database) collection(name string) ([]map[string]any, error) {
	if name == "tasks" {
		return database.data.Tasks, nil
	}
	if name == "contacts" {
		return database.data.Contacts, nil
	}
	return nil, errors.New("未知数据集合")
}

func (database *Database) handleCollection(name string, writer http.ResponseWriter, request *http.Request) {
	database.mu.Lock()
	defer database.mu.Unlock()
	collection, err := database.collection(name)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusNotFound)
		return
	}
	if request.Method == http.MethodGet {
		writer.Header().Set("Content-Type", "application/json; charset=utf-8")
		_ = json.NewEncoder(writer).Encode(collection)
		return
	}
	if request.Method != http.MethodPost {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var item map[string]any
	if err := json.NewDecoder(request.Body).Decode(&item); err != nil {
		http.Error(writer, "JSON 数据无效", http.StatusBadRequest)
		return
	}
	item["updatedAt"] = time.Now().Format(time.RFC3339)
	updated := false
	if id, ok := item["id"].(string); ok && id != "" {
		for index, existing := range collection {
			if existing["id"] == id {
				collection[index] = item
				updated = true
				break
			}
		}
	}
	if !updated {
		collection = append([]map[string]any{item}, collection...)
	}
	if name == "tasks" {
		database.data.Tasks = collection
	} else {
		database.data.Contacts = collection
	}
	if err := database.saveLocked(); err != nil {
		http.Error(writer, "保存本地数据失败", http.StatusInternalServerError)
		return
	}
	writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(writer).Encode(item)
}

func main() {
	database, err := openDatabase()
	if err != nil {
		panic(err)
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/api/health", func(writer http.ResponseWriter, _ *http.Request) { _, _ = writer.Write([]byte(`{"ok":true}`)) })
	mux.HandleFunc("/api/tasks", func(writer http.ResponseWriter, request *http.Request) {
		database.handleCollection("tasks", writer, request)
	})
	mux.HandleFunc("/api/contacts", func(writer http.ResponseWriter, request *http.Request) {
		database.handleCollection("contacts", writer, request)
	})
	port := os.Getenv("QSL_MAIL_PORT")
	if port == "" {
		port = "38765"
	}
	if err := http.ListenAndServe("127.0.0.1:"+port, mux); err != nil {
		panic(err)
	}
}
