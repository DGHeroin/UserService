package main

import (
    "UserService/HTTPAPI"
    "UserService/config"
    "UserService/dao"
    "UserService/dao/tables"
    "UserService/logger"
    "UserService/logger/ginlog"
    "flag"
    "github.com/gin-gonic/gin"
)

func init() {
    flag.StringVar(&config.HTTPAddress, "addr", ":8081", "http service address")
    flag.StringVar(&config.JWTKey, "jwtKey", "uSERsERVER-jwtkEY", "user service jwt key")
    flag.BoolVar(&config.UseEmailVerify, "useEmailVerify", false, "user email verify")
    flag.Parse()
}

func main() {
    { // 初始化 DB
        db, err := dao.NewDB("sqlite", "test.db")
        if err != nil {
            panic(err)
        }

        if err := tables.InitChain(db); err != nil {
            panic(err)
        }
        // set shared db
        dao.SharedDB(db)
    }

    { // 启动 http 服务
        gin.SetMode(gin.ReleaseMode)
        r := gin.Default()
        r.Use(ginlog.GetLogger(logger.Default()))

        HTTPAPI.InitGinAPI(r)

        if err := r.Run(config.HTTPAddress); err != nil {
            panic(err)
        }
    }

}
