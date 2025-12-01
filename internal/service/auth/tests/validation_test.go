package tests

import (
	"testing"

	"theb-backend/internal/service/auth/utils"
)

func TestValidatePhoneNumber(t *testing.T) {
	tests := []struct {
		name    string
		phone   string
		wantErr bool
	}{
		{"Valid Jordan number", "+962791234567", false},
		{"Valid US number", "+12025551234", false},
		{"Invalid - no plus", "962791234567", true},
		{"Invalid - empty", "", true},
		{"Invalid - too short", "+1234", true},
		{"Invalid - letters", "+96279ABCDEFG", true},
		{"Invalid - starts with zero", "+0791234567", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := utils.ValidatePhoneNumber(tt.phone)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidatePhoneNumber() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNormalizePhoneNumber(t *testing.T) {
	tests := []struct {
		name  string
		phone string
		want  string
	}{
		{"With spaces", "+962 79 123 4567", "+962791234567"},
		{"With dashes", "+962-79-123-4567", "+962791234567"},
		{"With parentheses", "+962(79)1234567", "+962791234567"},
		{"Without plus", "962791234567", "+962791234567"},
		{"Already clean", "+962791234567", "+962791234567"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := utils.NormalizePhoneNumber(tt.phone); got != tt.want {
				t.Errorf("NormalizePhoneNumber() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidateName(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"Valid English name", "John Doe", false},
		{"Valid Arabic name", "أحمد علي", false},
		{"Valid with numbers", "Ali123", false},
		{"Too short", "A", true},
		{"Empty", "", true},
		{"Only spaces", "   ", true},
		{"Too long", "abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyz", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := utils.ValidateName(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateName() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGenerateOTP(t *testing.T) {
	tests := []struct {
		name       string
		length     int
		wantErr    bool
		wantLength int
	}{
		{"6 digits", 6, false, 6},
		{"4 digits", 4, false, 4},
		{"8 digits", 8, false, 8},
		{"Invalid length", 0, true, 0},
		{"Negative length", -1, true, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := utils.GenerateOTP(tt.length)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateOTP() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && len(got) != tt.wantLength {
				t.Errorf("GenerateOTP() length = %v, want %v", len(got), tt.wantLength)
			}
			// Check if OTP contains only digits
			if !tt.wantErr {
				for _, c := range got {
					if c < '0' || c > '9' {
						t.Errorf("GenerateOTP() contains non-digit character: %c", c)
					}
				}
			}
		})
	}
}
