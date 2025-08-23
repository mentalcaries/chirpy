package auth

import (
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestHashedPassword(t *testing.T) {
	password1 := "correctPassword123!"
	password2 := "anotherPassword456!"
	hash1, _ := HashPassword(password1)
	hash2, _ := HashPassword(password2)

	tests := []struct {
		name     string
		password string
		hash     string
		wantErr  bool
	}{
		{
			name:     "Correct password",
			password: password1,
			hash:     hash1,
			wantErr:  false,
		},
		{
			name:     "Incorrect password",
			password: "wrongPassword",
			hash:     hash1,
			wantErr:  true,
		},
		{
			name:     "Password doesn't match different hash",
			password: password1,
			hash:     hash2,
			wantErr:  true,
		},
		{
			name:     "Empty password",
			password: "",
			hash:     hash1,
			wantErr:  true,
		},
		{
			name:     "Invalid hash",
			password: password1,
			hash:     "invalidhash",
			wantErr:  true,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			err := CheckPasswordHash(testCase.password, testCase.hash)
			if (err != nil) != testCase.wantErr {
				t.Errorf("CheckPasswordHash() error = %v, wantErr %v", err, testCase.wantErr)
			}
		})
	}
}

func TestJWT(t *testing.T) {
	uuid1, _ := uuid.Parse("2de4cae3-a60b-44bc-90c7-ae4ebb1e3dc5")
	uuid2, _ := uuid.Parse("8005f023-9ce3-4f53-8ddd-63efc7dc935b")
	expiresIn := 5 * time.Second

	secretKey := "supersecretkey"

	jwt1, _ := MakeJWT(uuid1, secretKey, expiresIn)
	jwt2, _ := MakeJWT(uuid2, secretKey, expiresIn)

	tests := []struct {
		name        string
		uuid        uuid.UUID
		decodedUUID uuid.UUID
		JWT         string
		wantErr     bool
	}{
		{
			name:        "Correctly decode JWT1",
			uuid:        uuid1,
			decodedUUID: uuid1,
			JWT:         jwt1,
			wantErr:     false,
		},
		{
			name:        "Correctly decode JWT2",
			uuid:        uuid2,
			decodedUUID: uuid2,
			JWT:         jwt2,
			wantErr:     false,
		},
		{
			name:        "Invalid JWT decoded",
			uuid:        uuid2,
			decodedUUID: uuid1,
			JWT:         jwt1,
			wantErr:     true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			decodedUUID, _ := ValidateJWT(test.JWT, secretKey)

			if (decodedUUID != test.uuid) != test.wantErr {
				t.Errorf("decoding failed")
			}
		})
	}
}

func TestValidateJWT(t *testing.T) {
	userID := uuid.New()
	validToken, _ := MakeJWT(userID, "secret", time.Hour)

	tests := []struct {
		name        string
		tokenString string
		tokenSecret string
		wantUserID  uuid.UUID
		wantErr     bool
	}{
		{
			name:        "Valid token",
			tokenString: validToken,
			tokenSecret: "secret",
			wantUserID:  userID,
			wantErr:     false,
		},
		{
			name:        "Invalid token",
			tokenString: "invalid.token.string",
			tokenSecret: "secret",
			wantUserID:  uuid.Nil,
			wantErr:     true,
		},
		{
			name:        "Wrong secret",
			tokenString: validToken,
			tokenSecret: "wrong_secret",
			wantUserID:  uuid.Nil,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotUserID, err := ValidateJWT(tt.tokenString, tt.tokenSecret)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateJWT() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotUserID != tt.wantUserID {
				t.Errorf("ValidateJWT() gotUserID = %v, want %v", gotUserID, tt.wantUserID)
			}
		})
	}
}

func TestExtractBearerToken(t *testing.T){
    token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"
    testHeader := http.Header{}
    testHeader.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c")

    t.Run("extract token", func(t *testing.T) {
        extractedToken, err := GetBearerToken(testHeader)
        if (err != nil) != false {
            t.Errorf("extraction failed: %s", err)
        }
        if extractedToken != token {
            t.Errorf("failed extraction: %v", extractedToken)
        }
    })

}

func TestGetBearerToken(t *testing.T) {
	tests := []struct {
		name      string
		headers   http.Header
		wantToken string
		wantErr   bool
	}{
		{
			name: "Valid Bearer token",
			headers: http.Header{
				"Authorization": []string{"Bearer valid_token"},
			},
			wantToken: "valid_token",
			wantErr:   false,
		},
		{
			name:      "Missing Authorization header",
			headers:   http.Header{},
			wantToken: "",
			wantErr:   true,
		},
		{
			name: "Malformed Authorization header",
			headers: http.Header{
				"Authorization": []string{"InvalidBearer token"},
			},
			wantToken: "",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotToken, err := GetBearerToken(tt.headers)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetBearerToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotToken != tt.wantToken {
				t.Errorf("GetBearerToken() gotToken = %v, want %v", gotToken, tt.wantToken)
			}
		})
	}
}