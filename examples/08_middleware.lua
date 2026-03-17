local app = lexgo.new({
	env = "production",
	allowedOrigins = { "t1", "t2" },
})

app.use(lexgo.middlewares.cors)
app.use(lexgo.middlewares.logger)

app.use(function(req, res, next)
	print("Request received for:", req.url)
	next()
end)

app.get("/", function(req, res)
	res.raw("Middleware test")
end)

app.post("/", function(req, res)
	res.raw("Middleware test 2")
end)

app.listen(3000)
