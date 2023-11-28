/*
 * @Description:
 * @Author: mali
 * @Date: 2023-11-23 14:36:18
 * @LastEditTime: 2023-11-28 14:59:06
 * @LastEditors: VSCode
 * @Reference:
 */
package kafka

import (
	"errors"
	"fmt"
	"sync"

	"github.com/IBM/sarama"
)

type KafkaConsumer struct {
	Hosts    []string        // Kafka主机IP:端口,例如:192.168.201.206:9092
	Ctopic   string          // topic名称
	Consumer sarama.Consumer // 消费者对象
	Kchan    chan string     // 消费者通道
}

func NewKafkaConsumer(hosts []string, topic string, Kchan chan string) *KafkaConsumer {
	consumer := &KafkaConsumer{
		Hosts:  hosts,
		Ctopic: topic,
		Kchan:  Kchan,
	}
	consumer.kafkaInit()
	return consumer
}

func (k *KafkaConsumer) kafkaInit() {
	// 定义配置选项
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.Version = sarama.V0_10_2_0

	// 初始化一个消费对象
	consumer, err := sarama.NewConsumer(k.Hosts, config)
	if err != nil {
		err = errors.New("NewConsumer错误,原因:" + err.Error())
		fmt.Println(err.Error())
		return
	}

	// 获取所有Topic
	topics, err := consumer.Topics()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// 判断是否有自定义的Topic
	var topicsName = ""
	for _, e := range topics {
		if e == k.Ctopic {
			topicsName = e
			break
		}
	}

	// 没有自定义的Topic则报错
	if topicsName == "" {
		err = errors.New("找不到topics内容")
		fmt.Println(err.Error())
		return
	}
	// 将消费对象保存到结构体以备后面使用
	k.Consumer = consumer
}

func (k *KafkaConsumer) Read() {
	var wg sync.WaitGroup
	// 遍历指定Topic分区持续监控消息
	Partitions, _ := k.Consumer.Partitions(k.Ctopic)
	for _, subPartitions := range Partitions {
		pc, err := k.Consumer.ConsumePartition(k.Ctopic, subPartitions, sarama.OffsetNewest)
		if err != nil {
			continue
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			// 这里进入另一个函数可以过滤消息内容
			k.processPartition(pc)
		}()
	}
	wg.Wait()
}

func (k *KafkaConsumer) processPartition(pc sarama.PartitionConsumer) {
	defer pc.AsyncClose()
	for msg := range pc.Messages() {
		// 这里将获取到的Topic信息发送到通道
		k.Kchan <- string(msg.Value)
	}
}
