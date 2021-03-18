package metrics

import "github.com/prometheus/client_golang/prometheus"

var MetricDescriptions = make(map[string]*prometheus.Desc)

var (
	Height = newDesc(
		prometheus.BuildFQName("hsc", "chain", "block_height"),
		"Current block height",
		[]string{"chain_id", "nodename"})

	TimePerBlock = newDesc(
		prometheus.BuildFQName("hsc", "chain", "block_time"),
		"Histogram metric of block duration",
		[]string{"chain_id", "nodename"})

	UnconfirmedTransactions = newDesc(
		prometheus.BuildFQName("hsc", "transactions", "in_mempool"),
		"Current depth of the mempool",
		[]string{"chain_id", "nodename"})

	TxPerBlock = newDesc(
		prometheus.BuildFQName("hsc", "transactions", "per_block"),
		"Histogram metric of transactions per block",
		[]string{"chain_id", "nodename"})

	TotalPeers = newDesc(
		prometheus.BuildFQName("hsc", "peers", "total"),
		"Current total peers",
		[]string{"chain_id", "nodename"})

	InboundPeers = newDesc(
		prometheus.BuildFQName("hsc", "peers", "inbound"),
		"Current inbound peers",
		[]string{"chain_id", "nodename"})

	OutboundPeers = newDesc(
		prometheus.BuildFQName("hsc", "peers", "outbound"),
		"Current outbound peers",
		[]string{"chain_id", "nodename"})

	Contracts = newDesc(
		prometheus.BuildFQName("hsc", "accounts", "contracts"),
		"Current contracts on the chain",
		[]string{"chain_id", "nodename"})

	Users = newDesc(
		prometheus.BuildFQName("hsc", "accounts", "users"),
		"Current users on the chain",
		[]string{"chain_id", "nodename"})
)

func newDesc(fqName, help string, variableLabels []string) *prometheus.Desc {
	desc := prometheus.NewDesc(fqName, help, variableLabels, nil)
	MetricDescriptions[fqName] = desc
	return desc
}
