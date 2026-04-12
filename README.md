# SSE-HUB — A lightweight pub/sub hub for Server-Sent Events

![sse-hub](https://github.com/user-attachments/assets/5cbb229a-3518-4064-8218-2b5c4d0bad2e)

Send and receive real-time updates without writing any SSE server-side logic — just POST to `sse-hub` to publish, and connect to it to subscribe.

## Quick Start

Build the binary:
```bash
make build
```

Show available commands:
```bash
./sse-hub
```

Start the server:
```bash
./sse-hub serve
```

Generate an API key placeholder for future `X-API-Key` support:
```bash
./sse-hub keygen
```

## Configuration

Configure via environment variables:

- `PORT` - Server port (default: 8080)
- `CLIENT_BUFFER` - Message buffer size per client (default: 10)
- `CORS_ORIGIN` - CORS origin (default: *)

**CLIENT_BUFFER explained:**
Each connected client has a message buffer. If a client is slow (poor connection, browser lag), messages queue in this buffer. When full, the server drops messages or disconnects the client. Larger buffer = more tolerance for slow clients but uses more memory.

Example:
```bash
PORT=3000 CLIENT_BUFFER=50 ./sse-hub serve
```

## API Summary

- `GET /events` opens an SSE stream
- `POST /publish` broadcasts a JSON payload to all connected clients

Current behavior:
- messages are broadcast to all connected clients
- the server keeps no event history
- there are no topics yet
- `X-API-Key` authentication is not implemented yet

See [docs/api.md](docs/api.md) for the current API description.

## Live Demo

Open `docs/examples/frontend/demo.html` in your browser for a small demo focused on notifications and read acknowledgements.

### Demo Setup

1. Build the binary:
   ```bash
   make build
   ```

2. Start the server:
   ```bash
   ./sse-hub serve
   ```

3. Start the read API example:
   ```bash
   php -S localhost:8081 -t docs/examples/backend/php
   ```

4. Open `docs/examples/frontend/demo.html` in your browser.

5. Publish a sample notification from PHP if needed:
   ```bash
   php docs/examples/backend/php/publisher.php
   ```

### Simple Frontend Example

```html
<!DOCTYPE html>
<html>
<head>
    <title>SSE Example</title>
</head>
<body>
    <div id="notifications"></div>
    
    <script>
        const eventSource = new EventSource('http://localhost:8080/events');
        
        eventSource.onmessage = function(event) {
            const data = JSON.parse(event.data);
            if (data.type === 'notification') {
                document.getElementById('notifications').innerHTML +=
                    '<div>' + data.title + ': ' + data.message + '</div>';
            }
        };
    </script>
</body>
</html>
```

How to publish a notification from backend (example in PHP)

```php
<?php

$payload = [
    'id' => 'notif-1001',
    'type' => 'notification',
    'title' => 'Order created',
    'message' => 'Order #12345 is ready for processing',
    'level' => 'info',
    'read' => false,
    'timestamp' => time(),
];

$ch = curl_init('http://127.0.0.1:8080/publish');

curl_setopt_array($ch, [
    CURLOPT_POST => true,
    CURLOPT_RETURNTRANSFER => true,
    CURLOPT_HTTPHEADER => ['Content-Type: application/json'],
    CURLOPT_POSTFIELDS => json_encode($payload, JSON_THROW_ON_ERROR),
]);

$response = curl_exec($ch);

if ($response === false) {
    fwrite(STDERR, 'cURL error: ' . curl_error($ch) . PHP_EOL);
}

curl_close($ch);

```

## Project Layout

- `cmd/` — CLI entrypoint and commands
- `internal/config/` — environment configuration
- `internal/sse/` — SSE server and HTTP handlers
- `docs/` — API notes, examples, and release planning
