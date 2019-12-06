// "use strict";
// version 0.1

//require the express and morgan packages
import express from "express";
import morgan from "morgan";
import { MongoClient, Collection, Db } from "mongodb";
import * as mongo from "./mongo_handlers";
import { Message } from "./message";
import { Channel } from "./channel";

import * as Amqp from "amqp-ts"
// import * as amqp from "amqplib";
// import Bluebird from "bluebird";

//create a new express application
const app = express();

const addr = process.env.ADDR || "80";
//split host and port using destructuring
// const [host, port] = addr.split(":");
// let portNum = parseInt(port);

//add JSON request body parsing middleware
app.use(express.json());
//add the request logging middleware
app.use(morgan("dev"));

// Connection URL
const url = 'mongodb://mongodb:27017/mongodb';

// Database Name
const dbName = 'mongodb';
var db: Db;

var messages: Collection;
var channels: Collection;

class RabbitObject {
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

// Reasoning for refactor: 
// https://bit.ly/342jCtj

// Create a new MongoClient
const mc = new MongoClient(url, { useUnifiedTopology: true });

const createConnection = async (): Promise<MongoClient> => {
    let client: MongoClient;
    try {
        client = await mc.connect();
    } catch (e) {
        console.log("Cannot connect to mongo: MongoNetworkError: failed to connect to server");
        console.log("Restarting Messaging server");
        process.exit(1);
    }
    return client!;
}

// const createMQConnection = async (): Bluebird<amqp.Connection> => {
//     let client: amqp.Connection;
//     try {
//         client = await amqp.connect("amqp://localhost");;
//     } catch (e) {
//         console.log("Cannot connect to RabbitMQ: failed to connect to server ", e);
//         process.exit(1);
//     }
//     return client!;
// }

// const createMQChannel = async (conn: amqp.Connection): Bluebird<amqp.Channel> => {
//     let channel: amqp.Channel;
//     try {
//         channel = await conn.createChannel();
//     } catch (e) {
//         console.log("Cannot create channel on RabbitMQ ", e);
//         process.exit(1);
//     }
//     return channel!;
// }

function sendObjectToQueue(q: Amqp.Queue, ob: RabbitObject) {
    let message = new Amqp.Message(ob)
    // let json = JSON.stringify(message)
    q.send(message)
    console.log("Sent out the message");
}


const main = async () => {
    let client = await createConnection();
    db = client.db(dbName);
    channels = db.collection("channels");
    messages = db.collection("messages");

    // let rabbitConn: Amqp.Connection;
    // rabbitConn = new Amqp.Connection("amqp://rabbitMQ");
    // let queue = rabbitConn.declareQueue("helloQueue");

    // let mqClient = await createMQConnection();
    // let mqChannel = await createMQChannel(mqClient);

    // mqClient.then(function (conn) {
    //     return conn.createChannel();
    // }).then(function (ch) {
    //     return ch.assertQueue(q).then(function (ok) {
    //         return ch.sendToQueue(q, Buffer.from('something to do'));
    //     });
    // })

    // mqClient.createChannel(function (error1: any, channel: any) {
    //     if (error1) {
    //         throw error1;
    //     }
    //     var queue = 'hello';
    //     var msg = 'Hello world';

    //     channel.assertQueue(queue, {
    //         durable: false
    //     });

    //     channel.sendToQueue(queue, Buffer.from(msg));
    //     console.log(" [x] Sent %s", msg);
    // });

    // TODO: We should do a test of the mongo helper methods

    app.listen(+addr, "", () => {
        //callback is executed once server is listening
        console.log(`server is listening at http://:${addr}...`);
        // console.log("port : " + port);
        // console.log("host : " + host);
    });





    app.use("/v1/channels", (req: any, res: any) => {
        // Check that the user is authenticated
        console.log("Inside Index.ts /channels")
        console.log(req.headers)
        if (req.headers['x-user'] == null) {
            console.log("pych null")
            res.status(401);
            res.send()
            return;
        }
        console.log("Passes x user ")
        console.log(req.body)
        switch (req.method) {
            case 'GET':
                // QUERY for all channels here
                let result = mongo.getAllChannels(channels, res);
                if (result.successMessage == "") {
                    res.status(500);
                    res.send()
                    break;
                }
                res.set("Content-Type", "application/json");
                res.json(result.resultJSON);
                res.send();
                break;

            case 'POST':
                if (req.body.name == null) {
                    res.status(500);
                    res.send()
                    //do something about the name property being null
                }
                // Call database to INSERT this new channel
                let newChannel = createChannel(req);

                mongo.insertNewChannel(channels, newChannel).then(insertResult => {
                    if (insertResult.errString == "duplicate") {
                        res.status(400).send()
                        return;
                    }

                    if (insertResult.errString == "Error inserting new channel") {
                        res.status(500);
                        res.send()
                        return
                    }
                    let insertChannel = insertResult.newChannel;
                    res.status(201)
                    res.set("Content-Type", "application/json");
                    res.json(insertChannel);
                    res.send();
                })
                break;
                //probably cant do this >>> .send("success");

                // // add to rabbitMQ queue
                // let obj = new RabbitObject('channel-new', insertChannel, null,
                //     insertChannel.members, null, null)
                // sendObjectToQueue(queue, obj)
            default:
                res.status(405);
                res.send()
                break;
        }
    });

    // Specific channel handler
    app.use("/v1/channels/:channelID", (req: any, res: any) => {
        // Check that the user is authenticated
        if (req.headers['X-User'] == null) {
            res.status(401);
            res.send()
            return;
        }
        // QUERY for the channel based on req.params.channelID
        if (req.params.channelID == null) {
            res.status(404);
            res.send()
            return;
        }
        let result = mongo.getChannelByID(channels, req.params.channelID);
        if (result.errString.length() > 0) {
            res.status(500);
            res.send()
            return;
        }
        let resultChannel = result.finalChannel;
        switch (req.method) {
            case 'GET':
                if (!isChannelMember(resultChannel, req.Header.Xuser)) {
                    res.status(403);
                    res.send()
                    break;
                }
                let returnedMessages;
                // QUERY for last 100 messages here
                if (req.params.before != null) {
                    returnedMessages = mongo.last100SpecificMessages(messages, resultChannel._id, req.params.before, res);
                    if (returnedMessages == null) {
                        res.status(500);
                        res.send()
                        break;
                    }
                } else {
                    returnedMessages = mongo.last100Messages(messages, resultChannel._id, res);
                    if (returnedMessages == null) {
                        res.status(500);
                        res.send()
                        break;
                    }
                }
                res.set("Content-Type", "application/json");
                // write last 100 messages to the client, encoded in JSON 
                // We already did this in the helper function
                res.json(returnedMessages);
                break;

            case 'POST':
                if (resultChannel.private && !isChannelMember(resultChannel, req.Header.Xuser)) {
                    res.status(403);
                    res.send()
                    break;
                }
                // Create a new message
                // Call database to INSERT a new message to the channel
                let newMessage = createMessage(req);
                let insertedResult = mongo.insertNewMessage(messages, newMessage);
                if (insertedResult.errString.length > 0) {
                    res.status(500);
                    res.send()
                }
                let insertedMessage = insertedResult.newMessage;
                res.set("Content-Type", "application/json");
                res.json(insertedMessage);
                res.status(201);
                  // probably cant do this >>> .send("success");

                // // add to rabbitMQ queue
                // let PostObj = new RabbitObject('message-new', null, insertedMessage,
                //     resultChannel.members, null, null)
                // sendObjectToQueue(queue, PostObj)
                res.send()

                break;
            case 'PATCH':
                if (!isChannelCreator(resultChannel, req.Header.Xuser)) {
                    res.status(403);
                    res.send()
                    break;
                }
                // Call database to UPDATE the channel name and/or description
                let updateResult = mongo.updateChannel(channels, resultChannel, req);
                if (updateResult.errString.length > 0) {
                    res.status(500);
                    res.send()
                    break;
                }
                let updatedChannel = updateResult.existingChannel;
                res.set("Content-Type", "application/json");
                res.json(updatedChannel);

                // add to rabbitMQ queue
                // let PatchObj = new RabbitObject('channel-update', updatedChannel, null,
                //     updatedChannel.members, null, null)
                // sendObjectToQueue(queue, PatchObj)
                res.send()
                break;
            case 'DELETE':
                if (!isChannelCreator(resultChannel, req.Header.Xuser)) {
                    res.status(403);
                    res.send()
                    break;
                }
                // Call database to DELETE this channel
                let result = mongo.deleteChannel(channels, messages, resultChannel);
                if (result.length > 0) {
                    res.status(500);
                    res.send()

                }

                // add to rabbitMQ queue
                // let obj = new RabbitObject('channel-delete', null, null, resultChannel.members,
                //     resultChannel._id, null)
                // sendObjectToQueue(queue, obj)
                res.set("Content-Type", "text/plain");
                res.send("Channel was successfully deleted");
                break;
            default:
                res.status(405);
                res.send()
                break;
        }
    });

    // Adding and removing members from your channel
    app.use("/v1/channels/:channelID/members", (req: any, res: any) => {
        // Check that the user is authenticated
        if (req.headers['X-User'] == null) {
            res.status(401);
            res.send()
            return;
        }
        // QUERY for the channel based on req.params.channelID
        if (req.params.channelID == null) {
            res.status(404);
            res.send()
            return;
        }
        let result = mongo.getChannelByID(channels, req.params.channelID);
        if (result.errString.length() > 0) {
            res.status(500);
            res.send()
            return;
        }
        let resultChannel = result.finalChannel;
        switch (req.method) {
            case 'POST':
                if (!isChannelCreator(resultChannel, req.Header.Xuser)) {
                    res.status(403);
                    res.send()
                    break;
                }
                // Call database to UPDATE the current channel
                let addResult = mongo.addChannelMember(channels, resultChannel, req);
                if (addResult.length > 0) {
                    res.status(500);
                    res.send()

                    break;
                }
                res.set("Content-Type", "application/json");
                res.status(201).send(req.user.ID + " was added to your channel");
                break;
            case 'DELETE':
                if (!isChannelCreator(resultChannel, req.Header['X-user'])) {
                    res.status(403)
                    res.send()

                    break;
                }
                // database to UPDATE the current channel members
                mongo.removeChannelMember(channels, resultChannel, req).then(errResult => {
                    if (errResult.length > 0) {
                        res.status(500);
                        res.send()
                        return;
                    }
                    res.set("Content-Type", "text/plain");
                    res.status(201).send(req.user.ID + " was removed from your channel");
                    return;
                })
                break;
            default:
                res.status(405);
                res.send()
                break;
        }
    });

    // Editing the body of or deleting a message
    app.use("/v1/messages/:messageID", (req: any, res: any) => {
        // Check that the user is authenticated
        if (req.headers['X-User'] == null) {
            res.status(401);
            res.send()

            return;
        }
        if (req.params.messageID == null) {
            res.status(404);
            res.send()

            return;
        }
        let result = mongo.getMessageByID(messages, req.params.messageID);
        if (result.errString.length() > 0) {
            res.status(500);
            res.send()

            return;
        }
        let resultMessage = result.finalMessage;
        switch (req.method) {
            case 'PATCH':
                if (!isMessageCreator(resultMessage, req.Header.Xuser)) {
                    res.status(403);
                    res.send()

                    break;
                }
                // TODO: Call the database to UPDATE the message in the database using the messageID
                mongo.updateMessage(messages, resultMessage, req).then(updatedResult => {
                    if (updatedResult.errString.length > 0) {
                        res.status(500);
                        res.send()
                        return;
                    }
                    let updatedMessage = updatedResult.existingMessage;
                    res.set("Content-Type", "application/json");
                    res.json(updatedMessage);
    
                    let resultChannel = mongo.getChannelByID(channels, updatedMessage.channelID)
                    // add to rabbitMQ queue
                    // let pobj = new RabbitObject('message-update', null, updatedMessage,
                    //     resultChannel.finalChannel, null, null)
                    // sendObjectToQueue(queue, pobj)
                    return;
                })
                break;
            case 'DELETE':
                if (!isMessageCreator(resultMessage, req.Header.Xuser)) {
                    res.status(403)
                    res.send()
                    break;
                }
                // Call database to DELETE the specified message using the messageID
                // Call database to DELETE this channel
                mongo.deleteMessage(messages, resultMessage).then(result => {
                    if (result.length > 0) {
                        res.status(500);
                        res.send()
    
                    }
                    // add to rabbitMQ queue
                    // let PostObj = new RabbitObject('message-delete', null, null,
                    //     deleteMessageChannel.finalChannel.members, null, resultMessage._id)
                    // sendObjectToQueue(queue, PostObj)

                    res.set("Content-Type", "text/plain");
                    res.send("Message deleted");

                    mongo.getChannelByID(channels, resultMessage.channelID)
                    return
                });
                break;
            default:
                res.status(405);
                res.send()
                break;
        }
    });

    function createChannel(req: any): Channel {
        let c = req.body;
        return new Channel(c.name, c.description, c.private,
            c.members, c.createdAt, c.creator, c.editedAt);
    }

    function createMessage(req: any): Message {
        let m = req.body;
        return new Message(req.params.ChannelID, m.createdAt, m.body,
            m.creator, m.editedAt);
    }

    function isChannelMember(channel: Channel, userID: any): boolean {
        let isMember = false;
        if (channel.private) {
            for (let i = 0; i < channel.members.length; i++) {
                if (channel.members[i] == userID) {
                    isMember = true;
                    break;
                }
            }
        } else {
            isMember = true;
        }
        return isMember;
    }

    function isChannelCreator(channel: Channel, userID: any): boolean {
        return channel.creator.ID == userID;
    }

    function isMessageCreator(message: Message, userID: any): boolean {
        return message.creator.ID == userID;
    }

    // error handler that will be called if
    // any handler earlier in the chain throws
    // an exception or passes an error to next()
    app.use((err: any, req: any, res: any) => {
        //write a stack trace to standard out,
        //which writes to the server's log
        console.error(err.stack)

        //but only report the error message
        //to the client, with a 500 status code
        res.set("Content-Type", "text/plain");
        res.status(500).send(err.message);
    });

}

main();