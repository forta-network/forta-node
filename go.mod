module github.com/forta-network/forta-node

go 1.22

require (
	github.com/creasty/defaults v1.7.0
	github.com/ethereum/go-ethereum v1.13.13
	github.com/fatih/color v1.16.0
	github.com/go-playground/validator/v10 v10.18.0
	github.com/go-redis/redis v6.15.9+incompatible
	github.com/goccy/go-json v0.10.2
	github.com/golang-jwt/jwt/v4 v4.5.0
	github.com/golang/mock v1.6.0
	github.com/golang/protobuf v1.5.3
	github.com/google/uuid v1.6.0
	github.com/gorilla/mux v1.8.1
	github.com/hashicorp/go-multierror v1.1.1
	github.com/ipfs/go-cid v0.4.1
	github.com/ipfs/go-ipfs-api v0.7.0
	github.com/libp2p/go-libp2p v0.33.0
	github.com/nats-io/nats.go v1.33.1
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/opencontainers/image-spec v1.1.0 // indirect
	github.com/patrickmn/go-cache v2.1.0+incompatible
	github.com/rs/cors v1.10.1
	github.com/shopspring/decimal v1.3.1
	github.com/sirupsen/logrus v1.9.3
	github.com/spf13/cobra v1.8.0
	github.com/spf13/viper v1.18.2
	github.com/stretchr/testify v1.8.4
	github.com/tylertreat/BoomFilters v0.0.0-20210315201527-1a82519a3e43
	golang.org/x/sync v0.6.0
	golang.org/x/time v0.5.0
	google.golang.org/grpc v1.62.0
	gopkg.in/yaml.v3 v3.0.1
)

replace github.com/docker/docker => github.com/moby/moby v20.10.25+incompatible

require (
	github.com/docker/docker v1.6.2
	github.com/docker/go-connections v0.5.0
	github.com/forta-network/forta-core-go v0.0.0-20240223172717-b36e5de0c586
	github.com/prometheus/client_golang v1.18.0
	github.com/prometheus/client_model v0.6.0
	github.com/prometheus/common v0.48.0
	golang.org/x/exp v0.0.0-20240222234643-814bf88cf225
)

require (
	bazil.org/fuse v0.0.0-20200117225306-7b5117fecadc // indirect
	github.com/Jorropo/jsync v1.0.1 // indirect
	github.com/Khan/genqlient v0.6.0 // indirect
	github.com/Microsoft/go-winio v0.6.1 // indirect
	github.com/alecthomas/units v0.0.0-20231202071711-9a357b53e9c9 // indirect
	github.com/asaskevich/govalidator v0.0.0-20230301143203-a9d515a09cc2 // indirect
	github.com/benbjohnson/clock v1.3.5 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/bits-and-blooms/bitset v1.13.0 // indirect
	github.com/bits-and-blooms/bloom v2.0.3+incompatible // indirect
	github.com/blang/semver/v4 v4.0.0 // indirect
	github.com/btcsuite/btcd v0.22.0-beta // indirect
	github.com/btcsuite/btcd/btcec/v2 v2.2.0 // indirect
	github.com/cenkalti/backoff/v4 v4.2.1 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/consensys/bavard v0.1.13 // indirect
	github.com/consensys/gnark-crypto v0.12.1 // indirect
	github.com/containerd/cgroups v1.1.0 // indirect
	github.com/coreos/go-systemd/v22 v22.5.0 // indirect
	github.com/crackcomm/go-gitignore v0.0.0-20231225121904-e25f5bc08668 // indirect
	github.com/crate-crypto/go-kzg-4844 v0.7.0 // indirect
	github.com/cskr/pubsub v1.0.2 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/davidlazar/go-crypto v0.0.0-20200604182044-b73af7476f6c // indirect
	github.com/deckarep/golang-set/v2 v2.6.0 // indirect
	github.com/decred/dcrd/dcrec/secp256k1/v4 v4.2.0 // indirect
	github.com/distribution/reference v0.5.0 // indirect
	github.com/docker/distribution v2.8.3+incompatible // indirect
	github.com/docker/go-units v0.5.0 // indirect
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/elastic/gosigar v0.14.2 // indirect
	github.com/ethereum/c-kzg-4844 v0.4.0 // indirect
	github.com/facebookgo/atomicfile v0.0.0-20151019160806-2de1f203e7d5 // indirect
	github.com/flynn/noise v1.1.0 // indirect
	github.com/forta-network/go-multicall v0.0.0-20230701154355-9467c4ddaa83 // indirect
	github.com/francoispqt/gojay v1.2.13 // indirect
	github.com/fsnotify/fsnotify v1.7.0 // indirect
	github.com/gabriel-vasile/mimetype v1.4.3 // indirect
	github.com/go-logr/logr v1.4.1 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-ole/go-ole v1.3.0 // indirect
	github.com/go-openapi/analysis v0.22.2 // indirect
	github.com/go-openapi/errors v0.21.0 // indirect
	github.com/go-openapi/jsonpointer v0.20.2 // indirect
	github.com/go-openapi/jsonreference v0.20.4 // indirect
	github.com/go-openapi/loads v0.21.5 // indirect
	github.com/go-openapi/runtime v0.27.1 // indirect
	github.com/go-openapi/spec v0.20.14 // indirect
	github.com/go-openapi/strfmt v0.22.0 // indirect
	github.com/go-openapi/swag v0.22.9 // indirect
	github.com/go-openapi/validate v0.23.0 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-task/slim-sprig v0.0.0-20230315185526-52ccab3ef572 // indirect
	github.com/godbus/dbus/v5 v5.1.0 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/google/gopacket v1.1.19 // indirect
	github.com/google/pprof v0.0.0-20240207164012-fb44976bdcd5 // indirect
	github.com/gorilla/websocket v1.5.1 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.19.1 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/golang-lru v1.0.2 // indirect
	github.com/hashicorp/golang-lru/v2 v2.0.7 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/holiman/uint256 v1.2.4 // indirect
	github.com/huin/goupnp v1.3.0 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/ipfs/bbloom v0.0.4 // indirect
	github.com/ipfs/boxo v0.18.0 // indirect
	github.com/ipfs/go-bitfield v1.1.0 // indirect
	github.com/ipfs/go-block-format v0.2.0 // indirect
	github.com/ipfs/go-cidutil v0.1.0 // indirect
	github.com/ipfs/go-datastore v0.6.0 // indirect
	github.com/ipfs/go-ds-measure v0.2.0 // indirect
	github.com/ipfs/go-fs-lock v0.0.7 // indirect
	github.com/ipfs/go-ipfs-delay v0.0.1 // indirect
	github.com/ipfs/go-ipfs-ds-help v1.1.1 // indirect
	github.com/ipfs/go-ipfs-pq v0.0.3 // indirect
	github.com/ipfs/go-ipfs-redirects-file v0.1.1 // indirect
	github.com/ipfs/go-ipfs-util v0.0.3 // indirect
	github.com/ipfs/go-ipld-cbor v0.1.0 // indirect
	github.com/ipfs/go-ipld-format v0.6.0 // indirect
	github.com/ipfs/go-ipld-legacy v0.2.1 // indirect
	github.com/ipfs/go-log v1.0.5 // indirect
	github.com/ipfs/go-log/v2 v2.5.1 // indirect
	github.com/ipfs/go-metrics-interface v0.0.1 // indirect
	github.com/ipfs/go-peertaskqueue v0.8.1 // indirect
	github.com/ipfs/go-unixfsnode v1.9.0 // indirect
	github.com/ipfs/kubo v0.26.0 // indirect
	github.com/ipld/go-car/v2 v2.13.1 // indirect
	github.com/ipld/go-codec-dagpb v1.6.0 // indirect
	github.com/ipld/go-ipld-prime v0.21.0 // indirect
	github.com/jackpal/go-nat-pmp v1.0.2 // indirect
	github.com/jbenet/go-temp-err-catcher v0.1.0 // indirect
	github.com/jbenet/goprocess v0.1.4 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/klauspost/compress v1.17.7 // indirect
	github.com/klauspost/cpuid/v2 v2.2.7 // indirect
	github.com/koron/go-ssdp v0.0.4 // indirect
	github.com/leodido/go-urn v1.4.0 // indirect
	github.com/libp2p/go-buffer-pool v0.1.0 // indirect
	github.com/libp2p/go-cidranger v1.1.0 // indirect
	github.com/libp2p/go-doh-resolver v0.4.0 // indirect
	github.com/libp2p/go-flow-metrics v0.1.0 // indirect
	github.com/libp2p/go-libp2p-asn-util v0.4.1 // indirect
	github.com/libp2p/go-libp2p-kad-dht v0.25.2 // indirect
	github.com/libp2p/go-libp2p-kbucket v0.6.3 // indirect
	github.com/libp2p/go-libp2p-pubsub v0.10.0 // indirect
	github.com/libp2p/go-libp2p-pubsub-router v0.6.0 // indirect
	github.com/libp2p/go-libp2p-record v0.2.0 // indirect
	github.com/libp2p/go-libp2p-routing-helpers v0.7.3 // indirect
	github.com/libp2p/go-libp2p-xor v0.1.0 // indirect
	github.com/libp2p/go-msgio v0.3.0 // indirect
	github.com/libp2p/go-nat v0.2.0 // indirect
	github.com/libp2p/go-netroute v0.2.1 // indirect
	github.com/libp2p/go-reuseport v0.4.0 // indirect
	github.com/libp2p/go-yamux/v4 v4.0.1 // indirect
	github.com/libp2p/zeroconf/v2 v2.2.0 // indirect
	github.com/magiconair/properties v1.8.7 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/marten-seemann/tcp v0.0.0-20210406111302-dfbc87cc63fd // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/miekg/dns v1.1.58 // indirect
	github.com/mikioh/tcpinfo v0.0.0-20190314235526-30a79bb1804b // indirect
	github.com/mikioh/tcpopt v0.0.0-20190314235656-172688c1accc // indirect
	github.com/minio/sha256-simd v1.0.1 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/mmcloughlin/addchain v0.4.0 // indirect
	github.com/moby/term v0.5.0 // indirect
	github.com/montanaflynn/stats v0.7.1 // indirect
	github.com/morikuni/aec v1.0.0 // indirect
	github.com/mr-tron/base58 v1.2.0 // indirect
	github.com/multiformats/go-base32 v0.1.0 // indirect
	github.com/multiformats/go-base36 v0.2.0 // indirect
	github.com/multiformats/go-multiaddr v0.12.2 // indirect
	github.com/multiformats/go-multiaddr-dns v0.3.1 // indirect
	github.com/multiformats/go-multiaddr-fmt v0.1.0 // indirect
	github.com/multiformats/go-multibase v0.2.0 // indirect
	github.com/multiformats/go-multicodec v0.9.0 // indirect
	github.com/multiformats/go-multihash v0.2.3 // indirect
	github.com/multiformats/go-multistream v0.5.0 // indirect
	github.com/multiformats/go-varint v0.0.7 // indirect
	github.com/nats-io/nkeys v0.4.7 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	github.com/nxadm/tail v1.4.11 // indirect
	github.com/oklog/ulid v1.3.1 // indirect
	github.com/onsi/ginkgo/v2 v2.15.0 // indirect
	github.com/opencontainers/runtime-spec v1.2.0 // indirect
	github.com/opentracing/opentracing-go v1.2.0 // indirect
	github.com/openzipkin/zipkin-go v0.4.2 // indirect
	github.com/pbnjay/memory v0.0.0-20210728143218-7b4eea64cf58 // indirect
	github.com/pelletier/go-toml/v2 v2.1.1 // indirect
	github.com/petar/GoLLRB v0.0.0-20210522233825-ae3b015fd3e9 // indirect
	github.com/pion/datachannel v1.5.5 // indirect
	github.com/pion/dtls/v2 v2.2.10 // indirect
	github.com/pion/ice/v2 v2.3.14 // indirect
	github.com/pion/interceptor v0.1.25 // indirect
	github.com/pion/logging v0.2.2 // indirect
	github.com/pion/mdns v0.0.12 // indirect
	github.com/pion/randutil v0.1.0 // indirect
	github.com/pion/rtcp v1.2.13 // indirect
	github.com/pion/rtp v1.8.3 // indirect
	github.com/pion/sctp v1.8.12 // indirect
	github.com/pion/sdp/v3 v3.0.6 // indirect
	github.com/pion/srtp/v2 v2.0.18 // indirect
	github.com/pion/stun v0.6.1 // indirect
	github.com/pion/transport/v2 v2.2.4 // indirect
	github.com/pion/turn/v2 v2.1.5 // indirect
	github.com/pion/webrtc/v3 v3.2.28 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/polydawn/refmt v0.89.0 // indirect
	github.com/prometheus/procfs v0.12.0 // indirect
	github.com/quic-go/qpack v0.4.0 // indirect
	github.com/quic-go/quic-go v0.41.0 // indirect
	github.com/quic-go/webtransport-go v0.6.0 // indirect
	github.com/raulk/go-watchdog v1.3.0 // indirect
	github.com/sagikazarmark/locafero v0.4.0 // indirect
	github.com/sagikazarmark/slog-shim v0.1.0 // indirect
	github.com/samber/lo v1.39.0 // indirect
	github.com/shirou/gopsutil v3.21.11+incompatible // indirect
	github.com/showwin/speedtest-go v1.6.10 // indirect
	github.com/sourcegraph/conc v0.3.0 // indirect
	github.com/spaolacci/murmur3 v1.1.0 // indirect
	github.com/spf13/afero v1.11.0 // indirect
	github.com/spf13/cast v1.6.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/subosito/gotenv v1.6.0 // indirect
	github.com/supranational/blst v0.3.11 // indirect
	github.com/tidwall/gjson v1.17.1 // indirect
	github.com/tidwall/match v1.1.1 // indirect
	github.com/tidwall/pretty v1.2.1 // indirect
	github.com/tklauser/go-sysconf v0.3.13 // indirect
	github.com/tklauser/numcpus v0.7.0 // indirect
	github.com/ucarion/urlpath v0.0.0-20200424170820-7ccc79b76bbb // indirect
	github.com/vektah/gqlparser/v2 v2.5.11 // indirect
	github.com/wealdtech/go-ens/v3 v3.6.0 // indirect
	github.com/wealdtech/go-multicodec v1.4.0 // indirect
	github.com/whyrusleeping/base32 v0.0.0-20170828182744-c30ac30633cc // indirect
	github.com/whyrusleeping/cbor v0.0.0-20171005072247-63513f603b11 // indirect
	github.com/whyrusleeping/cbor-gen v0.1.0 // indirect
	github.com/whyrusleeping/chunker v0.0.0-20181014151217-fe64bd25879f // indirect
	github.com/whyrusleeping/go-keyspace v0.0.0-20160322163242-5b898ac5add1 // indirect
	github.com/whyrusleeping/multiaddr-filter v0.0.0-20160516205228-e903e4adabd7 // indirect
	github.com/willf/bitset v1.1.11 // indirect
	github.com/yusufpapurcu/wmi v1.2.4 // indirect
	go.mongodb.org/mongo-driver v1.14.0 // indirect
	go.opencensus.io v0.24.0 // indirect
	go.opentelemetry.io/otel v1.24.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.24.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.24.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp v1.24.0 // indirect
	go.opentelemetry.io/otel/exporters/stdout/stdouttrace v1.24.0 // indirect
	go.opentelemetry.io/otel/exporters/zipkin v1.24.0 // indirect
	go.opentelemetry.io/otel/metric v1.24.0 // indirect
	go.opentelemetry.io/otel/sdk v1.24.0 // indirect
	go.opentelemetry.io/otel/trace v1.24.0 // indirect
	go.opentelemetry.io/proto/otlp v1.1.0 // indirect
	go.uber.org/dig v1.17.1 // indirect
	go.uber.org/fx v1.20.1 // indirect
	go.uber.org/mock v0.4.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.uber.org/zap v1.27.0 // indirect
	go4.org v0.0.0-20230225012048-214862532bf5 // indirect
	golang.org/x/crypto v0.19.0 // indirect
	golang.org/x/mod v0.15.0 // indirect
	golang.org/x/net v0.21.0 // indirect
	golang.org/x/sys v0.17.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	golang.org/x/tools v0.18.0 // indirect
	golang.org/x/xerrors v0.0.0-20231012003039-104605ab7028 // indirect
	gonum.org/v1/gonum v0.14.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20240221002015-b0ce06bbee7c // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240221002015-b0ce06bbee7c // indirect
	google.golang.org/protobuf v1.32.0 // indirect
	gopkg.in/ini.v1 v1.67.0 // indirect
	gotest.tools/v3 v3.4.0 // indirect
	lukechampine.com/blake3 v1.2.1 // indirect
	rsc.io/tmplfunc v0.0.3 // indirect
)
