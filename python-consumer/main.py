from connection import get_connection, get_channel, get_redis
from config import REDIS_URL
from consumer import start_consuming

def main():
    conn = get_connection()
    ch = get_channel(conn)
    redis = get_redis()
    try:
        start_consuming(ch, redis)
    finally:
        try:
            conn.close()
        except Exception:
            pass

if __name__ == "__main__":
    main()