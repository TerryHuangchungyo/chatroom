module main

go 1.15

require (
	github.com/gin-gonic/gin v1.6.3
	github.com/go-sql-driver/mysql v1.5.0
	github.com/gorilla/websocket v1.4.2
	local/config v0.0.0
	local/controller v0.0.0
	local/wsservice v0.0.0
	local/model v0.0.0
)

replace (
	local/config => ./config
	local/controller => ./controller
	local/wsservice => ./wsservice
	local/model => ./model
)
