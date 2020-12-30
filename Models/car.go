package Models

import "github.com/jinzhu/gorm"

//Car Model
type Car struct {
	gorm.Model

	Year      int
	Make      string
	ModelName string
	DriverID  int
}
