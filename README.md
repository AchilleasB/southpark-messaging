# South Park Messaging

Overview
- Small distributed system that lets South Park characters send/receive messages asynchronously.
- Components:
  - go-api-service —> Go HTTP API (producer) publishes JSON messages to RabbitMQ.
  - python-consumer —> Python service consumes RabbitMQ messages and persists them to Redis.
  - streamlit-web-app —> Streamlit UI reads from Redis and displays messages.
  - RabbitMQ — message broker (amqp).
  - Redis — simple persistence for the UI.

Architecture & flow
1. Client (Streamlit UI or any HTTP client) POSTs JSON to Go API `/messages`.
2. Go API validates and publishes the message to RabbitMQ queue `southpark_messages` via a Publisher port (adapter = rabbitmq).
3. Python consumer listens on `southpark_messages`, consumes messages, stores each JSON entry into a Redis list (key `southpark_messages`).
4. Streamlit UI reads from Redis list and displays messages in the browser.

Run everything
From the repo root:
```powershell

docker-compose up -d --build

docker-compose logs -f rabbitmq go-api python-consumer streamlit-web-app
```