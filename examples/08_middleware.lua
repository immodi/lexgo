local app = lexgo.new()

-- app.use(lexgo.middlewares.cors)
-- app.use(lexgo.middlewares.logger)

app.use(function(req, res, next)
	print("Request received for:", req.url)
	next()
end)

app.get("/", function(req, res)
	res.raw("Middleware test")
end)

app.listen(3000)
