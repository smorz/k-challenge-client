package challenge

import (
	"errors"
	"math/rand"
	"time"
)

const (
	Highest    float64 = 999999999 //The highest possible price for an instrument
	TimeLayout         = `'2006-01-02'`

	// Determine the ranges of random numbers.
	// In the future, if necessary, they can be received from the input.
	IDStart    = 1  // Instruments ID range start
	IDEnd      = 4  // Instruments ID range end
	Percentage = 90 // Average percentage of instruments traded per day
	Variation  = 9  // Percentage variation range

)

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

// generateDays creates a sufficient amount of days and determines which instrument was traded on each day.
// Is execute at TradeGenerator creating
func (t *TradeGenerator) generateDays() {
	//count of instrumnets
	instruCount := IDEnd - IDStart
	min := (instruCount * (Percentage - Variation)) / 100 // Minimum count of instruments traded per day
	max := (instruCount * (Percentage + Variation)) / 100 // Maximum count of instruments traded per day
	if max == 0 {
		max = 1
	}
	difference := max - min

	dayOffset := 0
	counter := 0
	dayTradeCount := min
	for {
		if difference != 0 {
			dayTradeCount = int(rand.Int63n(int64(difference))) + min
		}
		// If today's trade are  more than the required (today is last day)
		if counter+dayTradeCount > t.recordsCount {
			dayTradeCount = t.recordsCount - counter
		}
		//A slice of pseudo-random permutation of the offsets of instrument IDs from the start of the ID range.
		IDOffsets := rand.Perm(instruCount)

		for j := 0; j < dayTradeCount; j++ {

			t.days <- instrumentDay{
				instrumentID: IDStart + IDOffsets[j],
				deyOffset:    dayOffset,
			}
		}

		if counter += dayTradeCount; counter == t.recordsCount {
			break

		}
		dayOffset++
	}
}

// NewTradeGenerator Create a new TradeGenerator
func NewTradeGenerator(firstDay time.Time, recordsCount int) (*TradeGenerator, error) {
	if recordsCount <= 0 {
		err := errors.New("recordsCount is not posetive")
		return nil, err
	}

	t := TradeGenerator{
		recordsCount: recordsCount,
		generated:    0,
		firstDay:     firstDay,
		days:         make(chan instrumentDay, recordsCount),
	}
	go t.generateDays()
	return &t, nil
}
