package tables

import "gorm.io/gorm"

var (
    _chains []interface{}
)

func Add(p interface{}) {
    _chains = append(_chains, p)
}
func InitChain(db *gorm.DB) error {
    return db.AutoMigrate(_chains...)
}
