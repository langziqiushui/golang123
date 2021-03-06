package main

import (
	"fmt"
	"os"
	"time"
	"strconv"
    _ "github.com/jinzhu/gorm/dialects/mysql"
    "github.com/asaskevich/govalidator"
	"github.com/jinzhu/gorm"
	"github.com/kataras/iris"
	"github.com/kataras/iris/sessions"
	"github.com/kataras/iris/middleware/logger"
	"golang123/config"
	"golang123/model"
	"golang123/route"
	"golang123/sessmanager"
)

func init() {
	db, err := gorm.Open(config.DBConfig.Dialect, config.DBConfig.URL)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}
	
	if config.ServerConfig.Env == model.DevelopmentMode {
		db.LogMode(true)
	}

	db.DB().SetMaxIdleConns(config.DBConfig.MaxIdleConns);
	db.DB().SetMaxOpenConns(config.DBConfig.MaxOpenConns)

	model.DB = db;

	sess := sessions.New(sessions.Config{
		Cookie: config.ServerConfig.SessionID,
		Expires: time.Minute * time.Duration(config.ServerConfig.SessionTimeout),
	})
	sessmanager.Sess = sess

	govalidator.SetFieldsRequiredByDefault(true)
}

func main() {
	app := iris.New()

	app.Configure(iris.WithConfiguration(iris.Configuration{
		Charset: "UTF-8",	
	}))

	app.Use(logger.New())

	route.Route(app)

	app.OnErrorCode(iris.StatusNotFound, func(ctx iris.Context) {
		ctx.JSON(iris.Map{
			"errNo" : model.ErrorCode.NotFound,
			"msg"   : "Not Found",
			"data"  : iris.Map{},
		})
	})

	app.OnErrorCode(iris.StatusInternalServerError, func(ctx iris.Context) {
		ctx.JSON(iris.Map{
			"errNo" : model.ErrorCode.ERROR,
			"msg"   : "error",
			"data"  : iris.Map{},
		})
	})

	addr := iris.Addr(":" + strconv.Itoa(config.ServerConfig.Port))
	if config.ServerConfig.Env == model.DevelopmentMode {
		app.Run(addr)
	} else {
		app.Run(addr, iris.WithoutVersionChecker)
	}
}
