package kafka

import (
	configApp "payment-service/config"

	"github.com/IBM/sarama"
	"github.com/sirupsen/logrus"
)

type Kafka struct {
	brokers []string
}

type IKafka interface {
	ProduceMessage(string, []byte) error
}

func NewkafkaProducer(brokers []string) IKafka {
	return &Kafka{
		brokers: brokers,
	}
}

func (k *Kafka) ProduceMessage(topic string, data []byte) error {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = configApp.Config.Kafka.MaxRetry
	producer, err := sarama.NewSyncProducer(k.brokers, config)
	if err != nil {
		logrus.Errorf("failed to create producer: %v", err)
		return err
	}
	defer func(producer sarama.SyncProducer) {
		err := producer.Close()
		if err != nil {
			logrus.Errorf("Failed to close producer: %v", err)
			return
		}
	}(producer)

	message := &sarama.ProducerMessage{
		Topic:   topic,
		Headers: nil,
		Value:   sarama.ByteEncoder(data),
	}

	partition, offset, err := producer.SendMessage(message)
	if err != nil {
		logrus.Errorf("failed to produce message to kafka: %v", err)
		return err
	}

	logrus.Infof("Message is stored in topic(%s)/partition(%d)/offset(%d)\n", topic, partition, offset)
	return nil
}
