package river

import (
	"bytes"
	"html/template"
	"io"
	"io/ioutil"
	"strings"
	"sync"
)

type ViewConfig struct {
	Prefix     string
	Suffix     string
	Dir        string
	DelimLeft  string
	DelimRight string
	funcMap    map[string]interface{}
	Cache      bool
}

func (conf *ViewConfig) AddFunc(name string, fun interface{}) {
	conf.funcMap[name] = fun
}

func (conf *ViewConfig) GetViewFileName(name string) string {
	return strings.Join([]string{conf.Dir, conf.Prefix, name, conf.Suffix}, "")
}

type ViewEngine interface {
	Render(wr io.Writer, view string, data interface{})
}

type GoHtmlViewEngine struct {
	Config        *ViewConfig
	cacheTemplates map[string]*template.Template
	sync.Mutex
}

func (engine *GoHtmlViewEngine) newTemplate(fileName string) (*template.Template,error)  {
	text, textErr := ioutil.ReadFile(fileName)
	if textErr != nil {
		return nil,textErr
	}
	tpl, tplErr := template.New(fileName).
		Delims(engine.Config.DelimLeft, engine.Config.DelimRight).
		Funcs(engine.Config.funcMap).
		Parse(string(text))

	if tplErr != nil {
		return nil,tplErr
	}
	return tpl,nil
}

func (engine *GoHtmlViewEngine) GetTemplate(fileName string) (*template.Template,error) {
	engine.Lock()
	defer engine.Unlock()
	if engine.Config.Cache {
		tpl,exists :=engine.cacheTemplates[fileName]
		if exists {
			return tpl,nil
		}
		tpl,err := engine.newTemplate(fileName)
		if err != nil {
			return nil,err
		}
		engine.cacheTemplates[fileName] = tpl
		return tpl,nil
	}
	return engine.newTemplate(fileName)

}

func (engine *GoHtmlViewEngine) Render(wr io.Writer, name string, data interface{}) {

	fileName := engine.Config.GetViewFileName(name)

	tpl, tplErr := engine.GetTemplate(fileName)

	if tplErr != nil {
		_, _ = wr.Write([]byte(tplErr.Error()))
		return
	}
	err := tpl.Execute(wr, data)
	if err != nil {
		_, _ = wr.Write([]byte(err.Error()))
		return
	}

}

func NewGoHtmlViewEngine(config *ViewConfig) ViewEngine {

	engine := &GoHtmlViewEngine{
		Config:config,
		cacheTemplates: map[string]*template.Template{},
	}
	config.AddFunc("include",goHtmlViewEngineIncludeFun(engine))
	return engine
}

func goHtmlViewEngineIncludeFun(engine *GoHtmlViewEngine) func(fileName string,data ...interface{}) interface{} {

	return func(fileName string, data ...interface{}) interface{} {

		var viewData interface{}
		if len(data) > 0 {
			viewData = data[0]
		}
		if strings.HasSuffix(fileName,engine.Config.Suffix) {
			fileName = strings.TrimSuffix(fileName,engine.Config.Suffix)
		}
		tpl,err :=engine.GetTemplate(engine.Config.GetViewFileName(fileName))
		if err != nil {
			return err.Error()
		}
		buffer := new(bytes.Buffer)
		tplErr := tpl.Execute(buffer, viewData)
		if tplErr != nil {
			return  tplErr.Error()
		}
		v,rErr := ioutil.ReadAll(buffer)
		if rErr != nil {
			return  rErr.Error()
		}
		return template.HTML(string(v))
	}

}
