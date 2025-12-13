# PHP Examples for SSE Server

This directory contains PHP examples for integrating with the SSE server.

## Files

### `publisher.php`
Basic publisher class with examples of different message types.

### `stress_test.php`
Stress testing script that sends continuous messages with system metrics.

## Usage

### Basic Publisher
```bash
php publisher.php
```

### Stress Test
```bash
# Start SSE server first
PORT=8090 go run cmd/server/main.go

# In another terminal, run stress test
php stress_test.php
```

### Monitor Results
Open `docs/examples/frontend/monitor.html` in your browser to see real-time updates.

## Example Integration

### Laravel Integration
```php
// In a Laravel controller
use App\Services\SSEPublisher;

class OrderController extends Controller 
{
    public function store(Request $request) 
    {
        $order = Order::create($request->validated());
        
        // Notify via SSE
        $publisher = new SSEPublisher('http://localhost:8090');
        $publisher->publish([
            'type' => 'new_order',
            'order_id' => $order->id,
            'customer' => $order->customer_name,
            'amount' => $order->total,
            'timestamp' => time()
        ]);
        
        return response()->json($order);
    }
}
```

### WordPress Plugin
```php
// In a WordPress plugin
add_action('wp_insert_post', function($post_id) {
    if (get_post_type($post_id) === 'product') {
        $publisher = new SSEPublisher('http://localhost:8090');
        $publisher->publish([
            'type' => 'new_product',
            'product_id' => $post_id,
            'title' => get_the_title($post_id),
            'timestamp' => time()
        ]);
    }
});
```

### Symfony Integration
```php
// In a Symfony service
use Symfony\Component\EventDispatcher\EventSubscriberInterface;

class SSENotificationSubscriber implements EventSubscriberInterface
{
    private SSEPublisher $publisher;
    
    public function __construct(SSEPublisher $publisher) 
    {
        $this->publisher = $publisher;
    }
    
    public function onUserRegistered(UserRegisteredEvent $event): void
    {
        $this->publisher->publish([
            'type' => 'user_registered',
            'user_id' => $event->getUser()->getId(),
            'email' => $event->getUser()->getEmail(),
            'timestamp' => time()
        ]);
    }
    
    public static function getSubscribedEvents(): array
    {
        return [
            UserRegisteredEvent::class => 'onUserRegistered',
        ];
    }
}
```

## Performance Tips

1. **Connection Pooling**: Reuse cURL handles for multiple requests
2. **Async Publishing**: Use ReactPHP or Swoole for non-blocking requests
3. **Error Handling**: Always handle cURL errors gracefully
4. **Timeouts**: Set appropriate timeouts for production use

## Error Handling

```php
function publishWithRetry($data, $maxRetries = 3) {
    $attempt = 0;
    
    while ($attempt < $maxRetries) {
        try {
            $publisher = new SSEPublisher('http://localhost:8090');
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