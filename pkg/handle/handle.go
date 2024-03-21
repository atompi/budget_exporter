package handle

import (
	"net/http"

	"github.com/atompi/budget_exporter/pkg/options"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

func Handle(opts options.Options) {
	http.Handle("/metrics", promhttp.Handler())
	err := http.ListenAndServe(opts.Web.Listen, nil)
	if err != nil {
		zap.L().Sugar().Error(err)
	} else {
		zap.L().Sugar().Infof("exporter startup listening: %s", opts.Web.Listen)
	}
}
