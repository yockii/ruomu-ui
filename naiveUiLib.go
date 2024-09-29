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
				database.DB.Where(&model.MaterialVersionComponent{
					ComponentID:  modelComponent.ID,
					LibVersionID: libVersion.ID,
				}).Attrs(&model.MaterialVersionComponent{
					ID: util.SnowflakeId(),
				}).FirstOrCreate(new(model.MaterialVersionComponent))
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
			// 按钮
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
					Slots: []domain.SlotInfo{
						{
							Name: "default",
						},
						{
							Name: "icon",
						},
					},
				},
			},
			// 卡片
			{
				MaterialComponent: model.MaterialComponent{
					Name:        "卡片",
					TagName:     "NCard",
					Description: "卡片NCard",
				},
				MetaInfo: domain.ComponentMetaInfo{
					IsContainer: true,
					Props: []domain.PropGroup{
						{
							GroupName: "基础属性",
							Properties: []domain.Property{
								{
									Label:        "标题",
									Description:  "卡片标题",
									Name:         "title",
									Type:         "string",
									DefaultValue: "",
									Widget: domain.Widget{
										Component: domain.WidgetComponentInput,
									},
								},
								{
									Label:        "显示边框",
									Description:  "是否显示卡片边框",
									Name:         "bordered",
									Type:         "boolean",
									DefaultValue: true,
									Widget: domain.Widget{
										Component: domain.WidgetComponentSwitch,
										Inline:    true,
									},
								},
								{
									Label:        "允许关闭",
									Description:  "是否允许关闭",
									Name:         "closable",
									Type:         "boolean",
									DefaultValue: false,
									Widget: domain.Widget{
										Component: domain.WidgetComponentSwitch,
										Inline:    true,
									},
								},
								{
									Label:        "开启嵌入效果",
									Description:  "使用更深的背景色展现嵌入效果，只对亮色主题生效",
									Name:         "embeddable",
									Type:         "boolean",
									DefaultValue: false,
									Widget: domain.Widget{
										Component: domain.WidgetComponentSwitch,
										Inline:    true,
									},
								},
								{
									Label:        "可悬浮",
									Description:  "开启后，鼠标悬浮在卡片上时，会显示阴影效果",
									Name:         "hoverable",
									Type:         "boolean",
									DefaultValue: false,
									Widget: domain.Widget{
										Component: domain.WidgetComponentSwitch,
										Inline:    true,
									},
								},
								{
									Label:        "尺寸",
									Description:  "卡片尺寸",
									Name:         "size",
									Type:         "string",
									DefaultValue: "medium",
									Widget: domain.Widget{
										Component: domain.WidgetComponentSelect,
										Props: map[string]interface{}{
											"options": []map[string]string{
												{
													"label": "small",
													"value": "small",
												},
												{
													"label": "medium",
													"value": "medium",
												},
												{
													"label": "large",
													"value": "large",
												},
												{
													"label": "huge",
													"value": "huge",
												},
											},
										},
									},
								},
							},
						},
					},
					Events: map[string]domain.Event{
						"onClose": {
							Label:       "关闭事件",
							Description: "点击卡片关闭图标时的回调",
							Name:        "onClose",
							Params:      []string{},
						},
					},
					Slots: []domain.SlotInfo{
						{Name: "default"},
						{Name: "header"},
						{Name: "header-extra"},
						{Name: "cover"},
						{Name: "footer"},
						{Name: "action"},
					},
				},
			},
		},
	},
	{
		MaterialComponentGroup: model.MaterialComponentGroup{
			Name:        "数据录入",
			Description: "20个数据录入组件",
		},
		Components: []domain.Component{
			// 表单
			{
				MaterialComponent: model.MaterialComponent{
					Name:        "表单",
					TagName:     "NForm",
					Description: "表单NForm",
				},
				MetaInfo: domain.ComponentMetaInfo{
					IsContainer: true,
					Props: []domain.PropGroup{
						{
							GroupName: "基础属性",
							Properties: []domain.Property{
								{
									Label:        "禁用",
									Description:  "是否禁用表单",
									Name:         "disabled",
									Type:         "boolean",
									DefaultValue: false,
									Widget: domain.Widget{
										Component: domain.WidgetComponentSwitch,
									},
								},
								{
									Label:        "行内表单",
									Description:  "是否行内表单",
									Name:         "inline",
									Type:         "boolean",
									DefaultValue: false,
									Widget: domain.Widget{
										Component: domain.WidgetComponentSwitch,
									},
								},
								{
									Label:        "标签对齐方式",
									Description:  "标签对齐方式",
									Name:         "labelPlacement",
									Type:         "string",
									DefaultValue: "left",
									Widget: domain.Widget{
										Component: domain.WidgetComponentSelect,
										Props: map[string]interface{}{
											"options": []map[string]string{
												{
													"label": "left",
													"value": "left",
												},
												{
													"label": "top",
													"value": "top",
												},
											},
										},
									},
								},
								{
									Label:        "标签宽度",
									Description:  "标签宽度，像素/auto",
									Name:         "labelWidth",
									Type:         "string",
									DefaultValue: "auto",
									Widget: domain.Widget{
										Component: domain.WidgetComponentInput,
									},
								},
								{
									Label:        "标签对齐方式",
									Description:  "标签对齐方式",
									Name:         "labelAlign",
									Type:         "string",
									DefaultValue: "left",
									Widget: domain.Widget{
										Component: domain.WidgetComponentSelect,
										Props: map[string]interface{}{
											"options": []map[string]string{
												{
													"label": "left",
													"value": "left",
												},
												{
													"label": "right",
													"value": "right",
												},
											},
										},
									},
								},
								{
									Label:        "表单尺寸",
									Description:  "表单尺寸",
									Name:         "size",
									Type:         "string",
									DefaultValue: "medium",
									Widget: domain.Widget{
										Component: domain.WidgetComponentSelect,
										Props: map[string]interface{}{
											"options": []map[string]string{
												{
													"label": "small",
													"value": "small",
												},
												{
													"label": "medium",
													"value": "medium",
												},
												{
													"label": "large",
													"value": "large",
												},
											},
										},
									},
								},
								{
									Label:        "是否显示标签",
									Description:  "是否显示标签",
									Name:         "showLabel",
									Type:         "boolean",
									DefaultValue: true,
									Widget: domain.Widget{
										Component: domain.WidgetComponentSwitch,
									},
								},
								{
									Label:        "必填星号位置",
									Description:  "必填星号位置",
									Name:         "requireMarkPlacement",
									Type:         "string",
									DefaultValue: "right",
									Widget: domain.Widget{
										Component: domain.WidgetComponentSelect,
										Props: map[string]interface{}{
											"options": []map[string]string{
												{
													"label": "left",
													"value": "left",
												},
												{
													"label": "right",
													"value": "right",
												},
												{
													"label": "right-hanging",
													"value": "right-hanging",
												},
											},
										},
									},
								},
							},
						},
					},
					Slots: []domain.SlotInfo{{Name: "default"}},
				},
			},
			// 表单项
			{
				MaterialComponent: model.MaterialComponent{
					Name:        "表单项",
					TagName:     "NFormItem",
					Description: "表单项NFormItem",
				},
				MetaInfo: domain.ComponentMetaInfo{
					IsContainer: true,
					Props: []domain.PropGroup{
						{
							GroupName: "基础属性",
							Properties: []domain.Property{
								{
									Label:        "标签",
									Description:  "标签",
									Name:         "label",
									Type:         "string",
									DefaultValue: "",
									Widget: domain.Widget{
										Component: domain.WidgetComponentInput,
									},
								},
								{
									Label:        "标签对齐方式",
									Description:  "标签对齐方式",
									Name:         "labelAlign",
									Type:         "string",
									DefaultValue: "left",
									Widget: domain.Widget{
										Component: domain.WidgetComponentSelect,
										Props: map[string]interface{}{
											"options": []map[string]string{
												{
													"label": "left",
													"value": "left",
												},
												{
													"label": "right",
													"value": "right",
												},
											},
										},
									},
								},
								{
									Label:        "标签宽度",
									Description:  "标签宽度，像素/auto",
									Name:         "labelWidth",
									Type:         "string",
									DefaultValue: "auto",
									Widget: domain.Widget{
										Component: domain.WidgetComponentInput,
									},
								},
								{
									Label:        "标签对齐方式",
									Description:  "标签对齐方式",
									Name:         "labelPlacement",
									Type:         "string",
									DefaultValue: "",
									Widget: domain.Widget{
										Component: domain.WidgetComponentSelect,
										Props: map[string]interface{}{
											"options": []map[string]string{
												{
													"label": "left",
													"value": "left",
												},
												{
													"label": "top",
													"value": "top",
												},
											},
										},
									},
								},
							},
						},
					},
					Slots: []domain.SlotInfo{{Name: "default"}},
				},
			},
			// 文字输入
			{
				MaterialComponent: model.MaterialComponent{
					Name:        "文字输入",
					TagName:     "NInput",
					Description: "文字输入框",
				},
				MetaInfo: domain.ComponentMetaInfo{
					IsContainer: false,
					Props: []domain.PropGroup{
						{
							GroupName: "基础属性",
							Properties: []domain.Property{
								{
									Label:        "合法性校验",
									Description:  "校验当前的输入是否合法，如果返回 false 输入框便不会响应此次的输入",
									Name:         "allowInput",
									Type:         "function",
									DefaultValue: true,
									Widget: domain.Widget{
										Component: domain.WidgetComponentFunction,
										Props: map[string]interface{}{
											"params": []string{"value"},
										},
									},
								},
								{
									Label:        "自动聚焦",
									Description:  "自动聚焦",
									Name:         "autofocus",
									Type:         "boolean",
									DefaultValue: false,
									Widget: domain.Widget{
										Component: domain.WidgetComponentSwitch,
										Inline:    true,
									},
								},
								{
									Label:        "自适应高度",
									Description:  "是否自动增高输入框高度，只对 type=textarea 有效",
									Name:         "autosize",
									Type:         "boolean",
									DefaultValue: false,
									Widget: domain.Widget{
										Component: domain.WidgetComponentSwitch,
										Inline:    true,
									},
								},
								{
									Label:        "清除按钮",
									Description:  "是否显示清除按钮",
									Name:         "clearable",
									Type:         "boolean",
									DefaultValue: false,
									Widget: domain.Widget{
										Component: domain.WidgetComponentSwitch,
										Inline:    true,
									},
								},
								{
									Label:        "默认值",
									Description:  "默认值",
									Name:         "defaultValue",
									Type:         "string",
									DefaultValue: "",
									Widget: domain.Widget{
										Component: domain.WidgetComponentInput,
									},
								},
								{
									Label:        "计算输入字数",
									Description:  "是否显示输入字数统计",
									Name:         "countGraphemes",
									Type:         "function",
									DefaultValue: false,
									Widget: domain.Widget{
										Component: domain.WidgetComponentFunction,
									},
								},
								{
									Label:        "禁用",
									Description:  "是否禁用",
									Name:         "disabled",
									Type:         "boolean",
									DefaultValue: false,
									Widget: domain.Widget{
										Component: domain.WidgetComponentSwitch,
										Inline:    true,
									},
								},
								{
									Label:       "最大长度",
									Description: "最大输入长度",
									Name:        "maxlength",
									Type:        "number",
									Widget: domain.Widget{
										Component: domain.WidgetComponentNumber,
									},
								},
								{
									Label:       "最小长度",
									Description: "最小输入长度",
									Name:        "minlength",
									Type:        "number",
									Widget: domain.Widget{
										Component: domain.WidgetComponentNumber,
									},
								},
								{
									Label:        "占位符",
									Description:  "占位符",
									Name:         "placeholder",
									Type:         "string",
									DefaultValue: "",
									Widget: domain.Widget{
										Component: domain.WidgetComponentInput,
									},
								},
								{
									Label:        "只读",
									Description:  "是否只读",
									Name:         "readonly",
									Type:         "boolean",
									DefaultValue: false,
									Widget: domain.Widget{
										Component: domain.WidgetComponentSwitch,
										Inline:    true,
									},
								},
								{
									Label:        "圆角",
									Description:  "是否显示边框圆角",
									Name:         "round",
									Type:         "boolean",
									DefaultValue: false,
									Widget: domain.Widget{
										Component: domain.WidgetComponentSwitch,
										Inline:    true,
									},
								},
								{
									Label:        "行数",
									Description:  "输入框行数，对 type=\"textarea\" 有效",
									Name:         "rows",
									Type:         "number",
									DefaultValue: 3,
									Widget: domain.Widget{
										Component: domain.WidgetComponentNumber,
									},
								},
								{
									Label:        "字数统计",
									Description:  "是否显示输入字数统计",
									Name:         "showCount",
									Type:         "boolean",
									DefaultValue: false,
									Widget: domain.Widget{
										Component: domain.WidgetComponentSwitch,
										Inline:    true,
									},
								},
								{
									Label:        "显示密码时机",
									Description:  "密码明文的显示时机",
									Name:         "showPasswordOn",
									Type:         "string",
									DefaultValue: "focus",
									Widget: domain.Widget{
										Component: domain.WidgetComponentSelect,
										Props: map[string]interface{}{
											"options": []map[string]interface{}{
												{
													"label": "click",
													"value": "click",
												},
												{
													"label": "mousedown",
													"value": "mousedown",
												},
											},
										},
									},
								},
								{
									Label:        "尺寸",
									Description:  "输入框尺寸",
									Name:         "size",
									Type:         "string",
									DefaultValue: "medium",
									Widget: domain.Widget{
										Component: domain.WidgetComponentSelect,
										Props: map[string]interface{}{
											"options": []map[string]interface{}{
												{
													"label": "tiny",
													"value": "tiny",
												},
												{
													"label": "small",
													"value": "small",
												},
												{
													"label": "medium",
													"value": "medium",
												},
												{
													"label": "large",
													"value": "large",
												},
											},
										},
									},
								},
								{
									Label:        "类型",
									Description:  "输入框类型",
									Name:         "type",
									Type:         "string",
									DefaultValue: "text",
									Widget: domain.Widget{
										Component: domain.WidgetComponentSelect,
										Props: map[string]interface{}{
											"options": []map[string]interface{}{
												{
													"label": "text",
													"value": "text",
												},
												{
													"label": "textarea",
													"value": "textarea",
												},
												{
													"label": "password",
													"value": "password",
												},
											},
										},
									},
								},
								{
									Label:        "验证状态",
									Description:  "输入框验证状态",
									Name:         "status",
									Type:         "string",
									DefaultValue: "",
									Widget: domain.Widget{
										Component: domain.WidgetComponentSelect,
										Props: map[string]interface{}{
											"options": []map[string]interface{}{
												{
													"label": "error",
													"value": "error",
												},
												{
													"label": "warning",
													"value": "warning",
												},
												{
													"label": "success",
													"value": "success",
												},
											},
										},
									},
								},
								{
									Label:       "文本输入值",
									Description: "文本输入值",
									Name:        "value",
									Type:        "string",
									Widget: domain.Widget{
										Component: domain.WidgetComponentInput,
									},
								},
							},
						},
					},
					Events: map[string]domain.Event{
						"onBlur": {
							Label:       "失去焦点",
							Name:        "onBlur",
							Description: "失去焦点时触发",
							Params:      []string{},
						},
						"onChange": {
							Label:       "值改变",
							Name:        "onChange",
							Description: "值改变时触发",
							Params:      []string{"value"},
						},
						"onClear": {
							Label:       "清空",
							Name:        "onClear",
							Description: "清空时触发",
							Params:      []string{},
						},
						"onFocus": {
							Label:       "获得焦点",
							Name:        "onFocus",
							Description: "获得焦点时触发",
							Params:      []string{},
						},
						"onInput": {
							Label:       "输入时触发",
							Name:        "onInput",
							Description: "输入时触发",
							Params:      []string{"value"},
						},
						"onUpdateValue": {
							Label:       "更新值",
							Name:        "onUpdateValue",
							Description: "更新值时触发",
							Params:      []string{"value"},
						},
					},
					Slots: []domain.SlotInfo{
						{Name: "prefix"},
						{Name: "suffix"},
						{Name: "clear-icon"},
						{Name: "password-invisible-icon"},
						{Name: "password-visible-icon"},
						{Name: "separator"},
						{Name: "count"},
					},
				},
			},
		},
	},
}
