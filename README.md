# HTTP Logger

**HTTP Logger** is a lightweight HTTP server written in Go. It listens on a given port, logs all incoming HTTP requests, and optionally mirrors the request back in the response. It's ideal for debugging webhooks, clients, and API integrations.

## 📦 Install

Make sure you have [Go installed](https://golang.org/dl/), then:

```bash
go install https://github.com/sebps/http-logger@latest
```

## 🚀 Features

- Supports **all HTTP methods** and **any path**
- Logs request method, path, headers, body, and query parameters
- Optional **mirror mode** to return the request details in the response
- Clean CLI with `--help`, `--port`, and `--mirror` flags

## 🛠️ Usage

```bash
http-logger [options]
```

### Options

| Flag            | Description                                          | Default |
|-----------------|------------------------------------------------------|---------|
| `-p`, `--port`  | Port to run the server on                            | `8080`  |
| `-m`, `--mirror`| Enable mirror mode to echo back request information  | `false` |
| `-h`, `--help`  | Show help message and exit                           |         |

## 🧪 Examples

Start on default port (8080):

```bash
./http-logger
```

Start on custom port:

```bash
./http-logger --port 3000
```

Enable mirror mode:

```bash
./http-logger --mirror
```

Mirror mode + custom port:

```bash
./http-logger -p 5000 -m
```

## 🔁 Example Mirror Mode Response

```json
{
  "method": "POST",
  "path": "/submit",
  "body": "{\"name\": \"Bob\"}",
  "headers": "Content-Type: application/json",
  "params": "debug=true"
}
```