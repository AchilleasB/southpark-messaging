# streamlit-web-app

Purpose
- Simple Streamlit UI to send messages to the Go API and display messages persisted in Redis.

Structure
- `app.py` — Streamlit app:
  - Page config (must be the first Streamlit call).
  - Creates Redis client via `st.cache_resource`.
  - Sidebar: input form that POSTs to Go API (`/messages`).
  - Main area: reads the Redis list `southpark_messages` and displays messages; includes refresh / auto-refresh behavior.

Run locally
```bash
python -m venv .venv
.venv\Scripts\Activate.ps1   # PowerShell on Windows
pip install -r requirements.txt
streamlit run app.py
# open http://localhost:8501
```

Docker / Compose
- The app is containerized in `streamlit-web-app/Dockerfile`. Compose maps port `8501:8501` and injects `GO_API_URL` and `REDIS_URL`.

Troubleshooting
  - Verify `REDIS_URL` and `REDIS_KEY` inside container: `docker-compose exec streamlit-web-app env`
  - Confirm Redis has entries: `docker exec -it redis redis-cli LRANGE southpark_messages 0 -1`
  - Confirm app uses `st.cache_resource` for the Redis client and clears cache on refresh when needed.