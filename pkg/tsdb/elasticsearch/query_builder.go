package elasticsearch

import (
	"bytes"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/grafana/grafana/pkg/components/simplejson"
	"github.com/grafana/grafana/pkg/tsdb"

	"text/template"
)

type TemplateQueryModel struct {
	TimeRange *tsdb.TimeRange
	Model     *ElasticsearchRequestModel
}

var queryTemplate = `
{
  "size": 0,
  "query": {
    "bool": {
      "filter": [
        {
          "range": {{ . | formatTimeRange }}
        },
        {
          "query_string": {
            "analyze_wildcard": true,
            "query": "{{ .Model.Query }}"
          }
        }
      ]
    }
  },
  "aggs": {{ . | formatAggregates }}
}`

func convertTimeToUnixNano(rangeTime string) string {
	if rangeTime == "now" {
		rangeTime = "30s"
	}

	duration, err := time.ParseDuration(fmt.Sprintf("-%s", rangeTime))
	if err != nil {
		return err.Error()
	}

	return strconv.FormatInt(time.Now().Add(duration).UnixNano()/1000/1000, 10)
}

func formatTimeRange(data TemplateQueryModel) string {
	to := convertTimeToUnixNano(data.TimeRange.To)
	from := convertTimeToUnixNano(data.TimeRange.From)

	return fmt.Sprintf(`
    {
      "%s": {
        "gte":"%s",
        "lte":"%s",
        "format":"epoch_millis"
      }
    }`, data.Model.TimeField, from, to)
}

func formatAggregates(data TemplateQueryModel) string {
	aggregates := simplejson.New()

	for _, bAgg := range data.Model.BucketAggregates {
		bucket := simplejson.New()

		bucketAggregates := simplejson.New()
		bucketAggregates.Set("field", bAgg.Field)

		for key, value := range bAgg.Settings {
			if key == "trimEdges" {
				continue
			} else if key == "interval" {
				if value == "auto" {
					value = "5s"
				}
				bucketAggregates.Set(key, value)
			} else {
				bucketAggregates.Set(key, value)
			}
		}

		extendedBounds := simplejson.New()
		extendedBounds.Set("min", convertTimeToUnixNano(data.TimeRange.From))
		extendedBounds.Set("max", convertTimeToUnixNano(data.TimeRange.To))
		bucketAggregates.Set("extended_bounds", extendedBounds.MustMap())

		bucketAggregates.Set("format", "epoch_millis")
		bucket.Set(bAgg.Type, bucketAggregates.MustMap())

		metricAggregates := simplejson.New()
		for _, metric := range data.Model.Metrics {
			metricAggregate := simplejson.New()

			aggregate := simplejson.New()
			for key, value := range metric.Settings {
				aggregate.Set(key, value)
			}

			if metric.PipelineAggregate == "" {
				aggregate.Set("field", metric.Field)
			} else {
				aggregate.Set("buckets_path", metric.Field)
			}

			metricAggregate.Set(metric.Type, aggregate)
			metricAggregates.Set(metric.Id, metricAggregate.MustMap())
		}
		bucket.Set("aggs", metricAggregates.MustMap())

		aggregates.Set(bAgg.Id, bucket.MustMap())
	}

	aggString, err := aggregates.MarshalJSON()
	if err != nil {
		log.Printf("%s %s\n", string(aggString), err.Error())
	}
	return string(aggString)
}

func (model *ElasticsearchRequestModel) BuildQueryJson(timeRange *tsdb.TimeRange) (string, error) {

	templateQueryModel := TemplateQueryModel{
		TimeRange: timeRange,
		Model:     model,
	}

	funcMap := template.FuncMap{
		"formatTimeRange":  formatTimeRange,
		"formatAggregates": formatAggregates,
	}

	t, err := template.New("elasticsearchQuery").Funcs(funcMap).Parse(queryTemplate)
	if err != nil {
		return "", err
	}

	buffer := bytes.NewBufferString("")
	t.Execute(buffer, templateQueryModel)

	return string(buffer.Bytes()), nil
}
