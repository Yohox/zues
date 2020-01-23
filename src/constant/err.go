package constant

import "errors"

var (
	StructError = errors.New("TCP报文有误！")
	BodyOverFlowError = errors.New("报文超过限制！")
)