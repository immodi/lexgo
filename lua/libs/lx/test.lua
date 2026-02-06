local lx = require("lua.libs.lx.init")

print(lx.node("h1", nil, nil))
print("\n")

local n = lx.div {
    "Hello",
    id = "main",
    lx.span { "Nested" }
}

print(n)
