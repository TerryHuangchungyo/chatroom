module model

go 1.15

require (
    "github.com/go-sql-driver/mysql" v1.5.0
    "local/config" v0.0.0
)

replace (
    "local/config" => "../config"
)
