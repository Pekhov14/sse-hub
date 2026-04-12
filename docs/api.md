# SSE Hub API Documentation

## Overview

The current SSE hub provides two endpoints for real-time communication:
- `/events` - Server-Sent Events stream for clients
- `/publish` - HTTP endpoint for publishing messages

Current implementation notes:
- all messages are broadcast to all connected clients
- topics are not implemented yet
- message history and replay are not implemented
- `X-API-Key` authentication is planned, but not implemented yet

## Endpoints

### GET /events

Establishes a Server-Sent Events connection for receiving real-time messages.

**Headers:**
- `Content-Type: text/event-stream`
- `Cache-Control: no-cache`
- `Connection: keep-alive`
- `Access-Control-Allow-Origin: *`

**Response Format:**
```
data: {"message": "Hello", "timestamp": 1234567890}

data: {"type": "notification", "title": "Alert", "message": "Something happened"}

```

**Example:**
```javascript
const eventSource = new EventSource('http://localhost:8090/events');

eventSource.onmessage = function(event) {
    const data = JSON.parse(event.data);
    console.log('Received:', data);
};

eventSource.onerror = function(event) {
    console.error('Connection error');
};
```

### POST /publish

Publishes a message to all connected clients.

**Content-Type:** `application/json`

**Request Body:**
Any valid JSON object. Common patterns:

```json
{
    "message": "Hello World",
    "timestamp": 1234567890
}
```

```json
{
    "type": "notification",
    "title": "New Order",
    "message": "Order #12345 received",
    "user_id": 123
}
```

**Response:**
- **204 No Content** - Message published successfully
- **400 Bad Request** - Invalid JSON in request body
- **405 Method Not Allowed** - Non-POST request
- **500 Internal Server Error** - Server error

**Example:**
```javascript
fetch('http://localhost:8090/publish', {
    method: 'POST',
    headers: {
        'Content-Type': 'application/json'
    },
    body: JSON.stringify({
        message: 'Hello from client',
        timestamp: Date.now()
    })
});
```

## Configuration

The server is configured via environment variables:

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `8080` | Server port |
| `CLIENT_BUFFER` | `10` | Message buffer size per client |
| `CORS_ORIGIN` | `*` | CORS allowed origin |

**Example:**
```bash
PORT=3000 CLIENT_BUFFER=20 ./sse-hub serve
```

## Message Flow

1. **Client Connection**: Browser connects to `/events` endpoint
2. **Server Registration**: Server adds client to active connections list
3. **Message Publishing**: Backend service sends POST to `/publish`
4. **Broadcasting**: Server forwards message to all connected clients
5. **Client Disconnect**: Server automatically removes disconnected clients

## Error Handling

### Client-Side Errors
- **Connection Failed**: Check server URL and network connectivity
- **Message Parse Error**: Ensure server sends valid JSON

### Server-Side Errors
- **Invalid JSON**: Request body must be valid JSON
- **Client Buffer Full**: Client removed if message buffer overflows

## Best Practices

### Client Implementation
```javascript
class SSEClient {
    constructor(url) {
        this.url = url;
        this.eventSource = null;
        this.reconnectDelay = 1000;
    }
    
    connect() {
        this.eventSource = new EventSource(`${this.url}/events`);
        
        this.eventSource.onopen = () => {
            console.log('Connected');
            this.reconnectDelay = 1000; // Reset delay
        };
        
        this.eventSource.onmessage = (event) => {
            try {
                const data = JSON.parse(event.data);
                this.onMessage(data);
            } catch (e) {
                console.error('Invalid message format:', event.data);
            }
        };
        
        this.eventSource.onerror = () => {
            console.log('Connection lost, reconnecting...');
            this.reconnect();
        };
    }
    
    reconnect() {
        setTimeout(() => {
            this.connect();
            this.reconnectDelay = Math.min(this.reconnectDelay * 2, 30000);
        }, this.reconnectDelay);
    }
    
    onMessage(data) {
        // Override this method
        console.log('Received:', data);
    }
    
    disconnect() {
        if (this.eventSource) {
            this.eventSource.close();
        }
    }
}
```

### Publisher Implementation
```javascript
class SSEPublisher {
    constructor(url) {
        this.url = url;
    }
    
    async publish(data) {
        const response = await fetch(`${this.url}/publish`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(data)
        });
        
        if (!response.ok) {
            throw new Error(`HTTP ${response.status}: ${await response.text()}`);
        }
    }
}
```

## Scaling Considerations

- **Memory Usage**: Each client uses ~10KB memory (configurable buffer)
- **Connection Limits**: OS-dependent, typically 1000+ concurrent connections
- **Message Throughput**: Depends on message size and client count
- **Load Balancing**: Multi-instance fan-out is not implemented yet
