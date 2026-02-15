local app = lexgo.new({})
-- app.use(lexgo.middlewares.logger)

app.get("/", function(req, res)
	res.status(200)
	res.raw("Hello, world")
end)

app.get("/wild/:id/details", function(req, res)
	res.status(200)
	res.raw(req.params.id)
end)

app.get("/wild/*", function(req, res)
	res.status(200)
	res.html(string.format("<h1>%s</h1>", req.params["*"]))
end)

app.listen(3000)

-- TODO: NEED TO TEST IF IT WORKS
