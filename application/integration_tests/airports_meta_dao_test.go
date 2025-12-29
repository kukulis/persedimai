package integration_tests

import (
	"darbelis.eu/persedimai/di"
	"darbelis.eu/persedimai/internal/dao"
	"darbelis.eu/persedimai/internal/tables"
	"testing"
	"time"
)

func TestAirportsMetaDaoUpsert(t *testing.T) {
	// Setup database
	db, err := di.NewDatabase("test")
	if err != nil {
		t.Fatal(err)
	}

	// Create DAO
	airportsMetaDao := dao.NewAirportsMetaDao(db)

	// Create table
	err = airportsMetaDao.CreateTable()
	if err != nil {
		t.Fatal(err)
	}

	// Clear table before test
	if !ClearTestDatabase(db, "airports_meta") {
		t.Fatal("Failed to clear airports_meta table")
	}

	// Test 1: Insert new metadata record
	date1From := time.Date(2025, 12, 27, 0, 0, 0, 0, time.UTC)
	date1To := time.Date(2025, 12, 31, 0, 0, 0, 0, time.UTC)
	meta1 := &tables.AirportMeta{
		AirportCode:  "VNO",
		ImportedFrom: &date1From,
		ImportedTo:   &date1To,
	}

	err = airportsMetaDao.Upsert(meta1, true)
	if err != nil {
		t.Fatalf("First Upsert failed: %v", err)
	}

	// Verify the record was inserted
	retrieved1, err := airportsMetaDao.Get("VNO")
	if err != nil {
		t.Fatalf("Get failed after insert: %v", err)
	}
	if retrieved1 == nil {
		t.Fatal("Expected to retrieve inserted record, got nil")
	}
	if retrieved1.AirportCode != "VNO" {
		t.Errorf("Expected airport code 'VNO', got '%s'", retrieved1.AirportCode)
	}
	if retrieved1.ImportedFrom == nil || !retrieved1.ImportedFrom.Equal(date1From) {
		t.Errorf("Expected ImportedFrom to be %s, got %v", date1From, retrieved1.ImportedFrom)
	}
	if retrieved1.ImportedTo == nil || !retrieved1.ImportedTo.Equal(date1To) {
		t.Errorf("Expected ImportedTo to be %s, got %v", date1To, retrieved1.ImportedTo)
	}

	// Test 2: Update existing metadata record
	date2From := time.Date(2025, 12, 20, 0, 0, 0, 0, time.UTC)
	date2To := time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC)
	meta2 := &tables.AirportMeta{
		AirportCode:  "VNO",
		ImportedFrom: &date2From,
		ImportedTo:   &date2To,
	}

	err = airportsMetaDao.Upsert(meta2, true)
	if err != nil {
		t.Fatalf("Second Upsert (update) failed: %v", err)
	}

	// Verify the record was updated
	retrieved2, err := airportsMetaDao.Get("VNO")
	if err != nil {
		t.Fatalf("Get failed after update: %v", err)
	}
	if retrieved2 == nil {
		t.Fatal("Expected to retrieve updated record, got nil")
	}
	if retrieved2.ImportedFrom == nil || !retrieved2.ImportedFrom.Equal(date2From) {
		t.Errorf("Expected updated ImportedFrom to be %s, got %v", date2From, retrieved2.ImportedFrom)
	}
	if retrieved2.ImportedTo == nil || !retrieved2.ImportedTo.Equal(date2To) {
		t.Errorf("Expected updated ImportedTo to be %s, got %v", date2To, retrieved2.ImportedTo)
	}

	// Test 3: Insert metadata with null dates
	meta3 := &tables.AirportMeta{
		AirportCode:  "CDG",
		ImportedFrom: nil,
		ImportedTo:   nil,
	}

	err = airportsMetaDao.Upsert(meta3, true)
	if err != nil {
		t.Fatalf("Upsert with nil dates failed: %v", err)
	}

	// Verify the record with null dates
	retrieved3, err := airportsMetaDao.Get("CDG")
	if err != nil {
		t.Fatalf("Get failed for CDG: %v", err)
	}
	if retrieved3 == nil {
		t.Fatal("Expected to retrieve CDG record, got nil")
	}
	if retrieved3.AirportCode != "CDG" {
		t.Errorf("Expected airport code 'CDG', got '%s'", retrieved3.AirportCode)
	}
	if retrieved3.ImportedFrom != nil {
		t.Errorf("Expected ImportedFrom to be nil, got %v", retrieved3.ImportedFrom)
	}
	if retrieved3.ImportedTo != nil {
		t.Errorf("Expected ImportedTo to be nil, got %v", retrieved3.ImportedTo)
	}

	// Test 4: Insert another airport to verify multiple records
	date4From := time.Date(2025, 11, 1, 0, 0, 0, 0, time.UTC)
	date4To := time.Date(2025, 11, 30, 0, 0, 0, 0, time.UTC)
	meta4 := &tables.AirportMeta{
		AirportCode:  "JFK",
		ImportedFrom: &date4From,
		ImportedTo:   &date4To,
	}

	err = airportsMetaDao.Upsert(meta4, true)
	if err != nil {
		t.Fatalf("Upsert for JFK failed: %v", err)
	}

	// Verify JFK was inserted
	retrieved4, err := airportsMetaDao.Get("JFK")
	if err != nil {
		t.Fatalf("Get failed for JFK: %v", err)
	}
	if retrieved4 == nil {
		t.Fatal("Expected to retrieve JFK record, got nil")
	}
	if retrieved4.AirportCode != "JFK" {
		t.Errorf("Expected airport code 'JFK', got '%s'", retrieved4.AirportCode)
	}

	// Verify VNO still exists (to ensure we didn't overwrite it)
	retrievedVNO, err := airportsMetaDao.Get("VNO")
	if err != nil {
		t.Fatalf("Get failed for VNO after JFK insert: %v", err)
	}
	if retrievedVNO == nil {
		t.Fatal("Expected VNO record to still exist")
	}

	// Test 5: Upsert with updateDates=false should not update dates on duplicate
	date5From := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	date5To := time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC)
	meta5 := &tables.AirportMeta{
		AirportCode:  "JFK",
		ImportedFrom: &date5From,
		ImportedTo:   &date5To,
	}

	// This should NOT update the dates because updateDates=false
	err = airportsMetaDao.Upsert(meta5, false)
	if err != nil {
		t.Fatalf("Upsert with updateDates=false failed: %v", err)
	}

	// Verify dates were NOT updated
	retrieved5, err := airportsMetaDao.Get("JFK")
	if err != nil {
		t.Fatalf("Get failed for JFK after updateDates=false: %v", err)
	}
	if retrieved5 == nil {
		t.Fatal("Expected JFK record to exist")
	}
	// Dates should still be the original ones (date4From/date4To), not the new ones (date5From/date5To)
	if retrieved5.ImportedFrom == nil || !retrieved5.ImportedFrom.Equal(date4From) {
		t.Errorf("Expected ImportedFrom to remain %s (not updated), got %v", date4From, retrieved5.ImportedFrom)
	}
	if retrieved5.ImportedTo == nil || !retrieved5.ImportedTo.Equal(date4To) {
		t.Errorf("Expected ImportedTo to remain %s (not updated), got %v", date4To, retrieved5.ImportedTo)
	}
}

func TestAirportsMetaDaoGetFirstWithNullDates(t *testing.T) {
	// Setup database
	db, err := di.NewDatabase("test")
	if err != nil {
		t.Fatal(err)
	}

	// Create DAO
	airportsMetaDao := dao.NewAirportsMetaDao(db)

	// Create table
	err = airportsMetaDao.CreateTable()
	if err != nil {
		t.Fatal(err)
	}

	// Clear table before test
	if !ClearTestDatabase(db, "airports_meta") {
		t.Fatal("Failed to clear airports_meta table")
	}

	// Test 1: No records - should return nil
	result1, err := airportsMetaDao.GetFirstWithNullDates()
	if err != nil {
		t.Fatalf("GetFirstWithNullDates failed on empty table: %v", err)
	}
	if result1 != nil {
		t.Error("Expected nil for empty table, got a record")
	}

	// Test 2: Insert record with dates - should not be returned
	dateFrom := time.Date(2025, 12, 27, 0, 0, 0, 0, time.UTC)
	dateTo := time.Date(2025, 12, 31, 0, 0, 0, 0, time.UTC)
	meta1 := &tables.AirportMeta{
		AirportCode:  "VNO",
		ImportedFrom: &dateFrom,
		ImportedTo:   &dateTo,
	}
	err = airportsMetaDao.Upsert(meta1, true)
	if err != nil {
		t.Fatalf("Upsert failed: %v", err)
	}

	result2, err := airportsMetaDao.GetFirstWithNullDates()
	if err != nil {
		t.Fatalf("GetFirstWithNullDates failed: %v", err)
	}
	if result2 != nil {
		t.Error("Expected nil when no records have null dates, got a record")
	}

	// Test 3: Insert record with null dates - should be returned
	meta2 := &tables.AirportMeta{
		AirportCode:  "CDG",
		ImportedFrom: nil,
		ImportedTo:   nil,
	}
	err = airportsMetaDao.Upsert(meta2, true)
	if err != nil {
		t.Fatalf("Upsert with null dates failed: %v", err)
	}

	result3, err := airportsMetaDao.GetFirstWithNullDates()
	if err != nil {
		t.Fatalf("GetFirstWithNullDates failed: %v", err)
	}
	if result3 == nil {
		t.Fatal("Expected to get a record with null dates, got nil")
	}
	if result3.AirportCode != "CDG" {
		t.Errorf("Expected airport code 'CDG', got '%s'", result3.AirportCode)
	}
	if result3.ImportedFrom != nil {
		t.Errorf("Expected ImportedFrom to be nil, got %v", result3.ImportedFrom)
	}
	if result3.ImportedTo != nil {
		t.Errorf("Expected ImportedTo to be nil, got %v", result3.ImportedTo)
	}

	// Test 4: Insert another record with null dates - should return first one
	meta3 := &tables.AirportMeta{
		AirportCode:  "JFK",
		ImportedFrom: nil,
		ImportedTo:   nil,
	}
	err = airportsMetaDao.Upsert(meta3, true)
	if err != nil {
		t.Fatalf("Upsert second null record failed: %v", err)
	}

	result4, err := airportsMetaDao.GetFirstWithNullDates()
	if err != nil {
		t.Fatalf("GetFirstWithNullDates failed: %v", err)
	}
	if result4 == nil {
		t.Fatal("Expected to get a record with null dates, got nil")
	}
	// Should return one of them (CDG or JFK), we don't care which specific one
	if result4.AirportCode != "CDG" && result4.AirportCode != "JFK" {
		t.Errorf("Expected airport code 'CDG' or 'JFK', got '%s'", result4.AirportCode)
	}
}
