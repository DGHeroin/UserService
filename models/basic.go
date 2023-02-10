package models

import (
    "errors"
    "gorm.io/gorm"
)

type BasicModel struct {
    gorm.Model
}

var (
    ErrExist    = errors.New("data exist")
    ErrNotExist = errors.New("data not exist")
)
