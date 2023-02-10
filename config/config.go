package config

import "time"

var (
    HTTPAddress    string
    JWTKey         = `yes...yes...yes...`
    JWTExpire      = time.Hour
    UseEmailVerify = false
)
