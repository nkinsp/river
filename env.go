package river

import (
	"bufio"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

type envMap map[string]string

func (env envMap) Profile(filename string) {
	env["profile"] = filename
}

func (env envMap) ActiveProfile(active string) bool {
	return env.GetString("profile.active") == active
}

func (env envMap) Set(name string, v string) {
	env[name] = v
}

func (env envMap) Get(name string) string {
	return env[name]
}

func (env envMap) GetInt(name string) int {
	v, err := strconv.Atoi(env.GetString(name))
	if err != nil {
		return 0
	}
	return v
}

func (env envMap) GetString(name string) string {
	return env[name]
}

func (env envMap) GetBool(name string) bool {
	v, has := env[name]
	if !has {
		return false
	}
	b,err := strconv.ParseBool(v)
	if err != nil {
		return  false
	}
	return b

}

func (env envMap) Has(name string) bool {
	_, has := env[name]
	return has
}

var Env = &envMap{
	"profile.name":     "./application.conf",
	"profile.active":   "default",
	"server.addr":      ":8080",
	"static.enable":    "false",
	"static.dir":       "./static",
	"static.pattern":   "\\.[png|jpg|jpeg|gif|txt|html|js|css|ico]",
	"view.enable":      "false",
	"view.dir":         "./",
	"view.prefix":      "views/",
	"view.suffix":      ".html",
	"view.delimLeft":   "{{",
	"view.delimRight":  "}}",
	"view.cache":       "false",
	"session.enable":   "false",
	"session.name":     "session",
	"session.expire":   "1800",
	"session.domain":   "",
	"session.httpOnly": "true",
}


func loadConfigFile(filename string) {
	f, err := os.Open(filename)
	if err != nil {
		return
	}
	defer f.Close()
	br := bufio.NewReader(f)
	for {
		line, _, err := br.ReadLine()
		if err == io.EOF {
			break
		}
		lineString := strings.TrimSpace(string(line))
		if len(lineString) > 0 && !strings.HasPrefix(lineString,"##") {
		   lineAttr := strings.Split(lineString,"=")
		   if len(lineAttr) > 1{
			   key := strings.Replace(lineAttr[0]," ","",-1)
			   value := strings.Replace(lineAttr[1]," ","",-1)
			   Env.Set(key,value)
		   }
		}
	}



}

func initEnv() {

	//加载配置文件
	loadConfigFile(Env.Get("profile.name"))
	if !Env.ActiveProfile("default") {
		ncname := strings.Replace(Env.Get("profile.name"),".conf","-"+Env.Get("profile.active"),-1)
		loadConfigFile(ncname)
	}
	log.Println("[River] ","Active",Env.Get("profile.active"))
	if Env.Has("server.addr") {
		config.Addr = Env.GetString("server.addr")
	}
	if Env.GetBool("static.enable") {
		config.middlebrows = append(config.middlebrows, StaticFileMiddleware(Env.GetString("static.dir"), Env.GetString("static.pattern")))
	}
	if Env.GetBool("session.enable") {
		config.Session.Name = Env.GetString("session.name")
		config.Session.Path = Env.GetString("session.path")
		config.Session.HttpOnly = Env.GetBool("session.httpOnly")
		config.Session.ExpireTime = Env.GetInt("session.expireTime")
		config.Session.Domain = Env.GetString("session.domain")
		config.sessionManager = NewMemorySessionManager()
	}
	if Env.GetBool("view.enable") {
		config.viewConfig = &ViewConfig{
			Dir:        Env.GetString("view.dir"),
			Prefix:     Env.GetString("view.prefix"),
			Suffix:     Env.GetString("view.suffix"),
			DelimLeft:  Env.GetString("view.delimLeft"),
			DelimRight: Env.GetString("view.delimRight"),
			funcMap:    map[string]interface{}{},
			Cache:      Env.GetBool("view.cache"),
		}
		config.viewEngine = NewGoHtmlViewEngine(config.viewConfig)
	}

}
