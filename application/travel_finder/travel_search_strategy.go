package travel_finder

import "darbelis.eu/persedimai/data"

// TravelSearchStrategy defines the interface for different search implementations
type TravelSearchStrategy interface {
	// FindPath finds a sequence of travels from source to destination based on the filter criteria
	// Returns the travel path, or error if no path exists or search fails
	FindPath(filter *data.TravelFilter) ([]*TravelPath, error)

	// GetName returns the strategy name (for logging/debugging)
	GetName() string
}
