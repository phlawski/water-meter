package model

import (
	"time"

	"github.com/phlawski/water-meter/internal/db"
)

// Entry enriches a raw Reading with computed usage and cost.
type Entry struct {
	ID         int64
	ReadAt     time.Time
	ValueM3    float64
	PricePerM3 float64
	Notes      string
	UsageM3    *float64 // nil for the oldest reading (no previous to diff against)
	GrossCost  *float64 // water cost only
	NetCost    *float64 // water cost minus internet contribution
}

// BuildEntries takes readings ordered newest-first (as returned by ListReadings)
// and attaches usage/cost by diffing against the next-older reading.
// internetContribution is a fixed monthly PLN amount subtracted from each period's cost.
func BuildEntries(rows []db.Reading, internetContribution float64) []Entry {
	entries := make([]Entry, len(rows))
	for i, r := range rows {
		e := Entry{
			ID:         r.ID,
			ReadAt:     r.ReadAt,
			ValueM3:    r.ValueM3,
			PricePerM3: r.PricePerM3,
			Notes:      r.Notes,
		}
		// rows[i] is newer, rows[i+1] is older
		if i+1 < len(rows) {
			usage := r.ValueM3 - rows[i+1].ValueM3
			gross := usage * r.PricePerM3
			net := gross - internetContribution
			e.UsageM3 = &usage
			e.GrossCost = &gross
			e.NetCost = &net
		}
		entries[i] = e
	}
	return entries
}
