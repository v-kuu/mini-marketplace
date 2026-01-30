from locust import HttpUser, task, between
import random

class APIUser(HttpUser):
    wait_time = between(1, 3)

    @task(3)
    def get_products(self):
        self.client.get("/products")

    @task(2)
    def get_product(self):
        id = random.randint(1, 100)
        self.client.get(f"/products/{id}")

    @task(1)
    def create_product(self):
        id = random.randint(1, 100)
        self.client.post("/products", json={
            "id":"{id}",
            "name": "Test",
            "price": 499
        })

    @task(1)
    def health_check(self):
        self.client.get("/health")

