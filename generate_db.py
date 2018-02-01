import requests
from faker import Faker

fake = Faker()

for i in range(10000):
     r = requests.post('http://127.0.0.1:8000/games', data='{"title": "' + fake.name() + '"}')