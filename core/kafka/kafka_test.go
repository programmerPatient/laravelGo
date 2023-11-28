/*
 * @Description:kafka生产消费单元测试
 * @Author: mali
 * @Date: 2023-11-14 11:27:26
 * @LastEditTime: 2023-11-28 15:09:36
 * @LastEditors: VSCode
 * @Reference:
 */
package kafka

import (
	"testing"
	"time"

	"github.com/laravelGo/app/helper"
)

func TestKafkaProducer(t *testing.T) {
	hosts := []string{"192.168.21.23:9192", "192.168.21.23:9292", "192.168.21.23:9392"}
	topic := "test"
	// 调用初始化函数,并将上面的内容作为参数传进去
	nkm := NewKafkaProducer(hosts, topic)
	nkm.Sendmsg("ssss")
	time.Sleep(1 * time.Second)
}

func TestKafkaConsumer(t *testing.T) {
	hosts := []string{"192.168.21.23:9192", "192.168.21.23:9292", "192.168.21.23:9392"}
	topic := "test"
	response := make(chan string)
	// 调用初始化函数,并将上面的内容作为参数传进去
	nkm := NewKafkaConsumer(hosts, topic, response)
	nkm.Read()
	data := <-response
	helper.FormatPrint(data)
	time.Sleep(1 * time.Second)
}
