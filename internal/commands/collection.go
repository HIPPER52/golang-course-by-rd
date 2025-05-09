package commands

type CreateCollectionRequest struct {
	Name string `json:"name"`
}
type CreateCollectionResponse struct{}

type DeleteCollectionRequest struct {
	Name string `json:"name"`
}
type DeleteCollectionResponse struct{}

type ListCollectionsRequest struct{}
type ListCollectionsResponse struct {
	Collections []string `json:"collections"`
}
