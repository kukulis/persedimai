package tables

import "time"

// TransferSequence represents an ordered sequence of transfers forming a single path
type TransferSequence struct {
	Transfers []*Transfer
}

// NewTransferSequence creates a new TransferSequence from a slice of transfers
func NewTransferSequence(transfers []*Transfer) *TransferSequence {
	return &TransferSequence{Transfers: transfers}
}

// First returns the first transfer in the sequence, or nil if empty
func (ts *TransferSequence) First() *Transfer {
	if len(ts.Transfers) == 0 {
		return nil
	}
	return ts.Transfers[0]
}

// Last returns the last transfer in the sequence, or nil if empty
func (ts *TransferSequence) Last() *Transfer {
	if len(ts.Transfers) == 0 {
		return nil
	}
	return ts.Transfers[len(ts.Transfers)-1]
}

// Count returns the number of transfers in the sequence
func (ts *TransferSequence) Count() int {
	return len(ts.Transfers)
}

// TransferCount returns the number of transfers (same as Count for clarity)
func (ts *TransferSequence) TransferCount() int {
	return len(ts.Transfers)
}

// TotalDuration calculates the total journey time from first departure to last arrival
func (ts *TransferSequence) TotalDuration() time.Duration {
	if len(ts.Transfers) == 0 {
		return 0
	}
	return ts.Last().Arrival.Sub(ts.First().Departure)
}

// IsValid verifies that transfers are properly connected in sequence
// Returns true if:
// - Each transfer's destination matches the next transfer's origin
// - Each transfer departs after the previous one arrives
func (ts *TransferSequence) IsValid() bool {
	if len(ts.Transfers) == 0 {
		return true
	}

	for i := 0; i < len(ts.Transfers)-1; i++ {
		current := ts.Transfers[i]
		next := ts.Transfers[i+1]

		// Check if destination of current matches origin of next
		if current.To != next.From {
			return false
		}

		// Check if next transfer departs after current arrives
		if next.Departure.Before(current.Arrival) {
			return false
		}
	}

	return true
}

// ConnectionTime calculates waiting time at the given transfer index
// Returns 0 if index is out of bounds or if it's the last transfer
func (ts *TransferSequence) ConnectionTime(index int) time.Duration {
	if index < 0 || index >= len(ts.Transfers)-1 {
		return 0
	}

	return ts.Transfers[index+1].Departure.Sub(ts.Transfers[index].Arrival)
}

// TotalConnectionTime calculates the sum of all connection/waiting times
func (ts *TransferSequence) TotalConnectionTime() time.Duration {
	var total time.Duration
	for i := 0; i < len(ts.Transfers)-1; i++ {
		total += ts.ConnectionTime(i)
	}
	return total
}
