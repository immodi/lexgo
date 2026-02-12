local app = lexgo.new()

app.get("/html", function(req, res)
	res.html("<h1>Hello, this is an html response from Lua!</h1>")
end)

app.listen(3000)
