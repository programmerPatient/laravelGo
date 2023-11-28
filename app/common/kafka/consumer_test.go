/*
 * @Description:
 * @Author: mali
 * @Date: 2023-11-28 15:04:38
 * @LastEditTime: 2023-11-28 15:04:56
 * @LastEditors: VSCode
 * @Reference:
 */
package kafka

import "testing"

func ConsumerTest(T *testing.T) {
	read := make(chan string)
	go GetDefaultKafkaConsumer(read).Read()
	for {
		select {
		case msg := <-read:
			// 处理消息
			println(msg)
		}
	}
}
