# SSE-HUB — A lightweight pub/sub hub for Server-Sent Events

![sse-hub](https://github.com/user-attachments/assets/5cbb229a-3518-4064-8218-2b5c4d0bad2e)


Send and receive real-time updates without writing any SSE server-side logic — just POST to `sse-hub` to publish, and connect to it to subscribe.

Run sse-hub
`./sse-hub`

How to use in frontend

```html
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0"/>
    <title>SSE Monitor</title>
    <style>
        body {
            font-family: sans-serif;
            padding: 1rem;
        }
        .label {
            font-weight: bold;
        }
        .value {
            margin-left: 0.5rem;
        }
    </style>
</head>
<body>
    <div>
        <span class="label">Message:</span>
        <span id="message" class="value">-</span>
    </div>
    <div>
        <span class="label">Memory:</span>
        <span id="mem" class="value">-</span>
    </div>
    <div>
        <span class="label">CPU:</span>
        <span id="cpu" class="value">-</span>
    </div>

    <script>
        (() => {
            const eventSourceUrl = "http://127.0.0.1:8080/events";

            const elements = {
                message: document.getElementById("message"),
                mem: document.getElementById("mem"),
                cpu: document.getElementById("cpu"),
            };

            const updateUI = (data) => {
                elements.message.textContent = data.message ?? "-";
                elements.mem.textContent = data.mem ?? "-";
                elements.cpu.textContent = data.cpu ?? "-";
            };

            const eventSource = new EventSource(eventSourceUrl);

            eventSource.onmessage = (event) => {
                try {
                    const data = JSON.parse(event.data);
                    updateUI(data);
                } catch (error) {
                    console.error("Invalid JSON in SSE message:", error);
                }
            };

            eventSource.onerror = (error) => {
                console.error("SSE connection error:", error);
            };
        })();
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