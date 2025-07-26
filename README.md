### How to use in frontend

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