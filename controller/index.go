package controller

import (
	"encoding/json"

	"github.com/yockii/ruomu-ui/constant"
)

// Dispatch 注入点
func Dispatch(code string, headers map[string]string, value []byte) ([]byte, error) {
	switch code {
	case constant.InjectCodeMenuAdd:
		return wrapCall(value, MenuController.Add)
	case constant.InjectCodeMenuUpdate:
		return wrapCall(value, MenuController.Update)
	case constant.InjectCodeMenuDelete:
		return wrapCall(value, MenuController.Delete)
	case constant.InjectCodeMenuList:
		return wrapCall(value, MenuController.List)
	case constant.InjectCodePageAdd:
		return wrapCall(value, PageController.Add)
	case constant.InjectCodePageUpdate:
		return wrapCall(value, PageController.Update)
	case constant.InjectCodePageDelete:
		return wrapCall(value, PageController.Delete)
	case constant.InjectCodePageInstance:
		return wrapCall(value, PageController.Instance)
	case constant.InjectCodePageList:
		return wrapCall(value, PageController.List)
	}
	return nil, nil
}

func wrapCall(v []byte, f func([]byte) (any, error)) ([]byte, error) {
	r, err := f(v)
	if err != nil {
		return nil, err
	}
	bs, err := json.Marshal(r)
	return bs, err
}
