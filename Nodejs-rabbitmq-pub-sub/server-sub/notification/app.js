const amqp = require("amqplib");
const express = require('express');
const app = express();
const PORT = 3002;
const bodyParser = require("body-parser");

const messagesQueue = [];

async function consumeMessages() {
  const connection = await amqp.connect("amqp://localhost");
  const channel = await connection.createChannel();

  await channel.assertExchange("userExchange", "fanout");

  const q = await channel.assertQueue("notificationQueue");

  await channel.bindQueue(q.queue, "userExchange", "");

  channel.consume(q.queue, (msg) => {
    const data = JSON.parse(msg.content);
    messagesQueue.push(data);
    channel.ack(msg);
  });
}

consumeMessages();

app.use(bodyParser.json("application/json"));

app.get("/", (req, res) => {
  try {
    res.json(messagesQueue); 
  } catch (error) {
    console.error("Error retrieving messages:", error);
    res.status(500).send("Internal Server Error");
  }
});

app.listen(PORT, () => {
    console.log(`Server started on port ${PORT}`);
  });
