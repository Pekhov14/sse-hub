<?php

declare(strict_types=1);

header('Content-Type: application/json');
header('Access-Control-Allow-Origin: *');
header('Access-Control-Allow-Methods: POST, OPTIONS');
header('Access-Control-Allow-Headers: Content-Type');

if ($_SERVER['REQUEST_METHOD'] === 'OPTIONS') {
    http_response_code(204);
    exit;
}

if ($_SERVER['REQUEST_METHOD'] !== 'POST') {
    http_response_code(405);
    echo json_encode(['error' => 'Use POST']);
    exit;
}

$raw = file_get_contents('php://input');
$payload = json_decode($raw, true);

if (!is_array($payload) || empty($payload['id'])) {
    http_response_code(400);
    echo json_encode(['error' => 'Field "id" is required']);
    exit;
}

$response = [
    'ok' => true,
    'id' => (string) $payload['id'],
    'read' => true,
    'read_at' => gmdate(DATE_ATOM),
];

echo json_encode($response, JSON_PRETTY_PRINT);
