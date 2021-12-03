package RabbitMq

import (
	"fmt"
	"strconv"
	"testing"
	"time"
)

//简单模式
func TestSimplePublish(t *testing.T) { //生产者
	rabbitmq := NewRabbitMQSimple("Test")
	rabbitmq.PublishSimple("hello simple!")
	fmt.Println("发送成功！")
}

func TestSimpleReceive(t *testing.T) { //消费者
	rabbitmq := NewRabbitMQSimple("Test")
	rabbitmq.ConsumeSimple()
}

//工作模式
func TestWorkPublish(t *testing.T) {
	rabbitmq := NewRabbitMQSimple("Test")
	for i := 0; i < 20; i++ {
		rabbitmq.PublishSimple(fmt.Sprintln("hello work, number", i))
		time.Sleep(1 * time.Second)
	}
	fmt.Println("发送成功！")
}

func TestWorkReceive1(t *testing.T) {
	rabbitmq := NewRabbitMQSimple("Test")
	rabbitmq.ConsumeSimple()
}
func TestWorkReceive2(t *testing.T) {
	rabbitmq := NewRabbitMQSimple("Test")
	rabbitmq.ConsumeSimple()
}

//订阅模式
func TestPubSubPublish(t *testing.T) {
	rabbifmq := NewRabbitMQPubSub("newPub")
	for i := 0; i < 20; i++ {
		rabbifmq.PublishPub("pubsub data number " + strconv.Itoa(i))
		fmt.Println("product data number" + strconv.Itoa(i))
		time.Sleep(1 * time.Second)
	}
}

func TestPubSubReceive1(t *testing.T) {
	rabbitmq := NewRabbitMQPubSub("newPub")
	rabbitmq.ReceivePub()
}

func TestPubSubReceive2(t *testing.T) {
	rabbitmq := NewRabbitMQPubSub("newPub")
	rabbitmq.ReceivePub()
}

//路由模式
func TestRoutingPublish(t *testing.T) {
	rabbitmq1 := NewRabbitMQRouting("exRouting", "one")
	rabbitmq2 := NewRabbitMQRouting("exRouting", "two")
	for i := 0; i < 20; i++ {
		rabbitmq1.PublishRouting("Hello Routing one!" + strconv.Itoa(i))
		rabbitmq2.PublishRouting("Hello Routing two!" + strconv.Itoa(i))
		fmt.Println(i)
		time.Sleep(1 * time.Second)
	}
}

func TestRoutingReceive1(t *testing.T) {
	rabbitmq := NewRabbitMQRouting("exRouting", "one")
	rabbitmq.ReceiveRouting()
}

func TestRoutingReceive2(t *testing.T) {
	rabbitmq := NewRabbitMQRouting("exRouting", "two")
	rabbitmq.ReceiveRouting()
}

//主题模式
func TestTopicPublish(t *testing.T) {
	rabbitmq1 := NewRabbitMQTopic("exTopic", "test.topic.one")
	rabbitmq2 := NewRabbitMQTopic("exTopic", "test.topic.two")
	for i := 0; i < 20; i++ {
		rabbitmq1.PublishTopic("Hello test topic one! " + strconv.Itoa(i))
		rabbitmq2.PublishTopic("Hello test topic two! " + strconv.Itoa(i))
		fmt.Println(i)
		time.Sleep(1 * time.Second)
	}
}

func TestTopicReceive1(t *testing.T) {
	rabbitmq := NewRabbitMQTopic("exTopic", "#") //all
	rabbitmq.ReceiveTopic()
}

func TestTopicReceive2(t *testing.T) {
	rabbitmq := NewRabbitMQTopic("exTopic", "test.*.two")
	rabbitmq.ReceiveTopic()
}
