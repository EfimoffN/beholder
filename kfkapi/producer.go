package kfkapi

import (
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/IBM/sarama"
)

type KafkaProducer struct {
	brokerList []string
	producer   sarama.SyncProducer
}

func CreateKafkaProducer(brokcersList []string, log zerolog.Logger, topic string) (*KafkaProducer, error) {
	kafkaProducer := KafkaProducer{
		brokerList: brokcersList,
	}

	producer, err := kafkaProducer.newProducer()
	if err != nil {
		return nil, errors.Wrap(err, "create kafka producer")
	}

	kafkaProducer.producer = producer

	return &kafkaProducer, nil
}

func (kfk *KafkaProducer) ProducerClose() {
	kfk.producer.Close()
}

func (kfk *KafkaProducer) newProducer() (sarama.SyncProducer, error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 10
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(kfk.brokerList, config)
	if err != nil {
		return nil, errors.Wrap(err, "new sync prosucer")
	}

	return producer, nil
}

func (kfk *KafkaProducer) SendPost(post, topic string) error {
	_, _, err := kfk.producer.SendMessage(&sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(post),
	})

	if err != nil {
		return errors.Wrap(err, "send publication")
	}

	return nil
}
