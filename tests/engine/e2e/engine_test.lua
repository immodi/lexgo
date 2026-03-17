local app = lexgo.new({})

-- =========================
-- BASIC RESPONSE
-- =========================
app.get("/hello", function(req, res)
	res.raw("hello world")
end)

-- =========================
-- STATUS CODE
-- =========================
app.get("/status", function(req, res)
	res.status(201)
	res.raw("created")
end)

-- =========================
-- JSON RESPONSE
-- =========================
app.get("/json", function(req, res)
	res.json({
		ok = true,
		message = "json works",
	})
end)

-- =========================
-- PARAM PASSING
-- =========================
app.get("/user/:id", function(req, res)
	res.raw("user " .. req.params.id)
end)

-- =========================
-- REQUEST METHOD CHECK
-- =========================
app.post("/method", function(req, res)
	res.raw("POST OK")
end)

app.get("/method", function(req, res)
	res.raw("GET OK")
end)

-- =========================
-- LARGE BODY
-- =========================
app.get("/large", function(req, res)
	local s = ""
	for i = 1, 1000 do
		s = s .. "a"
	end
	res.raw(s)
end)

-- =========================
-- ERROR HANDLING
-- =========================
app.get("/error", function(req, res)
	error("boom")
end)

-- =========================
-- CUSTOM ERROR HANDLER
-- =========================
app.error(function(req, res, err)
	res.status(500)
	res.raw("engine error: " .. err)
end)

-- =========================
-- NOT FOUND
-- =========================
app.notFound(function(req, res)
	res.status(404)
	res.raw("not found")
end)

app.listen(8081)
