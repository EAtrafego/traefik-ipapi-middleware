package main

import (
    "context"
    "encoding/json"
    "net/http"
    "github.com/traefik/traefik/v3/pkg/middlewares"
)

// Config define as configurações do plugin
type Config struct {
    HeaderPrefix string `json:"headerPrefix"` // Prefixo para os cabeçalhos (ex.: "Geo-")
}

// CreateConfig inicializa a configuração
func CreateConfig() *Config {
    return &Config{
        HeaderPrefix: "Geo-", // Valor padrão
    }
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
        // Se a API falhar, prossegue sem adicionar cabeçalhos
        m.next.ServeHTTP(rw, req)
        return
    }
    defer resp.Body.Close()

    // Decodifica a resposta
    var result struct {
        Status      string `json:"status"`
        Country     string `json:"country"`
        CountryCode string `json:"countryCode"`
        Region      string `json:"region"`
        RegionName  string `json:"regionName"`
        City        string `json:"city"`
        Zip         string `json:"zip"`
        Lat         float64 `json:"lat"`
        Lon         float64 `json:"lon"`
        ISP         string `json:"isp"`
    }
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil || result.Status != "success" {
        // Se a resposta for inválida, prossegue sem adicionar cabeçalhos
        m.next.ServeHTTP(rw, req)
        return
    }

    // Adiciona cabeçalhos com informações de geolocalização
    req.Header.Set(m.config.HeaderPrefix+"Country", result.Country)
    req.Header.Set(m.config.HeaderPrefix+"Country-Code", result.CountryCode)
    req.Header.Set(m.config.HeaderPrefix+"Region", result.RegionName)
    req.Header.Set(m.config.HeaderPrefix+"City", result.City)
    req.Header.Set(m.config.HeaderPrefix+"Zip", result.Zip)
    req.Header.Set(m.config.HeaderPrefix+"Latitude", string(result.Lat))
    req.Header.Set(m.config.HeaderPrefix+"Longitude", string(result.Lon))
    req.Header.Set(m.config.HeaderPrefix+"ISP", result.ISP)

    // Prossegue com a requisição
    m.next.ServeHTTP(rw, req)
}