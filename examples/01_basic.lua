local app = lexgo.new({})

app.get("/", function(req, res)
	res.status(500)
	res.raw("Raw text")
end)

app.listen(3000)
