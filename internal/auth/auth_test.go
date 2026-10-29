package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestCreateJWT(t *testing.T) {
	// Test data
	userID := "123e4567-e89b-12d3-a456-426614174000"
	tokenSecret := "mysecretkey"
	expiresIn := 1 * time.Hour

	// Test MakeJWT
	token, err := MakeJWT(uuid.MustParse(userID), tokenSecret, expiresIn)
	if err != nil {
		t.Fatalf("MakeJWT failed: %v", err)
	}

	if token == "" {
		t.Fatal("MakeJWT returned an empty token")
	}
}

func TestValidateJWT(t *testing.T) {
	// Test data
	userID := "123e4567-e89b-12d3-a456-426614174000"
	tokenSecret := "mysecretkey"
	expiresIn := 1 * time.Hour

	// Create a valid token
	token, err := MakeJWT(uuid.MustParse(userID), tokenSecret, expiresIn)
	if err != nil {
		t.Fatalf("MakeJWT failed: %v", err)
	}

	// Test ValidateJWT with a valid token
	parsedUserID, err := ValidateJWT(token, tokenSecret)
	if err != nil {
		t.Fatalf("ValidateJWT failed: %v", err)
	}

	if parsedUserID.String() != userID {
		t.Fatalf("ValidateJWT returned wrong userID: got %v, want %v", parsedUserID, userID)
	}

}

func TestBearerToken(t *testing.T) {
	// Test data
	headers := make(map[string][]string)
	headers["Authorization"] = []string{"Bearer mytoken"}

	// Test GetBearerToken
	token, err := GetBearerToken(headers)
	if err != nil {
		t.Fatalf("GetBearerToken failed: %v", err)
	}

	if token != "mytoken" {
		t.Fatalf("GetBearerToken returned wrong token: got %v, want %v", token, "mytoken")
	}

}

func TestBearerTokenMissing(t *testing.T) {
	// Test data
	headers := make(map[string][]string)

	// Test GetBearerToken with missing Authorization header
	_, err := GetBearerToken(headers)

	if err == nil {
		t.Fatal("GetBearerToken did not return an error for missing Authorization header")
	}
}
