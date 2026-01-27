package kafka

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

// Consumer Kafka消费者
type Consumer struct {
	reader *kafka.Reader
	config *Config
}

// NewConsumer 创建新的消费者
func NewConsumer(config *Config) *Consumer {
	return &Consumer{
		reader: config.NewReader(),
		config: config,
	}
}

// Close 关闭消费者
func (c *Consumer) Close() error {
	return c.reader.Close()
}

func RunConsumer(config *Config) {
	ctx := context.Background()
	consumer := NewConsumer(config)
	defer consumer.Close()

	// 方式1: 使用内置处理
	consumer.StartConsuming(ctx)

	// 方式2: 使用自定义处理函数
	/*
		consumer.ConsumeWithHandler(ctx, func(msg *kafka.Message) {
			fmt.Printf("自定义处理 - 消息: %s\n", string(msg.Value))
			// 这里可以添加你的业务逻辑
		})
	*/
}

// ConsumeMessage 消费单条消息
func (c *Consumer) ConsumeMessage(ctx context.Context) (*kafka.Message, error) {
	msg, err := c.reader.ReadMessage(ctx)
	if err != nil {
		return nil, fmt.Errorf("消费消息失败: %w", err)
	}
	return &msg, nil
}

// StartConsuming 开始持续消费消息
func (c *Consumer) StartConsuming(ctx context.Context) {
	fmt.Printf("开始消费消息 -> Broker: %s, Topic: %s, Group: %s\n",
		c.config.Broker, c.config.Topic, c.config.GroupID)

	for {
		select {
		case <-ctx.Done():
			fmt.Println("停止消费消息")
			return
		default:
			// 设置读取超时，避免永久阻塞
			readCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
			msg, err := c.ConsumeMessage(readCtx)
			cancel()

			if err != nil {
				if err == context.DeadlineExceeded {
					continue
				}
				log.Printf("消费消息错误: %v", err)
				time.Sleep(5 * time.Second)
				continue
			}

			c.handleMessage(msg)
		}
	}
}

// handleMessage 处理收到的消息
func (c *Consumer) handleMessage(msg *kafka.Message) {
	fmt.Printf("✓ 收到消息: %s\n", string(msg.Value))
	fmt.Printf("  主题: %s, 分区: %d, 偏移量: %d, 键: %s\n",
		msg.Topic, msg.Partition, msg.Offset, string(msg.Key))
	fmt.Println("---")
}

// ConsumeWithHandler 使用自定义处理函数消费消息
func (c *Consumer) ConsumeWithHandler(ctx context.Context, handler func(*kafka.Message)) {
	fmt.Printf("开始消费消息(自定义处理) -> Broker: %s, Topic: %s\n",
		c.config.Broker, c.config.Topic)

	for {
		select {
		case <-ctx.Done():
			fmt.Println("停止消费消息")
			return
		default:
			readCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
			msg, err := c.ConsumeMessage(readCtx)
			cancel()

			if err != nil {
				if err == context.DeadlineExceeded {
					continue
				}
				log.Printf("消费消息错误: %v", err)
				time.Sleep(5 * time.Second)
				continue
			}

			handler(msg)
		}
	}
}
