local app = lexgo.new()

app.get("/throw", function(req, res)
	res.status(200)
	res.json({ message = "this is a json response" })

	-- should trigger framework protection
	res.html("<h1>This is illegal</h1>")
end)

app.listen(3000)
