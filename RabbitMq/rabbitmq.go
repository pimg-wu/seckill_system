package RabbitMq

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

//url格式：amqp://账号:密码@服务器地址:端口号/vhost
const MQURL = "amqp://testuser:123456@192.168.48.128:5672/test"

// const MQURL = "amqp://testuser:123456@127.0.0.1:5672/test"

type RabbitMQ struct {
	conn      *amqp.Connection
	channel   *amqp.Channel
	QueueName string //队列名称
	Exchange  string //交换机
	key       string //key
	Mqurl     string //连接信息
}

//创建RabbitMQ实例
func NewRabbitMQ(queueName string, exchange string, key string) *RabbitMQ {
	rabbitmq := &RabbitMQ{QueueName: queueName, Exchange: exchange, key: key, Mqurl: MQURL}
	var err error
	rabbitmq.conn, err = amqp.Dial(rabbitmq.Mqurl)
	rabbitmq.failOnErr(err, "创建连接错误！")
	rabbitmq.channel, err = rabbitmq.conn.Channel()
	rabbitmq.failOnErr(err, "获取channel失败！")
	return rabbitmq
}

//断开
func (r *RabbitMQ) Destory() {
	r.channel.Close()
	r.conn.Close()
}

//错误处理函数
func (r *RabbitMQ) failOnErr(err error, message string) {
	if err != nil {
		log.Fatalf("%s:%s", message, err)
		panic(fmt.Sprintf("%s:%s", message, err))
	}
}

//1、创建简单模式实例
func NewRabbitMQSimple(queueName string) *RabbitMQ {
	return NewRabbitMQ(queueName, "", "")
}

//2、简单模式生产代码
func (r *RabbitMQ) PublishSimple(message string) {
	//申请队列，若不存在自动创建，存在则跳过
	//保证队列存在，消息发送到队列
	_, err := r.channel.QueueDeclare(
		r.QueueName,
		false, //是否持久化
		false, //是否为自动删除
		false, //是否具有排他性
		false, //是否阻塞
		nil,   //额外属性
	)
	if err != nil {
		fmt.Println(err)
	}

	//发送消息到队列中
	r.channel.Publish(
		r.Exchange,
		r.QueueName,
		false, //若为true，根据exchange类型和routkey规则，若无法找到符合条件的队列会把发送的消息返回给发送者
		false, //若为true，当exchange发送消息到队列后发现队列上没有绑定消费者，则把消息返回给发送者
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		},
	)
}

//3、简单模式消费者
func (r *RabbitMQ) ConsumeSimple() {
	//申请队列，若不存在自动创建，存在则跳过
	//保证队列存在，消息发送到队列
	_, err := r.channel.QueueDeclare(
		r.QueueName,
		false, //是否持久化
		false, //是否为自动删除
		false, //是否具有排他性
		false, //是否阻塞
		nil,   //额外属性
	)
	if err != nil {
		fmt.Println(err)
	}

	//接收消息
	msgs, err := r.channel.Consume(
		r.QueueName,
		"",    //用来区分多个消费者
		true,  //是否自动应答
		false, //是否具有排他性
		false, //若为true，表示不能将同一个conn中发送的消息传播给消费者
		false, //队列消费是否阻塞
		nil,
	)
	if err != nil {
		fmt.Println(err)
	}

	forever := make(chan bool)
	go func() {
		for d := range msgs {
			//实现要处理的逻辑函数
			log.Printf("Received a message: %s", d.Body)
		}
	}()

	log.Printf("[*] waiting for message, To exit press CTRL+C")
	<-forever
}

//订阅模式
func NewRabbitMQPubSub(exchangeName string) *RabbitMQ {
	return NewRabbitMQ("", exchangeName, "")
}

func (r *RabbitMQ) PublishPub(message string) {
	err := r.channel.ExchangeDeclare(r.Exchange, "fanout", true, false, false, false, nil)
	r.failOnErr(err, "Failed to declare an exchange")

	err = r.channel.Publish(r.Exchange, "", false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(message)})
}

func (r *RabbitMQ) ReceivePub() {
	err := r.channel.ExchangeDeclare(r.Exchange, "fanout", true, false, false, false, nil)
	r.failOnErr(err, "Failed to declare an exchange")

	q, err := r.channel.QueueDeclare("", false, false, true, false, nil)
	r.failOnErr(err, "Failed to declare a queue")

	err = r.channel.QueueBind(q.Name, "", r.Exchange, false, nil)

	message, err := r.channel.Consume(q.Name, "", true, false, false, false, nil)

	forever := make(chan bool)

	go func() {
		for d := range message {
			log.Printf("Receive a message :%s", d.Body)
		}
	}()

	fmt.Println("Please press CTRL+C to quit")
	<-forever
}

//路由模式
func NewRabbitMQRouting(exchangeName string, routingKey string) *RabbitMQ {
	return NewRabbitMQ("", exchangeName, routingKey)
}

//路由模式发送消息
func (r *RabbitMQ) PublishRouting(message string) {
	err := r.channel.ExchangeDeclare(r.Exchange, "direct", true, false, false, false, nil)
	r.failOnErr(err, "Failed to declare an exchange")

	err = r.channel.Publish(r.Exchange, r.key, false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(message),
	})
}

func (r *RabbitMQ) ReceiveRouting() {
	err := r.channel.ExchangeDeclare(r.Exchange, "direct", true, false, false, false, nil)
	r.failOnErr(err, "Failed to declare an exchange")

	q, err := r.channel.QueueDeclare("", false, false, true, false, nil)
	r.failOnErr(err, "failed to declare a queue")

	//绑定队列到exchange中
	err = r.channel.QueueBind(q.Name, r.key, r.Exchange, false, nil)

	message, err := r.channel.Consume(q.Name, "", true, false, false, false, nil)

	forever := make(chan bool)

	go func() {
		for d := range message {
			log.Printf("Receive a message: %s", d.Body)
		}
	}()

	fmt.Println("Please press CTRL+C to quit")
	<-forever
}

//主题模式
func NewRabbitMQTopic(exchangeName string, routingKey string) *RabbitMQ {
	return NewRabbitMQ("", exchangeName, routingKey)
}

func (r *RabbitMQ) PublishTopic(message string) {
	err := r.channel.ExchangeDeclare(r.Exchange, "topic", true, false, false, false, nil)
	r.failOnErr(err, "failed to declare an exchang")

	err = r.channel.Publish(r.Exchange, r.key, false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(message),
	})
}

//主题模式接收消息，需注意key规则，“*”用于匹配一个单词，“#”用于匹配多个单词
func (r *RabbitMQ) ReceiveTopic() {
	err := r.channel.ExchangeDeclare(r.Exchange, "topic", true, false, false, false, nil)
	r.failOnErr(err, "failed to declare an exchang")

	q, err := r.channel.QueueDeclare("", false, false, true, false, nil)
	r.failOnErr(err, "failed to declare a queue")

	err = r.channel.QueueBind(q.Name, r.key, r.Exchange, false, nil)

	messages, err := r.channel.Consume(q.Name, "", true, false, false, false, nil)

	forever := make(chan bool)

	go func() {
		for d := range messages {
			log.Printf("Received a message: %s", d.Body)
		}
	}()

	fmt.Println("Please press CTRL+C to quit")
	<-forever
}
