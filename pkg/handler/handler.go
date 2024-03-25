package handler

import (
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/atompi/budget_exporter/pkg/options"
	csvutil "github.com/atompi/budget_exporter/pkg/util/csv"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

type metrics struct {
	total     *prometheus.GaugeVec
	based     *prometheus.GaugeVec
	increased *prometheus.GaugeVec
	left      *prometheus.GaugeVec
}

func RootHandler(opts options.Options) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"metrics_path": opts.Web.Path})
	}
}

func MetricsHandler(opts options.ScrapeOptions) gin.HandlerFunc {
	reg := prometheus.NewRegistry()
	m := newMetrics(reg)

	go func() {
		for {
			records, err := getCSV(opts.Address)
			if err != nil {
				zap.L().Sugar().Errorf("get csv records failed: %v", err)
				continue
			}
			for _, record := range *records {
				business := record[opts.LabelHeader.Business]
				provider := record[opts.LabelHeader.Provider]
				total, err := strconv.ParseFloat(record[opts.LabelHeader.Total], 64)
				if err != nil {
					zap.L().Sugar().Errorf("parse total failed: %v", err)
					continue
				}
				based, err := strconv.ParseFloat(record[opts.LabelHeader.Based], 64)
				if err != nil {
					zap.L().Sugar().Errorf("parse based failed: %v", err)
					continue
				}
				increased, err := strconv.ParseFloat(record[opts.LabelHeader.Increased], 64)
				if err != nil {
					zap.L().Sugar().Errorf("parse increased failed: %v", err)
					continue
				}
				left, err := strconv.ParseFloat(record[opts.LabelHeader.Left], 64)
				if err != nil {
					zap.L().Sugar().Errorf("parse left failed: %v", err)
					continue
				}
				m.total.WithLabelValues(business, provider).Set(total)
				m.based.WithLabelValues(business, provider).Set(based)
				m.increased.WithLabelValues(business, provider).Set(increased)
				m.left.WithLabelValues(business, provider).Set(left)
			}
			time.Sleep(time.Duration(opts.Interval) * time.Second)
		}
	}()

	h := promhttp.HandlerFor(
		reg,
		promhttp.HandlerOpts{
			Registry:          reg,
			EnableOpenMetrics: true,
		},
	)
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func newMetrics(reg *prometheus.Registry) *metrics {
	m := &metrics{
		total: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "budget_total",
				Help: "Budget total",
			},
			[]string{
				"business",
				"provider",
			},
		),
		based: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "budget_based",
				Help: "Budget based",
			},
			[]string{
				"business",
				"provider",
			}),
		increased: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "budget_increased",
				Help: "Budget increased",
			},
			[]string{
				"business",
				"provider",
			}),
		left: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "budget_left",
				Help: "Budget left",
			},
			[]string{
				"business",
				"provider",
			}),
	}
	reg.MustRegister(m.total, m.based, m.increased, m.left)
	return m
}

func getCSV(filePath string) (records *[]map[string]string, err error) {
	f, err := os.Open(filePath)
	if err != nil {
		return
	}
	defer f.Close()

	r := csvutil.BOMAwareCSVReader(f)
	res, err := r.ReadAll()
	if err != nil {
		return
	}

	data := &res
	records, err = csvutil.DataToMap(data)
	return
}
