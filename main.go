package main

import (
    "context"
    "encoding/json"
    "net/http"
    "github.com/traefik/traefik/v3/pkg/middlewares"
)

// Config define as configurações do plugin
type Config struct {
    BlockedCountries []string `json:"blockedCountries"`
}

// CreateConfig inicializa a configuração
func CreateConfig() *Config {
    return &Config{}
}

// IPAPIMiddleware é o struct do middleware
type IPAPIMiddleware struct {
    next   http.Handler
    name   string
    config *Config
}

// New cria uma nova instância do middleware
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
    return &IPAPIMiddleware{
        next:   next,
        name:   name,
        config: config,
    }, nil
}

// ServeHTTP processa a requisição
func (m *IPAPIMiddleware) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
    // Obtém o IP do cliente
    clientIP := req.RemoteAddr
    if forwardedFor := req.Header.Get("X-Forwarded-For"); forwardedFor != "" {
        clientIP = forwardedFor
    }

    // Consulta a API ip-api.com
    resp, err := http.Get("http://ip-api.com/json/" + clientIP)
    if err != nil {
        http.Error(rw, "Erro ao consultar ip-api.com", http.StatusInternalServerError)
        return
    }
    defer resp.Body.Close()

    // Decodifica a resposta
    var result struct {
        CountryCode string `json:"countryCode"`
        Status      string `json:"status"`
    }
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        http.Error(rw, "Erro ao processar resposta", http.StatusInternalServerError)
        return
    }

    // Verifica se o país está bloqueado
    if result.Status == "success" {
        for _, blocked := range m.config.BlockedCountries {
            if result.CountryCode == blocked {
                http.Error(rw, "Acesso bloqueado para este país", http.StatusForbidden)
                return
            }
        }
    }

    // Prossegue com a requisição
    m.next.ServeHTTP(rw, req)
}