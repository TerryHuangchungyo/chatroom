module controller

go 1.15

require (
    github.com/gin-gonic/gin v1.6.3
    local/wsservice v0.0.0
)

replace (
    local/wsservice => ../wsservice
)