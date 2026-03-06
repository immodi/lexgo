import requests
import sys

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


# =========================
# ROUTE TESTS
# =========================
def test_root():
    r = requests.get(BASE + "/", timeout=TIMEOUT)
    assert_equal(r.status_code, 200, "/")
    assert_equal(r.text, "root", "/ body")


def test_static_routes():
    r1 = requests.get(BASE + "/static/hello", timeout=TIMEOUT)
    assert_equal(r1.status_code, 200, "/static/hello")
    assert_equal(r1.text, "static hello", "/static/hello body")

    r2 = requests.get(BASE + "/static/nested/test", timeout=TIMEOUT)
    assert_equal(r2.status_code, 200, "/static/nested/test")
    assert_equal(r2.text, "static nested", "/static/nested/test body")


def test_param_routes():
    r1 = requests.get(BASE + "/user/123", timeout=TIMEOUT)
    assert_equal(r1.status_code, 200, "/user/:id")
    assert_equal(r1.text, "user 123", "/user/:id body")

    r2 = requests.get(BASE + "/user/42/john", timeout=TIMEOUT)
    assert_equal(r2.status_code, 200, "/user/:id/:name")
    assert_equal(r2.text, "user name john", "/user/:id/:name body")

    r3 = requests.get(BASE + "/user/abc/details", timeout=TIMEOUT)
    assert_equal(r3.status_code, 200, "/user/:id/details")
    assert_equal(r3.text, "details abc", "/user/:id/details body")


def test_wildcard_routes():
    r1 = requests.get(BASE + "/files/a.txt", timeout=TIMEOUT)
    assert_equal(r1.status_code, 200, "/files/*")
    assert_equal(r1.text, "files a.txt", "/files/* body")

    r2 = requests.get(BASE + "/files/a/b/c", timeout=TIMEOUT)
    assert_equal(r2.status_code, 200, "/files/* nested")
    assert_equal(r2.text, "files a/b/c", "/files/* nested body")


def test_priority_routes():
    r1 = requests.get(BASE + "/priority/test", timeout=TIMEOUT)
    assert_equal(r1.status_code, 200, "priority static")
    assert_equal(r1.text, "priority static", "priority static body")

    r2 = requests.get(BASE + "/priority/42", timeout=TIMEOUT)
    assert_equal(r2.status_code, 200, "priority param")
    assert_equal(r2.text, "priority param 42", "priority param body")

    r3 = requests.get(BASE + "/priority/a/b/c", timeout=TIMEOUT)
    assert_equal(r3.status_code, 200, "priority wildcard")
    assert_equal(r3.text, "priority wild a/b/c", "priority wildcard body")


def test_method_routes():
    r1 = requests.get(BASE + "/method", timeout=TIMEOUT)
    assert_equal(r1.status_code, 200, "GET /method")
    assert_equal(r1.text, "GET", "GET /method body")

    r2 = requests.post(BASE + "/method", timeout=TIMEOUT)
    assert_equal(r2.status_code, 200, "POST /method")
    assert_equal(r2.text, "POST", "POST /method body")

    r3 = requests.put(BASE + "/method", timeout=TIMEOUT)
    assert_equal(r3.status_code, 404, "PUT /method")


def test_not_found():
    r = requests.get(BASE + "/does-not-exist", timeout=TIMEOUT)
    assert_equal(r.status_code, 404, "/does-not-exist")
    assert_equal(r.text, "not found", "/does-not-exist body")


def main():
    runner = TestRunner()

    runner.test("ROOT /", test_root)
    runner.test("STATIC routes", test_static_routes)
    runner.test("PARAM routes", test_param_routes)
    runner.test("WILDCARD routes", test_wildcard_routes)
    runner.test("PRIORITY routes", test_priority_routes)
    runner.test("METHOD routes", test_method_routes)
    runner.test("NOT FOUND", test_not_found)

    runner.summary()


if __name__ == "__main__":
    main()
