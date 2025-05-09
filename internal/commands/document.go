package commands

type PutDocumentRequest struct {
	Collection string                 `json:"collection"`
	Document   map[string]interface{} `json:"document"`
}
type PutDocumentResponse struct{}

type GetDocumentRequest struct {
	Collection string `json:"collection"`
	Key        string `json:"key"`
}
type GetDocumentResponse struct {
	Document map[string]interface{} `json:"document"`
	Found    bool                   `json:"found"`
}

type DeleteDocumentRequest struct {
	Collection string `json:"collection"`
	Key        string `json:"key"`
}
type DeleteDocumentResponse struct {
	Success bool
}

type ListDocumentsRequest struct {
	Collection string `json:"collection"`
}
type ListDocumentsResponse struct {
	Documents []map[string]interface{} `json:"documents"`
}
