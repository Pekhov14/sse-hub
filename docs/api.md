# SSE Server API Documentation

## Overview

The SSE server provides two main endpoints for real-time communication:
- `/events` - Server-Sent Events stream for clients
- `/publish` - HTTP endpoint for publishing messages

## Authentication

Authentication is optional and controlled by the `API_KEY` environment variable.

### Methods
1. **Query Parameter**: `?key=your-api-key`
2. **HTTP Header**: `X-API-Key: your-api-key`

If no API key is configured, the server runs without authentication.

## Endpoints

### GET /events

Establishes a Server-Sent Events connection for receiving real-time messages.

**Headers:**
- `Content-Type: text/event-stream`
- `Cache-Control: no-cache`
- `Connection: keep-alive`
- `Access-Control-Allow-Origin: *`

**Authentication:**
- Query parameter: `GET /events?key=your-api-key`
- Not required if API key is not configured

**Response Format:**
```
data: {"message": "Hello", "timestamp": 1234567890}

data: {"type": "notification", "title": "Alert", "message": "Something happened"}

```

**Example:**
```javascript
const eventSource = new EventSource('/events?key=your-api-key');

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

**Authentication:**
- Header: `X-API-Key: your-api-key`
- Not required if API key is not configured

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
- **401 Unauthorized** - Invalid or missing API key
- **405 Method Not Allowed** - Non-POST request
- **500 Internal Server Error** - Server error

**Example:**
```javascript
fetch('/publish', {
    method: 'POST',
    headers: {
        'Content-Type': 'application/json',
        'X-API-Key': 'your-api-key'
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
| `API_KEY` | `` | API key for authentication (optional) |

**Example:**
```bash
PORT=3000 CLIENT_BUFFER=20 API_KEY=secret123 go run cmd/server/main.go
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
- **Authentication Failed**: Verify API key is correct
- **Message Parse Error**: Ensure server sends valid JSON

### Server-Side Errors
- **Invalid JSON**: Request body must be valid JSON
- **Authentication Error**: API key missing or incorrect
- **Client Buffer Full**: Client removed if message buffer overflows

## Best Practices

### Client Implementation
```javascript
class SSEClient {
    constructor(url, apiKey) {
        this.url = url;
        this.apiKey = apiKey;
        this.eventSource = null;
        this.reconnectDelay = 1000;
    }
    
    connect() {
        const url = this.apiKey ? 
            `${this.url}/events?key=${this.apiKey}` : 
            `${this.url}/events`;
            
        this.eventSource = new EventSource(url);
        
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
    constructor(url, apiKey) {
        this.url = url;
        this.apiKey = apiKey;
    }
    
    async publish(data) {
        const headers = {
            'Content-Type': 'application/json'
        };
        
        if (this.apiKey) {
            headers['X-API-Key'] = this.apiKey;
        }
        
        const response = await fetch(`${this.url}/publish`, {
            method: 'POST',
            headers,
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
- **Load Balancing**: Use sticky sessions or shared state (Redis) for multiple instances