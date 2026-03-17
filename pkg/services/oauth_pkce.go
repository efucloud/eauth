package services

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"strings"
)

func normalizePKCEParams(codeChallenge, codeChallengeMethod string) (normalizedChallenge, normalizedMethod string, err error) {
	codeChallenge = strings.TrimSpace(codeChallenge)
	codeChallengeMethod = strings.TrimSpace(codeChallengeMethod)
	if len(codeChallenge) == 0 {
		if len(codeChallengeMethod) > 0 {
			return "", "", fmt.Errorf("code_challenge is required when code_challenge_method is provided")
		}
		return "", "", nil
	}
	if len(codeChallengeMethod) == 0 {
		return codeChallenge, "plain", nil
	}
	if strings.EqualFold(codeChallengeMethod, "plain") {
		return codeChallenge, "plain", nil
	}
	if strings.EqualFold(codeChallengeMethod, "S256") {
		return codeChallenge, "S256", nil
	}
	return "", "", fmt.Errorf("unsupported code_challenge_method: %s", codeChallengeMethod)
}

func verifyPKCE(codeChallenge, codeChallengeMethod, codeVerifier string) error {
	if len(codeChallenge) == 0 {
		return nil
	}
	if len(codeVerifier) == 0 {
		return fmt.Errorf("code_verifier is required for pkce client")
	}
	if strings.EqualFold(codeChallengeMethod, "plain") {
		if codeVerifier != codeChallenge {
			return fmt.Errorf("pkce verification failed")
		}
		return nil
	}
	if strings.EqualFold(codeChallengeMethod, "S256") {
		hash := sha256.Sum256([]byte(codeVerifier))
		encoded := base64.RawURLEncoding.EncodeToString(hash[:])
		if encoded != codeChallenge {
			return fmt.Errorf("pkce verification failed")
		}
		return nil
	}
	return fmt.Errorf("unsupported code_challenge_method: %s", codeChallengeMethod)
}
