package models

import (
	// "gorm.io/gorm"

	"encoding/json"
	"fmt"
	"github.com/go-admin-team/go-admin-core/sdk"
	"github.com/go-admin-team/go-admin-core/sdk/runtime"
	"github.com/go-admin-team/go-admin-core/storage"
	"go-admin/common/models"
)

type SysApi struct {
	Id         int    `json:"id" gorm:"primaryKey;autoIncrement;comment:主键编码"`
	Handle     string `json:"handle" gorm:"type:varchar(128);comment:handle"`
	Title      string `json:"title" gorm:"type:varchar(128);comment:标题"`
	Permission string `json:"permission" gorm:"size:255;"`
	Path       string `json:"path" gorm:"type:varchar(128);comment:地址"`
	Paths      string `json:"paths" gorm:"type:varchar(128);comment:Paths"`
	Action     string `json:"action" gorm:"type:varchar(16);comment:类型"`
	ParentId   int    `json:"parentId" gorm:"comment:按钮id"`
	Sort       int    `json:"sort" gorm:"type:tinyint(4);comment:排序"`
	models.ModelTime
	models.ControlBy
}

func (SysApi) TableName() string {
	return "sys_api"
}

func (e *SysApi) Generate() models.ActiveRecord {
	o := *e
	return &o
}

func (e *SysApi) GetId() interface{} {
	return e.Id
}

func SaveSysApi(message storage.Messager) (err error) {
	var rb []byte
	rb, err = json.Marshal(message.GetValues())
	if err != nil {
		fmt.Errorf("json Marshal error, %s", err.Error())
		return err
	}

	var l runtime.Routers
	err = json.Unmarshal(rb, &l)
	if err != nil {
		fmt.Errorf("json Unmarshal error, %s", err.Error())
		return err
	}
	dbList := sdk.Runtime.GetDb()
	for _, d := range dbList {
		for _, v := range l.List {
			if v.HttpMethod != "HEAD" {
				err := d.Debug().Where(SysApi{Path: v.RelativePath, Action: v.HttpMethod}).
					Attrs(SysApi{Handle: v.Handler}).
					FirstOrCreate(&SysApi{}).Error
				if err != nil {
					err := fmt.Errorf("Models SaveSysApi error: %s \r\n ", err.Error())
					return err
				}
			}
		}
	}
	return nil
}