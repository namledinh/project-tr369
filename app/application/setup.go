package application

import (
	"fmt"
	"time"
	"usp-management-device-api/common/logging"

	"github.com/gomodule/redigo/redis"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/plain"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func mustConnectionSQL(host, user, pass, database, port, applicationName string) *gorm.DB {
	dsn := ""
	// add check if user and pass are empty, then user is root and not set password
	if user == "" && pass == "" {
		dsn = fmt.Sprintf("host=%s dbname=%s port=%s sslmode=disable application_name=%s user=root",
			host,
			database,
			port,
			applicationName,
		)
	} else {
		dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable application_name=%s",
			host,
			user,
			pass,
			database,
			port,
			applicationName,
		)
	}
	var err error

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		PrepareStmt:            true,
		TranslateError:         true,
		SkipDefaultTransaction: true,
		// Logger:                 logger.Default.LogMode(logger.),
	})

	if err != nil {
		panic(err)
	}

	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(100)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Minute * 15)
	sqlDB.SetConnMaxIdleTime(time.Minute * 15)

	return db
}

func kafkaConnect(hosts []string, user, pass string) *kafka.Writer {
	// kafka dialer setup
	mechanism := plain.Mechanism{
		Username: user,
		Password: pass,
	}
	sharedTransport := &kafka.Transport{
		SASL: mechanism,
	}

	kafkaWriter := kafka.NewWriter(kafka.WriterConfig{
		Brokers:      hosts,
		Balancer:     &kafka.LeastBytes{},
		BatchSize:    100,
		BatchTimeout: 50,
		RequiredAcks: int(kafka.RequireAll),
		Async:        false,
		WriteTimeout: 10 * time.Second,
	},
	)

	kafkaWriter.Transport = sharedTransport

	return kafkaWriter
}

func redisConnectPassword(host, pass string, database int) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     300,
		IdleTimeout: 120 * time.Second,
		MaxActive:   1000,
		Wait:        true,
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial("tcp",
				host,
				redis.DialPassword(pass),
				redis.DialDatabase(database))
			if err != nil {
				return nil, err
			}
			return conn, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Second {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}
}

func readerKafkaSetup(
	bootstrapServers []string,
	groupID string,
	topic string,
	username string,
	password string,
) *kafka.Reader {
	readerConfig := kafka.ReaderConfig{
		// Kafka configure
		Brokers: bootstrapServers,
		// Topic configure
		Topic:                  topic,
		GroupID:                groupID,
		WatchPartitionChanges:  true,
		PartitionWatchInterval: 5 * time.Second,

		// offset configure
		StartOffset: kafka.LastOffset,

		// Batching configure
		MinBytes:         1,
		MaxBytes:         256 * 1024 * 1024, // MiB
		MaxWait:          25 * time.Second,
		QueueCapacity:    1000,
		ReadBatchTimeout: 1 * time.Minute,
		CommitInterval:   time.Millisecond * 500,

		// Logger configure
		Logger:      kafka.LoggerFunc(logging.GetLogger("info")),
		ErrorLogger: kafka.LoggerFunc(logging.GetLogger("error")),

		// Safe default
		MaxAttempts:       3,
		HeartbeatInterval: 3 * time.Second,
	}

	if username != "" && password != "" {
		readerConfig.Dialer = &kafka.Dialer{
			Timeout:   10 * time.Second,
			DualStack: true,
			SASLMechanism: plain.Mechanism{
				Username: username,
				Password: password,
			},
		}
	}

	return kafka.NewReader(readerConfig)
}

func minioConnection(endpoint, accessKey, secretKey string) *minio.Client {
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: false,
	})
	if err != nil {
		panic(err)
	}
	return minioClient
}
