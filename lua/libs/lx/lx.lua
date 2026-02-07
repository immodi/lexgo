local node_module = require("lua.libs.lx.node")
local element_module = require("lua.libs.lx.element")

local lx = {
    node = node_module.node,
}

-- Create element with node injected
lx.element = element_module.element(lx.node)

setmetatable(lx, {
    __index = function(t, key)
        -- ignore internal lookups
        if type(key) ~= "string" then return nil end

        -- create the element factory
        local el = t.element(key)

        -- cache it so next access is fast
        rawset(t, key, el)

        return el
    end,
})

return lx
