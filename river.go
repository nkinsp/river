package river

type Request struct {

}
type HanderFun func()

type RequestBody interface {

}

type App struct {

	Router Router

}

type RouterFun func(router Router)

func (app *App) ConfigRouter(fun RouterFun)  {


}




func Create() int  {

	app := App{Router{}}

	app.ConfigRouter(func(router Router) {
		router.Any("",[]string{"",""}, func() {

		})
	})

	return 0
}