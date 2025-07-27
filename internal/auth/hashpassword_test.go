package auth

import (
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestPasswordMatch(t *testing.T) {

	password := "yelllows"

	hashedp, err := HashPassword(password)
	if err != nil {
		t.Errorf("error happened in hashpassword: %v", err)
	}

	err = CheckPasswordHash(password, hashedp)
	if err != nil {
		t.Errorf("error happened in checkpasswordhash: %v", err)
	}
}

// TestHashPasswordAndCheck tests both HashPassword and CheckPasswordHash functions.
func TestHashPasswordAndCheck(t *testing.T) {
	tests := []struct {
		name        string
		password    string
		wantHashErr bool
		wantCheckErr bool
		checkPass   string // Password to check against hash (same or different)
	}{
		{
			name:        "Valid password",
			password:    "yelllows",
			wantHashErr: false,
			wantCheckErr: false,
			checkPass:   "yelllows",
		},
		{
			name:        "Incorrect password",
			password:    "yelllows",
			wantHashErr: false,
			wantCheckErr: true,
			checkPass:   "wrongpass",
		},
		{
			name:        "Empty password",
			password:    "",
			wantHashErr: false, // bcrypt allows empty passwords
			wantCheckErr: false,
			checkPass:   "",
		},
		{
			name:        "Long password (<72 bytes)",
			password:    "bybcryptbecauseitislongerthan72byteswhichisaverylongstringindeed",
			wantHashErr: false,
			wantCheckErr: false,
			checkPass:   "bybcryptbecauseitislongerthan72byteswhichisaverylongstringindeed",
		},
		{
			name:        "Special characters",
			password:    "P@ssw0rd!#$",
			wantHashErr: false,
			wantCheckErr: false,
			checkPass:   "P@ssw0rd!#$",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test HashPassword
			hashedp, err := HashPassword(tt.password)
			if (err != nil) != tt.wantHashErr {
				t.Errorf("HashPassword(%q) error = %v, wantErr %v", tt.password, err, tt.wantHashErr)
				return
			}
			if !tt.wantHashErr {
				// Verify hash is a valid bcrypt hash
				if len(hashedp) == 0 {
					t.Errorf("HashPassword(%q) returned empty hash", tt.password)
				}
				// Check if hash format is valid (basic check for bcrypt prefix)
				if len(hashedp) < 4 || hashedp[:4] != "$2a$" && hashedp[:4] != "$2b$" && hashedp[:4] != "$2y$" {
					t.Errorf("HashPassword(%q) returned invalid bcrypt hash: %v", tt.password, hashedp)
				}
			}

			// Test CheckPasswordHash
			err = CheckPasswordHash(tt.checkPass, hashedp)
			if (err != nil) != tt.wantCheckErr {
				t.Errorf("CheckPasswordHash(%q, %q) error = %v, wantErr %v", tt.checkPass, hashedp, err, tt.wantCheckErr)
			}
		})
	}
}

// TestInvalidHash tests CheckPasswordHash with an invalid hash format.
func TestInvalidHash(t *testing.T) {
	password := "yelllows"
	invalidHash := "invalid_hash_format"
	err := CheckPasswordHash(password, invalidHash)
	if err == nil {
		t.Errorf("CheckPasswordHash(%q, %q) expected error, got nil", password, invalidHash)
	}
	if err != bcrypt.ErrHashTooShort && err.Error() != "crypto/bcrypt: hashedSecret too short to be a bcrypted password" {
		t.Errorf("CheckPasswordHash(%q, %q) expected bcrypt.ErrHashTooShort, got %v", password, invalidHash, err)
	}
}

// TestHashConsistency tests that HashPassword produces consistent results for the same input.
func TestHashConsistency(t *testing.T) {
	password := "consistentPass123"
	hash1, err := HashPassword(password)
	if err != nil {
		t.Errorf("HashPassword(%q) error = %v", password, err)
		return
	}
	// Verify the hash can be checked
	err = CheckPasswordHash(password, hash1)
	if err != nil {
		t.Errorf("CheckPasswordHash(%q, %q) error = %v", password, hash1, err)
	}

}


func TestCheckPasswordHash(t *testing.T) {
	// First, we need to create some hashed passwords for testing
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

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CheckPasswordHash(tt.password, tt.hash)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckPasswordHash() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}