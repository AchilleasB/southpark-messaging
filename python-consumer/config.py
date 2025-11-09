import os

RABBITMQ_URL = os.getenv("RABBITMQ_URL", "amqp://guest:guest@rabbitmq:5672/")
QUEUE_NAME = os.getenv("QUEUE_NAME", "southpark_messages")

REDIS_URL = os.getenv("REDIS_URL", "redis://redis:6379/0")
REDIS_KEY = os.getenv("REDIS_KEY", "southpark_messages") 