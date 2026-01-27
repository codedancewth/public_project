package kafka

import (
	"github.com/segmentio/kafka-go"
	"time"
)

// Config Kafka配置
type Config struct {
	Broker      string
	Topic       string
	GroupID     string
	DialTimeout time.Duration
	MaxAttempts int
}

// DefaultConfig 默认配置
func DefaultConfig() *Config {
	return &Config{
		Broker:  "xxxxx:port", // 你的云服务器地址
		Topic:   "quickstart-events",
		GroupID: "go-consumer-group",
	}
}

// NewWriter 创建Kafka生产者
func (c *Config) NewWriter() *kafka.Writer {
	return &kafka.Writer{
		Addr:  kafka.TCP(c.Broker),
		Topic: c.Topic,
		// 可以后续在研究
		Balancer:     &kafka.LeastBytes{},
		BatchTimeout: 10 * time.Millisecond,
		WriteTimeout: 10 * time.Second,
		RequiredAcks: kafka.RequireOne,
		Async:        false,
	}
}

// NewReader 创建Kafka消费者
func (c *Config) NewReader() *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{c.Broker},
		Topic:   c.Topic,
		GroupID: c.GroupID,
		// 可以后续在研究
		MinBytes:       10e3,
		MaxBytes:       10e6,
		CommitInterval: time.Second,
		StartOffset:    kafka.FirstOffset,
	})
}
