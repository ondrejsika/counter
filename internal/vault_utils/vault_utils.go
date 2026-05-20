package vault_utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

func GetVaultEnv() (string, string, string, error) {
	vaultAddr := os.Getenv("VAULT_ADDR")
	if vaultAddr == "" {
		vaultAddr = "http://vault.vault:8200"
	}
	vaultAuthRole := os.Getenv("VAULT_KUBERNETES_AUTH_ROLE")
	if vaultAuthRole == "" {
		return "", "", "", fmt.Errorf("VAULT_KUBERNETES_AUTH_ROLE is not set")
	}
	vaultAuthPath := os.Getenv("VAULT_KUBERNETES_AUTH_PATH")
	if vaultAuthPath == "" {
		vaultAuthPath = "kubernetes"
	}
	return vaultAddr, vaultAuthRole, vaultAuthPath, nil
}

func GetVaultTokenFromKubernetes(vaultAddr, vaultAuthPath, vaultAuthRole, vaultSecretPath string) (string, error) {
	jwtBytes, err := os.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/token")
	if err != nil {
		return "", fmt.Errorf("read kubernetes service account token: %w", err)
	}

	body, err := json.Marshal(map[string]string{
		"jwt":  string(jwtBytes),
		"role": vaultAuthRole,
	})
	if err != nil {
		return "", fmt.Errorf("marshal vault login request: %w", err)
	}

	url := fmt.Sprintf("%s/v1/auth/%s/login", vaultAddr, vaultAuthPath)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return "", fmt.Errorf("vault kubernetes login: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("vault kubernetes login status %d: %s", resp.StatusCode, string(respBody))
	}

	var result struct {
		Auth struct {
			ClientToken string `json:"client_token"`
		} `json:"auth"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("decode vault login response: %w", err)
	}

	if result.Auth.ClientToken == "" {
		return "", fmt.Errorf("vault returned empty client token")
	}

	return result.Auth.ClientToken, nil
}

func GetVaultSecret(vaultAddr, vaultToken, vaultSecretPath string) (map[string]string, error) {
	url := fmt.Sprintf("%s/v1/%s", vaultAddr, vaultSecretPath)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("create vault request: %w", err)
	}
	req.Header.Set("X-Vault-Token", vaultToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("vault secret fetch: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("vault secret fetch status %d: %s", resp.StatusCode, string(respBody))
	}

	var envelope struct {
		Data json.RawMessage `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&envelope); err != nil {
		return nil, fmt.Errorf("decode vault secret response: %w", err)
	}

	// Detect KV v2 (data.data) vs KV v1 (data is the map directly)
	var raw map[string]interface{}
	if err := json.Unmarshal(envelope.Data, &raw); err != nil {
		return nil, fmt.Errorf("parse vault secret data: %w", err)
	}

	if nested, ok := raw["data"]; ok {
		if nestedMap, ok := nested.(map[string]interface{}); ok {
			return toStringMap(nestedMap), nil
		}
	}

	return toStringMap(raw), nil
}

func LoadEnvFromVaultSecret(vaultAddr, vaultToken, vaultSecretPath string) error {
	secret, err := GetVaultSecret(vaultAddr, vaultToken, vaultSecretPath)
	if err != nil {
		return err
	}
	for k, v := range secret {
		if _, exists := os.LookupEnv(k); exists {
			continue
		}
		if err := os.Setenv(k, v); err != nil {
			return fmt.Errorf("set env %s: %w", k, err)
		}
	}
	return nil
}

func toStringMap(m map[string]interface{}) map[string]string {
	out := make(map[string]string, len(m))
	for k, v := range m {
		if str, ok := v.(string); ok {
			out[k] = str
		}
	}
	return out
}
