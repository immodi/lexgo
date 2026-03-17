local app = lexgo.new({})

-- =========================
-- ROOT
-- =========================
app.get("/", function(req, res)
	res.raw("root")
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
	res.raw("user name " .. req.params.name)
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
-- METHOD MATCHING
-- =========================
app.post("/method", function(req, res)
	res.raw("POST")
end)

app.get("/method", function(req, res)
	res.raw("GET")
end)

-- =========================
-- NOT FOUND
-- =========================
app.notFound(function(req, res)
	res.status(404)
	res.raw("not found")
end)

app.listen(8081)
