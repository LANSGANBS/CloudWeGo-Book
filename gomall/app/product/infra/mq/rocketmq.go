package mq

import (
	"context"
	"fmt"
	"sync"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"github.com/apache/rocketmq-client-go/v2/rlog"
	"github.com/cloudwego/kitex/pkg/klog"
)

const (
	TopicStockDeduct   = "stock_deduct"
	TopicStockRestore  = "stock_restore"
	TopicStockSync     = "stock_sync"
	ConsumerGroupStock = "stock_consumer_group"
)

var (
	ProducerInstance rocketmq.Producer
	consumerInstance rocketmq.PushConsumer

	nameServerAddr = []string{"127.0.0.1:9876"}
	brokerAddr     = "127.0.0.1:10911"

	once sync.Once
)

func init() {
	rlog.SetLogLevel("error")
}

func InitRocketMQ() error {
	var initErr error
	once.Do(func() {
		p, err := rocketmq.NewProducer(
			producer.WithNameServer(nameServerAddr),
			producer.WithRetry(3),
			producer.WithQueueSelector(NewStockQueueSelector()),
		)
		if err != nil {
			initErr = fmt.Errorf("failed to create producer: %w", err)
			return
		}

		err = p.Start()
		if err != nil {
			initErr = fmt.Errorf("failed to start producer: %w", err)
			return
		}

		ProducerInstance = p
		klog.Info("RocketMQ producer started successfully")
	})

	return initErr
}

func StartConsumer(ctx context.Context, handler func(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error)) error {
	var err error
	consumerInstance, err = rocketmq.NewPushConsumer(
		consumer.WithNameServer(nameServerAddr),
		consumer.WithConsumerModel(consumer.Clustering),
		consumer.WithConsumerOrder(true),
		consumer.WithGroupName(ConsumerGroupStock),
	)
	if err != nil {
		return fmt.Errorf("failed to create consumer: %w", err)
	}

	err = consumerInstance.Subscribe(TopicStockDeduct, consumer.MessageSelector{}, handler)
	if err != nil {
		return fmt.Errorf("failed to subscribe topic: %w", err)
	}

	err = consumerInstance.Start()
	if err != nil {
		return fmt.Errorf("failed to start consumer: %w", err)
	}

	klog.Info("RocketMQ consumer started successfully")

	go func() {
		<-ctx.Done()
		klog.Info("Shutting down RocketMQ consumer...")
		consumerInstance.Shutdown()
	}()

	return nil
}

func Shutdown() {
	if ProducerInstance != nil {
		ProducerInstance.Shutdown()
		klog.Info("RocketMQ producer shutdown")
	}
	if consumerInstance != nil {
		consumerInstance.Shutdown()
		klog.Info("RocketMQ consumer shutdown")
	}
}

type StockQueueSelector struct{}

func NewStockQueueSelector() *StockQueueSelector {
	return &StockQueueSelector{}
}

func (s *StockQueueSelector) Select(message *primitive.Message, queue []*primitive.MessageQueue, shardingKey string) *primitive.MessageQueue {
	if len(queue) == 0 {
		return nil
	}

	if shardingKey == "" {
		shardingKey = message.GetShardingKey()
	}

	if shardingKey == "" {
		return queue[0]
	}

	var hash uint32
	for _, c := range shardingKey {
		hash = hash*31 + uint32(c)
	}

	index := hash % uint32(len(queue))
	return queue[index]
}
