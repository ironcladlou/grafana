package elasticsearch

import (
	"fmt"
	"testing"

	"github.com/grafana/grafana/pkg/tsdb"
	. "github.com/smartystreets/goconvey/convey"
)

var testResponseJson = `
{
  "took": 1588,
  "timed_out": false,
  "_shards": {
    "total": 2250,
    "successful": 2250,
    "failed": 0
  },
  "hits": {
    "total": 5,
    "max_score": 0,
    "hits": []
  },
  "aggregations": {
    "2": {
      "buckets": [
        {
          "key_as_string": "1487020955000",
          "key": 1487020955000,
          "doc_count": 0,
          "1": {
            "value": null
          }
        },
        {
          "key_as_string": "1487020985000",
          "key": 1487020985000,
          "doc_count": 1,
          "1": {
            "value": 0
          }
        },
        {
          "key_as_string": "1487021010000",
          "key": 1487021010000,
          "doc_count": 0,
          "1": {
            "value": null
          }
        },
        {
          "key_as_string": "1487021045000",
          "key": 1487021045000,
          "doc_count": 1,
          "1": {
            "value": 1234
          },
          "3": {
            "value": 123
          }
        },
        {
          "key_as_string": "1487021105000",
          "key": 1487021105000,
          "doc_count": 1,
          "1": {
            "value": 155
          },
          "3": {
            "value": 0
          }
        },
        {
          "key_as_string": "1487021165000",
          "key": 1487021165000,
          "doc_count": 1,
          "1": {
            "value": 0
          },
          "3": {
            "value": 0
          }
        },
        {
          "key_as_string": "1487021180000",
          "key": 1487021180000,
          "doc_count": 0,
          "1": {
            "value": null
          }
        },
        {
          "key_as_string": "1487021210000",
          "key": 1487021210000,
          "doc_count": 0,
          "1": {
            "value": null
          }
        },
        {
          "key_as_string": "1487021225000",
          "key": 1487021225000,
          "doc_count": 1,
          "1": {
            "value": 0
          },
          "3": {
            "value": 1000
          }
        }
      ]
    }
  }
}`

func TestElasticserachQueryParser(t *testing.T) {
	Convey("Elasticserach QueryBuilder query parsing", t, func() {

		Convey("Parse ElasticSearch Requry Results", func() {
			queryResult, err := parseQueryResult([]byte(testResponseJson))
			So(err, ShouldBeNil)
			So(queryResult, ShouldNotBeNil)

			qR := &tsdb.QueryResult{}
			fmt.Println(qR)
			So(len(queryResult.Series), ShouldEqual, 2)
		})
	})
}
