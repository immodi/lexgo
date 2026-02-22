local app = lexgo.new({})
app.use(lexgo.middlewares.logger)

-- =========================
-- ROOT
-- =========================
app.get("/", function(req, res)
	res.status(200)
	res.raw("root ok")
end)

-- =========================
-- STATIC ROUTES
-- =========================
app.get("/static/hello", function(req, res)
	res.raw("static hello")
end)

app.get("/static/nested/test", function(req, res)
	res.raw("static nested")
end)

-- =========================
-- PARAM ROUTES
-- =========================
app.get("/user/:id", function(req, res)
	res.raw("user " .. req.params.id)
end)

app.get("/user/:id/:name", function(req, res)
	res.raw("user + name " .. req.params.id)
end)

app.get("/user/:id/details", function(req, res)
	res.raw("details " .. req.params.id)
end)

-- =========================
-- WILDCARD ROUTES
-- =========================
app.get("/files/*", function(req, res)
	res.raw("files " .. req.params["*"])
end)

-- =========================
-- PRIORITY TEST
-- static > param > wildcard
-- =========================
app.get("/priority/test", function(req, res)
	res.raw("priority static")
end)

app.get("/priority/:id", function(req, res)
	res.raw("priority param " .. req.params.id)
end)

app.get("/priority/*", function(req, res)
	res.raw("priority wild " .. req.params["*"])
end)

-- =========================
-- METHOD TESTS
-- =========================
app.post("/method", function(req, res)
	res.raw("POST OK")
end)

app.get("/method", function(req, res)
	res.raw("GET OK")
end)

-- =========================
-- ERROR ROUTE
-- =========================
app.get("/boom", function(req, res)
	error("something exploded")
end)

-- =========================
-- NOT FOUND
-- =========================
app.notFound(function(req, res)
	res.status(404)
	res.raw("custom 404")
end)

-- =========================
-- ERROR HANDLER
-- =========================
app.error(function(req, res, err)
	print("ERROR:", err)
	res.status(500)
	res.raw("custom 500")
end)

app.listen(3000)
