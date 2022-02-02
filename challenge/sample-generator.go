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
func (t *TradeGenerator) table() string {
	return "trade"
}

// fields Returns fields name of the table
func (t *TradeGenerator) fields() []string {
	return []string{"id", "instrumentid", "dateen", "open", "high", "low", "close"}
}

// values Generates randmon values.
func (t *TradeGenerator) values() []interface{} {
	for {
		select {
		case day := <-t.days:
			open, high, low, close := t.generateOHLC()
			t.mu.Lock()
			defer t.mu.Unlock()
			t.generated++
			return []interface{}{1, day.instrumentID,
				t.firstDay.AddDate(0, 0, day.deyOffset).Format(TimeLayout),
				open, high, low, close,
			}
		default:
			if t.finished() {
				return nil
			}
			//fmt.Print(t.generated, t.recordsCount, len(t.days), " | ")
			continue
		}
	}
}

func (t *TradeGenerator) finished() bool {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.generated == t.recordsCount
}

// generateLowHigh Generates random low and high
func (t *TradeGenerator) generateLowHigh() (low, high float64) {
	a := rand.Float64() * Highest
	b := rand.Float64() * Highest
	if b < a {
		return b, a
	}
	return a, b
}

// generateOHLC Generates random open, high, low, and close
func (t *TradeGenerator) generateOHLC() (open, high, low, close float64) {
	low, high = t.generateLowHigh()
	difference := high - low
	open = low + rand.Float64()*difference
	close = low + rand.Float64()*difference
	return
}

func (t *TradeGenerator) generateDays() {

	instruCount := IDEnd - IDStart                        //count of instrumnets
	min := (instruCount * (Percentage - Variation)) / 100 // Minimum count of instruments traded per day
	max := (instruCount * (Percentage + Variation)) / 100 // Maximum count of instruments traded per day
	difference := max - min
	dayOffset := 0
	counter := 0
	for {
		tradeCount := int(rand.Int63n(int64(difference))) + min
		if counter+tradeCount > t.recordsCount {
			tradeCount = t.recordsCount - counter
		}
		//A piece of pseudo-random permutation of the offsets of instrument IDs from the start of the ID range.
		IDOffsets := rand.Perm(instruCount)

		for j := 0; j < tradeCount; j++ {

			t.days <- TradableDay{
				instrumentID: IDStart + IDOffsets[j],
				deyOffset:    dayOffset,
			}
		}

		if counter += tradeCount; counter == t.recordsCount {
			break

		}
		dayOffset++
	}
}

func NewTradeGenerator(recordsCount int) *TradeGenerator {
	var t TradeGenerator
	t.recordsCount = recordsCount
	t.days = make(chan TradableDay, t.recordsCount)
	go t.generateDays()
	return &t
}
