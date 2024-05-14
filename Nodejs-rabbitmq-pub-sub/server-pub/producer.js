const amqp = require("amqplib");
const config = require("./config");

class Producer {
  channel;

  async createChannel() {
    const connection = await amqp.connect(config.rabbitMQ.url);
    this.channel = await connection.createChannel();
  }

  async publishMessage(message) {
    if (!this.channel) {
      await this.createChannel();
    }

    const exchangeName = config.rabbitMQ.exchangeName;
    await this.channel.assertExchange(exchangeName, "fanout");

    const msg = {
      message: message,
      dateTime: new Date(),
    };
    await this.channel.publish(
      exchangeName,
      '',
      Buffer.from(JSON.stringify(msg))
    );

  }
}

module.exports = Producer;
