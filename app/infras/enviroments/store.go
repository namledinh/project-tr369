package enviroments

import (
	"reflect"
	"sync"

	"github.com/caarlos0/env"
)

var once sync.Once
var internalConfig config

//go:generate envdoc --output ../../../env-doc.md
type config struct {
	// Application name
	AppName string `env:"APP_NAME" envDefault:"kafka-configure-worker" json:"app_name"`
	// Application port
	AppPort string `env:"APP_PORT" envDefault:"8080" json:"app_port"`
	// Application log level
	AppLogLevel string `env:"APP_LOG_LEVEL" envDefault:"info" json:"app_log_level"`
	// App deploy environment
	AppDeployEnv string `env:"APP_DEPLOY_ENV" envDefault:"dev" json:"app_deploy_env"`

	// SQL host
	SQLHost string `env:"SQL_HOST,notEmpty,required" json:"sql_host"`
	// SQL port
	SQLPort string `env:"SQL_PORT,notEmpty,required" json:"sql_port"`
	// SQL USP User
	SQLUSPUser string `env:"SQL_USP_USER" json:"sql_usp_user"`
	// SQL USP Password
	SQLUSPPass string `env:"SQL_USP_PASS" json:"sql_usp_pass"`
	// SQL USP Database
	SQLUSPDB string `env:"SQL_USP_DB,notEmpty,required" json:"sql_usp_db"`

	// MinIO Endpoint
	MinIOEndpoint string `env:"MINIO_ENDPOINT,notEmpty,required" json:"minio_endpoint"`
	// MinIO Access Key
	MinIOAccessKey string `env:"MINIO_ACCESS_KEY,notEmpty,required" json:"minio_access_key"`
	// MinIO Secret Key
	MinIOSecretKey string `env:"MINIO_SECRET_KEY,notEmpty,required" json:"minio_secret_key"`
}

type store struct {
	config *config
}

func NewStore() *store {
	return &store{}
}

func (s *store) Load() {
	if reflect.DeepEqual(internalConfig, config{}) {
		once.Do(
			func() {
				err := env.Parse(&internalConfig)
				if err != nil {
					panic(err)
				}
			})
	}
	s.config = &internalConfig
}

func (s *store) Print() {
	if s.config == nil {
		s.Load()
	}
	println("Application Name:", s.GetAppName())
	println("Application Port:", s.GetAppPort())
	println("Application Log Level:", s.GetAppLogLevel())
	println("SQL Host:", s.GetSQLHost())
	println("SQL Port:", s.GetSQLPort())
	println("SQL USP User:", s.GetSQLUSPUser())
	println("SQL USP Password:", s.GetSQLUSPPass())
	println("SQL USP Database:", s.GetSQLUSPDB())
	println("MinIO Endpoint:", s.GetMinIOEndpoint())
	println("MinIO Access Key:", s.GetMinIOAccessKey())
	println("MinIO Secret Key:", s.GetMinIOSecretKey())
}

func (s *store) GetAppName() string        { return s.config.AppName }
func (s *store) GetAppPort() string        { return s.config.AppPort }
func (s *store) GetAppLogLevel() string    { return s.config.AppLogLevel }
func (s *store) GetAppDeployEnv() string   { return s.config.AppDeployEnv }
func (s *store) GetSQLHost() string        { return s.config.SQLHost }
func (s *store) GetSQLPort() string        { return s.config.SQLPort }
func (s *store) GetSQLUSPUser() string     { return s.config.SQLUSPUser }
func (s *store) GetSQLUSPPass() string     { return s.config.SQLUSPPass }
func (s *store) GetSQLUSPDB() string       { return s.config.SQLUSPDB }
func (s *store) GetMinIOEndpoint() string  { return s.config.MinIOEndpoint }
func (s *store) GetMinIOAccessKey() string { return s.config.MinIOAccessKey }
func (s *store) GetMinIOSecretKey() string { return s.config.MinIOSecretKey }
