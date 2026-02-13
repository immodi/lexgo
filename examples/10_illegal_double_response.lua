local app = lexgo.new({})

app.use(function(req, res, next)
	print("Request received for:", req.url)
	next()
end)

app.get("/throw", function(req, res)
	res.status(200)
	res.json({ message = "this is a json response" })

	-- should trigger framework protection
	res.html("<h1>This is illegal</h1>")
end)

app.error(function(err, res)
	res.raw(err)
end)

app.notFound(function(req, res)
	res.status(400)
	res.json({ message = "not found" })
end)

app.listen(3000)
