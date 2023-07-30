import random
import socket
import time

if __name__ == "__main__":
    s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    s.connect(("localhost", 8888))
    time.sleep(10)
    s.send(str.encode("".join([random.choice('abcdefghijklmnopqrstuvwxyz!@#$%^&*()') for i in range(2000)])+"\n"))
    time.sleep(10)
    s.send(str.encode("".join([random.choice('abcdefghijklmnopqrstuvwxyz!@#$%^&*()') for i in range(2000)])+"\n"))
    # data = s.recv(1024)
    # print(data)
    s.close()

