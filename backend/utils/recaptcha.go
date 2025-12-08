package utils

import (
	"encoding/json"
	"log"
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
		log.Println("⚠️ reCAPTCHA: No secret key configured, skipping verification")
		return true, nil
	}

	if token == "" {
		log.Println("❌ reCAPTCHA: Empty token received")
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

	// Log the result for debugging
	log.Printf("🔐 reCAPTCHA: success=%v, score=%.2f, action=%s, hostname=%s",
		result.Success, result.Score, result.Action, result.Hostname)

	// Check if verification was successful and score is acceptable
	// Score >= 0.5 is considered human (as recommended by Google)
	if result.Success && result.Score >= 0.5 {
		log.Println("✅ reCAPTCHA: Verification PASSED")
		return true, nil
	}

	log.Printf("❌ reCAPTCHA: Verification FAILED (score too low or errors: %v)", result.ErrorCodes)
	return false, nil
}
