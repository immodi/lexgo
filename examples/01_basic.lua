local app = lexgo.new({
	env = "dev",
	allowedOrigins = { "*" },
})

app.get("/", function(req, res)
	res.status(500)
	res.raw("Raw text hello")
end)

app.listen(3000)
