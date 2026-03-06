---@class Node
---@field tag string
---@field isSelfClosing boolean
---@field props table<string, any>
---@field children Node[]

local lx = {}

------------------------------------------------------------------
-- Parser
------------------------------------------------------------------

local function escape_attr(s)
    return tostring(s):gsub('"', "&quot;")
end

local function escape_text(s)
    s = tostring(s)
    s = s:gsub("&", "&amp;")
    s = s:gsub("<", "&lt;")
    s = s:gsub(">", "&gt;")
    return s
end

local function parseChild(c)
    local parts = {}

    for _, child in ipairs(c.children) do
        if type(child) == "string" and child ~= "" then
            parts[#parts + 1] = escape_text(child)
        elseif type(child) == "table" then
            parts[#parts + 1] = lx._parser(child)
        end
    end

    return table.concat(parts)
end

function lx._parser(t)
    local element = "<" .. t.tag

    if t.props then
        local props = {}
        for k, v in pairs(t.props) do
            if type(v) == "boolean" then
                if v then
                    props[#props + 1] = k
                end
            else
                props[#props + 1] =
                    string.format('%s="%s"', k, escape_attr(v))
            end
        end

        if #props > 0 then
            element = element .. " " .. table.concat(props, " ")
        end
    end

    if t.isSelfClosing then
        element = element .. " />"
    else
        element =
            element .. ">" .. parseChild(t) .. "</" .. t.tag .. ">"
    end

    return element
end

------------------------------------------------------------------
-- Node
------------------------------------------------------------------

local function fmt(v)
    if type(v) == "boolean" then
        return v and "true" or "false"
    elseif type(v) == "string" then
        return '"' .. v .. '"'
    else
        return tostring(v)
    end
end

function lx.node(tag, isSelfClosing, props, children)
    props = props or {}
    isSelfClosing = isSelfClosing or false
    children = children or {}

    local element = {
        tag = tag,
        isSelfClosing = isSelfClosing,
        props = props,
        children = children,
    }

    function element:render()
        if not self then
            error("can't render a nil value, you need to pass the lx element to render()")
        end
        return lx._parser(self)
    end

    return setmetatable(element, {
        __tostring = function(self)
            local props_str = "{"
            for k, v in pairs(self.props) do
                props_str =
                    props_str .. string.format(" %s=%s", k, fmt(v))
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

------------------------------------------------------------------
-- Element Factory (dependency injection kept internally)
------------------------------------------------------------------

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

lx.element = create_element(lx.node)

------------------------------------------------------------------
-- Lazy tag creation (lx.div, lx.h1, etc.)
------------------------------------------------------------------

setmetatable(lx, {
    __index = function(t, key)
        if type(key) ~= "string" then
            return nil
        end

        local el = t.element(key)
        rawset(t, key, el)
        return el
    end,
})

------------------------------------------------------------------

return lx
