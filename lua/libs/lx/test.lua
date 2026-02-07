local lx = require("lua.libs.lx.lx")

-- print(lx.node("h1", false, { class = "text-area" }, {}))
-- print("\n")

local n = lx.div {
    "Hello",
    id = "main",
    lx.span { "Nested" }
}

print(n:render())
