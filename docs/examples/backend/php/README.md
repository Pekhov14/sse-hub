# PHP Examples for SSE Hub

This directory contains a minimal PHP pair for a complete notification flow.

## Files

### `publisher.php`
Basic notification publisher for the `/publish` endpoint.

### `read.php`
Small endpoint that accepts read acknowledgements from the frontend.

## Usage

### Basic Publisher
```bash
php publisher.php
```

### Read API
```bash
php -S localhost:8081 -t docs/examples/backend/php
```

### Monitor Results
Open `docs/examples/frontend/demo.html` in your browser to see real-time updates.

## Notification Example

```php
$publisher = new SSEPublisher('http://localhost:8080');
$publisher->publish([
    'id' => 'notif-1001',
    'type' => 'notification',
    'title' => 'Order created',
    'message' => 'Order #12345 is ready for processing',
    'level' => 'info',
    'read' => false,
    'timestamp' => time(),
]);
```

## Read Receipt Example

```php
POST /read.php
{
    "id": "notif-1001",
    "read": true,
    "read_at": "2026-04-12T12:00:00Z"
}
```

## Error Handling

```php
function publishWithRetry($data, $maxRetries = 3) {
    $attempt = 0;
    
    while ($attempt < $maxRetries) {
        try {
            $publisher = new SSEPublisher('http://localhost:8080');
            $publisher->publish($data);
            return true;
        } catch (Exception $e) {
            $attempt++;
            if ($attempt >= $maxRetries) {
                error_log("Failed to publish SSE message after {$maxRetries} attempts: " . $e->getMessage());
                return false;
            }
            sleep(1); // Wait before retry
        }
    }
    
    return false;
}
```
