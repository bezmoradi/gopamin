package templates

import (
	_ "embed"
)

//go:embed files/env/env-dynamodb.tmpl
var dynamodbEnv []byte

func DynamodbEnvTemplate() ([]byte, string) {
	return dynamodbEnv, ".env"
}

//go:embed files/readme/readme-dynamodb.tmpl
var dynamodbReadme []byte

func DynamodbReadmeTemplate() ([]byte, string) {
	return dynamodbReadme, "README.md"
}

//go:embed files/internal/adapters/repositories/dynamodb/repository.tmpl
var dynamodbRepository []byte

func DynamodbRepositoryTemplate() ([]byte, string) {
	return dynamodbRepository, "internal/adapters/repositories/dynamodb/repository.go"
}

//go:embed files/env/env-mongodb.tmpl
var mongodbEnv []byte

func MongodbEnvTemplate() ([]byte, string) {
	return mongodbEnv, ".env"
}

//go:embed files/readme/readme-mongodb.tmpl
var mongodbReadme []byte

func MongodbReadmeTemplate() ([]byte, string) {
	return mongodbReadme, "README.md"
}

//go:embed files/internal/adapters/repositories/mongodb/repository.tmpl
var mongodbRepository []byte

func MongodbRepositoryTemplate() ([]byte, string) {
	return mongodbRepository, "internal/adapters/repositories/mongodb/repository.go"
}

//go:embed files/env/env-mysql.tmpl
var mysqlEnv []byte

func MysqlEnvTemplate() ([]byte, string) {
	return mysqlEnv, ".env"
}

//go:embed files/readme/readme-mysql.tmpl
var mysqlReadme []byte

func MysqlReadmeTemplate() ([]byte, string) {
	return mysqlReadme, "README.md"
}

//go:embed files/internal/adapters/repositories/mysql/repository.tmpl
var mysqlRepository []byte

func MysqlRepositoryTemplate() ([]byte, string) {
	return mysqlRepository, "internal/adapters/repositories/mysql/repository.go"
}

//go:embed files/env/env-postgres.tmpl
var postgresEnv []byte

func PostgresEnvTemplate() ([]byte, string) {
	return postgresEnv, ".env"
}

//go:embed files/readme/readme-postgres.tmpl
var postgresReadme []byte

func PostgresReadmeTemplate() ([]byte, string) {
	return postgresReadme, "README.md"
}

//go:embed files/internal/adapters/repositories/postgres/repository.tmpl
var postgresRepository []byte

func PostgresRepositoryTemplate() ([]byte, string) {
	return postgresRepository, "internal/adapters/repositories/postgres/repository.go"
}

//go:embed files/env/env-sqlite.tmpl
var sqliteEnv []byte

func SqliteEnvTemplate() ([]byte, string) {
	return sqliteEnv, ".env"
}

//go:embed files/readme/readme-sqlite.tmpl
var sqliteReadme []byte

func SqliteReadmeTemplate() ([]byte, string) {
	return sqliteReadme, "README.md"
}

//go:embed files/internal/adapters/repositories/sqlite/repository.tmpl
var sqliteRepository []byte

func SqliteRepositoryTemplate() ([]byte, string) {
	return sqliteRepository, "internal/adapters/repositories/sqlite/repository.go"
}

//go:embed files/env/env-badgerdb.tmpl
var badgerdbEnv []byte

func BadgerdbEnvTemplate() ([]byte, string) {
	return badgerdbEnv, ".env"
}

//go:embed files/readme/readme-badgerdb.tmpl
var badgerdbReadme []byte

func BadgerdbReadmeTemplate() ([]byte, string) {
	return badgerdbReadme, "README.md"
}

//go:embed files/internal/adapters/repositories/badgerdb/repository.tmpl
var badgerdbRepository []byte

func BadgerdbRepositoryTemplate() ([]byte, string) {
	return badgerdbRepository, "internal/adapters/repositories/badgerdb/repository.go"
}

//go:embed files/env/env-mariadb.tmpl
var mariadbEnv []byte

func MariadbEnvTemplate() ([]byte, string) {
	return mariadbEnv, ".env"
}

//go:embed files/readme/readme-mariadb.tmpl
var mariadbReadme []byte

func MariadbReadmeTemplate() ([]byte, string) {
	return mariadbReadme, "README.md"
}

//go:embed files/internal/adapters/repositories/mariadb/repository.tmpl
var mariadbRepository []byte

func MariadbRepositoryTemplate() ([]byte, string) {
	return mariadbRepository, "internal/adapters/repositories/mariadb/repository.go"
}

//go:embed files/env/env-cassandra.tmpl
var cassandraEnv []byte

func CassandraEnvTemplate() ([]byte, string) {
	return cassandraEnv, ".env"
}

//go:embed files/readme/readme-cassandra.tmpl
var cassandraReadme []byte

func CassandraReadmeTemplate() ([]byte, string) {
	return cassandraReadme, "README.md"
}

//go:embed files/internal/adapters/repositories/cassandra/repository.tmpl
var cassandraRepository []byte

func CassandraRepositoryTemplate() ([]byte, string) {
	return cassandraRepository, "internal/adapters/repositories/cassandra/repository.go"
}
