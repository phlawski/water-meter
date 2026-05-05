package handler

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/phlawski/water-meter/internal/db"
	"github.com/phlawski/water-meter/internal/model"
)

const (
	configKeyPrice               = "price_per_m3"
	configKeyInternetContrib     = "internet_contribution"
)

type indexData struct {
	Entries              []model.Entry
	Today                string
	PricePerM3           float64
	InternetContribution float64
}

func (h *Handler) getConfigFloat(ctx context.Context, key string) (float64, error) {
	val, err := h.q.GetConfigValue(ctx, key)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	return strconv.ParseFloat(val, 64)
}

func (h *Handler) Index(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	rows, err := h.q.ListReadings(ctx)
	if err != nil {
		http.Error(w, "db error", http.StatusInternalServerError)
		log.Printf("ListReadings: %v", err)
		return
	}

	price, err := h.getConfigFloat(ctx, configKeyPrice)
	if err != nil {
		http.Error(w, "db error", http.StatusInternalServerError)
		log.Printf("getConfigFloat price: %v", err)
		return
	}

	internetContrib, err := h.getConfigFloat(ctx, configKeyInternetContrib)
	if err != nil {
		http.Error(w, "db error", http.StatusInternalServerError)
		log.Printf("getConfigFloat internet: %v", err)
		return
	}

	data := indexData{
		Entries:              model.BuildEntries(rows, internetContrib),
		Today:                time.Now().Format("2006-01-02"),
		PricePerM3:           price,
		InternetContribution: internetContrib,
	}

	if err := tmpl.ExecuteTemplate(w, "index.html", data); err != nil {
		log.Printf("template: %v", err)
	}
}

func (h *Handler) CreateReading(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad form", http.StatusBadRequest)
		return
	}

	readAt, err := time.Parse("2006-01-02", r.FormValue("read_at"))
	if err != nil {
		http.Error(w, "invalid date", http.StatusBadRequest)
		return
	}

	valueM3, err := strconv.ParseFloat(r.FormValue("value_m3"), 64)
	if err != nil || valueM3 < 0 {
		http.Error(w, "invalid meter value", http.StatusBadRequest)
		return
	}

	pricePerM3, err := h.getConfigFloat(ctx, configKeyPrice)
	if err != nil {
		http.Error(w, "db error", http.StatusInternalServerError)
		log.Printf("getConfigFloat price: %v", err)
		return
	}

	_, err = h.q.InsertReading(ctx, db.InsertReadingParams{
		ReadAt:     readAt,
		ValueM3:    valueM3,
		PricePerM3: pricePerM3,
		Notes:      r.FormValue("notes"),
	})
	if err != nil {
		http.Error(w, "db error", http.StatusInternalServerError)
		log.Printf("InsertReading: %v", err)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *Handler) DeleteReading(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	if err := h.q.DeleteReading(r.Context(), id); err != nil {
		http.Error(w, "db error", http.StatusInternalServerError)
		log.Printf("DeleteReading: %v", err)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *Handler) SetPrice(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad form", http.StatusBadRequest)
		return
	}

	price, err := strconv.ParseFloat(r.FormValue("price_per_m3"), 64)
	if err != nil || price < 0 {
		http.Error(w, "invalid price", http.StatusBadRequest)
		return
	}

	if err := h.q.SetConfigValue(r.Context(), db.SetConfigValueParams{
		Key:   configKeyPrice,
		Value: strconv.FormatFloat(price, 'f', 4, 64),
	}); err != nil {
		http.Error(w, "db error", http.StatusInternalServerError)
		log.Printf("SetConfigValue price: %v", err)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *Handler) SetInternetContribution(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad form", http.StatusBadRequest)
		return
	}

	contrib, err := strconv.ParseFloat(r.FormValue("internet_contribution"), 64)
	if err != nil || contrib < 0 {
		http.Error(w, "invalid amount", http.StatusBadRequest)
		return
	}

	if err := h.q.SetConfigValue(r.Context(), db.SetConfigValueParams{
		Key:   configKeyInternetContrib,
		Value: strconv.FormatFloat(contrib, 'f', 4, 64),
	}); err != nil {
		http.Error(w, "db error", http.StatusInternalServerError)
		log.Printf("SetConfigValue internet: %v", err)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
