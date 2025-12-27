package aviation_edge

import (
	"bytes"
	"io"
	"net/http"
	"strings"
	"testing"
)

// mockTransport is a mock implementation of http.RoundTripper
// It allows us to mock HTTP responses without making actual network calls
type mockTransport struct {
	response *http.Response
	err      error
}

func (m *mockTransport) RoundTrip(*http.Request) (*http.Response, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.response, nil
}

// newMockClient creates a new AviationEdgeApiClient with a mocked HTTP client
func newMockClient(statusCode int, body string) *AviationEdgeApiClient {
	return &AviationEdgeApiClient{
		APIKey:  "test-api-key",
		BaseURL: BaseURL,
		HTTPClient: &http.Client{
			Transport: &mockTransport{
				response: &http.Response{
					StatusCode: statusCode,
					Body:       io.NopCloser(bytes.NewBufferString(body)),
					Header:     make(http.Header),
				},
			},
		},
	}
}

// TestGetAirports_Success tests a successful API response
func TestGetAirports_Success(t *testing.T) {
	responseBody := `[
		{
			"airportId": 1,
			"nameAirport": "Los Angeles International Airport",
			"codeIataAirport": "LAX",
			"codeIcaoAirport": "KLAX",
			"latitudeAirport": 33.9425,
			"longitudeAirport": -118.408,
			"geonameId": "5368361",
			"timezone": "America/Los_Angeles",
			"GMT": "-8",
			"phone": "+1 855-463-5252",
			"nameCountry": "United States",
			"codeIso2Country": "US",
			"codeIataCity": "LAX"
		},
		{
			"airportId": 2,
			"nameAirport": "John F Kennedy International Airport",
			"codeIataAirport": "JFK",
			"codeIcaoAirport": "KJFK",
			"latitudeAirport": 40.6398,
			"longitudeAirport": -73.7789,
			"geonameId": "5125738",
			"timezone": "America/New_York",
			"GMT": "-5",
			"phone": "+1 718-244-4444",
			"nameCountry": "United States",
			"codeIso2Country": "US",
			"codeIataCity": "NYC"
		}
	]`

	client := newMockClient(200, responseBody)

	airports, err := client.GetAirports(AirportsParams{
		CodeIso2Country: "US",
	})

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if len(airports) != 2 {
		t.Fatalf("Expected 2 airports, got %d", len(airports))
	}

	// Verify first airport data
	if airports[0].CodeIataAirport != "LAX" {
		t.Errorf("Expected airport code 'LAX', got '%s'", airports[0].CodeIataAirport)
	}
	if airports[0].NameAirport != "Los Angeles International Airport" {
		t.Errorf("Expected airport name 'Los Angeles International Airport', got '%s'", airports[0].NameAirport)
	}
	if airports[0].CodeIso2Country != "US" {
		t.Errorf("Expected country code 'US', got '%s'", airports[0].CodeIso2Country)
	}

	// Verify second airport data
	if airports[1].CodeIataAirport != "JFK" {
		t.Errorf("Expected airport code 'JFK', got '%s'", airports[1].CodeIataAirport)
	}
}

// TestGetAirports_ErrorJSON_Status200 tests error JSON response with 200 status code
// The Aviation Edge API sometimes returns error messages with a 200 status code
func TestGetAirports_ErrorJSON_Status200(t *testing.T) {
	errorBody := `{
		"success": false,
		"error": "Invalid API key",
		"message": "The API key provided is invalid or expired"
	}`

	client := newMockClient(200, errorBody)

	airports, err := client.GetAirports(AirportsParams{
		CodeIso2Country: "US",
	})

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if airports != nil {
		t.Errorf("Expected nil airports, got %v", airports)
	}

	expectedErrorMsg := "API error: The API key provided is invalid or expired"
	if !strings.Contains(err.Error(), "API error:") {
		t.Errorf("Expected error message to contain 'API error:', got: %v", err)
	}
	if !strings.Contains(err.Error(), "The API key provided is invalid or expired") {
		t.Errorf("Expected error message '%s', got: %v", expectedErrorMsg, err)
	}
}

// TestGetAirports_ErrorJSON_Status400 tests error JSON response with 400 status code
func TestGetAirports_ErrorJSON_Status400(t *testing.T) {
	errorBody := `{
		"error": "BadRequest",
		"message": "Invalid country code provided"
	}`

	client := newMockClient(400, errorBody)

	airports, err := client.GetAirports(AirportsParams{
		CodeIso2Country: "INVALID",
	})

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if airports != nil {
		t.Errorf("Expected nil airports, got %v", airports)
	}

	// Should contain the status code and error message
	if !strings.Contains(err.Error(), "400") {
		t.Errorf("Expected error message to contain status code '400', got: %v", err)
	}
	if !strings.Contains(err.Error(), "Invalid country code provided") {
		t.Errorf("Expected error message to contain 'Invalid country code provided', got: %v", err)
	}
}

// TestGetAirports_ErrorJSON_Status500 tests error JSON response with 500 status code
func TestGetAirports_ErrorJSON_Status500(t *testing.T) {
	errorBody := `{
		"error": "InternalServerError",
		"message": "An internal server error occurred"
	}`

	client := newMockClient(500, errorBody)

	airports, err := client.GetAirports(AirportsParams{
		CodeIataAirport: "LAX",
	})

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if airports != nil {
		t.Errorf("Expected nil airports, got %v", airports)
	}

	// Should contain the status code and error message
	if !strings.Contains(err.Error(), "500") {
		t.Errorf("Expected error message to contain status code '500', got: %v", err)
	}
	if !strings.Contains(err.Error(), "An internal server error occurred") {
		t.Errorf("Expected error message to contain 'An internal server error occurred', got: %v", err)
	}
}

// TestGetAirports_NonJSONResponse tests a non-JSON response from the API
// This could happen if the API returns HTML (e.g., maintenance page, 404 page)
func TestGetAirports_NonJSONResponse(t *testing.T) {
	htmlResponse := `<!DOCTYPE html>
<html>
<head><title>503 Service Unavailable</title></head>
<body>
<h1>Service Unavailable</h1>
<p>The service is temporarily unavailable. Please try again later.</p>
</body>
</html>`

	client := newMockClient(503, htmlResponse)

	airports, err := client.GetAirports(AirportsParams{
		CodeIso2Country: "US",
	})

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if airports != nil {
		t.Errorf("Expected nil airports, got %v", airports)
	}

	// Error should mention non-JSON response and the status code
	if !strings.Contains(err.Error(), "503") {
		t.Errorf("Expected error to contain status code '503', got: %v", err)
	}
	if !strings.Contains(err.Error(), "non-JSON response") {
		t.Errorf("Expected error to mention 'non-JSON response', got: %v", err)
	}
	// Should mention the log file path
	if !strings.Contains(err.Error(), "logged to") {
		t.Errorf("Expected error to mention log file, got: %v", err)
	}
}

// TestGetAirports_UnexpectedJSONStructure tests valid JSON with unexpected structure
// This tests the case where the API returns valid JSON with 200 status,
// but it doesn't match the expected response structure nor the error structure
func TestGetAirports_UnexpectedJSONStructure(t *testing.T) {
	unexpectedJSON := `{
		"data": {
			"airports": [],
			"metadata": {
				"version": "2.0",
				"timestamp": 1234567890
			}
		},
		"status": "ok"
	}`

	client := newMockClient(200, unexpectedJSON)

	airports, err := client.GetAirports(AirportsParams{
		CodeIso2Country: "US",
	})

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if airports != nil {
		t.Errorf("Expected nil airports, got %v", airports)
	}

	// Error should mention unexpected JSON format
	if !strings.Contains(err.Error(), "unexpected JSON format") {
		t.Errorf("Expected error to mention 'unexpected JSON format', got: %v", err)
	}
	// Should mention the log file path
	if !strings.Contains(err.Error(), "logged to") {
		t.Errorf("Expected error to mention log file, got: %v", err)
	}
}

// TestGetAirports_EmptyArray tests a successful response with empty array
func TestGetAirports_EmptyArray(t *testing.T) {
	emptyResponse := `[]`

	client := newMockClient(200, emptyResponse)

	airports, err := client.GetAirports(AirportsParams{
		CodeIataAirport: "NONEXISTENT",
	})

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if airports == nil {
		t.Fatal("Expected non-nil airports slice, got nil")
	}

	if len(airports) != 0 {
		t.Errorf("Expected 0 airports, got %d", len(airports))
	}
}

// TestGetAirports_PartiallyValidJSON tests JSON that can be parsed but has some invalid fields
func TestGetAirports_PartiallyValidJSON(t *testing.T) {
	partiallyValidJSON := `[
		{
			"airportId": 1,
			"nameAirport": "Test Airport",
			"codeIataAirport": "TST",
			"codeIcaoAirport": "KTST",
			"latitudeAirport": "invalid",
			"longitudeAirport": "invalid",
			"nameCountry": "Test Country",
			"codeIso2Country": "TC"
		}
	]`

	client := newMockClient(200, partiallyValidJSON)

	airports, err := client.GetAirports(AirportsParams{
		CodeIataAirport: "TST",
	})

	// This should fail during unmarshaling because latitude/longitude are strings instead of float64
	if err == nil {
		t.Fatal("Expected error due to invalid field types, got nil")
	}

	if airports != nil {
		t.Errorf("Expected nil airports, got %v", airports)
	}

	// Should fail on unmarshal and go through error handling
	if !strings.Contains(err.Error(), "parse airport response") && !strings.Contains(err.Error(), "unexpected JSON format") {
		t.Errorf("Expected error related to parsing, got: %v", err)
	}
}
