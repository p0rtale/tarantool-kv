package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	jsoniter "github.com/json-iterator/go"
	"github.com/tarantool/go-tarantool/v2"
)

type PostRequest struct {
	Key   string `json:"key"`
	Value any    `json:"value"`
}

type PutRequest struct {
	Value any `json:"value"`
}

func PostKV(conn *tarantool.Connection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req PostRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		resp, err := conn.Do(tarantool.NewSelectRequest("kv").Key([]any{req.Key})).Get()
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}
		if len(resp) > 0 {
			http.Error(w, "Key already exists", http.StatusConflict)
			return
		}

		_, err = conn.Do(tarantool.NewInsertRequest("kv").Tuple([]any{req.Key, req.Value})).Get()
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}

func PutKV(conn *tarantool.Connection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		key := vars["id"]

		var req PutRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		resp, err := conn.Do(tarantool.NewSelectRequest("kv").Key([]any{key})).Get()
		if err != nil || len(resp) == 0 {
			http.Error(w, "Key not found", http.StatusNotFound)
			return
		}

		_, err = conn.Do(
			tarantool.NewUpdateRequest("kv").
				Key([]any{key}).
				Operations(tarantool.NewOperations().Assign(1, req.Value)),
		).Get()
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func GetKV(conn *tarantool.Connection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		key := vars["id"]

		resp, err := conn.Do(tarantool.NewSelectRequest("kv").Key([]any{key})).Get()
		if err != nil || len(resp) == 0 {
			http.Error(w, "Key not found", http.StatusNotFound)
			return
		}

		value := resp[0].([]any)[1]
		w.Header().Set("Content-Type", "application/json")
		jsoniter.NewEncoder(w).Encode(map[string]any{"key": key, "value": value})
	}
}

func DeleteKV(conn *tarantool.Connection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		key := vars["id"]

		resp, err := conn.Do(tarantool.NewSelectRequest("kv").Key([]any{key})).Get()
		if err != nil || len(resp) == 0 {
			http.Error(w, "Key not found", http.StatusNotFound)
			return
		}

		_, err = conn.Do(tarantool.NewDeleteRequest("kv").Key([]any{key})).Get()
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
