local app = lexgo.new({})
local lx = lexgo.libs.lx

app.use(lexgo.middlewares.logger)

app.get("/lx", function(req, res)
	local html = lx.div({
		id = "container",
		class = "container",
		lx.h1({ "Hello from LX!" }),
		lx.section({
			class = "section-class",
			lx.p({ class = "text-span", "This is a span" }),
			lx.img({
				true,
				loading = "lazy",
				style = "width: 100px; height: 100px;",
				src = "https://www.w3schools.com/html/img_girl.jpg",
			}),
		}),
	}):render()

	res.status(200)
	res.html(html)
end)

app.listen(3000)
