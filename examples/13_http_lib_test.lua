local app = lexgo.new({})
local http = lexgo.libs.http

app.use(lexgo.middlewares.logger)

-- GET test (your existing one)
app.get("/http-test", function(req, res)
	local r, err = http.get("https://api.github.com")

	if err then
		res.status(500)
		res.html("<h1>HTTP Request Failed</h1><p>" .. err .. "</p>")
		return
	end

	local contentType = "unknown"
	if r.headers["Content-Type"] then
		contentType = r.headers["Content-Type"][1]
	end

	local html = "<h1>HTTP Request Success</h1>"
		.. "<p>Status: "
		.. r.status
		.. "</p>"
		.. "<p>Content-Type: "
		.. contentType
		.. "</p>"
		.. "<pre>"
		.. r.body
		.. "</pre>"

	res.status(200)
	res.html(html)
end)

-- POST test
app.get("/http-post-test", function(req, res)
	local payload = {
		name = "Ahmed",
		email = "ahmed@example.com",
	}

	local r, err = http.post("https://httpbin.org/post", "application/json", payload)

	if err then
		res.status(500)
		res.html("<h1>HTTP POST Failed</h1><p>" .. err .. "</p>")
		return
	end

	local contentType = "unknown"
	if r.headers["Content-Type"] then
		contentType = r.headers["Content-Type"][1]
	end

	local html = "<h1>HTTP POST Success</h1>"
		.. "<p>Status: "
		.. r.status
		.. "</p>"
		.. "<p>Content-Type: "
		.. contentType
		.. "</p>"
		.. "<pre>"
		.. r.body
		.. "</pre>"

	res.status(200)
	res.html(html)
end)

app.listen(3000)
