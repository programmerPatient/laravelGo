/*
 * @Description:
 * @Author: mali
 * @Date: 2023-11-28 11:03:40
 * @LastEditTime: 2023-11-28 15:10:51
 * @LastEditors: VSCode
 * @Reference:
 */
package kafka

import (
	"sync"

	"github.com/laravelGo/core/config"
	"github.com/laravelGo/core/kafka"
)

type KafkaProducer struct {
	hosts    []string
	topic    string
	Producer *kafka.KafkaProducer
}

var default_producer_once sync.Once
var default_product *KafkaProducer

func GetDefaultKafkaProducer() *KafkaProducer {
	default_producer_once.Do(func() {
		hosts := config.GetStringSlice("kafka.hosts")
		topic := config.GetString("kafka.topic")
		default_product = GetKafkaProducer(hosts, topic)
	})
	return default_product
}

func GetKafkaProducer(hosts []string, topic string) *KafkaProducer {
	return &KafkaProducer{
		hosts:    hosts,
		topic:    topic,
		Producer: kafka.NewKafkaProducer(hosts, topic),
	}
}

func (p *KafkaProducer) Sendmsg(msg string) {
	p.Producer.Sendmsg(msg)
}
