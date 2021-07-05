package fastlystats

type NRMetricType string

const (
	NRCount        = NRMetricType("count")
	NRDistribution = NRMetricType("distribution")
	NRGauge        = NRMetricType("gauge")
	NRSummary      = NRMetricType("summary")
	NRUniqueCount  = NRMetricType("uniqueCount")
)

type NewRelicMetricDescriptor struct {
	Name       string            `json:"name"`
	Value      interface{}       `json:"value"`
	Type       NRMetricType      `json:"type"`
	Timestamp  int64             `json:"timestamp"`
	Attributes map[string]string `json:"attributes"`
}

var NRMetricDescriptors = []NewRelicMetricDescriptor{
	{
		Name:       "requests",
		Type:       NRGauge,
		Attributes: map[string]string{"system": "fastly"},
	},
	{
		Name:       "hits",
		Type:       NRGauge,
		Attributes: map[string]string{"system": "fastly"},
	},
	{
		Name:       "hits_time",
		Type:       NRGauge,
		Attributes: map[string]string{"system": "fastly"},
	},
	{
		Name:       "miss",
		Type:       NRGauge,
		Attributes: map[string]string{"system": "fastly"},
	},
	{
		Name:       "miss_time",
		Type:       NRGauge,
		Attributes: map[string]string{"system": "fastly"},
	},
	{
		Name:       "pass",
		Type:       NRGauge,
		Attributes: map[string]string{"system": "fastly"},
	},
	{
		Name:       "pass_time",
		Type:       NRGauge,
		Attributes: map[string]string{"system": "fastly"},
	},
	{
		Name:       "synth",
		Type:       NRGauge,
		Attributes: map[string]string{"system": "fastly"},
	},
	{
		Name:       "errors",
		Type:       NRGauge,
		Attributes: map[string]string{"system": "fastly"},
	},
	{
		Name:       "restarts",
		Type:       NRGauge,
		Attributes: map[string]string{"system": "fastly"},
	},
	{
		Name:       "hit_ratio",
		Type:       NRGauge,
		Attributes: map[string]string{"system": "fastly"},
	},
	{
		Name:       "bandwidth",
		Type:       NRGauge,
		Attributes: map[string]string{"system": "fastly"},
	},
	{
		Name:       "req_body_bytes",
		Type:       NRGauge,
		Attributes: map[string]string{"system": "fastly"},
	},
	{
		Name:       "req_header_bytes",
		Type:       NRGauge,
		Attributes: map[string]string{"system": "fastly"},
	},
	{
		Name:       "resp_body_bytes",
		Type:       NRGauge,
		Attributes: map[string]string{"system": "fastly"},
	},
	{
		Name:       "resp_header_bytes",
		Type:       NRGauge,
		Attributes: map[string]string{"system": "fastly"},
	},
	{
		Name:       "bereq_body_bytes",
		Type:       NRGauge,
		Attributes: map[string]string{"system": "fastly"},
	},
	{
		Name:       "bereq_header_bytes",
		Type:       NRGauge,
		Attributes: map[string]string{"system": "fastly"},
	},
	{
		Name:       "uncachable",
		Type:       NRGauge,
		Attributes: map[string]string{"system": "fastly"},
	},
	{
		Name:       "pipe",
		Type:       NRGauge,
		Attributes: map[string]string{"system": "fastly"},
	},
	{
		Name:       "tls",
		Type:       NRGauge,
		Attributes: map[string]string{"system": "fastly"},
	},
	{
		Name:       "tls_v10",
		Type:       NRGauge,
		Attributes: map[string]string{"system": "fastly"},
	},
	{
		Name:       "tls_v11",
		Type:       NRGauge,
		Attributes: map[string]string{"system": "fastly"},
	},
	{
		Name:       "tls_v12",
		Type:       NRGauge,
		Attributes: map[string]string{"system": "fastly"},
	},
	{
		Name:       "tls_v13",
		Type:       NRGauge,
		Attributes: map[string]string{"system": "fastly"},
	},
	{
		Name:       "shield",
		Type:       NRGauge,
		Attributes: map[string]string{"system": "fastly"},
	},
	{
		Name:       "shield_resp_body_bytes",
		Type:       NRGauge,
		Attributes: map[string]string{"system": "fastly"},
	},
	{
		Name:       "shield_resp_header_bytes",
		Type:       NRGauge,
		Attributes: map[string]string{"system": "fastly"},
	},
	{
		Name:       "ipv6",
		Type:       NRGauge,
		Attributes: map[string]string{"system": "fastly"},
	},
	{
		Name:       "otfp",
		Type:       NRGauge,
		Attributes: map[string]string{"system": "fastly"},
	},
	{
		Name:       "video",
		Type:       NRGauge,
		Attributes: map[string]string{"system": "fastly"},
	},
	{
		Name: "pci",
		Type: NRGauge,

		Attributes: map[string]string{"system": "fastly"},
	},
	{
		Name:       "log",
		Type:       NRGauge,
		Attributes: map[string]string{"system": "fastly"},
	},
	{
		Name:       "http2",
		Type:       NRGauge,
		Attributes: map[string]string{"system": "fastly"},
	},
	{
		Name:       "waf_logged",
		Type:       NRGauge,
		Attributes: map[string]string{"system": "fastly"},
	},
	{
		Name:       "waf_blocked",
		Type:       NRGauge,
		Attributes: map[string]string{"system": "fastly"},
	},
	{
		Name:       "waf_passed",
		Type:       NRGauge,
		Attributes: map[string]string{"system": "fastly"},
	},
	{
		Name:       "attack_req_body_bytes",
		Type:       NRGauge,
		Attributes: map[string]string{"system": "fastly"},
	},
	{
		Name:       "attack_req_header_bytes",
		Type:       NRGauge,
		Attributes: map[string]string{"system": "fastly"},
	},
	{
		Name:       "attack_resp_synth_bytes",
		Type:       NRGauge,
		Attributes: map[string]string{"system": "fastly"},
	},
	{
		Name:       "imgopto",
		Type:       NRGauge,
		Attributes: map[string]string{"system": "fastly"},
	},
	{
		Name:       "status_200",
		Type:       NRGauge,
		Attributes: map[string]string{"system": "fastly"},
	},
	{
		Name:       "status_204",
		Type:       NRGauge,
		Attributes: map[string]string{"system": "fastly"},
	},
	{
		Name:       "status_206",
		Type:       NRGauge,
		Attributes: map[string]string{"system": "fastly"},
	},
	{
		Name:       "status_301",
		Type:       NRGauge,
		Attributes: map[string]string{"system": "fastly"},
	},
	{
		Name:       "status_302",
		Type:       NRGauge,
		Attributes: map[string]string{"system": "fastly"},
	},
	{
		Name:       "status_304",
		Type:       NRGauge,
		Attributes: map[string]string{"system": "fastly"},
	},
	{
		Name:       "status_400",
		Type:       NRGauge,
		Attributes: map[string]string{"system": "fastly"},
	},
	{
		Name: "status_401",
		Type: NRGauge,

		Attributes: map[string]string{"system": "fastly"},
	},
	{
		Name:       "status_403",
		Type:       NRGauge,
		Attributes: map[string]string{"system": "fastly"},
	},
	{
		Name:       "status_404",
		Type:       NRGauge,
		Attributes: map[string]string{"system": "fastly"},
	},
	{
		Name:       "status_416",
		Type:       NRGauge,
		Attributes: map[string]string{"system": "fastly"},
	},
	{
		Name:       "status_500",
		Type:       NRGauge,
		Attributes: map[string]string{"system": "fastly"},
	},
	{
		Name:       "status_501",
		Type:       NRGauge,
		Attributes: map[string]string{"system": "fastly"},
	},
	{
		Name:       "status_502",
		Type:       NRGauge,
		Attributes: map[string]string{"system": "fastly"},
	},
	{
		Name:       "status_503",
		Type:       NRGauge,
		Attributes: map[string]string{"system": "fastly"},
	},
	{
		Name:       "status_504",
		Type:       NRGauge,
		Attributes: map[string]string{"system": "fastly"},
	},
	{
		Name:       "status_505",
		Type:       NRGauge,
		Attributes: map[string]string{"system": "fastly"},
	},
	{
		Name:       "status_1xx",
		Type:       NRGauge,
		Attributes: map[string]string{"system": "fastly"},
	},
	{
		Name:       "status_2xx",
		Type:       NRGauge,
		Attributes: map[string]string{"system": "fastly"},
	},
	{
		Name: "status_3xx",
		Type: NRGauge,

		Attributes: map[string]string{"system": "fastly"},
	},
	{
		Name:       "status_4xx",
		Type:       NRGauge,
		Attributes: map[string]string{"system": "fastly"},
	},
	{
		Name:       "status_5xx",
		Type:       NRGauge,
		Attributes: map[string]string{"system": "fastly"},
	},
	{
		Name:       "object_size_1k",
		Type:       NRGauge,
		Attributes: map[string]string{"system": "fastly"},
	},
	{
		Name:       "object_size_10k",
		Type:       NRGauge,
		Attributes: map[string]string{"system": "fastly"},
	},
	{
		Name:       "object_size_100k",
		Type:       NRGauge,
		Attributes: map[string]string{"system": "fastly"},
	},
	{
		Name:       "object_size_1m",
		Type:       NRGauge,
		Attributes: map[string]string{"system": "fastly"},
	},
	{
		Name:       "object_size_10m",
		Type:       NRGauge,
		Attributes: map[string]string{"system": "fastly"},
	},
	{
		Name:       "object_size_100m",
		Type:       NRGauge,
		Attributes: map[string]string{"system": "fastly"},
	},
	{
		Name:       "object_size_1g",
		Type:       NRGauge,
		Attributes: map[string]string{"system": "fastly"},
	},
	{
		Name:       "billed_header_bytes",
		Type:       NRGauge,
		Attributes: map[string]string{"system": "fastly"},
	},
	{
		Name:       "billed_body_bytes",
		Type:       NRGauge,
		Attributes: map[string]string{"system": "fastly"},
	},
}
