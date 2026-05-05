package handler

import (
	"html/template"

	"github.com/phlawski/water-meter/internal/db"
)

var funcMap = template.FuncMap{
	"deref": func(f *float64) float64 {
		if f == nil {
			return 0
		}
		return *f
	},
}

var tmpl = template.Must(template.New("").Funcs(funcMap).ParseGlob("templates/*.html"))

type Handler struct {
	q *db.Queries
}

func New(q *db.Queries) *Handler {
	return &Handler{q: q}
}
