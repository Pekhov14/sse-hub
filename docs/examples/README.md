# SSE Server Examples

This directory contains practical examples of how to use the SSE server in different scenarios.

## Frontend Examples

### HTML/JavaScript Client (`frontend/`)
A complete web interface for testing the SSE server:
- Real-time connection to SSE server
- Message publishing interface
- Authentication support
- Message history display

**Usage:**
1. Start the SSE server: `go run cmd/server/main.go`
2. Open `frontend/index.html` in your browser
3. Connect to `http://localhost:8090`
4. Send and receive messages in real-time

## Backend Examples

### PHP Publisher (`backend/php/`)
Shows how to send messages from PHP applications:
- Simple message publishing
- Error handling
- Web form integration
- Multiple message types

**Usage:**
```bash
php backend/php/publisher.php
```

### Node.js Publisher (`backend/nodejs/`)
Demonstrates Node.js integration:
- Promise-based API
- Express.js middleware
- Periodic updates
- Real-time metrics

**Usage:**
```bash
node backend/nodejs/publisher.js
```

## Common Use Cases

### 1. Real-time Notifications
```javascript
// Send notification
await publisher.publish({
    type: 'notification',
    title: 'New Message',
    message: 'You have a new message from John',
    user_id: 123
});
```

### 2. Live Updates
```javascript
// Send data update
await publisher.publish({
    type: 'data_update',
    table: 'orders',
    action: 'created',
    record_id: 456
});
```

### 3. System Monitoring
```javascript
// Send system status
await publisher.publish({
    type: 'system_status',
    service: 'database',
    status: 'healthy',
    response_time: 45
});
```

### 4. Chat Messages
```javascript
// Send chat message
await publisher.publish({
    type: 'chat',
    room: 'general',
    user: 'alice',
    message: 'Hello everyone!'
});
```

## Authentication

If your SSE server uses API key authentication:

**Environment variable:**
```bash
API_KEY=your-secret-key go run cmd/server/main.go
```

**Client connection:**
```javascript
const eventSource = new EventSource('/events?key=your-secret-key');
```

**Publishing with key:**
```javascript
fetch('/publish', {
    method: 'POST',
    headers: {
        'X-API-Key': 'your-secret-key',
        'Content-Type': 'application/json'
    },
    body: JSON.stringify(data)
});
```

## Testing

### Quick Test
1. Start the SSE server: `PORT=8090 go run cmd/server/main.go`
2. Open `frontend/monitor.html` in your browser
3. Run stress test: `php backend/php/stress_test.php`
4. Watch real-time system metrics in the browser

### Full Test Suite
1. Start the SSE server
2. Open the frontend example in your browser
3. Run one of the backend examples
4. Watch messages appear in real-time in the browser

## Integration Tips

- **Error Handling**: Always wrap publish calls in try-catch blocks
- **Reconnection**: Implement client-side reconnection logic for production
- **Message Format**: Use consistent message structure across your application
- **Performance**: Consider message frequency and client count for scaling
- **Security**: Use API keys in production environments