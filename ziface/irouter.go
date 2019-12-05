package ziface

/*

 */
type IRouter interface {

	//处理conn业务之前的钩子方法
	PreHandle(request IRequest)

	//处理conn业务的方法
	Handle(rquest IRequest)

	//处理conn业务之后的钩子方法
	PostHandle(requesr IRequest)
}
