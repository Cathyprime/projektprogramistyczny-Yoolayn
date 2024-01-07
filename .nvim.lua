require("project").command({
    nargs = 1,
    command = "Docker",
    map = {
        up = "!docker compose up -d",
        down = "!docker compose down",
    }
})

require("project").registers({
    m = ":e server/src/main.go\n",
    u = ":!docker compose up -d\n",
    d = ":!docker compose down\n",
})
