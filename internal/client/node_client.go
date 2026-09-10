package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type NodeClient struct {
	baseURL    string
	jwtSecret  []byte
	httpClient *http.Client
}

func NewNodeClient(baseURL, jwtSecret string) *NodeClient {
	if baseURL == "" {
		baseURL = "http://node-api:4000"
	}
	return &NodeClient{
		baseURL:   baseURL,
		jwtSecret: []byte(jwtSecret),
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (n *NodeClient) generateToken() (string, error) {
	claims := jwt.MapClaims{
		"iss": "go-api",
		"aud": "node-api",
		"exp": time.Now().Add(1 * time.Minute).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(n.jwtSecret)
}

func (n *NodeClient) SendStats(payload any) (json.RawMessage, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("error serializando payload: %w", err)
	}

	token, err := n.generateToken()
	if err != nil {
		return nil, fmt.Errorf("error generando token JWT: %w", err)
	}

	req, err := http.NewRequestWithContext(
		context.Background(),
		http.MethodPost,
		n.baseURL+"/api/stats",
		bytes.NewReader(body),
	)
	if err != nil {
		return nil, fmt.Errorf("error creando petición http: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := n.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error conectando a node-api: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("node-api devolvió estado %d", resp.StatusCode)
	}

	var raw json.RawMessage
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, fmt.Errorf("error decodificando respuesta: %w", err)
	}

	return raw, nil
}