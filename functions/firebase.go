package functions

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type FirebaseTokenInfo struct {
	LocalID string `json:"localId"`
	Email   string `json:"email,omitempty"`
	Name    string `json:"displayName,omitempty"`
	Photo   string `json:"photoUrl,omitempty"`
}

func VerifyFirebaseToken(idToken string) (string, error) {
	url := fmt.Sprintf("https://identitytoolkit.googleapis.com/v1/accounts:lookup?key=%s", "AIzaSyBmSkCnePbHTi2BcngOIVekwP7CxJJ0SzQ")

	payload := map[string]interface{}{
		"idToken": idToken,
	}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result struct {
		Users []FirebaseTokenInfo `json:"users"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}
	if len(result.Users) == 0 {
		return "", fmt.Errorf("no user found")
	}

	return result.Users[0].LocalID, nil
}
