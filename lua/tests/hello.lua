app.get("/", function(req, res)
    res.send("<h1>Hello from Lua!</h1>")
end)


app.get("/ping", function(req, res)
    res.send("<h5>pong</h5>")
end)
