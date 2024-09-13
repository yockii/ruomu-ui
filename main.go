package main

import (
	logger "github.com/sirupsen/logrus"
	"github.com/yockii/ruomu-core/config"
	"github.com/yockii/ruomu-core/database"
	"github.com/yockii/ruomu-core/shared"
	"github.com/yockii/ruomu-core/util"
	moduleModel "github.com/yockii/ruomu-module/model"
	"github.com/yockii/ruomu-ui/constant"
	"github.com/yockii/ruomu-ui/controller"
	"github.com/yockii/ruomu-ui/model"
	"os"
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
		&model.ProjectMaterialLibVersion{},
	)

	_ = util.InitNode(1)

	return nil
}
func (UiCore) InjectCall(code string, headers map[string][]string, value []byte) ([]byte, error) {
	return controller.Dispatch(code, headers, value)
}

func init() {
	config.AddConfigPath("../conf")
	if err := config.DefaultInstance.ReadInConfig(); err != nil {
		logger.Warnf("No config file: %s ", err)
	}

	config.Set("moduleName", constant.ModuleCode)
	config.Set("logger.level", "debug")
	config.InitialLogger()
}
func main() {
	defer database.Close()

	// 如果微核调用启动，则启动服务监听 微核调用在cmd中带入 --mc 参数

	// 检查是否有启动参数 --mc
	args := os.Args
	runningInMicroCore := false
	for _, arg := range args {
		if arg == "--mc" {
			runningInMicroCore = true
			break
		}
	}
	if runningInMicroCore {
		shared.ModuleServe(constant.ModuleCode, &UiCore{})
	} else {
		registerModule()
		logger.Info("UI模块注册完成")
	}
}

func registerModule() {
	UiCore{}.Initial(map[string]string{})

	// 直接写表数据即可
	m := &moduleModel.Module{
		Name: constant.ModuleName,
	}
	database.DB.Where(&moduleModel.Module{
		Code: constant.ModuleCode,
	}).Attrs(&moduleModel.Module{
		ID:     util.SnowflakeId(),
		Name:   constant.ModuleName,
		Code:   constant.ModuleCode,
		Cmd:    "./plugins/ruomu-ui --mc",
		Status: 1,
	}).FirstOrCreate(m)

	md := new(moduleModel.ModuleDependency)
	database.DB.Where(&moduleModel.ModuleDependency{
		ModuleCode: constant.ModuleCode,
	}).Attrs(&moduleModel.ModuleDependency{
		ID:             util.SnowflakeId(),
		ModuleCode:     constant.ModuleCode,
		DependenceCode: "ruomu-uc",
	}).FirstOrCreate(md)

	// 注入信息
	{
		mjiList := []*moduleModel.ModuleInjectInfo{
			// render
			{
				ID:                util.SnowflakeId(),
				ModuleID:          m.ID,
				Name:              "获取主页面",
				Type:              11,
				InjectCode:        constant.InjectCodeIndexHtml,
				AuthorizationCode: "anon",
			},
			{
				ID:                util.SnowflakeId(),
				ModuleID:          m.ID,
				Name:              "获取配置页面",
				Type:              11,
				InjectCode:        constant.InjectCodeCanvasHtml,
				AuthorizationCode: "page:schema",
			},
			// json
			{
				ID:                util.SnowflakeId(),
				ModuleID:          m.ID,
				Name:              "新增项目",
				Type:              2,
				InjectCode:        constant.InjectCodeProjectAdd,
				AuthorizationCode: "project:add",
			},
			{
				ID:                util.SnowflakeId(),
				ModuleID:          m.ID,
				Name:              "修改项目",
				Type:              3,
				InjectCode:        constant.InjectCodeProjectUpdate,
				AuthorizationCode: "project:update",
			},
			{
				ID:                util.SnowflakeId(),
				ModuleID:          m.ID,
				Name:              "删除项目",
				Type:              4,
				InjectCode:        constant.InjectCodeProjectDelete,
				AuthorizationCode: "project:delete",
			},
			{
				ID:                util.SnowflakeId(),
				ModuleID:          m.ID,
				Name:              "获取项目列表",
				Type:              1,
				InjectCode:        constant.InjectCodeProjectList,
				AuthorizationCode: "project:list",
			},
			{
				ID:                util.SnowflakeId(),
				ModuleID:          m.ID,
				Name:              "获取项目详情",
				Type:              1,
				InjectCode:        constant.InjectCodeProjectInstance,
				AuthorizationCode: "project:instance",
			},
			{
				ID:                util.SnowflakeId(),
				ModuleID:          m.ID,
				Name:              "新增页面",
				Type:              2,
				InjectCode:        constant.InjectCodePageAdd,
				AuthorizationCode: "page:add",
			},
			{
				ID:                util.SnowflakeId(),
				ModuleID:          m.ID,
				Name:              "修改页面",
				Type:              3,
				InjectCode:        constant.InjectCodePageUpdate,
				AuthorizationCode: "page:update",
			},
			{
				ID:                util.SnowflakeId(),
				ModuleID:          m.ID,
				Name:              "删除页面",
				Type:              4,
				InjectCode:        constant.InjectCodePageDelete,
				AuthorizationCode: "page:delete",
			},
			{
				ID:                util.SnowflakeId(),
				ModuleID:          m.ID,
				Name:              "获取页面列表",
				Type:              1,
				InjectCode:        constant.InjectCodePageList,
				AuthorizationCode: "anon",
			},
			{
				ID:                util.SnowflakeId(),
				ModuleID:          m.ID,
				Name:              "获取页面详情",
				Type:              1,
				InjectCode:        constant.InjectCodePageInstance,
				AuthorizationCode: "page:instance",
			},
			{
				ID:                util.SnowflakeId(),
				ModuleID:          m.ID,
				Name:              "获取页面schema",
				Type:              1,
				InjectCode:        constant.InjectCodePageSchema,
				AuthorizationCode: "anon",
			},
			{
				ID:                util.SnowflakeId(),
				ModuleID:          m.ID,
				Name:              "新增组件库",
				Type:              2,
				InjectCode:        constant.InjectCodeMaterialLibAdd,
				AuthorizationCode: "materialLib:add",
			},
			{
				ID:                util.SnowflakeId(),
				ModuleID:          m.ID,
				Name:              "修改组件库",
				Type:              3,
				InjectCode:        constant.InjectCodeMaterialLibUpdate,
				AuthorizationCode: "materialLib:update",
			},
			{
				ID:                util.SnowflakeId(),
				ModuleID:          m.ID,
				Name:              "删除组件库",
				Type:              4,
				InjectCode:        constant.InjectCodeMaterialLibDelete,
				AuthorizationCode: "materialLib:delete",
			},
			{
				ID:                util.SnowflakeId(),
				ModuleID:          m.ID,
				Name:              "获取组件库列表",
				Type:              1,
				InjectCode:        constant.InjectCodeMaterialLibList,
				AuthorizationCode: "materialLib:list",
			},
			{
				ID:                util.SnowflakeId(),
				ModuleID:          m.ID,
				Name:              "获取组件库详情",
				Type:              1,
				InjectCode:        constant.InjectCodeMaterialLibInstance,
				AuthorizationCode: "materialLib:instance",
			},
			{
				ID:                util.SnowflakeId(),
				ModuleID:          m.ID,
				Name:              "新增组件库分组",
				Type:              2,
				InjectCode:        constant.InjectCodeMaterialComponentGroupAdd,
				AuthorizationCode: "materialComponentGroup:add",
			},
			{
				ID:                util.SnowflakeId(),
				ModuleID:          m.ID,
				Name:              "修改组件库分组",
				Type:              3,
				InjectCode:        constant.InjectCodeMaterialComponentGroupUpdate,
				AuthorizationCode: "materialComponentGroup:update",
			},
			{
				ID:                util.SnowflakeId(),
				ModuleID:          m.ID,
				Name:              "删除组件库分组",
				Type:              4,
				InjectCode:        constant.InjectCodeMaterialComponentGroupDelete,
				AuthorizationCode: "materialComponentGroup:delete",
			},
			{
				ID:                util.SnowflakeId(),
				ModuleID:          m.ID,
				Name:              "获取组件库分组列表",
				Type:              1,
				InjectCode:        constant.InjectCodeMaterialComponentGroupList,
				AuthorizationCode: "materialComponentGroup:list",
			},
			{
				ID:                util.SnowflakeId(),
				ModuleID:          m.ID,
				Name:              "获取组件库分组详情",
				Type:              1,
				InjectCode:        constant.InjectCodeMaterialComponentGroupInstance,
				AuthorizationCode: "materialComponentGroup:instance",
			},
			{
				ID:                util.SnowflakeId(),
				ModuleID:          m.ID,
				Name:              "新增组件库组件",
				Type:              2,
				InjectCode:        constant.InjectCodeMaterialComponentAdd,
				AuthorizationCode: "materialComponent:add",
			},
			{
				ID:                util.SnowflakeId(),
				ModuleID:          m.ID,
				Name:              "修改组件库组件",
				Type:              3,
				InjectCode:        constant.InjectCodeMaterialComponentUpdate,
				AuthorizationCode: "materialComponent:update",
			},
			{
				ID:                util.SnowflakeId(),
				ModuleID:          m.ID,
				Name:              "删除组件库组件",
				Type:              4,
				InjectCode:        constant.InjectCodeMaterialComponentDelete,
				AuthorizationCode: "materialComponent:delete",
			},
			{
				ID:                util.SnowflakeId(),
				ModuleID:          m.ID,
				Name:              "获取组件库组件列表",
				Type:              1,
				InjectCode:        constant.InjectCodeMaterialComponentList,
				AuthorizationCode: "materialComponent:list",
			},
			{
				ID:                util.SnowflakeId(),
				ModuleID:          m.ID,
				Name:              "获取组件库组件详情",
				Type:              1,
				InjectCode:        constant.InjectCodeMaterialComponentInstance,
				AuthorizationCode: "materialComponent:instance",
			},
			{
				ID:                util.SnowflakeId(),
				ModuleID:          m.ID,
				Name:              "新增项目关联物料库版本",
				Type:              1,
				InjectCode:        constant.InjectCodeProjectMaterialLibVersionAdd,
				AuthorizationCode: "projectMaterialLibVersion:list",
			},
			{
				ID:                util.SnowflakeId(),
				ModuleID:          m.ID,
				Name:              "修改项目关联物料库版本",
				Type:              2,
				InjectCode:        constant.InjectCodeProjectMaterialLibVersionUpdate,
				AuthorizationCode: "projectMaterialLibVersion:update",
			},
			{
				ID:                util.SnowflakeId(),
				ModuleID:          m.ID,
				Name:              "删除项目关联物料库版本",
				Type:              3,
				InjectCode:        constant.InjectCodeProjectMaterialLibVersionDelete,
				AuthorizationCode: "projectMaterialLibVersion:delete",
			},
			{
				ID:                util.SnowflakeId(),
				ModuleID:          m.ID,
				Name:              "获取项目关联物料库版本列表",
				Type:              1,
				InjectCode:        constant.InjectCodeProjectMaterialLibVersionList,
				AuthorizationCode: "projectMaterialLibVersion:list",
			},
		}

		for _, mji := range mjiList {
			t := new(moduleModel.ModuleInjectInfo)
			database.DB.Where(&moduleModel.ModuleInjectInfo{
				ModuleID:   mji.ModuleID,
				InjectCode: mji.InjectCode,
			}).Attrs(mji).FirstOrCreate(t)
		}
	}

	// 设置参数
	for k, v := range config.GetStringMapString("database") {
		s := &moduleModel.ModuleSettings{
			ID:       util.SnowflakeId(),
			ModuleID: m.ID,
			Code:     "database." + k,
			Value:    v,
		}
		t := new(moduleModel.ModuleSettings)
		database.DB.Where(&moduleModel.ModuleSettings{
			ModuleID: m.ID,
			Code:     s.Code,
		}).Attrs(s).FirstOrCreate(t)
	}

	// 界面数据
	{
		// 创建项目
		project := &model.Project{}
		database.DB.Where(&model.Project{
			Name: project.Name,
		}).Attrs(&model.Project{
			ID:          util.SnowflakeId(),
			Name:        "默认项目",
			Description: "若木平台基础项目",
			Status:      1,
		}).FirstOrCreate(project)

		naiveUiLibId := persistNaiveUiLib()

		// 物料版本
		naiveUiLibVersion := new(model.MaterialLibVersion)
		database.DB.Where(&model.MaterialLibVersion{
			LibID: naiveUiLibId,
		}).Attrs(&model.MaterialLibVersion{
			ID:            util.SnowflakeId(),
			LibID:         naiveUiLibId,
			Version:       "2.39.0",
			PluginUseName: "naive",
			CdnJsUrl:      "https://unpkg.com/naive-ui@2.39.0/dist/index.prod.js",
		}).FirstOrCreate(naiveUiLibVersion)

		// 项目关联物料库版本
		database.DB.Where(&model.ProjectMaterialLibVersion{
			ProjectID:    project.ID,
			LibID:        naiveUiLibId,
			LibVersionID: naiveUiLibVersion.ID,
		}).Attrs(&model.ProjectMaterialLibVersion{
			ID:           util.SnowflakeId(),
			ProjectID:    project.ID,
			LibID:        naiveUiLibId,
			LibVersionID: naiveUiLibVersion.ID,
		}).FirstOrCreate(new(model.ProjectMaterialLibVersion))

		config.Set("project.id", project.ID)
		_ = config.WriteConfig()
	}
}
