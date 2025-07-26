# sse-hub — A lightweight pub/sub hub for Server-Sent Events

Send and receive real-time updates without writing any SSE server-side logic — just POST to `sse-hub` to publish, and connect to it to subscribe.

 How to use in frontend

```html
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>SSE</title>
</head>
<body>
    Mem: <span id="mem"></span>
    <br>
    CPU: <span id="cpu"></span>

    <script>
        const eventSrc = new EventSource("http://127.0.0.1:8080/events");

        const mem = document.getElementById("mem");
        const cpu = document.getElementById("cpu");

        eventSrc.addEventListener("mem", (event) => {
            mem.textContent = event.data;
        });

        eventSrc.addEventListener("cpu", (event) => {
            cpu.textContent = event.data;
        });

        eventSrc.onerror = (err) => {
            console.log("sse error", err);
        };
    </script>
</body>
</html>
```