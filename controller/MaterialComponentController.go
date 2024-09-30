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
	"gorm.io/gorm"
	"sync"

	"github.com/yockii/ruomu-ui/model"
)

var MaterialComponentController = new(materialComponentController)

type materialComponentController struct{}

func (_ *materialComponentController) Add(value []byte) (any, error) {
	instance := new(model.MaterialComponent)
	if err := json.Unmarshal(value, instance); err != nil {
		logger.Errorln(err)
		return &server.CommonResponse{
			Code: server.ResponseCodeParamParseError,
			Msg:  server.ResponseMsgParamParseError,
		}, nil
	}

	// 处理必填
	if instance.Name == "" || instance.LibID == 0 || instance.GroupID == 0 {
		return &server.CommonResponse{
			Code: server.ResponseCodeParamNotEnough,
			Msg:  server.ResponseMsgParamNotEnough + " name / lib id / group id",
		}, nil
	}

	var c int64
	if err := database.DB.Model(&model.MaterialComponent{}).Where(&model.MaterialComponent{LibID: instance.LibID, Name: instance.Name}).Count(&c).Error; err != nil {
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

	// 尝试解析libVersionId
	libVersionId := gjson.GetBytes(value, "libVersionId").Uint()
	if libVersionId > 0 {
		libVersion := new(model.MaterialVersionComponent)
		if err := database.DB.Where(&model.MaterialVersionComponent{ComponentID: instance.ID, LibVersionID: libVersionId}).Attrs(&model.MaterialVersionComponent{
			ID: util.SnowflakeId(),
		}).FirstOrCreate(libVersion); err != nil {
			logger.Errorln(err)
		}
	}

	return &server.CommonResponse{
		Data: instance,
	}, nil
}

func (_ *materialComponentController) Update(value []byte) (any, error) {
	instance := new(model.MaterialComponent)
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

	if err := database.DB.Model(&model.MaterialComponent{ID: instance.ID}).Updates(&model.MaterialComponent{
		LibID:       instance.LibID,
		GroupID:     instance.GroupID,
		Name:        instance.Name,
		Description: instance.Description,
		TagName:     instance.TagName,
		Thumbnail:   instance.Thumbnail,
		Schema:      instance.Schema,
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

func (_ *materialComponentController) Delete(value []byte) (any, error) {
	instance := new(model.MaterialComponent)
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

func (_ *materialComponentController) List(value []byte) (any, error) {
	instance := new(model.MaterialComponent)
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

	condition := &model.MaterialComponent{
		ID:      instance.ID,
		LibID:   instance.LibID,
		GroupID: instance.GroupID,
	}
	if instance.Name != "" {
		tx.Where("name like ?", "%"+instance.Name+"%")
		instance.Name = ""
	}
	if instance.TagName != "" {
		tx.Where("tag_name like ?", "%"+instance.TagName+"%")
		instance.TagName = ""
	}

	var list []*model.MaterialComponent
	var total int64
	err := tx.Omit("schema").Find(&list, condition).Offset(-1).Count(&total).Error
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

func (_ *materialComponentController) ListWithMetaInfo(value []byte) (any, error) {
	instance := new(domain.MaterialComponentCondition)
	if err := json.Unmarshal(value, instance); err != nil {
		logger.Errorln(err)
		return &server.CommonResponse{
			Code: server.ResponseCodeParamParseError,
			Msg:  server.ResponseMsgParamParseError,
		}, nil
	}

	if instance.LibVersionID == 0 {
		return &server.CommonResponse{
			Code: server.ResponseCodeParamNotEnough,
			Msg:  server.ResponseMsgParamNotEnough + " libVersionId",
		}, nil
	}

	sub := database.DB.Model(&model.MaterialVersionComponent{}).Select("component_id").Where(&model.ProjectMaterialLibVersion{LibVersionID: instance.LibVersionID})

	tx := database.DB.Model(&model.MaterialComponent{}).Where("id in (?)", sub)

	//condition := &model.MaterialComponent{
	//	ID:      instance.ID,
	//	LibID:   instance.LibID,
	//	GroupID: instance.GroupID,
	//}
	//if instance.Name != "" {
	//	tx.Where("name like ?", "%"+instance.Name+"%")
	//	instance.Name = ""
	//}
	//if instance.TagName != "" {
	//	tx.Where("tag_name like ?", "%"+instance.TagName+"%")
	//	instance.TagName = ""
	//}

	var list []*model.MaterialComponent
	if err := tx.Find(&list).Error; err != nil {
		logger.Errorln(err)
		return &server.CommonResponse{
			Code: server.ResponseCodeDatabase,
			Msg:  server.ResponseMsgDatabase + err.Error(),
		}, nil
	}

	var result []*domain.Component
	var wg sync.WaitGroup
	for _, item := range list {
		c := &domain.Component{
			MaterialComponent: *item,
		}
		result = append(result, c)
		wg.Add(1)
		go func(dc *domain.Component) {
			defer wg.Done()
			// 把schema 解析成metaInfo
			metaInfo := new(domain.ComponentMetaInfo)
			if err := json.Unmarshal([]byte(dc.Schema), metaInfo); err != nil {
				logger.Errorln(err)
				return
			}
			dc.MetaInfo = *metaInfo
			dc.Schema = ""
		}(c)
	}
	wg.Wait()

	return &server.CommonResponse{
		Data: &server.Paginate{
			Total:  int64(len(list)),
			Offset: -1,
			Limit:  -1,
			Items:  result,
		},
	}, nil
}

func (_ *materialComponentController) Instance(value []byte) (any, error) {
	instance := new(model.MaterialComponent)
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

func (_ *materialComponentController) VersionAdd(value []byte) (any, error) {
	instance := new(model.MaterialVersionComponent)
	if err := json.Unmarshal(value, instance); err != nil {
		logger.Errorln(err)
		return &server.CommonResponse{
			Code: server.ResponseCodeParamParseError,
			Msg:  server.ResponseMsgParamParseError,
		}, nil
	}

	// 处理必填
	if instance.ComponentID == 0 || instance.LibVersionID == 0 {
		return &server.CommonResponse{
			Code: server.ResponseCodeParamNotEnough,
			Msg:  server.ResponseMsgParamNotEnough + " component id / lib id",
		}, nil
	}

	var c int64
	if err := database.DB.Model(&model.MaterialVersionComponent{}).Where(&model.MaterialVersionComponent{ComponentID: instance.ComponentID, LibVersionID: instance.LibVersionID}).Count(&c).Error; err != nil {
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

func (_ *materialComponentController) VersionDelete(value []byte) (any, error) {
	instance := new(model.MaterialVersionComponent)
	if err := json.Unmarshal(value, instance); err != nil {
		logger.Errorln(err)
		return &server.CommonResponse{
			Code: server.ResponseCodeParamParseError,
			Msg:  server.ResponseMsgParamParseError,
		}, nil
	}

	// 处理必填
	if instance.ID == 0 && (instance.ComponentID == 0 || instance.LibVersionID == 0) {
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
