package main

import (
	"encoding/json"
	logger "github.com/sirupsen/logrus"
	"github.com/yockii/ruomu-core/database"
	"github.com/yockii/ruomu-core/util"
	"github.com/yockii/ruomu-ui/domain"
	"github.com/yockii/ruomu-ui/model"
)

func persistNaiveUiLib() uint64 {
	// 创建物料库
	naiveUILib := &model.MaterialLib{
		Code: "naiveUi",
	}
	{
		database.DB.Where(&model.MaterialLib{
			Code: "naiveUi",
		}).Attrs(&model.MaterialLib{
			ID:           util.SnowflakeId(),
			Code:         naiveUILib.Code,
			Name:         "NaiveUI",
			Description:  "NaiveUI物料库",
			PackageName:  "naive-ui",
			Website:      "https://www.naiveui.com/",
			ThumbnailUrl: "https://www.naiveui.com/assets/naivelogo-BdDVTUmz.svg",
		}).FirstOrCreate(naiveUILib)
	}

	// version
	libVersion := &model.MaterialLibVersion{}
	{
		database.DB.Where(&model.MaterialLibVersion{
			LibID:   naiveUILib.ID,
			Version: "2.39.0",
		}).Attrs(&model.MaterialLibVersion{
			ID:            util.SnowflakeId(),
			LibID:         naiveUILib.ID,
			Version:       "2.39.0",
			PluginUseName: "naive",
			CdnJsUrl:      "https://unpkg.com/naive-ui@2.39.0/dist/index.js",
		}).FirstOrCreate(libVersion)

		// 更新默认版本
		database.DB.Model(model.MaterialLib{}).Where("id = ?", naiveUILib.ID).Update("active_version_id", libVersion.ID)
	}

	// 创建组
	{
		for _, group := range NaiveUiLibGroups {
			modelGroup := new(model.MaterialComponentGroup)
			database.DB.Where(&model.MaterialComponentGroup{
				LibCode: naiveUILib.Code,
				Name:    group.Name,
			}).Attrs(&model.MaterialComponentGroup{
				ID:          util.SnowflakeId(),
				LibID:       naiveUILib.ID,
				LibCode:     naiveUILib.Code,
				Name:        group.Name,
				Description: group.Description,
			}).FirstOrCreate(modelGroup)

			for _, component := range group.Components {
				modelComponent := new(model.MaterialComponent)
				// metaInfo 转string作为schema
				schemaBs, err := json.Marshal(component.MetaInfo)
				if err != nil {
					logger.Error("NaiveUiLibGroup Component MetaInfo Marshal Error", err)
				}
				schema := string(schemaBs)
				database.DB.Where(&model.MaterialComponent{
					LibID:   naiveUILib.ID,
					TagName: component.TagName,
				}).Attrs(&model.MaterialComponent{
					ID:          util.SnowflakeId(),
					LibID:       naiveUILib.ID,
					GroupID:     modelGroup.ID,
					TagName:     component.TagName,
					Name:        component.Name,
					Description: component.Description,
					Schema:      schema,
				}).FirstOrCreate(modelComponent)

				// component-version
				database.DB.Where(&model.MaterialComponentVersion{
					ComponentID:  modelComponent.ID,
					LibVersionID: libVersion.ID,
				}).Attrs(&model.MaterialComponentVersion{
					ID: util.SnowflakeId(),
				}).FirstOrCreate(new(model.MaterialComponentVersion))
			}
		}
	}

	return naiveUILib.ID
}

var NaiveUiLibGroups = []domain.Group{
	{
		MaterialComponentGroup: model.MaterialComponentGroup{
			Name:        "通用",
			Description: "15个通用组件",
		},
		Components: []domain.Component{
			{
				MaterialComponent: model.MaterialComponent{
					Name:        "按钮",
					TagName:     "NButton",
					Description: "按钮NButton",
				},
				MetaInfo: domain.ComponentMetaInfo{
					IsContainer: true,
					Props: []domain.PropGroup{
						{
							GroupName: "基础属性",
							Properties: []domain.Property{
								{
									Label:        "块级显示",
									Description:  "按钮是否显示为块级",
									Name:         "block",
									Type:         "boolean",
									DefaultValue: false,
									Widget: domain.Widget{
										Component: domain.WidgetComponentSwitch,
										Inline:    true,
									},
								},
								{
									Label:        "边框",
									Description:  "按钮是否显示 border",
									Name:         "bordered",
									Type:         "boolean",
									DefaultValue: false,
									Widget: domain.Widget{
										Component: domain.WidgetComponentSwitch,
										Inline:    true,
									},
								},
								{
									Label:        "圆形",
									Description:  "按钮是否为圆形",
									Name:         "circle",
									Type:         "boolean",
									DefaultValue: false,
									Widget: domain.Widget{
										Component: domain.WidgetComponentSwitch,
										Inline:    true,
									},
								},
								{
									Label:        "颜色",
									Description:  "按钮颜色（支持形如 #FFF， #FFFFFF， yellow，rgb(0, 0, 0) 的颜色）",
									Name:         "color",
									Type:         "string",
									DefaultValue: "",
									Widget: domain.Widget{
										Component: domain.WidgetComponentColor,
										Inline:    true,
									},
								},
								{
									Label:        "虚线边框",
									Description:  "按钮边框是否为虚线",
									Name:         "dashed",
									Type:         "boolean",
									DefaultValue: false,
									Widget: domain.Widget{
										Component: domain.WidgetComponentSwitch,
										Inline:    true,
									},
								},
								{
									Label:        "禁用",
									Description:  "按钮是否禁用",
									Name:         "disabled",
									Type:         "boolean",
									DefaultValue: false,
									Widget: domain.Widget{
										Component: domain.WidgetComponentSwitch,
										Inline:    true,
									},
								},
								{
									Label:        "加载中",
									Description:  "按钮是否为加载中状态",
									Name:         "loading",
									Type:         "boolean",
									DefaultValue: false,
									Widget: domain.Widget{
										Component: domain.WidgetComponentSwitch,
										Inline:    true,
									},
								},
							},
						},
					},
					Events: map[string]domain.Event{
						"onClick": {
							Label:       "点击事件",
							Description: "点击按钮时的事件",
							Name:        "onClick",
							Params:      []string{"mouseEvent"},
						},
					},
					Slots: []string{"default", "icon"},
				},
			},
		},
	},
}
