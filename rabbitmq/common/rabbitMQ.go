package common

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

// 连接信息: amqp://账户:密码@服务器地址:端口号/vhost
const MQURL = "amqp://imoocuser:imoocuser@127.0.0.1:5672/imooc"

// rabbitMQ结构体
type RabbitMQ struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	//队列名称
	QueueName string
	//交换机名称
	Exchange string
	//bind Key 名称
	Key string
	//连接信息
	Mqurl string
}

// NewRabbitMQ 创建结构体实例
func NewRabbitMQ(queueName string, exchange string, key string) *RabbitMQ {
	return &RabbitMQ{QueueName: queueName, Exchange: exchange, Key: key, Mqurl: MQURL}
}

// Destory 断开channel 和 connection
func (r *RabbitMQ) Destory() {
	r.channel.Close()
	r.conn.Close()
}

// failOnErr 错误处理函数
func (r *RabbitMQ) failOnErr(err error, message string) {
	if err != nil {
		log.Fatalf("%s:%s", message, err)
		panic(fmt.Sprintf("%s:%s", message, err))
	}
}

// NewRabbitMQSimple
// l1nkkk: 其他模式是根据参数变化（queueName、exchange、key）来组合成不同模式,
// Simple 模式下，通过queueName区别每一个Simple模式，
// Exchange和Key都为空，即Simple模式只需要传入queueName，
// 此时 Exchange 为 default，类型为direct，
// 创建Simple模式下RabbitMQ实例
func NewRabbitMQSimple(queueName string) *RabbitMQ {
	// 1.创建RabbitMQ实例
	rabbitmq := NewRabbitMQ(queueName, "", "")
	var err error

	// 下面代码实际上可以放到NewRabbitMQ里
	// 2.获取connection
	rabbitmq.conn, err = amqp.Dial(rabbitmq.Mqurl)
	rabbitmq.failOnErr(err, "failed to connect rabb"+
		"itmq!")

	// 3.获取channel
	rabbitmq.channel, err = rabbitmq.conn.Channel()
	rabbitmq.failOnErr(err, "failed to open a channel")
	return rabbitmq
}

// PublishSimple Simple 模式下生产者代码: 生产消息
func (r *RabbitMQ) PublishSimple(message string) {
	// 1.申请队列，如果队列不存在会自动创建，存在则跳过创建
	_, err := r.channel.QueueDeclare(
		// 队列名称
		r.QueueName,
		//是否持久化；为false的话，如果队列重启，消息就没了。
		false,
		//是否自动删除；当最后一个消费者断开连接的话，是否自动将消息从queue删除
		false,
		//是否具有排他性；若为true，创建一个只有自己可见的队列，其他用户不能访问
		false,
		//是否阻塞处理；发送消息的时候，是否要等待服务器的响应
		false,
		//额外的属性
		nil,
	)
	if err != nil {
		fmt.Println(err)
	}
	// 2.调用channel 发送消息到队列中
	r.channel.Publish(
		r.Exchange,
		r.QueueName,
		// l1nkkk: 如何返还？
		// mandatory, 如果为true，根据自身exchange类型和routekey规则，
		// 如果无法找到符合条件的队列，则会把消息返还给发送者
		false,
		// immediate, 如果为true，当exchange发送消息到队列后发现队列上没有消费者，
		// 则会把消息返还给发送者
		false,
		// 将要发送的信息
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})
}

// ConsumeSimple Simple 模式下消费者代码: 消费消息
func (r *RabbitMQ) ConsumeSimple() {
	//1.申请队列，如果队列不存在会自动创建，存在则跳过创建
	q, err := r.channel.QueueDeclare(
		r.QueueName,
		//是否持久化
		false,
		//是否自动删除
		false,
		//是否具有排他性
		false,
		//是否阻塞处理
		false,
		//额外的属性
		nil,
	)
	if err != nil {
		fmt.Println(err)
	}

	// 2. 接收消息;return <-chan Delivery, error
	msgs, err := r.channel.Consume(
		q.Name, // queue
		// 用来区分多个消费者
		"", // consumer
		// 是否自动应答，接收一个msg，是否告诉 rbmq 消息已经消费完成，
		// 这样 rbmq 就可以删除该msg；如果为false，就需要通过回调函数来通知rbmq
		true, // auto-ack
		// 是否具有排他性，即是否独有
		false, // exclusive
		// 如果设置为true，表示 不能将同一个Conenction中生产者发送的
		// 消息传递给这个Connection中的消费者
		false, // no-local
		// 队列是否阻塞，false表示设置为阻塞；阻塞是指消费完一个下个才进来
		// l1nkkk: （有点奇怪）？？
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		fmt.Println(err)
	}

	forever := make(chan bool)
	// 3. 启用协程处理消息
	go func() {
		for d := range msgs {
			// 消息逻辑处理，可以自行设计逻辑
			log.Printf("Received a message: %s", d.Body)

		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever

}



// NewRabbitMQPubSub 订阅模式创建RabbitMQ实例
func NewRabbitMQPubSub(exchangeName string) *RabbitMQ {
	// 1. 创建RabbitMQ实例
	// l1nkkk: 这种模式下，不需要设置 qName和key，但是需要exchangeName
	rabbitmq := NewRabbitMQ("", exchangeName, "")
	var err error
	// 2. 获取connection
	rabbitmq.conn, err = amqp.Dial(rabbitmq.Mqurl)
	rabbitmq.failOnErr(err, "failed to connect rabbitmq!")
	// 3. 获取channel
	rabbitmq.channel, err = rabbitmq.conn.Channel()
	rabbitmq.failOnErr(err, "failed to open a channel")
	return rabbitmq
}

// PublishPub 订阅模式生产
func (r *RabbitMQ) PublishPub(message string) {
	// 1.尝试创建交换机。如果交换机存在，就不管，如果交换机不存在，则创建
	err := r.channel.ExchangeDeclare(
		r.Exchange,
		"fanout",		// kind，订阅模式下交换机的类型，要定义成fanout（广播类型）
		true, 		// durable
		false,	// autoDelete
		// internal: true表示这个exchange不可以被client用来推送消息，
		// 仅用来进行exchange和exchange之间的绑定
		false,
		false,
		nil,
	)
	r.failOnErr(err, "Failed to declare an excha"+
		"nge")

	// 2.发送消息
	err = r.channel.Publish(
		r.Exchange,
		"",
		false,	// mandatort
		false,	// immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})
}

// RecieveSub 订阅模式消费端代码
func (r *RabbitMQ) RecieveSub() {
	// 1.试探性创建交换机
	err := r.channel.ExchangeDeclare(
		r.Exchange,
		//交换机类型
		"fanout",
		true,
		false,
		//YES表示这个exchange不可以被client用来推送消息，仅用来进行exchange和exchange之间的绑定
		false,
		false,
		nil,
	)
	r.failOnErr(err, "Failed to declare an exch"+
		"ange")

	// 2.试探性创建队列，这里注意队列名称不要写
	q, err := r.channel.QueueDeclare(
		"", // 随机生成队列名称
		false,
		false,
		true,
		false,
		nil,
	)
	r.failOnErr(err, "Failed to declare a queue")

	// 3. 绑定队列到 exchange 中
	err = r.channel.QueueBind(
		q.Name,
		// 在pub/sub模式下，这里的key要必须为空，否则不是订阅模式
		"",
		r.Exchange,
		false,
		nil)

	// 4.获取消费消息 chan
	messges, err := r.channel.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	forever := make(chan bool)

	go func() {
		for d := range messges {
			log.Printf("Received a message: %s", d.Body)
		}
	}()

	fmt.Println("退出请按 CTRL+C")
	<-forever
}

// NewRabbitMQRouting 创建路由模式的RabbitMQ实例
func NewRabbitMQRouting(exchangeName string, routingKey string) *RabbitMQ {
	// 创建RabbitMQ实例，注意这里的 routingKey
	rabbitmq := NewRabbitMQ("", exchangeName, routingKey)
	var err error
	// 获取connection
	rabbitmq.conn, err = amqp.Dial(rabbitmq.Mqurl)
	rabbitmq.failOnErr(err, "failed to connect rabbitmq!")
	// 获取channel
	rabbitmq.channel, err = rabbitmq.conn.Channel()
	rabbitmq.failOnErr(err, "failed to open a channel")
	return rabbitmq
}

// PublishRouting 路由模式发送消息
func (r *RabbitMQ) PublishRouting(message string) {
	// 1.尝试创建交换机
	err := r.channel.ExchangeDeclare(
		r.Exchange,
		// 要改成direct，不能是fanout，其他代码都和订阅模式相同
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)

	r.failOnErr(err, "Failed to declare an excha"+
		"nge")

	// 2.发送消息
	err = r.channel.Publish(
		r.Exchange,
		//要设置
		r.Key,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})
}

//路由模式接受消息
func (r *RabbitMQ) RecieveRouting() {
	// 1.试探性创建交换机
	err := r.channel.ExchangeDeclare(
		r.Exchange,
		// 交换机类型，direct
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	r.failOnErr(err, "Failed to declare an exch"+
		"ange")
	// 2.试探性创建队列，这里注意队列名称不要写
	q, err := r.channel.QueueDeclare(
		"", //随机生产队列名称
		false,
		false,
		true,
		false,
		nil,
	)
	r.failOnErr(err, "Failed to declare a queue")

	//绑定队列到 exchange 中
	err = r.channel.QueueBind(
		q.Name,
		//需要绑定key
		r.Key,
		r.Exchange,
		false,
		nil)

	//消费消息
	messges, err := r.channel.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	forever := make(chan bool)

	go func() {
		for d := range messges {
			log.Printf("Received a message: %s", d.Body)
		}
	}()

	//fmt.Println("退出请按 CTRL+C\n")
	<-forever
}

//话题模式
//创建RabbitMQ实例
func NewRabbitMQTopic(exchangeName string, routingKey string) *RabbitMQ {
	//创建RabbitMQ实例
	rabbitmq := NewRabbitMQ("", exchangeName, routingKey)
	var err error
	//获取connection
	rabbitmq.conn, err = amqp.Dial(rabbitmq.Mqurl)
	rabbitmq.failOnErr(err, "failed to connect rabbitmq!")
	//获取channel
	rabbitmq.channel, err = rabbitmq.conn.Channel()
	rabbitmq.failOnErr(err, "failed to open a channel")
	return rabbitmq
}

//话题模式发送消息
func (r *RabbitMQ) PublishTopic(message string) {
	//1.尝试创建交换机
	err := r.channel.ExchangeDeclare(
		r.Exchange,
		//要改成topic
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)

	r.failOnErr(err, "Failed to declare an excha"+
		"nge")

	//2.发送消息
	err = r.channel.Publish(
		r.Exchange,
		//要设置
		r.Key,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})
}

//话题模式接受消息
//要注意key,规则
//其中“*”用于匹配一个单词，“#”用于匹配多个单词（可以是零个）
//匹配 imooc.* 表示匹配 imooc.hello, 但是imooc.hello.one需要用imooc.#才能匹配到
func (r *RabbitMQ) RecieveTopic() {
	//1.试探性创建交换机
	err := r.channel.ExchangeDeclare(
		r.Exchange,
		//交换机类型
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	r.failOnErr(err, "Failed to declare an exch"+
		"ange")
	//2.试探性创建队列，这里注意队列名称不要写
	q, err := r.channel.QueueDeclare(
		"", //随机生产队列名称
		false,
		false,
		true,
		false,
		nil,
	)
	r.failOnErr(err, "Failed to declare a queue")

	//绑定队列到 exchange 中
	err = r.channel.QueueBind(
		q.Name,
		//在pub/sub模式下，这里的key要为空
		r.Key,
		r.Exchange,
		false,
		nil)

	//消费消息
	messges, err := r.channel.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	forever := make(chan bool)

	go func() {
		for d := range messges {
			log.Printf("Received a message: %s", d.Body)
		}
	}()

	fmt.Println("退出请按 CTRL+C")
	<-forever
}
