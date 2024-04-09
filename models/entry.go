package models

type Entry struct {
	Key      string `json:"key,omitempty" validate:"required"`
	Url      string `json:"url,omitempty" validate:"required"`
	Password string `json:"password,omitempty" validate:"omitempty"`
	Visits   int    `json:"visits,omitempty" validate:"omitempty"`
}
