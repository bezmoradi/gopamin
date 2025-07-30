package templates

import (
	_ "embed"
)

//go:embed files/env/env-dynamodb.tmpl
var dynamodbEnv []byte

func DynamodbEnvTemplate() ([]byte, string) {
	return dynamodbEnv, ".env"
}

//go:embed files/makefile/makefile-dynamodb.tmpl
var dynamodbMakefile []byte

func DynamodbMakefileTemplate() ([]byte, string) {
	return dynamodbMakefile, "Makefile"
}

//go:embed files/readme/readme-dynamodb.tmpl
var dynamodbReadme []byte

func DynamodbReadmeTemplate() ([]byte, string) {
	return dynamodbReadme, "README.md"
}

//go:embed files/docker-compose/docker-compose-dynamodb.tmpl
var dynamodbDockerCompose []byte

func DynamodbDockerComposeTemplate() ([]byte, string) {
	return dynamodbDockerCompose, "docker-compose.yml"
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

//go:embed files/makefile/makefile-mongodb.tmpl
var mongodbMakefile []byte

func MongodbMakefileTemplate() ([]byte, string) {
	return mongodbMakefile, "Makefile"
}

//go:embed files/readme/readme-mongodb.tmpl
var mongodbReadme []byte

func MongodbReadmeTemplate() ([]byte, string) {
	return mongodbReadme, "README.md"
}

//go:embed files/docker-compose/docker-compose-mongodb.tmpl
var mongodbDockerCompose []byte

func MongodbDockerComposeTemplate() ([]byte, string) {
	return mongodbDockerCompose, "docker-compose.yml"
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

//go:embed files/docker-compose/docker-compose-mysql.tmpl
var mysqlDockerCompose []byte

func MysqlDockerComposeTemplate() ([]byte, string) {
	return mysqlDockerCompose, "docker-compose.yml"
}

//go:embed files/makefile/makefile-mysql.tmpl
var mysqlMakefile []byte

func MysqlMakefileTemplate() ([]byte, string) {
	return mysqlMakefile, "Makefile"
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

//go:embed files/makefile/makefile-postgres.tmpl
var postgresMakefile []byte

func PostgresMakefileTemplate() ([]byte, string) {
	return postgresMakefile, "Makefile"
}

//go:embed files/readme/readme-postgres.tmpl
var postgresReadme []byte

func PostgresReadmeTemplate() ([]byte, string) {
	return postgresReadme, "README.md"
}

//go:embed files/docker-compose/docker-compose-postgres.tmpl
var postgresDockerCompose []byte

func PostgresDockerComposeTemplate() ([]byte, string) {
	return postgresDockerCompose, "docker-compose.yml"
}

//go:embed files/internal/adapters/repositories/postgres/repository.tmpl
var postgresRepository []byte

func PostgresRepositoryTemplate() ([]byte, string) {
	return postgresRepository, "internal/adapters/repositories/postgres/repository.go"
}

//go:embed files/env/env-redis.tmpl
var redisEnv []byte

func RedisEnvTemplate() ([]byte, string) {
	return redisEnv, ".env"
}

//go:embed files/makefile/makefile-redis.tmpl
var redisMakefile []byte

func RedisMakefileTemplate() ([]byte, string) {
	return redisMakefile, "Makefile"
}

//go:embed files/readme/readme-redis.tmpl
var redisReadme []byte

func RedisReadmeTemplate() ([]byte, string) {
	return redisReadme, "README.md"
}

//go:embed files/docker-compose/docker-compose-redis.tmpl
var redisDockerCompose []byte

func RedisDockerComposeTemplate() ([]byte, string) {
	return redisDockerCompose, "docker-compose.yml"
}

//go:embed files/internal/adapters/repositories/redis/repository.tmpl
var redisRepository []byte

func RedisRepositoryTemplate() ([]byte, string) {
	return redisRepository, "internal/adapters/repositories/redis/repository.go"
}

//go:embed files/env/env-sqlite.tmpl
var sqliteEnv []byte

func SqliteEnvTemplate() ([]byte, string) {
	return sqliteEnv, ".env"
}

//go:embed files/makefile/makefile-sqlite.tmpl
var sqliteMakefile []byte

func SqliteMakefileTemplate() ([]byte, string) {
	return sqliteMakefile, "Makefile"
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

//go:embed files/makefile/makefile-badgerdb.tmpl
var badgerdbMakefile []byte

func BadgerdbMakefileTemplate() ([]byte, string) {
	return badgerdbMakefile, "Makefile"
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

//go:embed files/makefile/makefile-mariadb.tmpl
var mariadbMakefile []byte

func MariadbMakefileTemplate() ([]byte, string) {
	return mariadbMakefile, "Makefile"
}

//go:embed files/readme/readme-mariadb.tmpl
var mariadbReadme []byte

func MariadbReadmeTemplate() ([]byte, string) {
	return mariadbReadme, "README.md"
}

//go:embed files/docker-compose/docker-compose-mariadb.tmpl
var mariadbDockerCompose []byte

func MariadbDockerComposeTemplate() ([]byte, string) {
	return mariadbDockerCompose, "docker-compose.yml"
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

//go:embed files/makefile/makefile-cassandra.tmpl
var cassandraMakefile []byte

func CassandraMakefileTemplate() ([]byte, string) {
	return cassandraMakefile, "Makefile"
}

//go:embed files/readme/readme-cassandra.tmpl
var cassandraReadme []byte

func CassandraReadmeTemplate() ([]byte, string) {
	return cassandraReadme, "README.md"
}

//go:embed files/docker-compose/docker-compose-cassandra.tmpl
var cassandraDockerCompose []byte

func CassandraDockerComposeTemplate() ([]byte, string) {
	return cassandraDockerCompose, "docker-compose.yml"
}

//go:embed files/internal/adapters/repositories/cassandra/repository.tmpl
var cassandraRepository []byte

func CassandraRepositoryTemplate() ([]byte, string) {
	return cassandraRepository, "internal/adapters/repositories/cassandra/repository.go"
}
