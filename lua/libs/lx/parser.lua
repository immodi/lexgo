local function escape_attr(s)
    return tostring(s):gsub('"', "&quot;")
end

local function parseChild(c)
    local parts = {}
    for _, child in ipairs(c.children) do
        if type(child) == "string" and child ~= "" then
            parts[#parts + 1] = child
        elseif type(child) == "table" then
            parts[#parts + 1] = parser(child)
        end
    end
    return table.concat(parts)
end

function parser(t)
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
        element = element .. ">" .. parseChild(t) .. "</" .. t.tag .. ">"
    end

    return element
end

return parser
