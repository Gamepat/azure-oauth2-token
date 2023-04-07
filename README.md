# Azure OAuth2 Token

Go module for retrieving bearer access-token for accessing azure-api

Reference: https://learn.microsoft.com/en-us/azure/active-directory/develop/v2-oauth2-auth-code-flow#request-an-access-token-with-a-client_secret

## Usage

### Require

Add the module to your `go.mod`:

```
require github.com/gamepat/azure-oauth2-token v0.1.0
```

### Examples

Get access-token:

```go
import (
    "log"
    auth "github.com/gamepat/azure-oauth2-token"
)

func main() {
    cfg := auth.AuthConfig{
        ClientID:     "<YOUR_CLIENT_ID>",
        ClientSecret: "<YOUR_CLIENT_SECRET>",
        ClientScope:  "<YOUR_CLIENT_SCOPE>",
    }

    token, err := auth.RequestAccessToken("<YOUR_AZ_TENANT_ID>", cfg)
    if err != nil {
        log.Fatalf("error: %v", err)
    }
    log.Println(token)
}
```

Get access-token infos:

```go
import (
    "log"
    auth "github.com/gamepat/azure-oauth2-token"
)

func main() {
    cfg := auth.AuthConfig{
        ClientID:     "<YOUR_CLIENT_ID>",
        ClientSecret: "<YOUR_CLIENT_SECRET>",
        ClientScope:  "<YOUR_CLIENT_SCOPE>",
    }

    tokenInfo, err := auth.RequestAccessTokenInfo("<YOUR_AZ_TENANT_ID>", cfg)
    if err != nil {
        log.Fatalf("error: %v", err)
    }
    log.Println(tokenInfo)
}
```
