package river



type IRoute interface {

	Any(path string,methods []string,handler HanderFun) IRoute
	Get(path string,handler HanderFun) IRoute
	Post(path string,handler HanderFun) IRoute
	Put(path string,handler HanderFun)IRoute
	Delete(path string,handler HanderFun)IRoute

}

type Router struct {

}

func (router *Router) Any(path string,method []string,handler HanderFun) Router {

	return  *router;
}





