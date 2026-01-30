local app = lexgo.new()

app.get("/", function(req, res)
	res.status(500)
	res.raw("Raw text")
end)

app.get("/html", function(req, res)
	res.html("<h1>Hello, this is an html response from Lua!</h1>")
end)

app.get("/json", function(req, res)
	res.status(200)
	res.json({ message = "this is a json response from lua" })
end)

app.post("/json", function(req, res)
	local response = req.body.message or nil

	res.status(200)
	res.json({ message = response })
end)

app.get("/search", function(req, res)
	local query = req.query
	local q = query.q and query.q[1] or "nothing"
	local page = query.page and query.page[1] or "1"

	res.status(200)
	res.json({
		message = "Query received",
		search = q,
		page = page,
	})
end)

app.post("/search", function(req, res)
	local query = req.query
	local q = query.q and query.q[1] or "nothing"
	local page = query.page and query.page[1] or "1"

	local body = req.body
	local filter = body.filter or "none"

	res.status(200)
	res.json({
		message = "POST query received",
		search = q,
		page = page,
		filter = filter,
	})
end)

app.get("/form", function(req, res)
	res.html([[
		<form action="/submit-form" method="POST">
			<label for="name">Name:</label>
			<input type="text" id="name" name="name"><br><br>
			<label for="age">Age:</label>
			<input type="number" id="age" name="age"><br><br>
			<input type="submit" value="Submit">
		</form>
	]])
end)

app.post("/submit-form", function(req, res)
	local body = req.body
	local name = body.name or "Anonymous"
	local age = body.age or "unknown"

	res.status(200)
	res.html("<h1>Form submitted!</h1><p>Name: " .. name .. "</p><p>Age: " .. age .. "</p>")
end)

app.notFound(function(req, res)
	res.status(404)
	res.html("<h1>Not Found</h1>")
end)
