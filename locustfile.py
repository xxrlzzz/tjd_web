#! /bin/python3

from locust import HttpUser, between, task


class WebsiteUser(HttpUser):
    wait_time = between(5, 15)

    def on_start(self):
        pass
#         self.client.post("/login", {
#             "username": "test_user",
#             "password": ""
#         })

    @task
    def token(self):
        self.client.post("/token", json={"username": "1"})
        self.client.post("/token", json={"username": "xxrl"})

    @task
    def ping(self):
        self.client.get("/ping/")