package authorization_checker

import "crypto/subtle"

type AuthorizationChecker struct {
	validCredentials map[string]string
}

func New(username, password string) *AuthorizationChecker {
	return &AuthorizationChecker{validCredentials: map[string]string{username: password}}
}

func (a *AuthorizationChecker) AreCredentialsValid(username, password string) bool {
	for validUsername, validPassword := range a.validCredentials {
		usernameMatch := (subtle.ConstantTimeCompare([]byte(username), []byte(validUsername)) == 1)
		passwordMatch := (subtle.ConstantTimeCompare([]byte(password), []byte(validPassword)) == 1)
		if usernameMatch && passwordMatch {
			return true
		}
	}
	return false
}
