import socket


BIND_ADDR = '::'


if __name__ == '__main__':
    sock = socket.socket(socket.AF_INET6, socket.SOCK_DGRAM)
    sock.bind((BIND_ADDR, 8000))
    while True:
        data, addr = sock.recvfrom(10)
        print(data, addr)
