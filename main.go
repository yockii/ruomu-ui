package main

import (
	"github.com/yockii/ruomu-core/config"
	"github.com/yockii/ruomu-core/database"
	"github.com/yockii/ruomu-core/shared"
	"github.com/yockii/ruomu-core/util"
	"github.com/yockii/ruomu-ui/constant"
	"github.com/yockii/ruomu-ui/controller"
	"github.com/yockii/ruomu-ui/model"
)

type UiCore struct {
	shared.UnimplementedCommunicateServer
}

func (UiCore) Initial(params map[string]string) error {
	for key, value := range params {
		config.Set(key, value)
	}

	database.Initial()
	// 同步表结构
	_ = database.DB.AutoMigrate(
		&model.Project{},
		&model.Page{},
		&model.MaterialLib{},
		&model.MaterialLibVersion{},
		&model.MaterialComponent{},
		&model.MaterialComponentGroup{},
		&model.MaterialComponentVersion{},
	)

	_ = util.InitNode(1)

	//TODO 初始化页面

	return nil
}
func (UiCore) InjectCall(code string, headers map[string][]string, value []byte) ([]byte, error) {
	return controller.Dispatch(code, headers, value)
}

func init() {
	config.Set("moduleName", constant.ModuleName)
	config.Set("logger.level", "debug")
	config.InitialLogger()
}
func main() {
	defer database.Close()

	shared.ModuleServe(constant.ModuleName, &UiCore{})
}
