package models

type Response struct {
	Status bool   `json:"status"`
	Error  string `json:"error,omitempty"`
}

type CreateCollectionRequest struct {
	CollectionName string `json:"collection_name"`
	PrimaryKey     string `json:"primary_key"`
}

type DeleteCollectionRequest struct {
	CollectionName string `json:"collection_name"`
}

type ListCollectionsResponse struct {
	Status      bool     `json:"status"`
	Collections []string `json:"collections,omitempty"`
	Error       string   `json:"error,omitempty"`
}

type PutDocumentRequest struct {
	CollectionName string                 `json:"collection_name"`
	Document       map[string]interface{} `json:"document"`
}

type GetDocumentRequest struct {
	CollectionName string `json:"collection_name"`
	Key            string `json:"key"`
	Value          string `json:"value"`
}

type GetDocumentResponse struct {
	Status   bool                   `json:"status"`
	Document map[string]interface{} `json:"document,omitempty"`
	Error    string                 `json:"error,omitempty"`
}

type DeleteDocumentRequest struct {
	CollectionName string `json:"collection_name"`
	Key            string `json:"key"`
	Value          string `json:"value"`
}

type ListDocumentsRequest struct {
	CollectionName string `json:"collection_name"`
}

type ListDocumentsResponse struct {
	Status    bool                     `json:"status"`
	Documents []map[string]interface{} `json:"documents,omitempty"`
	Error     string                   `json:"error,omitempty"`
}

type CreateIndexRequest struct {
	CollectionName string `json:"collection_name"`
	Field          string `json:"field"`
	Unique         bool   `json:"unique"`
}

type DeleteIndexRequest struct {
	CollectionName string `json:"collection_name"`
	Field          string `json:"field"`
}
