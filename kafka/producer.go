package kafka

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

// Producer Kafka生产者
type Producer struct {
	writer *kafka.Writer
	config *Config
}

// NewProducer 创建新的生产者
func NewProducer(config *Config) *Producer {
	return &Producer{
		writer: config.NewWriter(),
		config: config,
	}
}

// Close 关闭生产者
func (p *Producer) Close() error {
	return p.writer.Close()
}

func RunProducer(config *Config) {
	ctx := context.Background()
	producer := NewProducer(config)
	defer producer.Close()

	// 每2秒生产一条消息
	producer.StartProducing(ctx, 2*time.Second)
}

// SendMessage 发送单条消息
func (p *Producer) SendMessage(ctx context.Context, key, value string) error {
	msg := kafka.Message{
		Key:   []byte(key),
		Value: []byte(value),
	}

	err := p.writer.WriteMessages(ctx, msg)
	if err != nil {
		return fmt.Errorf("发送消息失败: %w", err)
	}

	fmt.Printf("✓ 消息发送成功: %s\n", value)
	return nil
}

// SendMessageWithRetry 发送消息带重试机制
func (p *Producer) SendMessageWithRetry(ctx context.Context, key, value string, maxRetries int) error {
	var err error
	for i := 0; i < maxRetries; i++ {
		err = p.SendMessage(ctx, key, value)
		if err == nil {
			return nil
		}

		log.Printf("发送失败，第 %d 次重试: %v", i+1, err)
		time.Sleep(time.Duration(i+1) * time.Second) // 指数退避
	}

	return fmt.Errorf("发送消息失败，重试 %d 次后放弃: %w", maxRetries, err)
}

// StartProducing 开始持续生产消息
func (p *Producer) StartProducing(ctx context.Context, interval time.Duration) {
	fmt.Printf("开始生产消息 -> Broker: %s, Topic: %s\n", p.config.Broker, p.config.Topic)

	messageCount := 0
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			fmt.Println("停止生产消息")
			return
		case <-ticker.C:
			messageCount++
			message := fmt.Sprintf("测试消息 %d - %s",
				messageCount, time.Now().Format("2006-01-02 15:04:05"))

			err := p.SendMessageWithRetry(ctx,
				fmt.Sprintf("key-%d", messageCount),
				message,
				3,
			)

			if err != nil {
				log.Printf("最终发送失败: %v", err)
			}
		}
	}
}
