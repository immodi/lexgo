local app = lexgo.new()

app.get("/search", function(req, res)
	local q = req.query.q and req.query.q[1] or "nothing"
	local page = req.query.page and req.query.page[1] or "1"

	res.status(200)
	res.json({
		message = "Query received",
		search = q,
		page = page,
	})
end)

app.post("/search", function(req, res)
	local q = req.query.q and req.query.q[1] or "nothing"
	local page = req.query.page and req.query.page[1] or "1"
	local filter = req.body.filter or "none"

	res.status(200)
	res.json({
		message = "POST query received",
		search = q,
		page = page,
		filter = filter,
	})
end)

app.listen(3000)
