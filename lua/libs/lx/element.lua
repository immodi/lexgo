-- This factory function uses dependency injection to avoid circular dependencies.
-- Instead of element.lua requiring init.lua (which would create a cycle since
-- init.lua already requires element.lua), we pass the node function as a parameter.
-- This allows element to use node without needing to know where it comes from.

---@param node_fn function
---@return function
local function create_element(node_fn)
    ---@param tag string
    ---@return table
    return function(tag)
        return setmetatable({}, {
            __call = function(_, arg)
                local props = {}
                local children = {}
                if type(arg) == "table" then
                    for k, v in pairs(arg) do
                        if type(k) == "string" then
                            props[k] = v
                        else
                            table.insert(children, v)
                        end
                    end
                elseif type(arg) == "string" then
                    table.insert(children, arg)
                end
                return node_fn(tag, props, children)
            end
        })
    end
end

return { element = create_element }
