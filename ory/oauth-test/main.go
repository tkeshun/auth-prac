package main

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	HydraPublicURL = "http://localhost:4444"
	ClientID       = "my-client"
	RedirectURI    = "http://127.0.0.1:5555/callback"
	Scope          = "openid profile"
)

// 乱数を用いた `state` 文字列の生成
func generateState() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	// Go 1.20 以降の推奨方法
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	b := make([]byte, 16) // 16 文字のランダムな state
	for i := range b {
		b[i] = charset[rng.Intn(len(charset))]
	}
	return string(b)
}

// PKCE `code_verifier` & `code_challenge` 生成
func generatePKCE() (string, string) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-._~"
	codeVerifier := "mock_verifier_string_0123456789"

	hash := sha256.Sum256([]byte(codeVerifier))
	codeChallenge := base64.RawURLEncoding.EncodeToString(hash[:])

	return codeVerifier, codeChallenge
}

// **1. 認可リクエスト**
func requestAuthorization(codeChallenge string) (string, error) {
	state := generateState() // state を生成
	authURL := fmt.Sprintf("%s/oauth2/auth?client_id=%s&response_type=code&redirect_uri=%s&scope=%s&state=%s&code_challenge=%s&code_challenge_method=S256",
		HydraPublicURL, ClientID, url.QueryEscape(RedirectURI), url.QueryEscape(Scope), state, codeChallenge)

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			// リダイレクト時にURLを取得し、処理を行う
			log.Println("Redirected to:", req.URL.String())

			// Hydra からログインプロバイダーにリダイレクトされた場合、手動でリクエストを送る
			if strings.Contains(req.URL.String(), "http://localhost:8080/login") {
				return handleLogin(req.URL.String())
			}

			// Hydra の Consent 画面へリダイレクトされた場合、無視
			if strings.Contains(req.URL.String(), "consent_challenge") {
				return http.ErrUseLastResponse // リダイレクトをキャンセル
			}

			return nil
		},
	}

	resp, err := client.Get(authURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// 最終的なリダイレクト先URLを取得
	redirectedURL := resp.Request.URL.String()
	fmt.Printf("redirect url: %s", redirectedURL)
	parsedURL, err := url.Parse(redirectedURL)
	if err != nil {
		return "", err
	}

	// 認可コード取得
	authCode := parsedURL.Query().Get("code")
	if authCode == "" {
		return "", fmt.Errorf("authorization code not found in response")
	}
	return authCode, nil
}

// **2. モックログインプロバイダーへアクセス**
func handleLogin(loginURL string) error {
	log.Println("Accessing Mock Login Provider:", loginURL)

	// モックログインプロバイダーにアクセス（ログインチャレンジを取得）
	resp, err := http.Get(loginURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// モックログインプロバイダーが Hydra にログイン成功を通知すると、リダイレクトURLが返る
	redirectTo := resp.Request.URL.String()
	log.Println("Mock Login Provider Redirected to:", redirectTo)

	// 取得したURLへリダイレクト（Hydra に戻る）
	resp, err = http.Get(redirectTo)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

// **3. トークンリクエスト**
func requestToken(authCode, codeVerifier string) (string, error) {
	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("code", authCode)
	data.Set("redirect_uri", RedirectURI)
	data.Set("client_id", ClientID)
	data.Set("code_verifier", codeVerifier)

	req, err := http.NewRequest("POST", HydraPublicURL+"/oauth2/token", strings.NewReader(data.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to get token: %s", body)
	}

	// レスポンスからアクセストークン取得
	var tokenResp struct {
		AccessToken string `json:"access_token"`
	}
	json.Unmarshal(body, &tokenResp)
	return tokenResp.AccessToken, nil
}

func main() {
	codeVerifier, codeChallenge := generatePKCE()
	log.Println("PKCE Code Verifier:", codeVerifier)
	log.Println("PKCE Code Challenge:", codeChallenge)

	// 認可コード取得
	authCode, err := requestAuthorization(codeChallenge)
	if err != nil {
		log.Fatalf("Authorization request failed: %v", err)
	}
	log.Println("Authorization Code:", authCode)

	// アクセストークン取得
	accessToken, err := requestToken(authCode, codeVerifier)
	if err != nil {
		log.Fatalf("Token request failed: %v", err)
	}
	log.Println("Access Token:", accessToken)
}
