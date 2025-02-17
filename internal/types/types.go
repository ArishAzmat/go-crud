package types

type Todo struct {
	Id          int
	Title       string `validate:"required"`
	Description string `validate:"required"`
}
