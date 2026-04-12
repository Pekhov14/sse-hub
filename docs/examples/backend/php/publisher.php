<?php

declare(strict_types=1);

/**
 * Notification publisher example for SSE Hub.
 */
class SSEPublisher {
    private $serverUrl;

    public function __construct($serverUrl) {
        $this->serverUrl = rtrim($serverUrl, '/');
    }

    public function publish($data) {
        $url = $this->serverUrl . '/publish';

        $jsonData = json_encode($data);
        if ($jsonData === false) {
            throw new Exception('Failed to encode notification payload');
        }

        $headers = [
            'Content-Type: application/json',
        ];

        $ch = curl_init();
        curl_setopt_array($ch, [
            CURLOPT_URL => $url,
            CURLOPT_POST => true,
            CURLOPT_POSTFIELDS => $jsonData,
            CURLOPT_HTTPHEADER => $headers,
            CURLOPT_RETURNTRANSFER => true,
            CURLOPT_TIMEOUT => 10,
        ]);

        $response = curl_exec($ch);
        $httpCode = curl_getinfo($ch, CURLINFO_HTTP_CODE);
        $error = curl_error($ch);
        curl_close($ch);

        if ($error) {
            throw new Exception("cURL error: " . $error);
        }

        if ($httpCode !== 204) {
            throw new Exception("HTTP error: " . $httpCode . " - " . $response);
        }

        return true;
    }
}

try {
    $sseServerUrl = 'http://localhost:8080';

    $publisher = new SSEPublisher($sseServerUrl);
    $notificationId = 'notif-' . date('YmdHis');

    $publisher->publish([
        'id' => $notificationId,
        'type' => 'notification',
        'title' => 'Order created',
        'message' => 'Order #12345 is ready for processing',
        'level' => 'info',
        'timestamp' => time(),
    ]);

    echo "Notification sent: {$notificationId}\n";
} catch (Exception $e) {
    echo "Error: " . $e->getMessage() . "\n";
}
