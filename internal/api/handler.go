package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"lesson_14/internal/models"
	"lesson_14/internal/mongodb"
)

type Handler struct {
	Store *mongodb.Store
}

func NewHandler(store *mongodb.Store) *Handler {
	return &Handler{Store: store}
}

func (h *Handler) PutDocument(w http.ResponseWriter, r *http.Request) {
	var req models.PutDocumentRequest
	if !decodeJSONBody(w, r, &req) {
		return
	}
	err := h.Store.PutDocument(r.Context(), req.CollectionName, req.Document)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to put document: %v", err), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) GetDocument(w http.ResponseWriter, r *http.Request) {
	var req models.GetDocumentRequest
	if !decodeJSONBody(w, r, &req) {
		return
	}
	doc, err := h.Store.GetDocument(r.Context(), req.CollectionName, req.Key, req.Value)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to get document: %v", err), http.StatusInternalServerError)
		return
	}
	if doc == nil {
		http.Error(w, "document not found", http.StatusNotFound)
		return
	}
	respondJSON(w, doc)
}

func (h *Handler) ListDocuments(w http.ResponseWriter, r *http.Request) {
	var req models.ListDocumentsRequest
	if !decodeJSONBody(w, r, &req) {
		return
	}
	docs, err := h.Store.ListDocuments(r.Context(), req.CollectionName)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to list documents: %v", err), http.StatusInternalServerError)
		return
	}
	respondJSON(w, docs)
}

func (h *Handler) DeleteDocument(w http.ResponseWriter, r *http.Request) {
	var req models.DeleteDocumentRequest
	if !decodeJSONBody(w, r, &req) {
		return
	}
	deleted, err := h.Store.DeleteDocument(r.Context(), req.CollectionName, req.Key, req.Value)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to delete document: %v", err), http.StatusInternalServerError)
		return
	}
	if !deleted {
		http.Error(w, "document not found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) CreateCollection(w http.ResponseWriter, r *http.Request) {
	var req models.CreateCollectionRequest
	if !decodeJSONBody(w, r, &req) {
		return
	}
	err := h.Store.CreateCollection(r.Context(), req.CollectionName)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to create collection: %v", err), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) DeleteCollection(w http.ResponseWriter, r *http.Request) {
	var req models.DeleteCollectionRequest
	if !decodeJSONBody(w, r, &req) {
		return
	}
	err := h.Store.DeleteCollection(r.Context(), req.CollectionName)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to delete collection: %v", err), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) ListCollections(w http.ResponseWriter, r *http.Request) {
	names, err := h.Store.ListCollections(r.Context())
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to list collections: %v", err), http.StatusInternalServerError)
		return
	}
	respondJSON(w, names)
}

func (h *Handler) CreateIndex(w http.ResponseWriter, r *http.Request) {
	var req models.CreateIndexRequest
	if !decodeJSONBody(w, r, &req) {
		return
	}
	err := h.Store.CreateIndex(r.Context(), req.CollectionName, req.Field, req.Unique)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to create index: %v", err), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) DeleteIndex(w http.ResponseWriter, r *http.Request) {
	var req models.DeleteIndexRequest
	if !decodeJSONBody(w, r, &req) {
		return
	}
	err := h.Store.DeleteIndex(r.Context(), req.CollectionName, req.Field)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to delete index: %v", err), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func respondJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}

func decodeJSONBody[T any](w http.ResponseWriter, r *http.Request, dst *T) bool {
	if err := json.NewDecoder(r.Body).Decode(dst); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return false
	}
	return true
}
