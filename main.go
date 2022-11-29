package main

import (
	"github.com/yockii/ruomu-core/config"
	"github.com/yockii/ruomu-core/database"
)

type UiCore struct{}

func (UiCore) Initial(params map[string]string) error {
	for key, value := range params {
		config.Set(key, value)
	}

	database.Initial()
	//TODO 同步表结构
	database.DB.Sync2()

	//TODO 初始化页面

	return nil
}
func (UiCore) InjectCall(code string, headers map[string]string, value []byte) ([]byte, error) {
	// TODO 注入点调用
	return nil, nil
}
func main() {

}
