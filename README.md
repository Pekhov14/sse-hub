# SSE-HUB — A lightweight pub/sub hub for Server-Sent Events

![sse-hub](https://github.com/user-attachments/assets/5cbb229a-3518-4064-8218-2b5c4d0bad2e)

Send and receive real-time updates without writing any SSE server-side logic — just POST to `sse-hub` to publish, and connect to it to subscribe.

## Quick Start

Run sse-hub:
```bash
./sse-hub
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
PORT=3000 CLIENT_BUFFER=50 ./sse-hub
```

## 🎯 Live Demo

Open `docs/examples/frontend/demo.html` in your browser for an interactive demo showcasing:

- **📊 System Monitoring** - Real-time server metrics (CPU, memory, users)
- **🔔 Notifications** - Push alerts and system notifications  
- **💬 Live Chat** - Real-time messaging
- **🔗 Connection Test** - Basic SSE functionality testing

### Quick Demo Setup

1. **Start the server:**
   ```bash
   PORT=8090 go run cmd/server/main.go
   ```

2. **Open the demo:**
   ```bash
   open docs/examples/frontend/demo.html
   ```

3. **Click "Connect to SSE" and test the features!**

4. **Generate test data (optional):**
   ```bash
   php docs/examples/backend/php/demo_publisher.php
   ```

### Simple Frontend Example

```html
<!DOCTYPE html>
<html>
<head>
    <title>SSE Example</title>
</head>
<body>
    <div id="messages"></div>
    
    <script>
        const eventSource = new EventSource('http://localhost:8090/events');
        
        eventSource.onmessage = function(event) {
            const data = JSON.parse(event.data);
            document.getElementById('messages').innerHTML += 
                '<div>' + data.message + '</div>';
        };
    </script>
</body>
</html>
```

How to use in backend (example in php)

```php
<?php declare(strict_types=1);

set_time_limit(0);


function publishSseMessage(int $index): void
{
    $data = [
        'message' => 'SSE message ' . $index,
        'mem'     => memory_get_usage(),
        'cpu'     => sys_getloadavg()[0],
    ];

    $ch = curl_init('http://127.0.0.1:8080/publish');

    curl_setopt_array($ch, [
        CURLOPT_POST           => true,
        CURLOPT_RETURNTRANSFER => true,
        CURLOPT_HTTPHEADER     => ['Content-Type: application/json'],
        CURLOPT_POSTFIELDS     => json_encode($data, JSON_THROW_ON_ERROR),
    ]);

    $response = curl_exec($ch);

    if ($response === false) {
        fwrite(STDERR, 'cURL error: ' . curl_error($ch) . PHP_EOL);
    }

    curl_close($ch);
}

function run(): void
{
    $iterations = 1_000_000;

    for ($i = 0; $i < $iterations; $i++) {
        sleep(1);
        publishSseMessage($i);
    }
}

run();

```