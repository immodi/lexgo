local app = lexgo.new()

app.get("/json", function(req, res)
	res.status(200)
	res.json({ message = "this is a json response from lua" })
end)

app.post("/json", function(req, res)
	local response = req.body.message or nil

	res.status(200)
	res.json({ message = response })
end)

app.put("/put-test", function(req, res)
	local message = req.body.message or "default PUT message"

	res.status(200)
	res.json({ method = "PUT", message = message })
end)

app.delete("/delete-test", function(req, res)
	res.status(200)
	res.json({ method = "DELETE", message = "Deleted successfully" })
end)

app.patch("/patch-test", function(req, res)
	local patchField = req.body.patch or "none"

	res.status(200)
	res.json({ method = "PATCH", patch = patchField })
end)

app.options("/options-test", function(req, res)
	res.status(200)
	res.setHeader("Allow", "GET, POST, PUT, DELETE, PATCH, OPTIONS")
	res.raw("")
end)

app.listen(3000)
