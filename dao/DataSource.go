package dao

import (
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"

    "gorm.io/driver/mysql"
    "gorm.io/driver/postgres"
)

func NewDB(dbType, dataSourceName string) (*gorm.DB, error) {
    switch dbType {
    case "mysql":
        // "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
        return gorm.Open(mysql.Open(dataSourceName), &gorm.Config{})
    case "postgres":
        // "host=localhost user=gorm password=gorm dbname=gorm port=9920 sslmode=disable TimeZone=Asia/Shanghai"
        return gorm.Open(postgres.Open(dataSourceName), &gorm.Config{})
    default:
        // "/data/my.db"
        return gorm.Open(sqlite.Open(dataSourceName), &gorm.Config{})
    }
}

var sharedDB *gorm.DB

func SharedDB(args ...*gorm.DB) *gorm.DB {
    if len(args) == 1 && args[0] != nil {
        sharedDB = args[0]
    }
    return sharedDB
}
