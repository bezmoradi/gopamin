package scaffolder

var databaseRecipes = map[string]databaseRecipe{
	"mysql":    {"mysql-repository", true, []string{"github.com/go-sql-driver/mysql", "github.com/google/uuid"}},
	"mariadb":  {"mariadb-repository", true, []string{"github.com/go-sql-driver/mysql", "github.com/google/uuid"}},
	"postgres": {"postgres-repository", true, []string{"github.com/jackc/pgx/v5/pgxpool", "github.com/google/uuid"}},
	"mongodb":  {"mongodb-repository", true, []string{"go.mongodb.org/mongo-driver/v2/mongo", "github.com/google/uuid"}},
	"dynamodb": {"dynamodb-repository", true, []string{
		"github.com/google/uuid",
		"github.com/aws/aws-sdk-go-v2",
		"github.com/aws/aws-sdk-go-v2/aws",
		"github.com/aws/aws-sdk-go-v2/config",
		"github.com/aws/aws-sdk-go-v2/credentials",
		"github.com/aws/aws-sdk-go-v2/service/dynamodb",
		"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue",
	}},
	"sqlite":    {"sqlite-repository", false, []string{"github.com/mattn/go-sqlite3", "github.com/google/uuid"}},
	"badgerdb":  {"badgerdb-repository", false, []string{"github.com/dgraph-io/badger/v4", "github.com/google/uuid"}},
	"cassandra": {"cassandra-repository", true, []string{"github.com/apache/cassandra-gocql-driver/v2", "github.com/google/uuid"}},
	"mock":      {"", false, []string{"github.com/google/uuid"}},
}
