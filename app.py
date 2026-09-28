from flask import Flask
import os
import socket
import redis

app = Flask(__name__)

redis_host = os.getenv("REDIS_HOST", "redis")
r = redis.Redis(host=redis_host, port=6379, decode_responses=True)

@app.route("/")
def hello():
    visits = r.incr("visits")
    hostname = socket.gethostname()
    return f"Hello from Docker! Hostname: {hostname}. Visits: {visits}\n"

if __name__ == "__main__":
    app.run(host="0.0.0.0", port=5000)
