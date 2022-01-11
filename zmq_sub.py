"""
subscribe data and print to stdout
usage:
python3 zmq_sub.py <port>
"""
import zmq
import sys

ctx = zmq.Context()
sock = ctx.socket(zmq.SUB)
sock.subscribe(b'')
print("connect to port", sys.argv[1])
sock.connect("tcp://localhost:" + sys.argv[1])
while True:
    data = sock.recv_multipart()
    print(data)
