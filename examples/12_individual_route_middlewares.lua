local app = lexgo.new({})

-- middlwares

-- app.use(lexgo.middlewares.logger)

local function indexLogger(req, res, next)
	print("request for:", req.url)
	next()
end

-- routes

local function index(req, res)
	res.status(200)
	res.json({ message = "this is a json response" })
end

local function home(req, res)
	res.status(200)
	res.html("<h1>HOME ROUTE</h1>")
end

local function notFound(req, res)
	res.status(404)
	res.raw("custom 404")
end

app.get("/", index, indexLogger)
app.get("/home", home)
app.notFound(notFound)

app.listen(3000)
