package integration_tests

import (
	"darbelis.eu/persedimai/di"
	"darbelis.eu/persedimai/internal/aviation_edge"
	"darbelis.eu/persedimai/internal/dao"
	"darbelis.eu/persedimai/internal/migrations"
	"testing"
)

func TestUpsertFlightSchedules(t *testing.T) {
	// Setup database
	db, err := di.NewDatabase("test")
	if err != nil {
		t.Fatal(err)
	}

	// Create table
	err = migrations.CreateFlightSchedulesTable(db)
	if err != nil {
		t.Fatal(err)
	}

	// Clear table before test
	if !ClearTestDatabase(db, "flight_schedules") {
		t.Fatal("Failed to clear flight_schedules table")
	}

	// Create DAO
	scheduleDao := dao.NewAviationEdgeFlightSchedulesDao(db)

	// Create first test schedule
	schedule1 := &aviation_edge.ScheduleResponse{
		Type:   "departure",
		Status: "scheduled",
		Departure: aviation_edge.Departure{
			IataCode:      "VNO",
			IcaoCode:      "EYVI",
			Terminal:      "1",
			Gate:          "A1",
			ScheduledTime: "2025-12-27T10:00:00.000",
		},
		Arrival: aviation_edge.Arrival{
			IataCode:      "CDG",
			IcaoCode:      "LFPG",
			Terminal:      "2E",
			Gate:          "K45",
			ScheduledTime: "2025-12-27T12:30:00.000",
		},
		Airline: aviation_edge.Airline{
			Name:     "Air France",
			IataCode: "AF",
			IcaoCode: "AFR",
		},
		Flight: aviation_edge.Flight{
			Number:     "1234",
			IataNumber: "AF1234",
			IcaoNumber: "AFR1234",
		},
		Aircraft: aviation_edge.Aircraft{
			RegNumber: "F-HBNK",
			IcaoCode:  "A320",
			ModelText: "Airbus A320",
		},
	}

	// Create second test schedule (different flight)
	schedule2 := &aviation_edge.ScheduleResponse{
		Type:   "departure",
		Status: "scheduled",
		Departure: aviation_edge.Departure{
			IataCode:      "VNO",
			IcaoCode:      "EYVI",
			Terminal:      "1",
			Gate:          "A2",
			ScheduledTime: "2025-12-27T14:00:00.000",
		},
		Arrival: aviation_edge.Arrival{
			IataCode:      "FRA",
			IcaoCode:      "EDDF",
			Terminal:      "1",
			Gate:          "B10",
			ScheduledTime: "2025-12-27T16:30:00.000",
		},
		Airline: aviation_edge.Airline{
			Name:     "Lufthansa",
			IataCode: "LH",
			IcaoCode: "DLH",
		},
		Flight: aviation_edge.Flight{
			Number:     "5678",
			IataNumber: "LH5678",
			IcaoNumber: "DLH5678",
		},
		Aircraft: aviation_edge.Aircraft{
			RegNumber: "D-AIZE",
			IcaoCode:  "A319",
			ModelText: "Airbus A319",
		},
	}

	// First call: insert one schedule
	err = scheduleDao.UpsertFlightSchedules([]*aviation_edge.ScheduleResponse{schedule1})
	if err != nil {
		t.Fatalf("First UpsertFlightSchedules failed: %v", err)
	}

	// Second call: insert two schedules (one duplicate, one new)
	// Update schedule1 with different gate to verify update behavior
	schedule1Updated := &aviation_edge.ScheduleResponse{
		Type:   "departure",
		Status: "scheduled",
		Departure: aviation_edge.Departure{
			IataCode:      "VNO",
			IcaoCode:      "EYVI",
			Terminal:      "1",
			Gate:          "A3", // Changed from A1
			ScheduledTime: "2025-12-27T10:00:00.000",
		},
		Arrival: aviation_edge.Arrival{
			IataCode:      "CDG",
			IcaoCode:      "LFPG",
			Terminal:      "2E",
			Gate:          "K45",
			ScheduledTime: "2025-12-27T12:30:00.000",
		},
		Airline: aviation_edge.Airline{
			Name:     "Air France",
			IataCode: "AF",
			IcaoCode: "AFR",
		},
		Flight: aviation_edge.Flight{
			Number:     "1234",
			IataNumber: "AF1234",
			IcaoNumber: "AFR1234",
		},
		Aircraft: aviation_edge.Aircraft{
			RegNumber: "F-HBNK",
			IcaoCode:  "A320",
			ModelText: "Airbus A320",
		},
	}

	err = scheduleDao.UpsertFlightSchedules([]*aviation_edge.ScheduleResponse{schedule1Updated, schedule2})
	if err != nil {
		t.Fatalf("Second UpsertFlightSchedules failed: %v", err)
	}

	// Verify results
	allSchedules, err := scheduleDao.GetAll()
	if err != nil {
		t.Fatalf("GetAll failed: %v", err)
	}

	// Should have exactly 2 schedules (not 3, because of upsert behavior)
	if len(allSchedules) != 2 {
		t.Errorf("Expected 2 schedules, got %d", len(allSchedules))
	}

	// Verify first schedule was updated (gate should be A3, not A1)
	foundUpdated := false
	for _, schedule := range allSchedules {
		if schedule.Flight.IataNumber == "AF1234" {
			foundUpdated = true
			if schedule.Departure.Gate != "A3" {
				t.Errorf("Expected first schedule to be updated with gate A3, got %s", schedule.Departure.Gate)
			}
			if schedule.Airline.Name != "Air France" {
				t.Errorf("Expected airline name 'Air France', got %s", schedule.Airline.Name)
			}
		}
	}
	if !foundUpdated {
		t.Error("First schedule (AF1234) not found in results")
	}

	// Verify second schedule was inserted
	foundSecond := false
	for _, schedule := range allSchedules {
		if schedule.Flight.IataNumber == "LH5678" {
			foundSecond = true
			if schedule.Airline.Name != "Lufthansa" {
				t.Errorf("Expected airline name 'Lufthansa', got %s", schedule.Airline.Name)
			}
			if schedule.Arrival.IataCode != "FRA" {
				t.Errorf("Expected arrival airport FRA, got %s", schedule.Arrival.IataCode)
			}
		}
	}
	if !foundSecond {
		t.Error("Second schedule (LH5678) not found in results")
	}
}
