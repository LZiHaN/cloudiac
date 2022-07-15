// Copyright (c) 2015-2022 CloudJ Technology Co., Ltd.

package apps

import (
	"cloudiac/portal/consts/e"
	"cloudiac/portal/libs/ctx"
	"cloudiac/portal/models/forms"
	"cloudiac/portal/models/resps"
	"cloudiac/portal/services"
)

// PlatformStatTodayOrg 当日新建组织数
func PlatformStatTodayOrg(c *ctx.ServiceContext, form *forms.PfStatForm) (interface{}, e.Error) {
	orgIds := parseOrgIds(form.OrgIds)

	count, err := services.GetTodayCreatedOrgs(c.DB(), orgIds)
	if err != nil {
		return nil, e.AutoNew(err, e.DBError)
	}

	return &resps.PfTodayStatResp{
		Name:  "当日新建组织数",
		Count: count,
	}, nil
}

// PlatformStatTodayProject 当日新建项目数
func PlatformStatTodayProject(c *ctx.ServiceContext, form *forms.PfStatForm) (interface{}, e.Error) {

	return nil, nil
}

// PlatformStatTodayStack 当日新建 Stack 数
func PlatformStatTodayStack(c *ctx.ServiceContext, form *forms.PfStatForm) (interface{}, e.Error) {

	return nil, nil
}

// PlatformStatTodayPG 当日新建合规策略组数量
func PlatformStatTodayPG(c *ctx.ServiceContext, form *forms.PfStatForm) (interface{}, e.Error) {

	return nil, nil
}

// PlatformStatTodayEnv 当日新建环境数
func PlatformStatTodayEnv(c *ctx.ServiceContext, form *forms.PfStatForm) (interface{}, e.Error) {

	return nil, nil
}

// PlatformStatTodayDestroyedEnv 当日销毁环境数
func PlatformStatTodayDestroyedEnv(c *ctx.ServiceContext, form *forms.PfStatForm) (interface{}, e.Error) {

	return nil, nil
}

// PlatformStatTodayResType 当日新建资源数：资源类型、数量
func PlatformStatTodayResType(c *ctx.ServiceContext, form *forms.PfStatForm) (interface{}, e.Error) {

	return nil, nil
}
