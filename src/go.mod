module main

go 1.15

require (
    github.com/gin-gonic/gin v1.6.3
    github.com/gorilla/websocket v1.4.2
    local/wsservice v0.0.0
    local/controller v0.0.0
)

replace (
    local/wsservice => ./wsservice
    local/controller => ./controller
)