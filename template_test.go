package gorm_template

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"testing"
)

func TestXml(t *testing.T) {

	dsn := "root:root@(127.0.0.1:3306)/test"

	cfg := &gorm.Config{}
	cfg.Logger = logger.Default
	cfg.Logger = cfg.Logger.LogMode(logger.Info)

	e, err := NewGormEngine(mysql.Open(dsn), cfg, "sql")
	if err != nil {
		panic(err)
	}

	m := map[string]interface{}{}
	err = e.QueryTpl("test.123", map[string]string{
		"Hello":"World",
	}, &m)
	if err != nil {
		panic(err)
	}

	fmt.Println(m)

}
