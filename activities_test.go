package iplocate_test

import (
	"io"
	"iplocate"
	"net/http"
	"strings"
	"testing"
)

type MockHTTPClient struct {
	Response *http.Response
	Err      error
}

func (m *MockHTTPClient) Get(url string) (*http.Response, error) {
	return m.Response, m.Err
}

// TestGetIP tests the GetIP activity with a mock server.
func TestGetIP(t *testing.T) {
	// Create a mock server that returns the fake IP address

	mockResponse := &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(strings.NewReader("127.0.0.1\n")),
	}

	ipActivities := iplocate.IPActivities{
		HTTPClient: &MockHTTPClient{Response: mockResponse},
	}

	// Call the GetIP function
	ip, err := ipActivities.GetIP()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Validate the returned IP
	expectedIP := "127.0.0.1"
	if ip != expectedIP {
		t.Fatalf("Expected IP to be '%s', but got '%s'", expectedIP, ip)
	}
}

// TestGetLocationInfo tests the GetLocationInfo activity with a mock server.
func TestGetLocationInfo(t *testing.T) {
	mockResponse := &http.Response{
		StatusCode: 200,
		Body: io.NopCloser(strings.NewReader(`{
            "city": "San Francisco",
            "regionName": "California",
            "country": "United States"
        }`)),
	}

	ipActivities := iplocate.IPActivities{
		HTTPClient: &MockHTTPClient{Response: mockResponse},
	}

	ip := "127.0.0.1"
	location, err := ipActivities.GetLocationInfo(ip)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	expectedLocation := "San Francisco, California, United States"
	if location != expectedLocation {
		t.Errorf("Expected location %v, got %v", expectedLocation, location)
	}
}
