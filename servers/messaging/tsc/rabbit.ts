import { Message } from "./message";
import { Channel } from "./channel";
import * as Amqp from "amqp-ts"
// import * as amqp from "amqplib"

export class RabbitObject {
    type: string;
    channel: Channel | any;
    message: Message | any;
    userIDs: string[] | any;
    channelID: string | any;
    messageID: string | any;
    constructor(t: string, c: Channel | any, m: Message | any, ids: string[] | any,
        cid: string | any, mid: string | any) {

        this.type = t;
        this.channel = c;
        this.message = m;
        this.userIDs = ids;
        this.channelID = cid;
        this.messageID = mid
    }
}

// export const createMQConnection = async (): Promise<amqp.Connection> => {
//     let client: amqp.Connection;
//     while (1) {
//         try {
//             client = await amqp.connect(mqURL);
//             break
//         } catch (e) {
//             console.log("Cannot connect to RabbitMQ: failed to connect to server ", e);
//             sleep(1)
//         }
//     }
//     return client!;
// }

// export const createMQChannel = async (conn: amqp.Connection): Promise<amqp.Channel> => {
//     let channel: amqp.Channel;
//     try {
//         channel = await conn.createChannel();
//     } catch (e) {
//         console.log("Cannot create channel on RabbitMQ ", e);
//         process.exit(1)
//     }
//     return channel!;
// }

export function sendObjectToQueue(q: Amqp.Queue, ob: RabbitObject) {
    const message = new Amqp.Message(ob)
    // let json = JSON.stringify(message)
    q.send(message)
    console.log("Sent out the message");
}