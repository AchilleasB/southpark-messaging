import pika
import config
from redis import Redis


def get_connection():
    params = pika.URLParameters(config.RABBITMQ_URL)
    return pika.BlockingConnection(params)

def get_channel(connection):
    return connection.channel()

def get_redis():
    return Redis.from_url(config.REDIS_URL, decode_responses=True)