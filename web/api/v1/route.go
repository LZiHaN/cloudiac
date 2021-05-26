package v1

import (
	"cloudiac/libs/ctrl"
	"cloudiac/web/api/v1/handlers"
	"cloudiac/web/middleware"

	"github.com/gin-gonic/gin"
)

func Register(g *gin.RouterGroup) {
	w := ctrl.GinRequestCtxWrap

	g.Any("/check", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"success": true,
		})
	})
	auth := g.Group("/auth")
	auth.POST("/login", w(handlers.User{}.Login))

	o := g.Group("/", w(middleware.Auth))
	o.GET("/org/search", w(handlers.Organization{}.Search))
	o.GET("/org/detail", w(handlers.Organization{}.Detail))
	o.GET("/user/getUserInfo", w(handlers.User{}.GetUserByToken))
	o.PUT("/user/updateSelf", w(handlers.User{}.Update))
	o.GET("/systemStatus/search", w(handlers.PortalSystemStatusSearch))
	o.PUT("/consulTags/update", w(handlers.ConsulTagUpdate))

	// IaC管理员权限
	sys := g.Group("/", w(middleware.Auth), w(middleware.IsAdmin))
	{
		sys.POST("/org/create", w(handlers.Organization{}.Create))
		sys.PUT("/org/update", w(handlers.Organization{}.Update))
		sys.PUT("/org/changeStatus", w(handlers.Organization{}.ChangeOrgStatus))

		ctrl.Register(sys.Group("system"), &handlers.SystemConfig{})
		ctrl.Register(sys.Group("token"), &handlers.Token{})
	}

	root := g.Group("/", w(middleware.Auth), w(middleware.AuthOrgId))
	{
		owner := root.Group("/", w(middleware.IsOrgOwner))
		owner.GET("/user/search", w(handlers.User{}.Search))
		owner.GET("/user/detail", w(handlers.User{}.Detail))
		owner.POST("/user/create", w(handlers.User{}.Create))
		owner.PUT("/user/removeUserForOrg", w(handlers.User{}.RemoveUserForOrg))
		owner.PUT("/user/userPassReset", w(handlers.User{}.UserPassReset))

		root.PUT("/user/update", w(handlers.User{}.Update))

		root.GET("/gitlab/listRepos", w(handlers.GitLab{}.ListRepos))
		root.GET("/gitlab/listBranches", w(handlers.GitLab{}.ListBranches))
		root.GET("/gitlab/getReadme", w(handlers.GitLab{}.GetReadmeContent))
		ctrl.Register(root.Group("notification"), &handlers.Notification{})
		ctrl.Register(root.Group("resourceAccount"), &handlers.ResourceAccount{})
		ctrl.Register(root.Group("template"), &handlers.Template{})
		ctrl.Register(root.Group("task"), &handlers.Task{})
		ctrl.Register(root.Group("taskComment"), &handlers.TaskComment{})
		ctrl.Register(root.Group("webhook"), &handlers.AccessToken{})
		ctrl.Register(root.Group("vcs"), &handlers.Vcs{})

		root.GET("/template/overview", w(handlers.Template{}.Overview))
		root.GET("/template/stateSearch", w(handlers.Template{}.Overview))
		root.GET("/task/last", w(handlers.Task{}.LastTask))

		root.GET("/consulKv/search", w(handlers.ConsulKVSearch))
		root.GET("/runnerList/search", w(handlers.RunnerListSearch))
		root.GET("/templateTfvars/search", w(handlers.TemplateTfvarsSearch))
		root.GET("/vcs/listEnableVcs", w(handlers.ListEnableVcs))

	}

	root.GET("/sse/hello/:filename", w(handlers.HelloSse))
	root.GET("/sse/test", w(handlers.TestSSE))
	g.GET("/taskLog/sse", w(handlers.TaskLogSSE))
}
