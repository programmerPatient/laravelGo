/*
 * @Description:
 * @Author: mali
 * @Date: 2023-11-28 15:04:38
 * @LastEditTime: 2023-11-28 15:05:38
 * @LastEditors: VSCode
 * @Reference:
 */
package kafka

import (
	"fmt"
	"sync"
	"testing"
)

func ProducerTest(T *testing.T) {
	max_num := 100000
	max_channel := make(chan struct{}, 1000)
	wait := sync.WaitGroup{}
	for i := 1; i <= max_num; i++ {
		wait.Add(1)
		max_channel <- struct{}{}
		go func(k int) {
			defer func() {
				<-max_channel
				wait.Done()
			}()
			GetDefaultKafkaProducer().Sendmsg(fmt.Sprintf("%v", k))
		}(i)
	}
	wait.Wait()
	close(max_channel)
}
