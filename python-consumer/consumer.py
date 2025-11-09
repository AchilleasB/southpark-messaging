import json
from typing import Callable
from config import QUEUE_NAME, REDIS_KEY

def default_callback_factory(redis_client):
    def callback(ch, method, properties, body):
        try:
            msg = json.loads(body)
            author = msg.get("author", "<unknown>")
            body_text = msg.get("body", "")
            print(f"{author}: {body_text}")

            # Persist raw JSON to Redis list (right push)
            try:
                redis_client.rpush(REDIS_KEY, json.dumps(msg))
            except Exception as e:
                print(f"redis persist error: {e}")
        except Exception:
            print("received non-json message")
        ch.basic_ack(delivery_tag=method.delivery_tag)
    return callback

def start_consuming(channel, redis_client, callback: Callable = None):
    if callback is None:
        callback = default_callback_factory(redis_client)
    channel.queue_declare(queue=QUEUE_NAME, durable=True)
    channel.basic_qos(prefetch_count=1)
    channel.basic_consume(queue=QUEUE_NAME, on_message_callback=callback)
    print(f"consumer waiting for messages on queue '{QUEUE_NAME}'. To exit press CTRL+C")
    try:
        channel.start_consuming()
    except KeyboardInterrupt:
        channel.stop_consuming()