package scaffolder

var databaseRecipes = map[string]databaseRecipe{
	"mysql":    {"mysql-repository", []string{"github.com/go-sql-driver/mysql", "github.com/google/uuid"}},
	"mariadb":  {"mariadb-repository", []string{"github.com/go-sql-driver/mysql", "github.com/google/uuid"}},
	"postgres": {"postgres-repository", []string{"github.com/jackc/pgx/v5", "github.com/google/uuid"}},
	"mongodb":  {"mongodb-repository", []string{"go.mongodb.org/mongo-driver/v2", "github.com/google/uuid"}},
	"dynamodb": {"dynamodb-repository", []string{
		"github.com/google/uuid",
		"github.com/aws/aws-sdk-go-v2",
		"github.com/aws/aws-sdk-go-v2/config",
		"github.com/aws/aws-sdk-go-v2/credentials",
		"github.com/aws/aws-sdk-go-v2/service/dynamodb",
		"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue",
	}},
	"sqlite":    {"sqlite-repository", []string{"github.com/mattn/go-sqlite3", "github.com/golang-migrate/migrate/v4", "github.com/google/uuid"}},
	"badgerdb":  {"badgerdb-repository", []string{"github.com/dgraph-io/badger/v4", "github.com/google/uuid"}},
	"cassandra": {"cassandra-repository", []string{"github.com/apache/cassandra-gocql-driver/v2", "github.com/google/uuid"}},
	"mock":      {"", []string{"github.com/google/uuid"}},
}
