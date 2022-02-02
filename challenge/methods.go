package challenge

import (
	"math/rand"
)

const (
	Highest    float64 = 999999999 //The highest possible price for an instrument
	TimeLayout         = `'2006-01-02'`
)

func (g *Generator) GenerateDays() {}

func (t *Trade) table() string {
	return "trade"
}

func (t *Trade) fields() []string {
	return []string{"id", "instrumentid", "dateen", "open", "high", "low", "close"}
}

func (t *Trade) values() ([]interface{}, bool) {
	select {
	case day := <-t.days:
		open, high, low, close := t.generateOHLC()
		return []interface{}{1, day.instrumentID,
			t.firstDay.AddDate(0, 0, day.deyOffset).Format(TimeLayout),
			open, high, low, close,
		}, false
	default:
		return nil, t.generateDone
	}

}

// generateLowHigh Generates random low and high
func (t *Trade) generateLowHigh() (low, high float64) {
	a := rand.Float64() * Highest
	b := rand.Float64() * Highest
	if b < a {
		return b, a
	}
	return a, b
}

// generateOHLC Generates random open, high, low, and close
func (t *Trade) generateOHLC() (open, high, low, close float64) {
	low, high = t.generateLowHigh()
	difference := high - low
	open = low + rand.Float64()*difference
	close = low + rand.Float64()*difference
	return
}
