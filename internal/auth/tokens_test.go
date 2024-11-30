package auth

import (
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestCreateToken(t *testing.T) {
	token, err := MakeJWT(uuid.New(), "test", time.Second*2)
	if err != nil {
		t.Errorf("Failed to create token: %v", err)
	}

	t.Logf("Got token : %v", token)

}

func TestCheckingToken(t *testing.T) {
	userID := uuid.New()
	otherID := uuid.New()
	secret := "test"

	token, err := MakeJWT(userID, secret, time.Second*2)
	if err != nil {
		t.Errorf("Failed to create token: %v", err)
	}

	otherToken, err := MakeJWT(otherID, secret, time.Second*2)
	if err != nil {
		t.Errorf("Failed to create other token: %v", err)
	}

	tests := []struct {
		name        string
		token       string
		tokenSecret string
		userID      uuid.UUID
		wantERR     bool
	}{
		{
			name:        "Correct Token",
			token:       token,
			tokenSecret: secret,
			userID:      userID,
			wantERR:     false,
		},
		{
			name:        "Incorrect Token",
			token:       otherToken,
			tokenSecret: secret,
			userID:      userID,
			wantERR:     true,
		},
		{
			name:        "Incorrect Secret",
			token:       token,
			tokenSecret: "blah",
			userID:      userID,
			wantERR:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ID, err := ValidateJWT(tt.token, tt.tokenSecret)
			if err != nil && !tt.wantERR {
				t.Errorf("Failed validating token error: %v want error: %v", err, tt.wantERR)
			}
			if ID != tt.userID && !tt.wantERR {
				t.Errorf("UserID does not match token id: %v != %v", tt.userID, ID)
			}

		})
	}
}

func TestTokenExperation(t *testing.T) {
	userID := uuid.New()
	token, err := MakeJWT(userID, "test", time.Second*2)
	if err != nil {
		t.Errorf("Failed to create token: %v", err)
	}

	_, err = ValidateJWT(token, "test")
	if err != nil {
		t.Errorf("Invalid token: %v", err)
	}

	time.Sleep(time.Second * 3)

	_, err = ValidateJWT(token, "test")
	if err == nil {
		t.Error("Token was supposed to expire")
	}
}

func TestGetBearer(t *testing.T) {

	testHeader := http.Header{}
	testHeader.Set("Authorization", "Bearer sometokentogetfortheprogram")
	otherHeader := http.Header{}
	otherHeader.Set("Authorization", "Bearer ")

	tests := []struct {
		name    string
		header  http.Header
		token   string
		wantERR bool
	}{
		{
			name:    "Correct Authorization Header",
			header:  testHeader,
			token:   "sometokentogetfortheprogram",
			wantERR: false,
		},
		{
			name:    "Incorrect Authorization Header",
			header:  otherHeader,
			token:   "doesn'tmatter",
			wantERR: true,
		},
		{
			name:    "Missing Authorization Header",
			header:  http.Header{},
			token:   "",
			wantERR: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := GetBearerToken(tt.header)
			if err != nil && !tt.wantERR {
				t.Errorf("Failed to get correct Bearer Token | expecting: %v | but got error: %v", tt.token, err)
			}
			if token != tt.token && !tt.wantERR {
				t.Errorf("Token does not match expected | expecting: %v | got: %v", tt.token, token)
			}
		})
	}

}
