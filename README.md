# go-kafka-sample

The sample about kafka, which use these libraries:
- [segmentio/kafka-go](https://github.com/segmentio/kafka-go), [IBM/sarama](https://github.com/IBM/sarama) and [confluent](https://github.com/confluentinc/confluent-kafka-go)
- [kafka](https://github.com/core-go/kafka) to wrap [go-stomp](https://github.com/segmentio/kafka-go)
    - Simplify the way to initialize the consumer, publisher by configurations
        - Props: when you want to change the parameter of consumer or publisher, you can change the config file, and restart Kubernetes POD, do not need to change source code and re-compile.
- [core-go/mq](https://github.com/core-go/mq) to implement this flow, which can be considered a low code tool for message queue consumer:

  ![A common flow to consume a message from a message queue](https://cdn-images-1.medium.com/max/800/1*Y4QUN6QnfmJgaKigcNHbQA.png)

### Similar libraries for nodejs
We also provide these libraries to support nodejs:
- [kafka-plus](https://www.npmjs.com/package/kafka-plus) to support [kafkajs](https://www.npmjs.com/package/kafkajs), combine with [mq-one](https://www.npmjs.com/package/mq-one) for nodejs. Example is at [kafka-sample](https://github.com/typescript-tutorial/kafka-sample).
