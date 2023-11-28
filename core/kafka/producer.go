/*
 * @Description:
 * @Author: mali
 * @Date: 2023-11-23 14:36:18
 * @LastEditTime: 2023-11-28 15:09:40
 * @LastEditors: VSCode
 * @Reference:
 */
package kafka

import (
	"errors"
	"fmt"

	"github.com/IBM/sarama"
	"github.com/laravelGo/core/logger"
)

type KafkaProducer struct {
	Hosts         []string             // Kafka主机
	Ptopic        string               // Topic
	AsyncProducer sarama.AsyncProducer // Kafka生产者接口对象
}

func NewKafkaProducer(hosts []string, topic string) *KafkaProducer {
	KafkaProduct := &KafkaProducer{
		Hosts:  hosts,
		Ptopic: topic,
	}
	KafkaProduct.kafkaInit()
	return KafkaProduct
}

func (k *KafkaProducer) kafkaInit() {
	// 定义配置参数
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	config.Version = sarama.V0_10_2_0
	// 初始化一个生产者对象
	producer, err := sarama.NewAsyncProducer(k.Hosts, config)
	if err != nil {
		err = errors.New("NewAsyncProducer错误,原因:" + err.Error())
		fmt.Println(err.Error())
		return
	}

	// 保存对象到结构体
	k.AsyncProducer = producer
}

func (k *KafkaProducer) Sendmsg(data string) {
	msg := &sarama.ProducerMessage{
		Topic: k.Ptopic,
	}
	// 信息编码
	msg.Value = sarama.ByteEncoder([]byte(data))
	// 将信息发送给通道
	k.AsyncProducer.Input() <- msg
	select {
	case res := <-k.AsyncProducer.Successes():
		fmt.Printf("%+v", res)
	case err := <-k.AsyncProducer.Errors():
		logger.ErrorJSON("kafka", "producer", map[string]interface{}{
			"msg":  "投递消息失败",
			"data": k,
			"err":  err,
		})
	}
}
