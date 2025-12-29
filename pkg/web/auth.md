# Web Server Security & TLS

This document details how to secure the gRouter Web Server using TLS (Transport Layer Security).

## TLS Configuration

To enable HTTPS, configure the `tls` section in your `config.yaml`.

**Configuration**:
```yaml
web:
  port: 8443 # Common port for HTTPS
  tls:
    enabled: true
    cert_file: "./certs/server.pem"
    key_file: "./certs/server.key"
```

## Certificate Generation

Use `mkcert` to generate locally trusted development certificates.

1.  **Install mkcert**:
    ```bash
    brew install mkcert # or equivalent
    mkcert -install
    ```
2.  **Generate Certs**:
    ```bash
    mkdir -p certs
    cd certs
    mkcert -key-file server.key -cert-file server.pem localhost 127.0.0.1 ::1
    ```

## Server Setup

Ensure your application loads the configuration with TLS enabled.

**Running the Server**:
```bash
# Assuming you are running the natsdemosvc service or gRouter
go run services/natsdemosvc/cmd/natsdemosvc/main.go --config configs/config.yaml
```

*The logs should show:*
```text
INFO    Starting web server    {"port": 8443, "tls": true}
```

## Client Setup

### 1. cURL
To test the endpoint, use `curl`. If you used `mkcert -install`, your system might already trust the cert. If not, or for explicit verification:

```bash
curl -v https://localhost:8443/health/live
```

If using self-signed certs without `mkcert` trust, use `-k` (insecure):
```bash
curl -k -v https://localhost:8443/health/live
```

### 2. Go Client
Here is how to configure a Go `http.Client` to accept the custom CA or skip verification (for dev).

```go
package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"net/http"
	"os"
)

func main() {
	// A. For Dev/Self-Signed (Insecure)
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	// B. For Custom CA (Secure)
	// caCert, _ := os.ReadFile("./certs/rootCA.pem")
	// caCertPool := x509.NewCertPool()
	// caCertPool.AppendCertsFromPEM(caCert)
	// tr := &http.Transport{
	// 	TLSClientConfig: &tls.Config{RootCAs: caCertPool},
	// }

	client := &http.Client{Transport: tr}

	resp, err := client.Get("https://localhost:8443/health/live")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Println(string(body))
}
```

### 3. Browser
Simply navigate to `https://localhost:8443/health/live`.
- If using `mkcert`, it should load without warnings.
- Otherwise, you will see a "Not Secure" warning which you must accept to proceed.
