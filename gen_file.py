import random

count = 10**8

with open("data","wb+") as f:
    for i in range(count):
        f.write(random.randint(0,count).to_bytes(8,byteorder='big',signed=True))