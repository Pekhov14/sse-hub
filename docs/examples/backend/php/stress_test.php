<?php 
declare(strict_types=1);

set_time_limit(0);

function publishSseMessage(int $index): void
{
    $data = [
        'message' => 'SSE message ' . $index,
        'mem'     => memory_get_usage(),
        'cpu'     => sys_getloadavg()[0],
    ];

    $ch = curl_init('http://127.0.0.1:8090/publish');
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
    
    echo "Starting stress test with {$iterations} iterations...\n";
    echo "Sending messages to http://127.0.0.1:8090/publish\n";
    echo "Press Ctrl+C to stop\n\n";
    
    for ($i = 0; $i < $iterations; $i++) {
        sleep(1);
        publishSseMessage($i);
        
        if ($i % 10 === 0) {
            echo "Sent {$i} messages\n";
        }
    }
}

run();