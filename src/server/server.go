package main

import (
	"config"
	"flag"
	"fmt"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	tem "html/template"
	"log"
	"os"
	"strconv"
	"tools"
)

func Add1(a string, b int) int {
	aInt, _ := strconv.Atoi(a)
	return aInt + b
}
func main() {

	configFile := ""
	flag.StringVar(&configFile, "c", "server.conf", "config file")
	flag.Parse()
	//load config
	config.Load(configFile)

	conf := config.Get()

	//init DB
	for _, db := range conf.Dbs {
		err := tools.CreateDbConn(db.Database, db.Host, db.Port, db.User, db.Password, db.MaxConn, db.MaxIdle, db.LifeTime)
		tools.CheckErr(err)
		if db.Trace == 1 {
			tools.GetDb(db.Database).DbMap.TraceOn("[db]", log.New(os.Stdout, "myapp:", log.Lmicroseconds))
		}
	}

	template := tem.FuncMap{"Add1": Add1}
	m := martini.Classic()
	m.Use(martini.Static("public"))
	m.Use(render.Renderer(render.Options{
		Funcs:      []tem.FuncMap{template},
		Directory:  "./templates",     // Specify what path to load the templates from.
		Extensions: []string{".html"}, // Specify extensions to load for templates.
	}))

	//初始化路由
	InitRoute(m)

	m.RunOnAddr(fmt.Sprintf(":%d", config.Get().Http.Port))
}
