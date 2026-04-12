# SSE Hub Examples

This directory contains one complete, canonical example of how to use the current SSE hub.

## Frontend Examples

### HTML/JavaScript Demo (`frontend/demo.html`)
A minimal web interface focused on one use case:
- connect to the SSE hub
- publish a notification with an `id`
- receive and render notifications in real time
- mark a notification as read with a normal HTTP POST

**Usage:**
1. Build the binary: `make build`
2. Start the SSE hub: `./sse-hub serve`
3. Start the read API example: `php -S localhost:8081 -t docs/examples/backend/php`
4. Open `frontend/demo.html` in your browser
5. Connect to `http://localhost:8080`
6. Publish, receive, and mark notifications as read

## Backend Examples

### PHP Publisher (`backend/php/`)
Shows the same notification flow from a backend script:
- build a notification payload
- include a stable notification `id`
- send it to `POST /publish`
- handle publish errors

### PHP Read API (`backend/php/read.php`)
Receives the read acknowledgement from the frontend:
- accepts a notification `id`
- returns a JSON confirmation
- demonstrates that user actions go back through normal HTTP, not SSE

**Usage:**
```bash
php docs/examples/backend/php/publisher.php
```

## Canonical Payload

Use this payload shape as the simplest starting point:

```javascript
await publisher.publish({
    id: 'notif-1001',
    type: 'notification',
    title: 'New order',
    message: 'Order #12345 is ready for processing',
    level: 'info',
    read: false,
    timestamp: Date.now()
});
```

## Testing

### Quick Test
1. Start the SSE hub: `./sse-hub serve`
2. Start the read API example: `php -S localhost:8081 -t docs/examples/backend/php`
3. Open `docs/examples/frontend/demo.html` in your browser
4. Run `php docs/examples/backend/php/publisher.php`
5. Watch the notification appear in the browser
6. Click `Mark as read` and watch the read acknowledgement in the log

## Integration Tips

- **Error Handling**: Always wrap publish calls in try-catch blocks
- **Reconnection**: Implement client-side reconnection logic for production
- **Message Format**: Give notifications an `id` if the frontend will act on them
- **Directionality**: Use SSE for delivery and normal HTTP for user actions such as read acknowledgements
- **Security**: Add authentication before exposing `/publish` publicly
