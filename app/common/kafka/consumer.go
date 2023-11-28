/*
 * @Description:
 * @Author: mali
 * @Date: 2023-11-28 13:55:06
 * @LastEditTime: 2023-11-28 15:10:43
 * @LastEditors: VSCode
 * @Reference:
 */
package kafka

import (
	"sync"

	"github.com/laravelGo/core/config"
	"github.com/laravelGo/core/kafka"
)

type KafkaConsumer struct {
	hosts    []string
	topic    string
	kchan    chan string
	Consumer *kafka.KafkaConsumer
}

var default_consumer_once sync.Once
var default_consumer *KafkaConsumer

/**
 * @Author: mali
 * @Func:
 * @Description: 默认kafka消费端
 * @Param:
 * @Return:
 * @Example:
 * @param {chanstring} Kchan 接收数据的
 */
func GetDefaultKafkaConsumer(Kchan chan string) *KafkaConsumer {
	default_consumer_once.Do(func() {
		hosts := config.GetStringSlice("kafka.hosts")
		topic := config.GetString("kafka.topic")
		default_consumer = GetKafkaConsumer(hosts, topic, Kchan)
	})
	return default_consumer
}

func GetKafkaConsumer(hosts []string, topic string, Kchan chan string) *KafkaConsumer {
	return &KafkaConsumer{
		hosts:    hosts,
		topic:    topic,
		kchan:    Kchan,
		Consumer: kafka.NewKafkaConsumer(hosts, topic, Kchan),
	}
}

func (p *KafkaConsumer) Read() {
	// 调用初始化函数,并将上面的内容作为参数传进去
	p.Consumer.Read()
}
