package main

import (
    "fmt"

    "github.com/Shopify/sarama"
)

// 基于sarama第三方库开发的kafka client

func main() {
    config := sarama.NewConfig()
    config.Producer.RequiredAcks = sarama.WaitForAll          // 发送完数据需要leader和follow都确认
    config.Producer.Partitioner = sarama.NewRandomPartitioner // 新选出一个partition
    config.Producer.Return.Successes = true                   // 成功交付的消息将在success channel返回

    // 连接kafka
    fmt.Printf("config=%v\n", config)
    fmt.Printf("kafka version=%s\n", config.Version)
    client, err := sarama.NewAsyncProducer([]string{"192.168.20.45:6669"}, config)
    if err != nil {
        fmt.Println("producer closed, err:", err)
        return
    }
    defer client.Close()

    // 构造一个消息
    msg := &sarama.ProducerMessage{
        Topic: "snmp",
        Value: sarama.ByteEncoder("this is a test log"),
        }

    // 发送消息
    client.Input() <- msg

}



