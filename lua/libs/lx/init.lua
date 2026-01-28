local node_module = require("lua.libs.lx.node")
local element_module = require("lua.libs.lx.element")

local lx = {
    node = node_module.node,
}

-- Create element with node injected
lx.element = element_module.element(lx.node)


-- define standard elements
for _, tag in ipairs({ "div", "span" }) do
    lx[tag] = lx.element(tag)
end

return lx
