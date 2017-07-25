import socket
import time


SERVER_ADDR = '2001:412:abcd:2:0013:A200:4147:8C2B'


if __name__ == '__main__':
    c = 0
    sock = socket.socket(socket.AF_INET6, socket.SOCK_DGRAM)
    while True:
        sock.sendto(b'HELLO %i' % c, (SERVER_ADDR, 8000))
        c += 1
        time.sleep(1)
