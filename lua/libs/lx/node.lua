---@class Node
---@field tag string
---@field props table<string, any>
---@field children Node[]

local lx = {}

---@param tag string
---@param props table<string, any>|nil
---@param children Node[]|nil
---@return Node
function lx.node(tag, props, children)
    props = props or {}
    children = children or {}

    return setmetatable({
        tag = tag,
        props = props,
        children = children,
    }, {
        ---@param self Node
        __tostring = function(self)
            local props_str = "{"
            for k, v in pairs(self.props) do
                props_str = props_str .. string.format(" %s=%s", k, tostring(v))
            end
            props_str = props_str .. " }"

            local children_str = "["
            for i, child in ipairs(self.children) do
                children_str = children_str .. tostring(child)
                if i < #self.children then
                    children_str = children_str .. ", "
                end
            end
            children_str = children_str .. "]"

            return string.format("{ tag = %s, props = %s, children = %s }", self.tag, props_str, children_str)
        end,
    })
end

return lx
