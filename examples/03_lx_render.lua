local lx = require("lua.libs.lx.lx")

local app = lexgo.new({})

app.get("/lx-test", function(req, res)
	local html = lx.div({
		id = "main",
		class = "container",
		lx.h1({ "Hello from lx!" }),
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
