# Stock Hub

A Go application with WebSocket support for real-time stock data.

## Project Structure

```
stock_hub/
├── cmd/
│   └── server/          # Application entry point
├── internal/
│   ├── handler/         # HTTP handlers
│   ├── service/         # Business logic
│   ├── model/           # Data models
│   └── config/          # Configuration
├── pkg/
│   ├── websocket/       # WebSocket utilities
│   └── utils/           # Utility functions
├── api/                 # API definitions
├── web/
│   ├── static/          # Static files (CSS, JS)
│   └── templates/       # HTML templates
├── migrations/          # Database migrations
├── scripts/             # Build/deploy scripts
└── docs/                # Documentation
```

## Dependencies

- **gorilla/websocket**: WebSocket implementation

## Getting Started

1. Install dependencies:
```bash
go mod download
```

2. Run the server:
```bash
go run cmd/server/main.go
```

3. Connect to WebSocket:
```
ws://localhost:8080/ws
```

## Development

```bash
# Run server
go run cmd/server/main.go

# Build
go build -o bin/server cmd/server/main.go

# Run tests
go test ./...
```

![Screenshot 2026-02-05 at 20.16.03.png](../../../../var/folders/q5/x2zsxtkx22ddmcpn8f3xn4y40000gp/T/TemporaryItems/NSIRD_screencaptureui_x2yToe/Screenshot%202026-02-05%20at%2020.16.03.png)

![Screenshot 2026-02-05 at 20.16.38.png](../../../../var/folders/q5/x2zsxtkx22ddmcpn8f3xn4y40000gp/T/TemporaryItems/NSIRD_screencaptureui_uePrnw/Screenshot%202026-02-05%20at%2020.16.38.png)