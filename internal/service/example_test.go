package service

import (
	"fmt"
	"github.com/aiechoic/admin/internal/ioc"
	"os"
	"path/filepath"
)

func ExampleGetViper() {
	c := ioc.NewContainer()

	_ = os.Remove(filepath.Join("configs", "testing", "viper.yaml"))

	v, err := GetViper("viper", `ni: "hao"`, c) // create new config file
	if err != nil {
		panic(err)
	}
	value := v.GetString("ni")
	fmt.Println(value)

	v, err = GetViper("viper", `ni: "..."`, c) // read from existing config file
	if err != nil {
		panic(err)
	}
	value = v.GetString("ni")
	fmt.Println(value)

	// Output:
	// hao
	// hao
}

func ExampleGetDefaultRedisClient() {
	c := ioc.NewContainer()

	_ = os.Remove(filepath.Join("configs", "testing", "redis.yaml"))

	client, err := GetRedisClient("redis", c) // create new config file

	fmt.Printf("client == nil: %v\n", client == nil)
	fmt.Printf("err == nil: %v\n", err == nil)

	// Output:
	// client == nil: true
	// err == nil: false
}

func ExampleGetGormDB() {
	c := ioc.NewContainer()

	_ = os.Remove(filepath.Join("configs", "testing", "gorm.yaml"))

	db, err := GetGormDB("gorm", c) // create new config file

	fmt.Printf("db == nil: %v\n", db == nil)
	fmt.Printf("err == nil: %v\n", err == nil)

	// Output:
	// db == nil: true
	// err == nil: false
}

func ExampleGetEmailSender() {
	c := ioc.NewContainer()

	_ = os.Remove(filepath.Join("configs", "testing", "email.yaml"))

	sender, err := GetEmailSender("email", c) // create new config file

	if err != nil {
		panic(err)
	}

	err = sender.Send("user@gmail.com", map[string]interface{}{
		"code":            "123456",
		"expireInMinutes": 5,
	})

	fmt.Printf("err == nil: %v\n", err == nil)

	// Output:
	// err == nil: false
}
