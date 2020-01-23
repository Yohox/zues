package server

import (
	"zues/src/logger"
)

var InternalServer = &Server{
	Start: func() {
		logger.Infof("正在开启后端服务器")
		logger.Infof("开启后端服务器成功")
	},
	Stop: func() {
		logger.Infof("正在关闭后端服务器")
	},
}
