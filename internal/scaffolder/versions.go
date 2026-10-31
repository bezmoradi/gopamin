package scaffolder

import (
	"fmt"
	"sort"
)

// modulePins is the exact version of every module gopamin adds to a generated
// project's go.mod. Generation runs `go get module@version` (never `-u`), so one
// gopamin release always produces the same dependency graph instead of whatever
// was latest on the day it ran. Modules that release together (the OTel and AWS
// families) are bumped together.
//
// To bump: change the versions here, then generate a cross-section of projects
// and run `go build`, `go vet`, `go test`, `make lint` and `make arch-check` on
// each before releasing.
var modulePins = map[string]string{
	"github.com/joho/godotenv": "v1.5.1",

	"github.com/labstack/echo/v5":         "v5.4.4",
	"github.com/go-chi/chi/v5":            "v5.5.0",
	"github.com/gin-gonic/gin":            "v1.12.0",
	"github.com/julienschmidt/httprouter": "v1.3.0",
	"github.com/gorilla/mux":              "v1.8.1",
	"github.com/graphql-go/graphql":       "v0.8.1",

	"github.com/google/uuid":                      "v1.6.0",
	"github.com/go-sql-driver/mysql":              "v1.10.1",
	"github.com/jackc/pgx/v5":                     "v5.12.0",
	"go.mongodb.org/mongo-driver/v2":              "v2.9.2",
	"github.com/mattn/go-sqlite3":                 "v1.14.52",
	"github.com/golang-migrate/migrate/v4":        "v4.20.1",
	"github.com/dgraph-io/badger/v4":              "v4.9.6",
	"github.com/apache/cassandra-gocql-driver/v2": "v2.1.2",

	"github.com/aws/aws-sdk-go-v2":                                 "v1.47.5",
	"github.com/aws/aws-sdk-go-v2/config":                          "v1.33.13",
	"github.com/aws/aws-sdk-go-v2/credentials":                     "v1.20.13",
	"github.com/aws/aws-sdk-go-v2/service/dynamodb":                "v1.72.3",
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue": "v1.21.17",

	"github.com/sirupsen/logrus": "v1.10.2",
	"go.uber.org/zap":            "v1.28.0",

	"github.com/segmentio/kafka-go":  "v0.4.51",
	"github.com/rabbitmq/amqp091-go": "v1.16.0",
	"github.com/redis/go-redis/v9":   "v9.23.0",

	"go.opentelemetry.io/otel":                                          "v1.47.0",
	"go.opentelemetry.io/otel/sdk":                                      "v1.47.0",
	"go.opentelemetry.io/otel/sdk/metric":                               "v1.47.0",
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc":   "v1.47.0",
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc": "v1.47.0",
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp":     "v0.72.0",
	"go.opentelemetry.io/contrib/instrumentation/runtime":               "v0.72.0",

	"github.com/golang-jwt/jwt/v5": "v5.3.1",
	"github.com/swaggest/swgui":    "v1.8.9",

	"golang.org/x/net":           "v0.62.0",
	"golang.org/x/crypto":        "v0.58.0",
	"golang.org/x/sys":           "v0.49.0",
	"golang.org/x/text":          "v0.43.0",
	"google.golang.org/grpc":     "v1.85.0",
	"google.golang.org/protobuf": "v1.36.12",
	// gin's HTTP/3 support; v0.59.1 is the smallest fix for GO-2026-5676.
	"github.com/quic-go/quic-go": "v0.59.1",
}

// securityFloor lists transitive modules that the pinned direct modules pull
// in at whatever minimum version they declare, which can be one with a known
// vulnerability. Every project raises them to their pinned version; `go mod
// tidy` then drops the ones the project doesn't use.
var securityFloor = []string{
	"golang.org/x/net",
	"golang.org/x/crypto",
	"golang.org/x/sys",
	"golang.org/x/text",
	"google.golang.org/grpc",
	"google.golang.org/protobuf",
	"github.com/quic-go/quic-go",
}

// pinned turns module paths into `module@version` arguments for `go get`. An
// unpinned module is a gopamin bug, so it is reported rather than fetched at
// whatever version happens to be latest.
func pinned(modules []string) ([]string, error) {
	args := make([]string, 0, len(modules))
	var missing []string
	for _, m := range modules {
		v, ok := modulePins[m]
		if !ok {
			missing = append(missing, m)
			continue
		}
		args = append(args, m+"@"+v)
	}
	if len(missing) > 0 {
		sort.Strings(missing)
		return nil, fmt.Errorf("no pinned version for %v (add it to modulePins)", missing)
	}
	return args, nil
}
