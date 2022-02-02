package challenge

import (
	"math/rand"
)

const (
	Highest    float64 = 999999999 //The highest possible price for an instrument
	TimeLayout         = `'2006-01-02'`

	IDStart    = 1  // Instruments ID range start
	IDEnd      = 50 // Instruments ID range end
	Percentage = 90 // Average percentage of instruments traded per day
	Variation  = 9  // Percentage variation range

)

func (g *Generator) GenerateDays() {}

// table Returns table name to copy in
func (t *Trade) table() string {
	return "trade"
}

// fields Returns fields name of the table
func (t *Trade) fields() []string {
	return []string{"id", "instrumentid", "dateen", "open", "high", "low", "close"}
}

// values Generates randmon values.
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

func (t *Trade) generateDays(count int) {
	instruCount := IDEnd - IDStart                        //count of instrumnets
	min := (instruCount * (Percentage - Variation)) / 100 // Minimum count of instruments traded per day
	max := (instruCount * (Percentage + Variation)) / 100 // Maximum count of instruments traded per day
	difference := max - min
	dayOffset := 0
	counter := 0
	for {
		tradeCount := int(rand.Int63n(int64(difference))) + min
		if counter+tradeCount > count {
			tradeCount = count - counter
		}
		//A piece of pseudo-random permutation of the offsets of instrument IDs from the start of the ID range.
		IDOffsets := rand.Perm(instruCount)

		for j := 0; j < tradeCount; j++ {
			t.days <- TradableDay{
				instrumentID: IDStart + IDOffsets[j],
				deyOffset:    dayOffset,
			}
		}
		if counter += tradeCount; counter == count {
			break

		}
		dayOffset++
	}
	t.generateDone = true

}

func NewTrade(recordCount int) *Trade {
	var t Trade
	go t.generateDays(recordCount)
	t.days = make(chan TradableDay, recordCount)
	return &t
}
