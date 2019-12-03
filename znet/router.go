package znet

import (
	"czinx/ziface"
)

type BaseRoute struct{}

func (br *BaseRoute) PreHandle(req ziface.IRequest) {}

func (br *BaseRoute) Handle(req ziface.IRequest) {}

func (br *BaseRoute) PostHandle(req ziface.IRequest) {}
