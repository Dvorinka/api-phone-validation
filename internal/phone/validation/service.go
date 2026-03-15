package validation

import (
	"strings"

	"github.com/nyaruka/phonenumbers"
)

type Service struct{}

func NewService() *Service {
	return &Service{}
}

type ValidateInput struct {
	Phone         string `json:"phone"`
	DefaultRegion string `json:"default_region"`
}

type ValidateResult struct {
	Input              string `json:"input"`
	E164               string `json:"e164,omitempty"`
	National           string `json:"national,omitempty"`
	International      string `json:"international,omitempty"`
	IsValid            bool   `json:"is_valid"`
	IsPossible         bool   `json:"is_possible"`
	CountryCode        int    `json:"country_code,omitempty"`
	Region             string `json:"region,omitempty"`
	LineType           string `json:"line_type"`
	Carrier            string `json:"carrier,omitempty"`
	IsVoIP             bool   `json:"is_voip"`
	IsDisposableVoIP   bool   `json:"is_disposable_voip"`
	ValidationError    string `json:"validation_error,omitempty"`
	NormalizedInputRaw string `json:"normalized_input_raw"`
}

func (s *Service) Validate(input ValidateInput) ValidateResult {
	raw := strings.TrimSpace(input.Phone)
	region := strings.ToUpper(strings.TrimSpace(input.DefaultRegion))

	result := ValidateResult{
		Input:              raw,
		LineType:           "unknown",
		NormalizedInputRaw: compact(raw),
	}

	if raw == "" {
		result.ValidationError = "phone is required"
		return result
	}

	number, err := phonenumbers.Parse(raw, region)
	if err != nil {
		result.ValidationError = err.Error()
		return result
	}

	result.IsPossible = phonenumbers.IsPossibleNumber(number)
	result.IsValid = phonenumbers.IsValidNumber(number)
	result.CountryCode = int(number.GetCountryCode())
	result.Region = phonenumbers.GetRegionCodeForNumber(number)

	numberType := phonenumbers.GetNumberType(number)
	result.LineType = phoneTypeToString(numberType)
	result.IsVoIP = numberType == phonenumbers.VOIP
	carrier, carrierErr := phonenumbers.GetCarrierForNumber(number, "en")
	if carrierErr == nil {
		result.Carrier = carrier
	}
	result.IsDisposableVoIP = isDisposableVoIP(result.IsVoIP, result.Carrier)

	if result.IsPossible {
		result.E164 = phonenumbers.Format(number, phonenumbers.E164)
		result.National = phonenumbers.Format(number, phonenumbers.NATIONAL)
		result.International = phonenumbers.Format(number, phonenumbers.INTERNATIONAL)
	}

	return result
}

func phoneTypeToString(numberType phonenumbers.PhoneNumberType) string {
	switch numberType {
	case phonenumbers.FIXED_LINE:
		return "landline"
	case phonenumbers.MOBILE:
		return "mobile"
	case phonenumbers.FIXED_LINE_OR_MOBILE:
		return "fixed_line_or_mobile"
	case phonenumbers.TOLL_FREE:
		return "toll_free"
	case phonenumbers.PREMIUM_RATE:
		return "premium_rate"
	case phonenumbers.SHARED_COST:
		return "shared_cost"
	case phonenumbers.VOIP:
		return "voip"
	case phonenumbers.PERSONAL_NUMBER:
		return "personal_number"
	case phonenumbers.PAGER:
		return "pager"
	case phonenumbers.UAN:
		return "uan"
	case phonenumbers.VOICEMAIL:
		return "voicemail"
	default:
		return "unknown"
	}
}

func isDisposableVoIP(isVoIP bool, carrier string) bool {
	if !isVoIP {
		return false
	}

	lower := strings.ToLower(strings.TrimSpace(carrier))
	if lower == "" {
		return true
	}

	// Lightweight heuristic until we integrate commercial intelligence feeds.
	highRiskKeywords := []string{
		"virtual",
		"voip",
		"disposable",
		"temp",
		"anonymous",
	}

	for _, keyword := range highRiskKeywords {
		if strings.Contains(lower, keyword) {
			return true
		}
	}
	return false
}

func compact(value string) string {
	value = strings.ReplaceAll(value, " ", "")
	value = strings.ReplaceAll(value, "-", "")
	value = strings.ReplaceAll(value, "(", "")
	value = strings.ReplaceAll(value, ")", "")
	return value
}
