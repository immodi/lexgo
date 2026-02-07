local app = lexgo.new()

app.get("/boom", function(req, res)
	error("something exploded")
end)

app.notFound(function(req, res)
	res.status(404)
	res.html("<h1>Not Found</h1>")
end)

app.error(function(err, res)
	print(err)
	res.status(500)
	res.html('<h1 style="color: red;">Internal Server Error</h1>')
end)

app.listen(3000)
