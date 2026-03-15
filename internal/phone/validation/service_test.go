package validation

import "testing"

func TestValidateValidNumber(t *testing.T) {
	service := NewService()

	result := service.Validate(ValidateInput{
		Phone:         "+14155552671",
		DefaultRegion: "US",
	})

	if !result.IsPossible {
		t.Fatalf("expected number to be possible")
	}
	if !result.IsValid {
		t.Fatalf("expected number to be valid")
	}
	if result.E164 == "" {
		t.Fatalf("expected e164 format")
	}
	if result.CountryCode != 1 {
		t.Fatalf("expected country code 1, got %d", result.CountryCode)
	}
}

func TestValidateInvalidNumber(t *testing.T) {
	service := NewService()

	result := service.Validate(ValidateInput{
		Phone:         "123",
		DefaultRegion: "US",
	})

	if result.IsValid {
		t.Fatalf("expected number to be invalid")
	}
	if result.LineType == "" {
		t.Fatalf("expected line type response")
	}
}
