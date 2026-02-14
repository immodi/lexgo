local app = lexgo.new({})

app.get("/params/:id/:name", function(req, res)
	local id = req.params.id or "1"
	local name = req.params.name or "Name"

	res.html(string.format("<p>id: %s</p><p>name: %s</p>", id, name))
end)

app.listen(3000)
