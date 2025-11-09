import os
import json
import time
import requests
import redis
import streamlit as st

# ----------------------------
# IMPORTANT: Page config must be the FIRST Streamlit command
# ----------------------------
st.set_page_config(page_title="South Park Messaging", layout="centered")

# ----------------------------
# Imports & Global Connections
# ----------------------------
GO_API_URL = os.getenv("GO_API_URL", "http://localhost:8081")
REDIS_URL = os.getenv("REDIS_URL", "redis://localhost:6379/0")
REDIS_KEY = os.getenv("REDIS_KEY", "southpark_messages")

@st.cache_resource
def get_redis_client():
    return redis.Redis.from_url(REDIS_URL, decode_responses=True)

r = get_redis_client()

# ----------------------------
# Sending Function
# ----------------------------
def send_message(author: str, body: str) -> tuple[int, str]:
    """
    POST message to Go API /messages.
    Returns (status_code, text) or (0, error_msg) on failure.
    """
    url = f"{GO_API_URL.rstrip('/')}/messages"
    payload = {"author": author, "body": body}
    headers = {"Content-Type": "application/json"}
    try:
        resp = requests.post(url, json=payload, headers=headers, timeout=5)
        return resp.status_code, resp.text
    except Exception as e:
        return 0, str(e)

# ----------------------------
# Helper: Read messages from Redis
# ----------------------------
@st.cache_data(ttl=2)
def fetch_messages_from_redis() -> list[dict]:
    """
    Read all messages from Redis list (right-pushed JSON strings).
    Returns list of parsed message dicts in chronological order (oldest first).
    """
    try:
        raw = r.lrange(REDIS_KEY, 0, -1) or []
    except Exception:
        return []
    msgs = []
    for item in raw:
        try:
            msgs.append(json.loads(item))
        except Exception:
            # store raw fallback
            msgs.append({"author": "<unknown>", "body": str(item)})
    return msgs

# ----------------------------
# Main UI Logic
# ----------------------------
st.title("South Park Messaging")

# Sidebar: input form and send action
with st.sidebar.form("send_form", clear_on_submit=True):
    st.header("Send a message")
    author = st.text_input("Author", value="Cartman")
    body = st.text_area("Message", height=100)
    submit = st.form_submit_button("Send")

    if submit:
        if not author.strip() or not body.strip():
            st.error("Both author and message body are required.")
        else:
            code, txt = send_message(author.strip(), body.strip())
            if code == 202 or (100 <= code < 300):
                st.success("Message sent.")
            elif code == 0:
                st.error(f"Failed to send: {txt}")
            else:
                st.error(f"Send failed ({code}): {txt}")

# Main area: display messages and refresh control
st.subheader("Messages")

col1, col2 = st.columns([1, 3])
with col1:
    refresh = st.button("Refresh")
    auto_refresh = st.checkbox("Auto-refresh (every 2s)", value=True)

with col2:
    st.write("Messages stored in Redis (most recent last)")

# Fetch and display
if refresh:
    # clear cache once to force immediate fetch
    fetch_messages_from_redis.clear()

messages = fetch_messages_from_redis()

if not messages:
    st.info("No messages found.")
else:
    # show messages in reverse chronological order (newest first)
    for msg in reversed(messages):
        author = msg.get("author", "<unknown>")
        body = msg.get("body", "")
        st.markdown(f"**{st.session_state.get('last_author', author)}**" if False else f"**{author}**")
        st.write(body)
        st.markdown("---")

# Auto-refresh logic: simple rerun when enabled
if auto_refresh:
    time.sleep(2)
    # invalidate cache and rerun to update UI
    fetch_messages_from_redis.clear()
    st.experimental_rerun()