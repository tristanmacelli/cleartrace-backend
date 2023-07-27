import { Message } from "./message";
import { Channel } from "./channel";
import { sleep } from "./mongo_handlers";
import * as amqp from "amqplib";

interface TransactionResult {
  type: string;
  userIDs: number[];
}

export interface MessageTransaction extends TransactionResult {
  message?: Message;
  channelID?: string;
  messageID?: string;
}

export interface ChannelTransaction extends TransactionResult {
  channel?: Channel;
  channelID?: string;
}

export type MessagingTransaction = MessageTransaction | ChannelTransaction;

const mqURL = "amqp://userMessageQueue";
const mqName = "messageLoopbackQueue";

export const createMQConnection = async (): Promise<amqp.Connection> => {
  let retryInterval: number = 1;
  let client: amqp.Connection;

  while (1) {
    try {
      client = await amqp.connect(mqURL);
      break;
    } catch (e) {
      console.log("Cannot connect to RabbitMQ: failed to connect to server ", e);
      await sleep(retryInterval);
      retryInterval *= 2;
    }
  }
  return client!;
};

export const createMQChannel = async (conn: amqp.Connection): Promise<amqp.Channel> => {
  let channel: amqp.Channel;
  try {
    channel = await conn.createChannel();
  } catch (e) {
    console.log("Cannot create channel on RabbitMQ ", e);
    process.exit(1);
  }
  return channel!;
};

export const sendObjectToQueue = (channel: amqp.Channel, ob: MessagingTransaction) => {
  const json = JSON.stringify(ob);
  channel.sendToQueue(mqName, Buffer.from(json));
  // TODO: Remove the following output once tested & working
  console.log("Sent out the message");
};
