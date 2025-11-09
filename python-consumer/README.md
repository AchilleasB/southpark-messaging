# python-consumer

Purpose
- Consumes messages from RabbitMQ queue `southpark_messages`.
- Prints messages to console and persists JSON messages into Redis list key (default `southpark_messages`).

Structure
- `main.py` — entrypoint: creates RabbitMQ connection and Redis client, starts consumer.
- `connection.py` — RabbitMQ connection / channel helpers.
- `redis_client.py` — builds Redis client from URL.
- `worker.py` — consumer callback, message handling, persistence.
- `config.py` — environment defaults (RABBITMQ_URL, REDIS_URL, REDIS_KEY).

Run locally
```bash
python -m venv .venv
.venv\Scripts\Activate.ps1   # PowerShell on Windows
pip install -r requirements.txt
python main.py
```

Docker
- Image built via `python-consumer/Dockerfile`.
- Compose sets envs; ensure `REDIS_URL` and `RABBITMQ_URL` are provided.