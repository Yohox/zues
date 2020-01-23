package service

import (
	"strings"
	"zues/src/logger"
)



type FirstShakeParam struct {
	Content string
}

type SecondShakeParam struct {

}

type ThirdShakeParam struct {

}

var supportsProtocols = []string{"ECC"}

func getSupportProtocol(protocols []string) string {
	for _, pro := range protocols {
		for _, supportsProtocol := range supportsProtocols {
			if pro == supportsProtocol {
				return pro
			}
		}
	}
	return ""
}

type shakeError struct {
	error
}

func (p *shakeError) Error() string {
	return "握手失败！"
}

func (p ClientService) FirstShake(param []string) (bool, interface{}, error) {
	pro := getSupportProtocol(param)
	if pro == "" {
		logger.Errorf("不支持的协议！%s", strings.Join(param, ","))
		return false, nil, &shakeError{}
	}
	logger.Debugf("启用握手协议:%s, 正在进行第一次握手", pro)
	return true, pro, nil
}

func (p ClientService) SecondShake(param *SecondShakeParam) (bool, interface{}, error) {
	logger.Debugf("正在进行第二次握手")
	return false, nil, nil
}

func (p ClientService) ThirdShake(param *ThirdShakeParam) (bool, interface{}, error) {
	logger.Debugf("正在进行第三次握手")
	return false, nil ,nil
}
