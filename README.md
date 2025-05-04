# Traefik IP-API Middleware Plugin

Este plugin para o Traefik utiliza a API do [ip-api.com](https://ip-api.com/) para obter informações de geolocalização com base no IP do cliente e bloquear requisições de países especificados.

## Funcionalidades
- Consulta a geolocalização do IP via `ip-api.com`.
- Bloqueia requisições com base em códigos de países configurados.
- Suporta o cabeçalho `X-Forwarded-For` para IPs reais em proxies.

## Requisitos
- Traefik v3.0 ou superior.
- Go 1.21 ou superior (para desenvolvimento).
- Conexão com a internet para acessar `ip-api.com`.

## Instalação

### Configuração Estática
Adicione o plugin ao `traefik.yml`:

```yaml
experimental:
  plugins:
    ipapi:
      moduleName: github.com/seu-usuario/traefik-ipapi-middleware
      version: v0.1.0