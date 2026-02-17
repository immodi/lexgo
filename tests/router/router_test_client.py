import requests

BASE = "http://localhost:3000"

def test(method, path, expected_status, expected_text=None):
    try:
        response = requests.request(method, BASE + path)

        status_ok = response.status_code == expected_status
        text_ok = True

        if expected_text is not None:
            text_ok = expected_text in response.text

        success = status_ok and text_ok
        result = "PASS" if success else "FAIL"

        print(f"[{result}] {method} {path}")
        print(f"       Status: {response.status_code}")
        print(f"       Body:   {response.text}")
        print()

    except Exception as e:
        print(f"[ERROR] {method} {path} -> {e}")


def run():
    print("==== ROOT ====")
    test("GET", "/", 200, "root ok")

    print("==== STATIC ====")
    test("GET", "/static/hello", 200, "static hello")
    test("GET", "/static/nested/test", 200, "static nested")

    print("==== PARAM ====")
    test("GET", "/user/123", 200, "user 123")
    test("GET", "/user/abc/details", 200, "details abc")

    print("==== WILDCARD ====")
    test("GET", "/files/a.txt", 200, "files a.txt")
    test("GET", "/files/a/b/c", 200, "files a/b/c")

    print("==== PRIORITY ====")
    test("GET", "/priority/test", 200, "priority static")
    test("GET", "/priority/42", 200, "priority param 42")
    test("GET", "/priority/a/b/c", 200, "priority wild a/b/c")

    print("==== METHOD ====")
    test("GET", "/method", 200, "GET OK")
    test("POST", "/method", 200, "POST OK")
    test("PUT", "/method", 404)

    print("==== ERROR HANDLING ====")
    test("GET", "/boom", 500, "custom 500")

    print("==== NOT FOUND ====")
    test("GET", "/does-not-exist", 404, "custom 404")


if __name__ == "__main__":
    run()
