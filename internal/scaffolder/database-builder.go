package scaffolder

type mysqlBuilder struct {
	project *Project
}

func (m *mysqlBuilder) build() {
	fileGenerator([]string{"mock-repository"}, m.project)
	fileGenerator([]string{"user-test"}, m.project)
	fileGenerator([]string{"user"}, m.project)
	fileGenerator([]string{"user-repository-interface"}, m.project)
	fileGenerator([]string{"user-service-interface"}, m.project)
	fileGenerator([]string{"user-service"}, m.project)
	fileGenerator([]string{"user-service-test"}, m.project)

	fileGenerator([]string{"mysql-repository"}, m.project)
	fileGenerator([]string{"mysql-docker-compose"}, m.project)
	goGetPackages(m.project.Path, []string{"github.com/go-sql-driver/mysql", "github.com/google/uuid"})
}

type cassandraBuilder struct {
	project *Project
}

func (m *cassandraBuilder) build() {
	fileGenerator([]string{"mock-repository"}, m.project)
	fileGenerator([]string{"user-test"}, m.project)
	fileGenerator([]string{"user"}, m.project)
	fileGenerator([]string{"user-repository-interface"}, m.project)
	fileGenerator([]string{"user-service-interface"}, m.project)
	fileGenerator([]string{"user-service"}, m.project)
	fileGenerator([]string{"user-service-test"}, m.project)

	fileGenerator([]string{"cassandra-repository"}, m.project)
	fileGenerator([]string{"cassandra-docker-compose"}, m.project)
	goGetPackages(m.project.Path, []string{"github.com/gocql/gocql", "github.com/google/uuid"})
}

type mariadbBuilder struct {
	project *Project
}

func (m *mariadbBuilder) build() {
	fileGenerator([]string{"mock-repository"}, m.project)
	fileGenerator([]string{"user-test"}, m.project)
	fileGenerator([]string{"user"}, m.project)
	fileGenerator([]string{"user-repository-interface"}, m.project)
	fileGenerator([]string{"user-service-interface"}, m.project)
	fileGenerator([]string{"user-service"}, m.project)
	fileGenerator([]string{"user-service-test"}, m.project)

	fileGenerator([]string{"mariadb-repository"}, m.project)
	fileGenerator([]string{"mariadb-docker-compose"}, m.project)
	goGetPackages(m.project.Path, []string{"github.com/go-sql-driver/mysql", "github.com/google/uuid"})
}

type postgresBuilder struct {
	project *Project
}

func (p *postgresBuilder) build() {
	fileGenerator([]string{"mock-repository"}, p.project)
	fileGenerator([]string{"user-test"}, p.project)
	fileGenerator([]string{"user"}, p.project)
	fileGenerator([]string{"user-repository-interface"}, p.project)
	fileGenerator([]string{"user-service-interface"}, p.project)
	fileGenerator([]string{"user-service"}, p.project)
	fileGenerator([]string{"user-service-test"}, p.project)

	fileGenerator([]string{"postgres-repository"}, p.project)
	fileGenerator([]string{"postgres-docker-compose"}, p.project)
	goGetPackages(p.project.Path, []string{"github.com/jackc/pgx/v5", "github.com/google/uuid"})
}

type mongodbBuilder struct {
	project *Project
}

func (m *mongodbBuilder) build() {
	fileGenerator([]string{"mock-repository"}, m.project)
	fileGenerator([]string{"user-test"}, m.project)
	fileGenerator([]string{"user"}, m.project)
	fileGenerator([]string{"user-repository-interface"}, m.project)
	fileGenerator([]string{"user-service-interface"}, m.project)
	fileGenerator([]string{"user-service"}, m.project)
	fileGenerator([]string{"user-service-test"}, m.project)

	fileGenerator([]string{"mongodb-repository"}, m.project)
	fileGenerator([]string{"mongodb-docker-compose"}, m.project)
	goGetPackages(m.project.Path, []string{"go.mongodb.org/mongo-driver/mongo", "github.com/google/uuid"})
}

type dynamodbBuilder struct {
	project *Project
}

func (m *dynamodbBuilder) build() {
	fileGenerator([]string{"mock-repository"}, m.project)
	fileGenerator([]string{"user-test"}, m.project)
	fileGenerator([]string{"user"}, m.project)
	fileGenerator([]string{"user-repository-interface"}, m.project)
	fileGenerator([]string{"user-service-interface"}, m.project)
	fileGenerator([]string{"user-service"}, m.project)
	fileGenerator([]string{"user-service-test"}, m.project)

	fileGenerator([]string{"dynamodb-repository"}, m.project)
	fileGenerator([]string{"dynamodb-docker-compose"}, m.project)
	goGetPackages(m.project.Path, []string{
		"github.com/google/uuid",
		"github.com/aws/aws-sdk-go-v2",
		"github.com/aws/aws-sdk-go-v2/aws",
		"github.com/aws/aws-sdk-go-v2/config",
		"github.com/aws/aws-sdk-go-v2/service/dynamodb",
		"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue",
	})
}

type redisBuilder struct {
	project *Project
}

func (r *redisBuilder) build() {
	fileGenerator([]string{"mock-repository"}, r.project)
	fileGenerator([]string{"user-test"}, r.project)
	fileGenerator([]string{"user"}, r.project)
	fileGenerator([]string{"user-repository-interface"}, r.project)
	fileGenerator([]string{"user-service-interface"}, r.project)
	fileGenerator([]string{"user-service"}, r.project)
	fileGenerator([]string{"user-service-test"}, r.project)

	fileGenerator([]string{"redis-repository"}, r.project)
	fileGenerator([]string{"redis-docker-compose"}, r.project)
	goGetPackages(r.project.Path, []string{
		"github.com/redis/go-redis/v9",
		"github.com/google/uuid",
	})
}

type badgerdbBuilder struct {
	project *Project
}

func (r *badgerdbBuilder) build() {
	fileGenerator([]string{"mock-repository"}, r.project)
	fileGenerator([]string{"user-test"}, r.project)
	fileGenerator([]string{"user"}, r.project)
	fileGenerator([]string{"user-repository-interface"}, r.project)
	fileGenerator([]string{"user-service-interface"}, r.project)
	fileGenerator([]string{"user-service"}, r.project)
	fileGenerator([]string{"user-service-test"}, r.project)

	fileGenerator([]string{"badgerdb-repository"}, r.project)
	goGetPackages(r.project.Path, []string{
		"github.com/dgraph-io/badger/v4",
		"github.com/google/uuid",
	})
}

type sqliteBuilder struct {
	project *Project
}

func (s *sqliteBuilder) build() {
	fileGenerator([]string{"mock-repository"}, s.project)
	fileGenerator([]string{"user-test"}, s.project)
	fileGenerator([]string{"user"}, s.project)
	fileGenerator([]string{"user-repository-interface"}, s.project)
	fileGenerator([]string{"user-service-interface"}, s.project)
	fileGenerator([]string{"user-service"}, s.project)
	fileGenerator([]string{"user-service-test"}, s.project)

	fileGenerator([]string{"sqlite-repository"}, s.project)
	goGetPackages(s.project.Path, []string{"github.com/mattn/go-sqlite3", "github.com/google/uuid"})
}

type mockBuilder struct {
	project *Project
}

func (s *mockBuilder) build() {
	fileGenerator([]string{"mock-repository"}, s.project)
	fileGenerator([]string{"user-test"}, s.project)
	fileGenerator([]string{"user"}, s.project)
	fileGenerator([]string{"user-repository-interface"}, s.project)
	fileGenerator([]string{"user-service-interface"}, s.project)
	fileGenerator([]string{"user-service"}, s.project)
	fileGenerator([]string{"user-service-test"}, s.project)
	goGetPackages(s.project.Path, []string{"github.com/google/uuid"})
}

type databaseBuilderFactory func(p *Project) boilerplateBuilder

var databaseBuilderMap = map[string]databaseBuilderFactory{
	"mysql": func(p *Project) boilerplateBuilder {
		return &mysqlBuilder{p}
	},
	"mariadb": func(p *Project) boilerplateBuilder {
		return &mariadbBuilder{p}
	},
	"postgres": func(p *Project) boilerplateBuilder {
		return &postgresBuilder{p}
	},
	"mongodb": func(p *Project) boilerplateBuilder {
		return &mongodbBuilder{p}
	},
	"redis": func(p *Project) boilerplateBuilder {
		return &redisBuilder{p}
	},
	"dynamodb": func(p *Project) boilerplateBuilder {
		return &dynamodbBuilder{p}
	},
	"sqlite": func(p *Project) boilerplateBuilder {
		return &sqliteBuilder{p}
	},
	"badgerdb": func(p *Project) boilerplateBuilder {
		return &badgerdbBuilder{p}
	},
	"cassandra": func(p *Project) boilerplateBuilder {
		return &cassandraBuilder{p}
	},
	"mock": func(p *Project) boilerplateBuilder {
		return &mockBuilder{p}
	},
}
