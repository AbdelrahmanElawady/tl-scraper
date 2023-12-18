package endpoints

import (
	"net/url"
)

type Caller struct {
	key       string
	email     string
	accountID string
	baseURL   string
}

const defaultBaseURL = "https://api.testlodge.com/v1/account/"

func NewCaller(key, email string, accountID string) (*Caller, error) {
	baseURL, err := url.JoinPath(defaultBaseURL, accountID)
	if err != nil {
		return nil, err
	}
	return &Caller{
		key:       key,
		email:     email,
		accountID: accountID,
		baseURL:   baseURL,
	}, nil
}
