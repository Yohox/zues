package server

import (
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"net"
	"reflect"
	"unsafe"
	"zues/src/config"
	"zues/src/logger"
	"zues/src/protocol"
	"zues/src/service"
	"zues/src/utils"
)

var ClientServer = &Server{
	Start: doStart,
	Stop: doStop,
}

const messageSize = unsafe.Sizeof(protocol.ClientRequestMessage{})


var stop = false

func doStart(){
	logger.Infof("正在开启客户端服务器")
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", config.Cfg.ClientConfig.Ip, config.Cfg.ClientConfig.Port))
	if err != nil {
		logger.Errorf("客户端开启失败！", err.Error())
		panic("客户端开启失败！")
	}
	logger.Infof("开启客户端服务器成功")
	for !stop {
		conn, err := listener.Accept()
		if err != nil {
			logger.Errorf("监听失败！", err.Error())
			continue
		}
		go doAccept(conn)
	}
}

func doStop(){
	stop = true
}

func doAccept(conn net.Conn){
	var (
		clientRequestMessage *protocol.ClientRequestMessage
		clientResponseMessage *protocol.ClientResponseMessage
		err error
	)
	for {
		clientRequestMessage, err = protocol.ReadClientRequestMessage(conn)
		if err == io.EOF {
			continue
		} else if err != nil {
			err = errors.New("读取客户端消息头失败！ " + err.Error())
			break
		}
		clientResponseMessage, err = doDispatch(conn, clientRequestMessage)
		if err != nil {
			err = errors.New("分发消息失败！ " + err.Error())
			break
		}
		err = doResponse(conn, clientResponseMessage)
		if err != nil {
			err = errors.New("回复消息失败！ " + err.Error())
			break
		}
	}
	if err != nil {
		log.Errorf(err.Error())
	}
	_ = conn.Close()
}

type DispatchStruct struct {
	Name string
	Param []interface{}
}


func doDispatch(conn net.Conn, requestMessage *protocol.ClientRequestMessage) (*protocol.ClientResponseMessage, error) {
	dispatchStruct, err := getDispatchStruct(requestMessage)
	if err != nil {
		return nil, err
	}
	res, err := invokeService(dispatchStruct.Name, dispatchStruct.Param)
	if err != nil {
		return nil, err
	}
	if requestMessage.ShakeStep > 2 {
		res, err = utils.EncodeBody(res)
		if err != nil {
			return nil, err
		}
	}
	return &protocol.ClientResponseMessage{
		Len: uint32(len(res)),
		Body: res,
	}, nil
}

func getDispatchStruct(requestMessage *protocol.ClientRequestMessage) (*DispatchStruct, error) {
	dispatchStruct := &DispatchStruct{}
	if requestMessage.ShakeStep < 3 {
		switch requestMessage.ShakeStep {
		case 0:
			dispatchStruct.Name = "FirstShake"
		case 1:
			dispatchStruct.Name = "SecondShake"
		case 2:
			dispatchStruct.Name = "ThirdShake"
		default:
			return nil, errors.New("握手函数错误！")
		}
		err := json.Unmarshal([]byte(requestMessage.Body), &dispatchStruct.Param)
		if err != nil {
			return nil, err
		}
	} else {
		decodeBody, err := utils.DecodeBody(requestMessage.Body)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(decodeBody, dispatchStruct)
		if err != nil {
			return nil, err
		}
	}
	return dispatchStruct, nil
}

func invokeService(name string, params []interface{}) ([]byte, error) {
	clientService := service.ClientService{}
	of := reflect.ValueOf(clientService)
	method := of.MethodByName(name)
	if !method.IsValid() {
		return nil, errors.New(fmt.Sprintf("找不到方法：%s", name))
	}
	args, err := parseParam(method, params)
	if err != nil {
		return nil, err
	}
	res := method.Call(args)
	response, err := getClientServiceResponse(res)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("调用函数:%s %s", name, err.Error()))
	}
	bytes, err := json.Marshal(response)
	if err != nil {
		return nil, err
	} else {
		return bytes, nil
	}
}

func getClientServiceResponse(methodRes []reflect.Value) (*service.ClientServiceResponse, error) {
	if len(methodRes) != 3 {
		return nil, errors.New("返回结果数量有误！")
	}
	response := &service.ClientServiceResponse{}
	if methodRes[0].Kind() != reflect.Bool || methodRes[1].Kind() != reflect.Interface || methodRes[2].Kind() != reflect.Interface {
		return nil, errors.New("返回结果类型不符！")
	}
	response.Status = methodRes[0].Bool()
	bytes, err := json.Marshal(methodRes[1].Interface())
	if err != nil {
		return nil, errors.New("序列化结果失败！")
	}
	response.Res = bytes
	response.Message = methodRes[2].String()
	return response, nil
}

func parseParam(method reflect.Value, params []interface{}) ([]reflect.Value, error) {
	methodType := reflect.TypeOf(method.Interface())
	argNum := methodType.NumIn()
	if argNum != len(params) {
		return nil, errors.New("参数数量不一样！")
	}
	res := make([]reflect.Value, argNum)
	for i := 0 ; i < argNum; i++ {
		value, err := solveParam(methodType.In(i), params[i])
		if err != nil {
			return nil, err
		}
		res[i] = value
	}
	return res, nil
}

func solveParam(paramType reflect.Type, param interface{}) (reflect.Value, error) {
	switch paramType.Kind() {
	case reflect.Int:
		i, ok := param.(int)
		if !ok {
			return reflect.ValueOf(0), errors.New("参数类型转换失败！")
		}
		return reflect.ValueOf(i), nil
	case reflect.String:
		i, ok := param.(string)
		if !ok {
			return reflect.ValueOf(""), errors.New("参数类型转换失败！")
		}
		return reflect.ValueOf(i), nil
	case reflect.Slice:
		params, ok := param.([]interface{})
		if !ok {
			return reflect.ValueOf(""), errors.New("参数类型转换失败！")
		}
		slice := reflect.MakeSlice(paramType, len(params), len(params))
		for i:= 0; i < len(params); i++ {
			value, err := solveParam(paramType.Elem(), params[i])
			if err != nil {
				return reflect.ValueOf(""), errors.New("参数类型转换失败！")
			}
			slice.Index(i).Set(value)
		}
		return slice, nil
	default:
		return reflect.ValueOf(""), errors.New("参数类型错误！")
	}
}

func doResponse(conn net.Conn, responseMessage *protocol.ClientResponseMessage) error {
	bytes := append([]byte("head|"), responseMessage.Body...)
	_, err := conn.Write(bytes)
	if err != nil {
		return err
	}
	return nil
}