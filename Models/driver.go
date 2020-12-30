package Models

import "github.com/jinzhu/gorm"

type Driver struct {
	gorm.Model
	Name    string
	License string
	Cars    []Car
}
