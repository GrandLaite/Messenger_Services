package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	secret string
}

func NewAuthService(sec string) *AuthService {
	return &AuthService{secret: sec}
}

func (s *AuthService) GenerateToken(username, role string) (string, error) {
	claims := jwt.MapClaims{
		"username": username,
		"role":     role,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString([]byte(s.secret))
}

func (s *AuthService) LoginUser(username, password string) (string, error) {
	url := os.Getenv("USER_SERVICE_URL")
	if url == "" {
		url = "http://localhost:8082"
	}
	postURL := url + "/users/checkpassword"
	var reqData struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	reqData.Username = username
	reqData.Password = password
	buf, _ := json.Marshal(reqData)
	rq, err := http.NewRequest("POST", postURL, bytes.NewBuffer(buf))
	if err != nil {
		return "", err
	}
	rq.Header.Set("Content-Type", "application/json")
	cl := &http.Client{}
	resp, err := cl.Do(rq)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return "", errors.New(string(b))
	}
	var res struct {
		Username string `json:"username"`
		Role     string `json:"role"`
	}
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return "", err
	}
	return s.GenerateToken(res.Username, res.Role)
}

func HashPassword(pw string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func CheckPassword(hash, pw string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(pw))
}
