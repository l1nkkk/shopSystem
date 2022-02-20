package common

import (
	"fmt"
	"strconv"
	"testing"
	"time"
)

func TestTest(t *testing.T) {
	fmt.Println("test test")
}

// ========== Simple & Work mode
func TestRabbitMQ_ConsumeSimple(t *testing.T) {
	rabbitmq := NewRabbitMQSimple("" +
		"imoocSimple")
	rabbitmq.ConsumeSimple()
}

func TestRabbitMQ_PublishSimple(t *testing.T) {
	rabbitmq := NewRabbitMQSimple("" +
		"imoocSimple")
	rabbitmq.PublishSimple("Hello imooc!")
	fmt.Println("发送成功！")
}

// ========== Publish/Subscribe mode

// run 2 process for test
func TestRabbitMQ_RecieveSub(t *testing.T) {
	rabbitmq := NewRabbitMQPubSub("" +
		"newProduct")
	rabbitmq.RecieveSub()
}

func TestRabbitMQ_PublishPub(t *testing.T) {
	rabbitmq := NewRabbitMQPubSub("" +
		"newProduct")
	for i := 0; i < 100; i++ {
		rabbitmq.PublishPub("订阅模式生产第" +
			strconv.Itoa(i) + "条" + "数据")
		fmt.Println("订阅模式生产第" +
			strconv.Itoa(i) + "条" + "数据")
		time.Sleep(1 * time.Second)
	}
}

// ========== Routing mode

// Receive One
func TestRabbitMQ_RecieveRouting_1(t *testing.T) {
	imoocOne := NewRabbitMQRouting("exImooc", "imooc_one")
	imoocOne.RecieveRouting()
}

// Receive Two
func TestRabbitMQ_RecieveRouting_2(t *testing.T) {
	imoocOne := NewRabbitMQRouting("exImooc", "imooc_two")
	imoocOne.RecieveRouting()
}

func TestRabbitMQ_PublishRouting(t *testing.T) {
	// 由生产者控制可以接收的是哪些消费者，定义两个 MQ 实例
	imoocOne := NewRabbitMQRouting("exImooc", "imooc_one")
	imoocTwo := NewRabbitMQRouting("exImooc", "imooc_two")

	for i := 0; i <= 10; i++ {
		imoocOne.PublishRouting("Hello imooc one!" + strconv.Itoa(i))
		imoocTwo.PublishRouting("Hello imooc Two!" + strconv.Itoa(i))
		time.Sleep(1 * time.Second)
		fmt.Println(i)
	}
}

// ========== Topic mode

// receive all
func TestRabbitMQ_RecieveTopic_1(t *testing.T) {
	imoocOne := NewRabbitMQTopic("exImoocTopic", "#")
	imoocOne.RecieveTopic()
}

// receive two
func TestRabbitMQ_RecieveTopic_2(t *testing.T) {

	imoocOne := NewRabbitMQTopic("exImoocTopic", "imooc.*.two")
	imoocOne.RecieveTopic()
}

func TestRabbitMQ_PublishTopic(t *testing.T) {
	imoocOne := NewRabbitMQTopic("exImoocTopic", "imooc.topic.one")
	imoocTwo := NewRabbitMQTopic("exImoocTopic", "imooc.topic.two")
	for i := 0; i <= 10; i++ {
		imoocOne.PublishTopic("Hello imooc topic one!" + strconv.Itoa(i))
		imoocTwo.PublishTopic("Hello imooc topic Two!" + strconv.Itoa(i))
		time.Sleep(1 * time.Second)
		fmt.Println(i)
	}
}
