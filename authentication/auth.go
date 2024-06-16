package authentication

import (
	"encoding/json"
	"fmt"
	// "log"
	"net/http"
	"time"
	"crypto/rand"
    "encoding/base64"
	"errors"

	"github.com/dgrijalva/jwt-go"
	// "github.com/gorilla/mux"
	// "github.com/segmentio/ksuid"
)

type User struct {
	PhoneNumber string
	OTP         string
}

var users = make(map[string]User)

func generateRandomKey(length int) (string, error) {
	randomBytes := make([]byte, length)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(randomBytes), nil
}

var key, err = generateRandomKey(32)
var jwtSecret = []byte(key)

func SendOtpHandler(w http.ResponseWriter, r *http.Request) {
	var request struct {
		PhoneNumber string `json:"phone_number"`
	}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	phoneNumber := request.PhoneNumber
	if phoneNumber == "" {
		http.Error(w, "Phone number is required", http.StatusBadRequest)
		return
	}

	if len(phoneNumber) != 10 {
		http.Error(w, "Phone number must be 10 digit.", http.StatusBadRequest)
		return
	}

	otp := generateOTP()

	fmt.Printf("OTP for %s: %s\n", phoneNumber, otp)
	users[phoneNumber] = User{PhoneNumber: phoneNumber, OTP: otp}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OTP sent successfully"))
}

func VerifyOtpHandler(w http.ResponseWriter, r *http.Request){
	var request struct {
		PhoneNumber string `json:"phone_number"`
		OTP 		string `json:"otp"`
	}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	phoneNumber := request.PhoneNumber
	otp := request.OTP

	user, ok := users[phoneNumber]
	if !ok {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	if user.OTP != otp {
		http.Error(w, "Incorrect OTP", http.StatusUnauthorized)
		return
	}

	accessToken,expToken, err := generateJWTToken(phoneNumber)
	if err != nil {
		http.Error(w, "Failed to create token.", http.StatusInternalServerError)
		return
	}	
	// fmt.Printf(token)

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"access_token": accessToken, "expiry_token": expToken})
}

func GenerateAccessTokenHandler(w http.ResponseWriter, r *http.Request){
	var request struct{
		ExpToken     	string `json:"exp_token"`
		PhoneNumber		string `json:"phone_number"`
	}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	expiryToken := request.ExpToken
	phoneNumber := request.PhoneNumber
	claims, err := verifyExpToken(expiryToken)
	if err != nil {
		http.Error(w, "Invalid expiration token.", http.StatusUnauthorized)
		return
	}

	if claims["phone_number"] != phoneNumber {
		http.Error(w, "Mismatched phone number.", http.StatusForbidden)
		return
	}

	// accessExpiresAt := time.Now().Add(time.Minute * 15).Unix()
	// accessToken, err := generateAccessToken(phoneNumber, accessExpiresAt)
	accessToken,expToken, err := generateJWTToken(phoneNumber)
	if err != nil {
		http.Error(w, "Failed to generate access token", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"access_token": accessToken,"expiry_token": expToken})
}

func VerifyAccessTokenHandler(w http.ResponseWriter, r *http.Request) {
	var request struct {
		AccessToken string `json:"access_token"`
	}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	accessToken := request.AccessToken

	claims, err := verifyAccessToken(accessToken)
	if err != nil {
		http.Error(w, "Invalid access token", http.StatusUnauthorized)
		return
	}

	// Extract phone number from claims
	phoneNumber, ok := claims["phone_number"].(string)
	if !ok {
		http.Error(w, "Invalid access token claims", http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"phone_number": phoneNumber})
}

func generateOTP() string {
	// For simplicity, generating a random 6-digit number as OTP
	randomBytes := make([]byte, 3) // 3 bytes for a 24-bit random number
	rand.Read(randomBytes)
	return fmt.Sprintf("%06d", int(randomBytes[0])<<16|int(randomBytes[1])<<8|int(randomBytes[2]))
}

func generateJWTToken(phoneNumber string) (string,string ,error) {
	// Create a new token with a claims section.
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	// Set token claims
	accessExpiresAt := time.Now().Add(time.Hour * 1).Unix() // Access token expires in 1 hour
	expTokenExpiresAt := time.Now().Add(time.Hour * 24 * 7).Unix() // Exp token expires in 1 week

	// Generate the access token separately
	accessToken, err := generateAccessToken(phoneNumber, accessExpiresAt)
	if err != nil {
		return "","", err
	}
	expToken, err := generateExpiryToken(phoneNumber, expTokenExpiresAt)
	if err != nil {
		return "","", err
	}
	// Set the claims with the generated access token
	claims["phone_number"] = phoneNumber
	claims["access_token"] = accessToken
	claims["exp_token"] = expToken

	// Sign the token with the secret key
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "","", err
	}
	fmt.Printf(tokenString)
	return accessToken,expToken, nil
}

func generateAccessToken(phoneNumber string, expiresAt int64) (string, error) {
	// Create a new access token with a claims section.
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	// Set access token claims
	claims["phone_number"] = phoneNumber
	claims["exp"] = expiresAt

	// Sign the access token with the secret key
	accessTokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return accessTokenString, nil
}
func generateExpiryToken(phoneNumber string, expiresAt int64) (string, error) {
	// Create a new expiry token with a claims section.
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	// Set expiry token claims
	claims["phone_number"] = phoneNumber
	claims["exp"] = expiresAt

	// Sign the expiry token with the secret key
	expTokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return expTokenString, nil
}
func verifyExpToken(expToken string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(expToken, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method and key
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("Invalid token")
	}

	return claims, nil
}

func verifyAccessToken(accessToken string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method and key
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("Invalid token")
	}

	return claims, nil
}