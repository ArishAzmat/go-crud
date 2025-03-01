package types

type Todo struct {
	Id          int64  `json:"id"`
	Title       string `json:"title" validate:"required"`
	Description string `json:"description" validate:"required"`
}
