package backend_kafka

import (
	"fmt"

	"github.com/IBM/sarama"
	"github.com/rs/zerolog/log"
)

const partition = 0

func CreateTopicKafka(peers []string, topic, hostname string) error {
	admin, err := sarama.NewClusterAdmin(peers, nil)
	if err != nil {
		log.Error().Str("hostname", hostname).Msg(fmt.Sprintf("error=%s", err))
		return err
	}
	defer admin.Close()

	err = admin.CreateTopic(topic, &sarama.TopicDetail{
		NumPartitions:     1,
		ReplicationFactor: 1,
	}, false)
	if err != nil {
		kafkaErr, ok := err.(*sarama.TopicError)
		if ok && kafkaErr.Err == sarama.ErrTopicAlreadyExists {
			return nil
		}
		log.Error().Str("hostname", hostname).Msg(fmt.Sprintf("error=%s", err))
		return err
	}

	return nil
}

func DoCountKafka(peers []string, topic, hostname string) (int, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(peers, config)
	if err != nil {
		log.Error().Str("hostname", hostname).Msg(fmt.Sprintf("error=%s", err))
		return -1, err
	}
	defer producer.Close()

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder("1"),
	}

	_, offset, err := producer.SendMessage(msg)
	if err != nil {
		log.Error().Str("hostname", hostname).Msg(fmt.Sprintf("error=%s", err))
		return -1, err
	}

	return int(offset) + 1, nil
}

func GetCountKafka(peers []string, topic, hostname string) (int, error) {
	client, err := sarama.NewClient(peers, nil)
	if err != nil {
		log.Error().Str("hostname", hostname).Msg(fmt.Sprintf("error=%s", err))
		return -1, err
	}
	defer client.Close()

	offset, err := client.GetOffset(topic, partition, sarama.OffsetNewest)
	if err != nil {
		log.Error().Str("hostname", hostname).Msg(fmt.Sprintf("error=%s", err))
		return -1, err
	}

	return int(offset), nil
}
