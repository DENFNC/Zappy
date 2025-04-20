package dto

type ListResult[T any] struct {
	Items         []*T
	NextPageToken string
}
