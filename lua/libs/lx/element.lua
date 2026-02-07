-- This factory function uses dependency injection to avoid circular dependencies.
-- Instead of element.lua requiring init.lua (which would create a cycle since
-- init.lua already requires element.lua), we pass the node function as a parameter.
-- This allows element to use node without needing to know where it comes from.

---@param node_fn function
---@return function
local function create_element(node_fn)
    return function(tag)
        if type(tag) ~= "string" or tag == "element" then
            error("can't use 'lx.element' as a DOM element")
        end
        return setmetatable({}, {
            __call = function(_, arg)
                local props = {}
                local children = {}
                local isSelfClosing = false

                for k, v in pairs(arg) do
                    if type(k) == "number" and type(v) == "string" then
                        table.insert(children, v)
                    elseif type(k) == "number" and type(v) == "boolean" then
                        isSelfClosing = v
                    elseif type(k) == "number" and type(v) == "table" then
                        table.insert(children, v)
                    elseif type(k) == "string" then
                        props[k] = v
                    end
                end

                return node_fn(tag, isSelfClosing, props, children)
            end,
        })
    end
end

return { element = create_element }
