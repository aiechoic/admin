package gorm_test

import (
	"github.com/aiechoic/admin/core/gorm"
	"github.com/aiechoic/admin/core/ioc"
)

func ExampleGetDB() {
	c := ioc.NewContainer()

	db, err := gorm.GetDB("gorm", c)
	if err != nil {
		panic(err)
	}
	_ = db
}
