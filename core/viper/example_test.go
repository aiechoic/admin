package viper_test

import (
	"fmt"
	"github.com/aiechoic/admin/core/ioc"
	"github.com/aiechoic/admin/core/viper"
	"os"
	"path/filepath"
)

func ExampleGetViper() {
	c := ioc.NewContainer()

	_ = os.Remove(filepath.Join("configs", "testing", "viper.yaml"))

	v, err := viper.GetViper("viper", `ni: "hao"`, c) // create new config file
	if err != nil {
		panic(err)
	}
	value := v.GetString("ni")
	fmt.Println(value)

	v, err = viper.GetViper("viper", `ni: "..."`, c) // read from existing config file
	if err != nil {
		panic(err)
	}
	value = v.GetString("ni")
	fmt.Println(value)

	// Output:
	// hao
	// hao
}
