package controller

import (
	"encoding/json"
	"errors"
	logger "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"github.com/yockii/ruomu-core/database"
	"github.com/yockii/ruomu-core/server"
	"github.com/yockii/ruomu-core/util"
	"github.com/yockii/ruomu-ui/domain"
	"github.com/yockii/ruomu-ui/uiutil"
	"gorm.io/gorm"

	"github.com/yockii/ruomu-ui/model"
)

var MaterialLibVersionController = new(materialLibVersionController)

type materialLibVersionController struct{}

func (_ *materialLibVersionController) Add(value []byte) (any, error) {
	instance := new(model.MaterialLibVersion)
	if err := json.Unmarshal(value, instance); err != nil {
		logger.Errorln(err)
		return &server.CommonResponse{
			Code: server.ResponseCodeParamParseError,
			Msg:  server.ResponseMsgParamParseError,
		}, nil
	}

	// 处理必填
	if instance.LibID == 0 || instance.Version == "" || instance.PluginUseName == "" {
		return &server.CommonResponse{
			Code: server.ResponseCodeParamNotEnough,
			Msg:  server.ResponseMsgParamNotEnough + " libId / version / pluginUseName",
		}, nil
	}

	var c int64
	if err := database.DB.Model(&model.MaterialLibVersion{}).Where(&model.MaterialLibVersion{LibID: instance.LibID, Version: instance.Version}).Count(&c).Error; err != nil {
		logger.Errorln(err)
		return &server.CommonResponse{
			Code: server.ResponseCodeDatabase,
			Msg:  server.ResponseMsgDatabase + err.Error(),
		}, nil
	} else if c > 0 {
		return &server.CommonResponse{
			Code: server.ResponseCodeDuplicated,
			Msg:  server.ResponseMsgDuplicated,
		}, nil
	}

	instance.ID = util.SnowflakeId()
	if err := database.DB.Create(instance).Error; err != nil {
		logger.Errorln(err)
		return &server.CommonResponse{
			Code: server.ResponseCodeDatabase,
			Msg:  server.ResponseMsgDatabase + err.Error(),
		}, nil
	}
	return &server.CommonResponse{
		Data: instance,
	}, nil
}

func (_ *materialLibVersionController) Update(value []byte) (any, error) {
	instance := new(model.MaterialLibVersion)
	if err := json.Unmarshal(value, instance); err != nil {
		logger.Errorln(err)
		return &server.CommonResponse{
			Code: server.ResponseCodeParamParseError,
			Msg:  server.ResponseMsgParamParseError,
		}, nil
	}
	// 处理必填
	if instance.ID == 0 {
		return &server.CommonResponse{
			Code: server.ResponseCodeParamNotEnough,
			Msg:  server.ResponseMsgParamNotEnough + " id",
		}, nil
	}

	if err := database.DB.Model(&model.MaterialLibVersion{ID: instance.ID}).Updates(&model.MaterialLibVersion{
		LibID:         instance.LibID,
		Version:       instance.Version,
		PluginUseName: instance.PluginUseName,
		CdnJsUrl:      instance.CdnJsUrl,
		CdnCssUrl:     instance.CdnCssUrl,
	}).Error; err != nil {
		logger.Errorln(err)
		return &server.CommonResponse{
			Code: server.ResponseCodeDatabase,
			Msg:  server.ResponseMsgDatabase + err.Error(),
		}, nil
	}
	return &server.CommonResponse{
		Data: true,
	}, nil
}

func (_ *materialLibVersionController) Delete(value []byte) (any, error) {
	instance := new(model.MaterialLibVersion)
	if err := json.Unmarshal(value, instance); err != nil {
		logger.Errorln(err)
		return &server.CommonResponse{
			Code: server.ResponseCodeParamParseError,
			Msg:  server.ResponseMsgParamParseError,
		}, nil
	}
	// 处理必填
	if instance.ID == 0 {
		return &server.CommonResponse{
			Code: server.ResponseCodeParamNotEnough,
			Msg:  server.ResponseMsgParamNotEnough + " id",
		}, nil
	}

	if err := database.DB.Where(instance).Delete(instance).Error; err != nil {
		logger.Errorln(err)
		return &server.CommonResponse{
			Code: server.ResponseCodeDatabase,
			Msg:  server.ResponseMsgDatabase + err.Error(),
		}, nil
	}
	return &server.CommonResponse{
		Data: true,
	}, nil
}

func (_ *materialLibVersionController) List(value []byte) (any, error) {
	instance := new(domain.MaterialLibVersionCondition)
	if err := json.Unmarshal(value, instance); err != nil {
		logger.Errorln(err)
		return &server.CommonResponse{
			Code: server.ResponseCodeParamParseError,
			Msg:  server.ResponseMsgParamParseError,
		}, nil
	}

	vj := gjson.ParseBytes(value)
	paginate := new(server.Paginate)
	if vj.Get("limit").Exists() {
		paginate.Limit = int(vj.Get("limit").Int())
	}
	if vj.Get("offset").Exists() {
		paginate.Offset = int(vj.Get("offset").Int())
	}

	if paginate.Limit <= 0 && paginate.Limit != -1 {
		paginate.Limit = 10
	}

	tx := database.DB.Limit(paginate.Limit).Offset(paginate.Offset)

	condition := &model.MaterialLibVersion{
		ID:    instance.ID,
		LibID: instance.LibID,
	}
	if instance.Version != "" {
		tx.Where("version like ?", "%"+instance.Version+"%")
		instance.Version = ""
	}

	var list []*model.MaterialLibVersion
	var total int64
	err := tx.Find(&list, condition).Offset(-1).Count(&total).Error
	if err != nil {
		logger.Errorln(err)
		return &server.CommonResponse{
			Code: server.ResponseCodeDatabase,
			Msg:  server.ResponseMsgDatabase + err.Error(),
		}, nil
	}
	return &server.CommonResponse{
		Data: &server.Paginate{
			Total:  total,
			Offset: paginate.Offset,
			Limit:  paginate.Limit,
			Items:  list,
		},
	}, nil
}

func (_ *materialLibVersionController) Instance(value []byte) (any, error) {
	instance := new(model.MaterialLibVersion)
	if err := json.Unmarshal(value, instance); err != nil {
		logger.Errorln(err)
		return &server.CommonResponse{
			Code: server.ResponseCodeParamParseError,
			Msg:  server.ResponseMsgParamParseError,
		}, nil
	}
	if err := database.DB.Where(instance).Take(instance).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &server.CommonResponse{}, nil
		}
		logger.Errorln(err)
		return &server.CommonResponse{
			Code: server.ResponseCodeDatabase,
			Msg:  server.ResponseMsgDatabase + err.Error(),
		}, nil
	}
	return &server.CommonResponse{
		Data: instance,
	}, nil
}

func (_ *materialLibVersionController) Export(value []byte) (any, error) {
	instance := new(model.MaterialLibVersion)
	if err := json.Unmarshal(value, instance); err != nil {
		logger.Errorln(err)
		return &server.CommonResponse{
			Code: server.ResponseCodeParamParseError,
			Msg:  server.ResponseMsgParamParseError,
		}, nil
	}
	if err := database.DB.Where(instance).Take(instance).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &server.CommonResponse{}, nil
		}
		logger.Errorln(err)
		return &server.CommonResponse{
			Code: server.ResponseCodeDatabase,
			Msg:  server.ResponseMsgDatabase + err.Error(),
		}, nil
	}

	lib := new(model.MaterialLib)
	if err := database.DB.Where(&model.MaterialLib{ID: instance.LibID}).Take(lib).Error; err != nil {
		logger.Errorln(err)
		return &server.CommonResponse{
			Code: server.ResponseCodeDatabase,
			Msg:  server.ResponseMsgDatabase + err.Error(),
		}, nil
	}
	exportResult := new(domain.MaterialLibVersionImportExport)
	exportResult.Lib = lib
	exportResult.Version = instance

	if err := database.DB.Where(&model.MaterialComponentGroup{LibID: instance.LibID}).Find(&exportResult.Groups).Error; err != nil {
		logger.Errorln(err)
		return &server.CommonResponse{
			Code: server.ResponseCodeDatabase,
			Msg:  server.ResponseMsgDatabase + err.Error(),
		}, nil
	}

	if err := database.DB.Where(&model.MaterialComponent{LibID: instance.LibID}).Find(&exportResult.Components).Error; err != nil {
		logger.Errorln(err)
		return &server.CommonResponse{
			Code: server.ResponseCodeDatabase,
			Msg:  server.ResponseMsgDatabase + err.Error(),
		}, nil
	}

	// 将结果转换为JSON
	jsonData, err := json.Marshal(exportResult)
	if err != nil {
		logger.Errorln(err)
		return &server.CommonResponse{
			Code: server.ResponseCodeGeneration,
			Msg:  server.ResponseMsgGeneration + err.Error(),
		}, nil
	}
	result, err := uiutil.CompressBytes(jsonData)
	if err != nil {
		logger.Errorln(err)
		return &server.CommonResponse{
			Code: server.ResponseCodeGeneration,
			Msg:  server.ResponseMsgGeneration + err.Error(),
		}, nil
	}
	return &server.CommonResponse{
		Data: result,
	}, nil
}

func (_ *materialLibVersionController) Import(value []byte) (any, error) {
	decompressed, err := uiutil.DecompressBytes(value)
	if err != nil {
		logger.Errorln(err)
		return &server.CommonResponse{
			Code: server.ResponseCodeParamParseError,
			Msg:  server.ResponseMsgParamParseError + err.Error(),
		}, nil
	}
	if len(decompressed) == 0 {
		return &server.CommonResponse{
			Code: server.ResponseCodeParamParseError,
			Msg:  server.ResponseMsgParamParseError,
		}, nil
	}
	importRes := new(domain.MaterialLibVersionImportExport)
	if err := json.Unmarshal(decompressed, importRes); err != nil {
		logger.Errorln(err)
		return &server.CommonResponse{
			Code: server.ResponseCodeParamParseError,
			Msg:  server.ResponseMsgParamParseError + err.Error(),
		}, nil
	}
	importLib := importRes.Lib
	libInstance := new(model.MaterialLib)
	libInstance.Code = importLib.Code
	// 检查是否存在
	if err := database.DB.Where(libInstance).Take(libInstance).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 不存在就新增
			importLib.ID = util.SnowflakeId()
			if err := database.DB.Create(importLib).Error; err != nil {
				logger.Errorln(err)
				return &server.CommonResponse{
					Code: server.ResponseCodeDatabase,
					Msg:  server.ResponseMsgDatabase + err.Error(),
				}, nil
			}
		} else {
			logger.Errorln(err)
			return &server.CommonResponse{
				Code: server.ResponseCodeDatabase,
				Msg:  server.ResponseMsgDatabase + err.Error(),
			}, nil
		}
	} else {
		// 则将得到的libInstance的ID赋值给importLib
		importLib.ID = libInstance.ID
	}
	// 检查version是否存在
	versionInstance := new(model.MaterialLibVersion)
	versionInstance.LibID = importLib.ID
	versionInstance.Version = importRes.Version.Version
	if err := database.DB.Where(versionInstance).Take(versionInstance).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 不存在就新增
			importRes.Version.ID = util.SnowflakeId()
			if err := database.DB.Create(importRes.Version).Error; err != nil {
				logger.Errorln(err)
				return &server.CommonResponse{
					Code: server.ResponseCodeDatabase,
					Msg:  server.ResponseMsgDatabase + err.Error(),
				}, nil
			}
		} else {
			logger.Errorln(err)
			return &server.CommonResponse{
				Code: server.ResponseCodeDatabase,
				Msg:  server.ResponseMsgDatabase + err.Error(),
			}, nil
		}
	} else {
		importRes.Version.ID = versionInstance.ID
	}

	groupOldIdWithNewId := make(map[uint64]uint64)
	for _, group := range importRes.Groups {
		oldID := group.ID
		groupInstance := &model.MaterialComponentGroup{
			LibID: importLib.ID,
			Name:  group.Name,
		}
		if err := database.DB.Where(groupInstance).Take(groupInstance).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				group.ID = util.SnowflakeId()
				if err := database.DB.Create(group).Error; err != nil {
					logger.Errorln(err)
					return &server.CommonResponse{
						Code: server.ResponseCodeDatabase,
						Msg:  server.ResponseMsgDatabase + err.Error(),
					}, nil
				}
				groupOldIdWithNewId[oldID] = group.ID
			} else {
				logger.Errorln(err)
				return &server.CommonResponse{
					Code: server.ResponseCodeDatabase,
					Msg:  server.ResponseMsgDatabase + err.Error(),
				}, nil
			}
		} else {
			group.ID = groupInstance.ID
			groupOldIdWithNewId[oldID] = group.ID
		}
	}
	for _, component := range importRes.Components {
		componentInstance := &model.MaterialComponent{
			LibID:   importLib.ID,
			TagName: component.TagName,
		}
		if err := database.DB.Where(componentInstance).Take(componentInstance).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				component.ID = util.SnowflakeId()
				component.GroupID = groupOldIdWithNewId[component.GroupID]
				if err := database.DB.Create(component).Error; err != nil {
					logger.Errorln(err)
					return &server.CommonResponse{
						Code: server.ResponseCodeDatabase,
						Msg:  server.ResponseMsgDatabase + err.Error(),
					}, nil
				}
			} else {
				logger.Errorln(err)
				return &server.CommonResponse{
					Code: server.ResponseCodeDatabase,
					Msg:  server.ResponseMsgDatabase + err.Error(),
				}, nil
			}
		} else {
			// 更新
			component.ID = componentInstance.ID
			component.GroupID = groupOldIdWithNewId[component.GroupID]
			if err := database.DB.Model(componentInstance).Updates(component).Error; err != nil {
				logger.Errorln(err)
				return &server.CommonResponse{
					Code: server.ResponseCodeDatabase,
					Msg:  server.ResponseMsgDatabase + err.Error(),
				}, nil
			}
		}
	}

	return &server.CommonResponse{
		Data: true,
	}, nil
}
