package Models

import "github.com/jinzhu/gorm"

//Car Model
type Car struct {
	gorm.Model

	Year      string
	Make      string
	ModelName string
	DriverID  string
}
