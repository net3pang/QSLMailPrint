package main

import (
	"encoding/json"
	"net/http"
	"os"

	"qsl-mail/backend/store"
)

func main() {
	database, err := store.Open("")
	if err != nil {
		panic(err)
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/api/health", func(writer http.ResponseWriter, _ *http.Request) { _, _ = writer.Write([]byte(`{"ok":true}`)) })
	for _, name := range []string{"tasks", "contacts"} {
		collectionName := name
		mux.HandleFunc("/api/"+collectionName, func(writer http.ResponseWriter, request *http.Request) {
			if request.Method == http.MethodGet {
				items, err := database.Collection(collectionName)
				if err != nil {
					http.Error(writer, err.Error(), http.StatusNotFound)
					return
				}
				writer.Header().Set("Content-Type", "application/json; charset=utf-8")
				_ = json.NewEncoder(writer).Encode(items)
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
			saved, err := database.SaveRecord(collectionName, item)
			if err != nil {
				http.Error(writer, "保存本地数据失败", http.StatusInternalServerError)
				return
			}
			writer.Header().Set("Content-Type", "application/json; charset=utf-8")
			_ = json.NewEncoder(writer).Encode(saved)
		})
	}
	port := os.Getenv("QSL_MAIL_PORT")
	if port == "" {
		port = "38765"
	}
	if err := http.ListenAndServe("127.0.0.1:"+port, mux); err != nil {
		panic(err)
	}
}
