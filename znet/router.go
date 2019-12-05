package znet

import (
	"czinx/ziface"
)

//定义基础路由
type BaseRoute struct{}

func (br *BaseRoute) PreHandle(req ziface.IRequest) {}

func (br *BaseRoute) Handle(req ziface.IRequest) {}

func (br *BaseRoute) PostHandle(req ziface.IRequest) {}
