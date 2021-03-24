package models

import (
	"cloudiac/libs/db"
)

const (
	OrgEnable  = "enable"
	OrgDisable = "disable"
	OrgDeleted = "deleted"
)

type Organization struct {
	SoftDeleteModel

	Name         string    `json:"name" gorm:"size:32;not null;comment:'组织名称'"`
	Guid         string    `json:"guid" gorm:"size:32;not null;comment:'组织GUID'"`
	Description  string    `json:"description" gorm:"size:255;comment:'组织描述'"`
	Status       string    `json:"status" gorm:"type:enum('enable','disable','deleted');default:'enable';comment:'组织状态'"`
	VcsType      string    `json:"vcsType" gorm:"type:enum('gitlab','vmware','openstack');default:'gitlab';comment:'vcs类型'"`
	VcsVersion   string    `json:"vcsVersion" gorm:"size:32;comment:'vcs版本'"`
	VcsAuthInfo  string    `json:"vcsAuthInfo" gorm:"size:128;comment:'vcs认证信息'"`
	UserId       uint      `json:"userId" gorm:"not:null;comment:'创建人'"`
}

func (Organization) TableName() string {
	return "iac_org"
}

func (o Organization) Migrate(sess *db.Session) (err error) {
	err = o.AddUniqueIndex(sess, "unique__guid", "guid")
	if err != nil {
		return err
	}

	return nil
}
