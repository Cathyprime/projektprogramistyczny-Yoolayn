require("project").command({
    nargs = 1,
    command = "Docker",
    map = {
        up = "!docker compose up -d",
        down = "!docker compose down",
    }
})

require("project").registers({
    m = ":e cmd/main.go\n",
    j = [[0wye$a<Space>`json:"<C-R>""`<Esc>bb~]],
    b = ":e internal/handlers/board_handlers.go\n",
    u = ":e internal/handlers/user_handlers.go\n",
})
