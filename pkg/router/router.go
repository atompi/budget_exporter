package router

import (
	"github.com/atompi/budget_exporter/pkg/handler"
	"github.com/atompi/budget_exporter/pkg/options"
	"github.com/gin-gonic/gin"
)

type RouterGroupFunc func(*gin.RouterGroup, options.Options)

func RootRouter(routerGroup *gin.RouterGroup, opts options.Options) {
	routerGroup.GET("", handler.RootHandler(opts))
}

func MetricsRouter(routerGroup *gin.RouterGroup, opts options.Options) {
	routerGroup.GET(opts.Web.Path, handler.MetricsHandler(opts.Scrape))
}

func Register(e *gin.Engine, opts options.Options) {
	rootRouterGroup := e.Group("/")
	routerGroups := []RouterGroupFunc{}
	routerGroups = append(routerGroups, RootRouter, MetricsRouter)

	for _, routerGroup := range routerGroups {
		routerGroup(rootRouterGroup, opts)
	}
}
