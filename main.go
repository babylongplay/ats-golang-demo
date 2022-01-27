package atsgolangdemo

import "fmt"

const (
	apiURL        = `https://api-prod.autosystem.io`
	operatorToken = `TOKEN`
	secretKey     = `KEY`
)

func main() {
	// Request Client
	request := NewRequest(apiURL, operatorToken, secretKey)

	// create user
	payload := map[string]interface{}{
		"userName":     "",
		"name":         "",
		"registeredIp": "",
	}

	resp, err := request.Post("/api/user/create", payload, 120)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s: API status %d", "Create User", resp.Code)
}
