package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

var hydraAdminURL = "http://localhost:4445"

func main() {
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/consent", consentHandler)
	http.HandleFunc("/callback", callbackHandler)

	log.Println("Mock login provider running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// ログイン処理
func loginHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		// ログイン画面のHTMLを表示
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, `
			<form method="POST">
				<label>ユーザーID: <input name="username"></label><br>
				<button type="submit">ログイン</button>
			</form>
		`)
	case http.MethodPost:
		// ユーザーがフォームを送信したと仮定
		loginChallenge := r.URL.Query().Get("login_challenge")
		if loginChallenge == "" {
			http.Error(w, "login_challenge is missing", http.StatusBadRequest)
			return
		}

		// Hydra にログイン成功を通知
		data := map[string]string{
			"challenge": loginChallenge,
			"subject":   "mock-user",
		}
		jsonData, _ := json.Marshal(data)

		resp, err := http.Post(hydraAdminURL+"/oauth2/auth/requests/login/accept",
			"application/json", bytes.NewReader(jsonData)) // ← 修正
		if err != nil {
			http.Error(w, "Failed to communicate with Hydra", http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		var response map[string]string
		json.NewDecoder(resp.Body).Decode(&response)

		// Hydra のリダイレクトURLへ移動
		http.Redirect(w, r, response["redirect_to"], http.StatusFound)
	}
}

// 同意（Consent）処理
func consentHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		// 同意画面のHTMLを表示
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, `
			<form method="POST">
				<p>このアプリが次の情報を取得します: openid, profile</p>
				<button type="submit">同意する</button>
			</form>
		`)
	case http.MethodPost:
		consentChallenge := r.URL.Query().Get("consent_challenge")
		if consentChallenge == "" {
			http.Error(w, "consent_challenge is missing", http.StatusBadRequest)
			return
		}

		// Hydra に同意を通知
		data := map[string]interface{}{
			"challenge":                   consentChallenge,
			"grant_scope":                 []string{"openid", "profile"},
			"grant_access_token_audience": []string{"test-client"},
			"session": map[string]interface{}{
				"id_token": map[string]string{
					"sub":   "mock-user",
					"name":  "Test User",
					"email": "mock@example.com",
				},
			},
		}
		jsonData, _ := json.Marshal(data)

		resp, err := http.Post(hydraAdminURL+"/oauth2/auth/requests/consent/accept",
			"application/json", bytes.NewReader(jsonData)) // ← 修正
		if err != nil {
			http.Error(w, "Failed to communicate with Hydra", http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		var response map[string]string
		json.NewDecoder(resp.Body).Decode(&response)

		// Hydra のリダイレクトURLへ移動
		http.Redirect(w, r, response["redirect_to"], http.StatusFound)
	}
}

// 認可コードを受け取るエンドポイント
func callbackHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")

	fmt.Fprintf(w, "Authorization Code: %s\nState: %s", code, state)
}
