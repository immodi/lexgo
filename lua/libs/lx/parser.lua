local parser = {}

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

function parser.parseChild(c)
    local parts = {}
    for _, child in ipairs(c.children) do
        if type(child) == "string" and child ~= "" then
            parts[#parts + 1] = escape_text(child)
        elseif type(child) == "table" then
            parts[#parts + 1] = parser.parser(child)
        end
    end
    return table.concat(parts)
end

---@param t table
---@return string
function parser.parser(t)
    local element = "<" .. t.tag

    if t.props then
        local props = {}
        for k, v in pairs(t.props) do
            if type(v) == "boolean" then
                if v then
                    props[#props + 1] = k
                end
            else
                props[#props + 1] = string.format('%s="%s"', k, escape_attr(v))
            end
        end
        if #props > 0 then
            element = element .. " " .. table.concat(props, " ")
        end
    end

    if t.isSelfClosing then
        element = element .. " />"
    else
        element = element .. ">" .. parser.parseChild(t) .. "</" .. t.tag .. ">"
    end

    return element
end

return parser.parser
