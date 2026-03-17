import requests
import sys
import subprocess
import time

BASE = "http://localhost:8081"
TIMEOUT = 2

class TestRunner:
    def __init__(self):
        self.total = 0
        self.failed = 0

    def test(self, name, fn):
        self.total += 1
        try:
            fn()
            print(f"PASS  {name}")
        except Exception as e:
            self.failed += 1
            print(f"FAIL  {name} -> {e}")

    def summary(self):
        print("\n======================")
        print(f"Total: {self.total}")
        print(f"Passed: {self.total - self.failed}")
        print(f"Failed: {self.failed}")
        print("======================")
        sys.exit(0 if self.failed == 0 else 1)


def assert_equal(a, b, msg):
    if a != b:
        raise Exception(f"{msg} (expected {b}, got {a})")


def test_hello():
    r = requests.get(BASE + "/hello", timeout=TIMEOUT)
    assert_equal(r.status_code, 200, "status")
    assert_equal(r.text, "hello world", "body")


def test_status():
    r = requests.get(BASE + "/status", timeout=TIMEOUT)
    assert_equal(r.status_code, 201, "status")
    assert_equal(r.text, "created", "body")


def test_json():
    r = requests.get(BASE + "/json", timeout=TIMEOUT)
    assert_equal(r.status_code, 200, "status")
    data = r.json()
    assert_equal(data["ok"], True, "json ok")
    assert_equal(data["message"], "json works", "json message")


def test_params():
    r = requests.get(BASE + "/user/42", timeout=TIMEOUT)
    assert_equal(r.status_code, 200, "status")
    assert_equal(r.text, "user 42", "param extraction")


def test_method_get():
    r = requests.get(BASE + "/method", timeout=TIMEOUT)
    assert_equal(r.status_code, 200, "status")
    assert_equal(r.text, "GET OK", "GET handler")


def test_method_post():
    r = requests.post(BASE + "/method", timeout=TIMEOUT)
    assert_equal(r.status_code, 200, "status")
    assert_equal(r.text, "POST OK", "POST handler")


def test_large():
    r = requests.get(BASE + "/large", timeout=TIMEOUT)
    assert_equal(r.status_code, 200, "status")
    assert_equal(len(r.text), 1000, "large response size")


def test_not_found():
    r = requests.get(BASE + "/missing", timeout=TIMEOUT)
    assert_equal(r.status_code, 404, "status")
    assert_equal(r.text, "not found", "body")


def test_error():
    r = requests.get(BASE + "/error", timeout=TIMEOUT)
    assert_equal(r.status_code, 500, "status")
    if not r.text.startswith("engine error"):
        raise Exception("error handler not executed")


def main():
    runner = TestRunner()
    runner.test("GET /hello", test_hello)
    runner.test("GET /status", test_status)
    runner.test("GET /json", test_json)
    runner.test("GET /user/:id", test_params)
    runner.test("GET /method", test_method_get)
    runner.test("POST /method", test_method_post)
    runner.test("GET /large", test_large)
    runner.test("GET /missing", test_not_found)
    runner.test("GET /error", test_error)

    runner.summary()


if __name__ == "__main__":
    main()
