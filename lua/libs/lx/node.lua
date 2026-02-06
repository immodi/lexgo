---@class Node
---@field tag string
---@field isSelfClosing boolean
---@field props table<string, any>
---@field children Node[]

local lx = {}

---@param v any
---@return string
local function fmt(v)
    if type(v) == "boolean" then
        return v and "true" or "false"
    elseif type(v) == "string" then
        return '"' .. v .. '"' -- wrap strings in quotes
    else
        return tostring(v)
    end
end

---@param tag string
---@param isSelfClosing boolean
---@param props table<string, any>|nil
---@param children Node[]|nil
---@return Node
function lx.node(tag, isSelfClosing, props, children)
    props = props or {}
    isSelfClosing = isSelfClosing or false
    children = children or {}

    return setmetatable({
        tag = tag,
        isSelfClosing = isSelfClosing,
        props = props,
        children = children,
    }, {
        ---@param self Node
        __tostring = function(self)
            local props_str = "{"
            for k, v in pairs(self.props) do
                props_str = props_str .. string.format(" %s=%s", k, fmt(v))
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

            return string.format(
                "{ tag=%s, isSelfClosing=%s, props=%s, children=%s }",
                fmt(self.tag),
                fmt(self.isSelfClosing),
                props_str,
                children_str
            )
        end,
    })
end

return lx
