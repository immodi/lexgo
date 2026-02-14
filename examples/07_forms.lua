local app = lexgo.new({})

app.get("/form", function(req, res)
	res.html([[
		<form action="/submit-form" method="POST">
			<label>Name:</label>
			<input type="text" name="name"><br><br>
			<label>Age:</label>
			<input type="number" name="age"><br><br>
			<input type="submit" value="Submit">
		</form>
	]])
end)

app.post("/submit-form", function(req, res)
	local name = req.body.name or "Anonymous"
	local age = req.body.age or "unknown"

	res.status(200)
	res.html("<h1>Form submitted!</h1><p>Name: " .. name .. "</p><p>Age: " .. age .. "</p>")
end)

app.listen(3000)
