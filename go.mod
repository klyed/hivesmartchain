module hivesmartchain

go 1.12

require (
	github.com/BurntSushi/toml v0.3.1
	github.com/OneOfOne/xxhash v1.2.8
	github.com/alecthomas/jsonschema v0.0.0-20210301060011-54c507b6f074
	github.com/btcsuite/btcd v0.21.0-beta
	github.com/cep21/xdgbasedir v0.0.0-20170329171747-21470bfc93b9
	github.com/cosmos/iavl v0.15.3
	github.com/eapache/channels v1.1.0
	github.com/elgs/gojq v0.0.0-20201120033525-b5293fef2759
	github.com/fatih/color v1.10.0
	github.com/go-interpreter/wagon v0.6.0
	github.com/go-kit/kit v0.10.0
	github.com/go-ozzo/ozzo-validation v3.6.0+incompatible
	github.com/gogo/protobuf v1.3.2
	github.com/golang/protobuf v1.4.3
	github.com/hashicorp/golang-lru v0.5.4
	github.com/howeyc/gopass v0.0.0-20190910152052-7cb4b85ec19c
	github.com/iancoleman/strcase v0.1.3
	github.com/imdario/mergo v0.3.12
	github.com/jawher/mow.cli v1.2.0
	github.com/jmoiron/sqlx v1.3.1
	github.com/klyed/hive-go v0.4.0
	github.com/klyed/hivesmartchain v0.0.0-20210327202532-eb77a8b5c479
	github.com/lib/pq v1.10.0
	github.com/mattn/go-sqlite3 v1.14.6
	github.com/monax/relic v2.0.0+incompatible
	github.com/perlin-network/life v0.0.0-20191203030451-05c0e0f7eaea
	github.com/pkg/errors v0.9.1
	github.com/powerman/rpc-codec v1.2.2 // indirect
	github.com/prometheus/client_golang v1.9.0
	github.com/prometheus/client_model v0.2.0
	github.com/prometheus/common v0.19.0
	github.com/spf13/viper v1.7.1
	github.com/streadway/simpleuuid v0.0.0-20130420165545-6617b501e485
	github.com/stretchr/testify v1.7.0
	github.com/syndtr/goleveldb v1.0.1-0.20200815110645-5c35d600f0ca
	github.com/tendermint/tendermint v0.34.8
	github.com/tendermint/tm-db v0.6.4
	github.com/test-go/testify v1.1.4
	github.com/tmthrgd/go-bitset v0.0.0-20190904054048-394d9a556c05
	github.com/tmthrgd/go-hex v0.0.0-20190904060850-447a3041c3bc
	github.com/xeipuuv/gojsonschema v1.2.0
	github.com/xlab/treeprint v1.1.0
	golang.org/x/crypto v0.0.0-20210317152858-513c2a44f670
	golang.org/x/net v0.0.0-20210316092652-d523dce5a7f4
	google.golang.org/grpc v1.36.0
	gopkg.in/yaml.v2 v2.4.0
)

replace /tendermint/tendermint => /klyed/tendermint@v0.34.7

replace /steem-go/rpc@v0.10.0 => ./hive-go

replace /steem-go/rpc-codec@v0.0.0 => ./hive-go

replace /klyed/hiverpc-go => ./hive-go
