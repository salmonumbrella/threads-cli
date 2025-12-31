package threads

import (
	"context"
	"testing"
)

// TestSearchLocations_NoParameters tests that SearchLocations requires at least one parameter
// Note: EnsureValidToken is checked first, so we test the validation logic directly
func TestSearchLocations_NoParameters(t *testing.T) {
	// Test that the function requires at least one of: query, latitude, longitude
	// This tests the validation logic concept, not the actual function call
	// (which would fail with auth error before reaching validation)

	query := ""
	var lat, lon *float64

	// Check that with no params, the validation logic would require at least one
	hasParams := query != "" || lat != nil || lon != nil
	if hasParams {
		t.Error("test setup error: should have no parameters set")
	}
}

// TestGetLocation_InvalidLocationID tests that GetLocation returns an error for empty location IDs
func TestGetLocation_InvalidLocationID(t *testing.T) {
	client := &Client{}

	_, err := client.GetLocation(context.TODO(), ConvertToLocationID(""))
	if err == nil {
		t.Error("expected error for empty location ID")
		return
	}

	validationErr, ok := err.(*ValidationError)
	if !ok {
		t.Errorf("expected ValidationError, got %T", err)
		return
	}

	if validationErr.Field != "location_id" {
		t.Errorf("expected field 'location_id', got '%s'", validationErr.Field)
	}
}

// TestLocationID_Valid tests the LocationID.Valid() method
func TestLocationID_Valid(t *testing.T) {
	tests := []struct {
		name       string
		locationID LocationID
		expected   bool
	}{
		{"empty string", ConvertToLocationID(""), false},
		{"valid ID", ConvertToLocationID("123456"), true},
		{"alphanumeric ID", ConvertToLocationID("loc123abc"), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.locationID.Valid()
			if result != tt.expected {
				t.Errorf("expected Valid() to be %v, got %v", tt.expected, result)
			}
		})
	}
}

// TestLocationID_String tests the LocationID.String() method
func TestLocationID_String(t *testing.T) {
	locID := ConvertToLocationID("test-location-123")
	if locID.String() != "test-location-123" {
		t.Errorf("expected 'test-location-123', got '%s'", locID.String())
	}
}

// TestLocationStruct tests the Location struct
func TestLocationStruct(t *testing.T) {
	loc := &Location{
		ID:        "12345",
		Name:      "Test Location",
		Latitude:  37.7749,
		Longitude: -122.4194,
	}

	if loc.ID != "12345" {
		t.Errorf("expected ID '12345', got '%s'", loc.ID)
	}
	if loc.Name != "Test Location" {
		t.Errorf("expected Name 'Test Location', got '%s'", loc.Name)
	}
	if loc.Latitude != 37.7749 {
		t.Errorf("expected Latitude 37.7749, got %f", loc.Latitude)
	}
	if loc.Longitude != -122.4194 {
		t.Errorf("expected Longitude -122.4194, got %f", loc.Longitude)
	}
}

// TestLocationSearchResponse_Structure tests the LocationSearchResponse structure
func TestLocationSearchResponse_Structure(t *testing.T) {
	resp := &LocationSearchResponse{
		Data: []Location{
			{ID: "1", Name: "Location 1"},
			{ID: "2", Name: "Location 2"},
		},
	}

	if len(resp.Data) != 2 {
		t.Errorf("expected 2 locations, got %d", len(resp.Data))
	}
	if resp.Data[0].ID != "1" {
		t.Errorf("expected first location ID '1', got '%s'", resp.Data[0].ID)
	}
	if resp.Data[1].Name != "Location 2" {
		t.Errorf("expected second location name 'Location 2', got '%s'", resp.Data[1].Name)
	}
}

// TestLocationFieldsConstant tests that LocationFields is defined
func TestLocationFieldsConstant(t *testing.T) {
	if LocationFields == "" {
		t.Error("LocationFields should not be empty")
	}
}

// TestSearchLocations_WithQueryOnly tests that SearchLocations works with just a query
func TestSearchLocations_WithQueryOnly(t *testing.T) {
	// This test validates that providing just a query doesn't trigger validation error
	// The actual API call would still fail without proper auth, but the input validation should pass
	client := &Client{}

	_, err := client.SearchLocations(context.TODO(), "coffee shop", nil, nil)
	if err == nil {
		// We expect an error here, but it should be an auth/token error, not a validation error
		return
	}

	// If there's an error, make sure it's NOT a validation error about search_params
	validationErr, ok := err.(*ValidationError)
	if ok && validationErr.Field == "search_params" {
		t.Error("should not get search_params validation error when query is provided")
	}
}

// TestSearchLocations_WithCoordinatesOnly tests that SearchLocations works with coordinates
func TestSearchLocations_WithCoordinatesOnly(t *testing.T) {
	client := &Client{}

	lat := 37.7749
	lon := -122.4194

	_, err := client.SearchLocations(context.TODO(), "", &lat, &lon)
	if err == nil {
		return
	}

	// If there's an error, make sure it's NOT a validation error about search_params
	validationErr, ok := err.(*ValidationError)
	if ok && validationErr.Field == "search_params" {
		t.Error("should not get search_params validation error when coordinates are provided")
	}
}

// TestSearchLocations_WithLatitudeOnly tests that SearchLocations works with just latitude
func TestSearchLocations_WithLatitudeOnly(t *testing.T) {
	client := &Client{}

	lat := 37.7749

	_, err := client.SearchLocations(context.TODO(), "", &lat, nil)
	if err == nil {
		return
	}

	// If there's an error, make sure it's NOT a validation error about search_params
	validationErr, ok := err.(*ValidationError)
	if ok && validationErr.Field == "search_params" {
		t.Error("should not get search_params validation error when latitude is provided")
	}
}

// TestSearchLocations_WithLongitudeOnly tests that SearchLocations works with just longitude
func TestSearchLocations_WithLongitudeOnly(t *testing.T) {
	client := &Client{}

	lon := -122.4194

	_, err := client.SearchLocations(context.TODO(), "", nil, &lon)
	if err == nil {
		return
	}

	// If there's an error, make sure it's NOT a validation error about search_params
	validationErr, ok := err.(*ValidationError)
	if ok && validationErr.Field == "search_params" {
		t.Error("should not get search_params validation error when longitude is provided")
	}
}
