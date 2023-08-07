jwt parser plugin in traefik
========================================

This plugin mainly parses the jwt in the request, writes the request information in the format of key and value into the request header, and passes it to the next layer

## configuration

```yaml
JWTParser:
  TokenKey: "jwt key, for example: Authorization"
  SecretKey: "jwt secret"
  TrustKeys:
    - "user_id"
    - "user_name"
    - "only the key configuration with this option will be parsed"
```

### plugin configuration

```yaml
experimental:
  plugins:
    JWTParser:
      moduleName: github.com/qxsugar/traefik-jwt-parser
      version: v1.0.2
```

### middleware configuration

```yaml
http:
  routers:
    homeassistant:
      entryPoints:
        - web
      rule: 'Host(`www.example.com`)'
      service: xxx
      middlewares:
        - JWTParser
  middlewares:
    JWTParser:
      plugin:
        JWTParser:
          TokenKey: Authorization
          SecretKey: "secret"
          TrustKeys:
            - "user_id"
            - "user_name"
```
