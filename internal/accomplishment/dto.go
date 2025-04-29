package accomplishment

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/project-ippl-dev/tanding-api/internal/db"
)

type request struct {
	Title       string                 `json:"title"`
	Level       db.AccomplishmentLevel `json:"level"`
	Ranking     string                 `json:"ranking"`
	Category    string                 `json:"category"`
	Sport       string                 `json:"sport"`
	Description string                 `json:"description"`
	FileURL     string                 `json:"file_url"`
	Month       int16                  `json:"month`
	Year        int16                  `json:"year"`
}

func (r request) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Title, validation.Required),
		validation.Field(&r.Level, validation.Required),
		validation.Field(&r.Ranking, validation.Required),
		validation.Field(&r.Sport, validation.Required),
		validation.Field(&r.FileURL, validation.Required, is.URL),
		validation.Field(&r.Year, validation.Required, validation.Min(1900), validation.Max(3000)),
	)

}

type response struct {
	Title       string                 `json:"title"`
	Level       db.AccomplishmentLevel `json:"level"`
	Ranking     string                 `json:"ranking"`
	Category    string                 `json:"category"`
	Sport       string                 `json:"sport"`
	Description string                 `json:"description"`
	FileURL     string                 `json:"file_url"`
	Date        string                 `json:"date`
}
