package utils

import (
	"encoding/json"
	"net/http"
	"net/url"
	"os"
	"time"
)

// RecaptchaResponse represents Google's reCAPTCHA verification response
type RecaptchaResponse struct {
	Success     bool      `json:"success"`
	Score       float64   `json:"score"`
	Action      string    `json:"action"`
	ChallengeTS time.Time `json:"challenge_ts"`
	Hostname    string    `json:"hostname"`
	ErrorCodes  []string  `json:"error-codes"`
}

// VerifyRecaptcha verifies the reCAPTCHA token with Google's API
// Returns true if verification passes (score >= 0.5)
func VerifyRecaptcha(token string) (bool, error) {
	secretKey := os.Getenv("RECAPTCHA_SECRET_KEY")
	if secretKey == "" {
		// If no secret key configured, skip verification (for development)
		return true, nil
	}

	if token == "" {
		return false, nil
	}

	// Call Google's verification API
	resp, err := http.PostForm("https://www.google.com/recaptcha/api/siteverify",
		url.Values{
			"secret":   {secretKey},
			"response": {token},
		},
	)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	// Parse response
	var result RecaptchaResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return false, err
	}

	// Check if verification was successful and score is acceptable
	// Score >= 0.5 is considered human (as recommended by Google)
	if result.Success && result.Score >= 0.5 {
		return true, nil
	}

	return false, nil
}
