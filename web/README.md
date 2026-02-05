# Stock Hub Web - React Application

This is the React frontend for the Stock Hub WebSocket communication application.

## Development

1. Install dependencies:
```bash
cd web
npm install
```

2. Start development server:
```bash
npm run dev
```

The dev server will run on http://localhost:3000 and proxy WebSocket connections to the Go server on port 8082.

## Build

Build the React app for production:
```bash
npm run build
```

This will output the built files to `../static/` directory which the Go server serves.
