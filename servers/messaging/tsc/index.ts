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
import { User } from "./user";

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
    });

    app.use("/v1/channels/:channelID/members", (req: any, res: any) => {
        // Check that the user is authenticated
        if (req.headers['x-user'] == null) {
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
        mongo.getChannelByID(channels, req.params.channelID).then((result) => {
            if (result.errString.length > 0) {
                res.status(500);
                res.send()
                return;
            }
            let resultChannel = result.finalChannel;
            switch (req.method) {
                case 'POST':
                    let user = JSON.parse(req.headers['x-user'])
                    if (!isChannelCreator(resultChannel, user.id)) {
                        res.status(403);
                        res.set("Content-Type", "text/plain");
                        res.send("Cannot change members")
                        break;
                    }
                    // Call database to UPDATE the current channel
                    mongo.addChannelMember(channels, resultChannel, req).then((err) => {
                        if (err.length > 0) {
                            res.status(500);
                            res.send()
                            return;
                        }
                        res.set("Content-Type", "application/json");
                        res.status(201).send(req.user.ID + " was added to your channel");
                        return;
                    })
                    break;
                case 'DELETE':
                    if (!isChannelCreator(resultChannel, req.Header['X-user'])) {
                        res.status(403)
                        res.set("Content-Type", "text/plain");
                        res.send("Cannot delete members")
    
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
                    res.set("Content-Type", "text/plain");
                    res.send("Incorrect HTTP Header")
                    break;
            }
            return;
        });
    });

// Specific channel handler
    app.use("/v1/channels/:channelID", (req: any, res: any) => {
            
        // Check that the user is authenticated
        if (req.headers['x-user'] == null) {
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
        // req.params
        mongo.getChannelByID(channels, req.params.channelID).then((result) => {
            if (result.errString.length > 0) {
                res.status(500);
                res.send()
                return;
            }
            let resultChannel = result.finalChannel;
            switch (req.method) {
                case 'GET':
                    let user = JSON.parse(req.headers['x-user'])
                    if (!isChannelMember(resultChannel, user.id)) {
                        res.status(403);
                        res.set("Content-Type", "text/plain");
                        res.send("Cannot get messages")
                        break;
                    }
                    // QUERY for last 100 messages here
                    if (req.params.before != null) {
                        mongo.last100SpecificMessages(messages, resultChannel.id, req.params.before, res).then((returnedMessages) => {
                            if (returnedMessages == null) {
                                res.status(500);
                                res.send()
                                return;
                            }
                            res.set("Content-Type", "application/json");
                            res.json(returnedMessages);
                            res.send()
                            return;
                        })
                        
                    } else {
                        mongo.last100Messages(messages, resultChannel.id, res).then((returnedMessages) => {
                            if (returnedMessages == null) {
                                res.status(500);
                                res.send()
                                return;
                            }
                            res.set("Content-Type", "application/json");
                            res.json(returnedMessages);
                            res.send()
                            return;
                        });
                        
                    }
                    break;    
                case 'POST':
                    user = JSON.parse(req.headers['x-user'])
                    if (resultChannel.private && !isChannelMember(resultChannel, user.id)) {
                        res.status(403);
                        res.set("Content-Type", "text/plain");
                        res.send("Cannot post message")
                        break;
                    }
                    // Create a new message
                    // Call database to INSERT a new message to the channel
                    let newUser = JSON.parse(req.headers['x-user'])
                    let newMessage = createMessage(req, newUser);
                    mongo.insertNewMessage(messages, newMessage).then((insertedResult) =>{
                        if (insertedResult.errString.length > 0) {
                            res.status(500);
                            res.send()
                            return;
                        }
                        let insertedMessage = insertedResult.newMessage;
                        res.status(201);
                        res.set("Content-Type", "application/json");
                        res.json(insertedMessage);
                        // // add to rabbitMQ queue
                        // let PostObj = new RabbitObject('message-new', null, insertedMessage,
                        //     resultChannel.members, null, null)
                        // sendObjectToQueue(queue, PostObj)
                        res.send()
                        return;
                    })
                    break;
                case 'PATCH':
                    user = JSON.parse(req.headers['x-user'])
                    
                    if (!isChannelCreator(resultChannel, user.id)) {
                        res.status(403);
                        res.set("Content-Type", "text/plain");
                        res.send("Cannot amend channel")
                        break;
                    }
                    // Call database to UPDATE the channel name and/or description
                    mongo.updateChannel(channels, resultChannel, req).then((updateResult) => {
                        if (updateResult.errString.length > 0) {
                            res.status(500);
                            res.send()
                            return;
                        }
                        let updatedChannel = updateResult.existingChannel;
                        res.set("Content-Type", "application/json");
                        res.json(updatedChannel);
    
                        // add to rabbitMQ queue
                        // let PatchObj = new RabbitObject('channel-update', updatedChannel, null,
                        //     updatedChannel.members, null, null)
                        // sendObjectToQueue(queue, PatchObj)
                        res.send()
                        return;
                    })
                    break;
                case 'DELETE':
                    user = JSON.parse(req.headers['x-user'])
                    if (!isChannelCreator(resultChannel, user.id)) {
                        res.status(403);
                        res.set("Content-Type", "text/plain");
                        res.send("You cannot delete this message")
                        break;
                    }
                    // Call database to DELETE this channel
                    mongo.deleteChannel(channels, messages, resultChannel).then((result)=> {
                        if (result.length > 0) {
                            res.status(500);
                            res.send()
                            return;
                        }
                        // add to rabbitMQ queue
                        // let obj = new RabbitObject('channel-delete', null, null, resultChannel.members,
                        //     resultChannel._id, null)
                        // sendObjectToQueue(queue, obj)
                        res.set("Content-Type", "text/plain");
                        res.send("Channel was successfully deleted");
                        return;
                    })
                    break;
                default:
                    res.status(405);
                    res.set("Content-Type", "text/plain");
                    res.send("Incorrect HTTP Header")
                    break;
            }
            return;
        }) 
    });


    app.use("/v1/channels", (req: any, res: any) => {
        // Check that the user is authenticated
        if (req.headers['x-user'] == null) {
            res.status(401);
            res.send();
            return;
        }
        switch (req.method) {
            case 'GET':
                // QUERY for all channels here
                mongo.getAllChannels(channels, res).then((result)=> {
                    if (result.successMessage == "") {
                        res.status(500);
                        res.send()
                        return;
                    }
                    res.set("Content-Type", "application/json");
                    res.json(result.resultJSON);
                    res.send();
                    return
                })
                
                break;

            case 'POST':
                if (req.body.name == null) {
                    res.status(500);
                    res.send()
                    //do something about the name property being null
                }
                let newUser = JSON.parse(req.headers['x-user'])
                
                let newChannel = createChannel(req, newUser);

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
                res.set("Content-Type", "text/plain");
                res.send("Incorrect HTTP Header")
                break;
        }
    });

    // Editing the body of or deleting a message
    app.use("/v1/messages/:messageID", (req: any, res: any) => {
        // Check that the user is authenticated
        if (req.headers['x-user'] == null) {
            res.status(401);
            res.send();
            return;
        }
        if (req.params.messageID == null) {
            res.status(404);
            res.send()

            return;
        }
        mongo.getMessageByID(messages, req.params.messageID).then((result) => {
            if (result.errString.length > 0) {
                res.status(500);
                res.send()
                return;
            }
            let resultMessage = result.finalMessage;
            switch (req.method) {
                case 'PATCH':
                    let user = JSON.parse(req.headers['x-user'])
                    if (!isMessageCreator(resultMessage, user.id)) {
                        res.status(403);
                        res.set("Content-Type", "text/plain");
                        res.send("Cannot update message");
                        break;
                    }
                    // TODO: Call the database to UPDATE the message in the database using the messageID
                    mongo.updateMessage(messages, resultMessage, req).then(updatedResult => {
                        if (updatedResult.errString.length > 0) {
                            res.status(500);
                            res.send();
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
                        res.send()
                        return;
                    })
                    break;
                case 'DELETE':
                    user = JSON.parse(req.headers['x-user'])
                    if (!isMessageCreator(resultMessage, user.id)) {
                        res.status(403)
                        res.set("Content-Type", "text/plain");
                        res.send("Cannot delete message");
                        break;
                    }
                    // Call database to DELETE the specified message using the messageID
                    // Call database to DELETE this channel
                    mongo.deleteMessage(messages, resultMessage).then(result => {
                        if (result.length > 0) {
                            res.status(500);
                            res.send();
                            return;
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
                    res.set("Content-Type", "text/plain");
                    res.send("Incorrect HTTP Header")
                    break;
            }
        })  
    });

    function createChannel(req: any, creator: User): Channel {
        let c = req.body;

        c.members.push(creator.id)
        return new Channel(c.name, c.description, c.private,
            c.members, c.createdAt, creator, c.editedAt);
    }

    function createMessage(req: any, creator: User): Message {
        let m = req.body;
        return new Message(req.params.channelID, m.createdAt, m.body,
            creator, m.editedAt);
    }

    function isChannelMember(channel: Channel, userID: number): boolean {
        let isMember = false;
        if (channel.private) { 
            for (let i = 0; i < channel.members.length; i++) {
                if (channel.members[i] === userID) {
                    isMember = true;
                    break;
                }
            }
        } else {
            isMember = true;
        }
        return isMember;
    }

    function isChannelCreator(channel: Channel, userID: number): boolean {
        return channel.creator.id === userID;
    }

    function isMessageCreator(message: Message, userID: number): boolean {
        return message.creator.id === userID;
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