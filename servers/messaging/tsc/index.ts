"use strict";
// version 0.1

//require the express and morgan packages
import express from "express";
import morgan from "morgan";
import { MongoClient, Collection, Db } from "mongodb";
import * as mongo from "./mongo_handlers";
import { Message, isMessageCreator } from "./message";
import { Channel, isChannelCreator, isChannelMember } from "./channel";
import { RabbitObject, sendObjectToQueue } from "./rabbit";
import { User } from "./user";

import * as Amqp from "amqp-ts"
// import * as amqp from "amqplib"

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

var db: Db;
var messages: Collection;
var channels: Collection;

// Rabbit MQ variables
const mqName = "helloQueue"
const mqURL = "amqp://rabbitMQ"

// Mongo DB variables
const dbName = 'mongodb';
const mongoURL = 'mongodb://mongodb:27017/mongodb';

// Create a new MongoClient
const mc = new MongoClient(mongoURL, { useUnifiedTopology: true });

async function createConnection() {
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

const main = async () => {
    const client = await createConnection();
    const db = client.db(dbName);
    var channels = db.collection("channels");
    var messages = db.collection("messages");

    let rabbitConn: Amqp.Connection;
    rabbitConn = new Amqp.Connection(mqURL);
    let queue = rabbitConn.declareQueue(mqName);

    // let mqClient = await createMQConnection();
    // let mqChannel = await createMQChannel(mqClient);
    // mqChannel.assertQueue(mqName)
    // mqChannel.sendToQueue(mqName, Buffer.from('something to do'));

    // channel.assertQueue(queue, {
    //     durable: false
    // });
    // channel.sendToQueue(queue, Buffer.from(msg));

    app.listen(+addr, "", (req, res) => {
        //callback is executed once server is listening
        console.log(`server is listening at http://:${addr}...`);
    });

    // function isAuthenticated(req: any) {
    //     return req.headers['x-user'] != null
    // }

    // app.all('*', function preflightCheck(req, res, next){
    //     if (!isAuthenticated(req)) {
    //         res.status(401);
    //         res.send()
    //         next(new Error("401 Unauthorized"))
    //     } else {
    //         next();
    //     }
    // });

    app.use("/v1/channels/:channelID/members", (req: any, res: any) => {
        // Check that the user is authenticated
        if (req.headers['x-user'] == null) {
            res.status(401);
            res.send()
            return;
        }
        const user = JSON.parse(req.headers['x-user'])
        // QUERY for the channel based on req.params.channelID
        if (req.params.channelID == null) {
            res.status(404);
            res.send()
            return;
        }
        mongo.getChannelByID(channels, req.params.channelID).then((result) => {
            if (result.err) {
                res.status(500);
                res.send()
                return;
            }
            const resultChannel = result.channel;
            switch (req.method) {
                case 'POST':
                    if (!isChannelCreator(resultChannel, user.ID)) {
                        res.status(403);
                        res.set("Content-Type", "text/plain");
                        res.send("Cannot change members")
                        break;
                    }
                    // Call database to UPDATE the current channel
                    mongo.addChannelMember(channels, resultChannel, req).then((err) => {
                        if (err) {
                            res.status(500);
                            res.send()
                            return;
                        }
                        res.set("Content-Type", "text/plain");
                        res.status(201)
                        res.send(user.ID + " was added to your channel");
                        return;
                    })
                    break;
                case 'DELETE':
                    if (!isChannelCreator(resultChannel, user.ID)) {
                        res.status(403)
                        res.set("Content-Type", "text/plain");
                        res.send("Cannot delete members")
                        break;
                    }
                    // database to UPDATE the current channel members
                    mongo.removeChannelMember(channels, resultChannel, req).then(err => {
                        if (err) {
                            res.status(500);
                            res.send()
                            return;
                        }
                        res.set("Content-Type", "text/plain");
                        res.status(201)
                        res.send(user.ID + " was removed from your channel");
                        return;
                    })
                    break;
                default:
                    res.status(405);
                    res.set("Content-Type", "text/plain");
                    res.send("Method Not Allowed")
                    break;
            }
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
        const user = JSON.parse(req.headers['x-user'])
        // QUERY for the channel based on req.params.channelID
        if (req.params.channelID == null) {
            res.status(404);
            res.send()
            return;
        }
        // req.params
        mongo.getChannelByID(channels, req.params.channelID).then((result) => {
            if (result.err) {
                res.status(500);
                res.send()
                return;
            }
            let resultChannel = result.channel;
            switch (req.method) {
                case 'GET':
                    if (!isChannelMember(resultChannel, user.ID)) {
                        res.status(403);
                        res.set("Content-Type", "text/plain");
                        res.send("Cannot get messages")
                        break;
                    }
                    // QUERY for last 100 messages here
                    if (req.params.before != null) {
                        mongo.last100SpecificMessages(messages, resultChannel.id, req.params.before).then((result) => {
                            if (result.err) {
                                res.status(500);
                                res.send()
                                return;
                            }
                            res.set("Content-Type", "application/json");
                            res.json(result.last100messages);
                            res.send()
                            return;
                        })
                    } else {
                        mongo.last100Messages(messages, resultChannel.id).then((result) => {
                            if (result.err) {
                                res.status(500);
                                res.send()
                                return;
                            }
                            res.set("Content-Type", "application/json");
                            res.json(result.last100messages);
                            res.send()
                            return;
                        });
                    }
                    break;
                case 'POST':
                    if (!isChannelMember(resultChannel, user.ID)) {
                        res.status(403);
                        res.set("Content-Type", "text/plain");
                        res.send("Cannot post message")
                        break;
                    }
                    // Create a new message
                    // Call database to INSERT a new message to the channel
                    let newMessage = createMessage(req, user);
                    mongo.insertNewMessage(messages, newMessage).then((result) => {
                        if (result.err) {
                            res.status(500);
                            res.send()
                            return;
                        }
                        let insertedMessage = result.newMessage;
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
                    if (!isChannelCreator(resultChannel, user.ID)) {
                        res.status(403);
                        res.set("Content-Type", "text/plain");
                        res.send("Cannot amend channel")
                        break;
                    }
                    // Call database to UPDATE the channel name and/or description
                    mongo.updateChannel(channels, resultChannel, req).then((result) => {
                        if (result.err) {
                            res.status(500);
                            res.send()
                            return;
                        }
                        let updatedChannel = result.existingChannel;
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
                    if (!isChannelCreator(resultChannel, user.ID)) {
                        res.status(403);
                        res.set("Content-Type", "text/plain");
                        res.send("You cannot delete this channel")
                        break;
                    }
                    // Call database to DELETE this channel
                    mongo.deleteChannel(channels, messages, resultChannel).then((err) => {
                        if (err) {
                            res.status(500);
                            res.send()
                            return;
                        }
                        // add to rabbitMQ queue
                        // let obj = new RabbitObject('channel-delete', null, null, resultChannel.members,
                        //     resultChannel.id, null)
                        // sendObjectToQueue(queue, obj)
                        res.set("Content-Type", "text/plain");
                        res.send("Channel was successfully deleted");
                        return;
                    })
                    break;
                default:
                    res.status(405);
                    res.set("Content-Type", "text/plain");
                    res.send("Method Not Allowed")
                    break;
            }
        })
    });


    app.use("/v1/channels", (req: any, res: any) => {
        // Check that the user is authenticated
        if (req.headers['x-user'] == null) {
            res.status(401);
            res.send();
            return;
        }
        const user = JSON.parse(req.headers['x-user'])
        switch (req.method) {
            case 'GET':
                // QUERY for all channels here
                mongo.getAllChannels(channels, user.ID).then((result) => {
                    if (result.err) {
                        res.status(500);
                        res.send()
                        return;
                    }
                    res.set("Content-Type", "application/json");
                    res.json(result.allChannels);
                    res.send();
                    return;
                })
                break;
            case 'POST':
                if (req.body.name == null) {
                    res.status(500);
                    res.send();
                    break;
                    //do something about the name property being null
                }
                let newChannel = createChannel(req, user);

                mongo.insertNewChannel(channels, newChannel).then(result => {
                    if (result.duplicates) {
                        res.status(400);
                        res.send();
                        return;
                    }

                    if (result.err) {
                        res.status(500);
                        res.send();
                        return;
                    }
                    let insertChannel = result.newChannel;
                    // // add to rabbitMQ queue
                    // let obj = new RabbitObject('channel-new', insertChannel, null,
                    //     insertChannel.members, null, null)
                    // sendObjectToQueue(queue, obj)
                    res.status(201);
                    res.set("Content-Type", "application/json");
                    res.json(insertChannel);
                    res.send();
                    return;
                })
                break;
            default:
                res.status(405);
                res.set("Content-Type", "text/plain");
                res.send("Method Not Allowed")
                break;
        }
    });

    // Editing the body of or deleting a message
    app.use("/v1/messages/:messageID", async (req: any, res: any) => {
        // Check that the user is authenticated
        if (req.headers['x-user'] == null) {
            res.status(401);
            res.send();
            return;
        }
        const user = JSON.parse(req.headers['x-user'])
        if (req.params.messageID == null) {
            res.status(404);
            res.send()
            return;
        }
        let result = await mongo.getMessageByID(messages, req.params.messageID)
        if (result.err) {
            res.status(500);
            res.send()
            return;
        }
        // Can we use this as a const?
        let resultMessage = result.message;
        switch (req.method) {
            case 'PATCH':
                if (!isMessageCreator(resultMessage, user.ID)) {
                    res.status(403);
                    res.set("Content-Type", "text/plain");
                    res.send("Cannot update message");
                    break;
                }
                // TODO: Call the database to UPDATE the message in the database using the messageID
                let result = await mongo.updateMessage(messages, resultMessage, req)
                if (result.err) {
                    res.status(500);
                    res.send();
                    return;
                }
                let updatedMessage = result.existingMessage;
                res.set("Content-Type", "application/json");
                res.json(updatedMessage);

                mongo.getChannelByID(channels, updatedMessage.channelID).then((result) => {
                    // add to rabbitMQ queue
                    let members = null
                    if (result.channel.private) {
                        members = result.channel.members
                    }
                    let post = new RabbitObject('message-update', null, updatedMessage,
                        result.channel.members, null, null)
                    sendObjectToQueue(queue, post)
                })
                res.send();
                return;
            case 'DELETE':
                if (!isMessageCreator(resultMessage, user.ID)) {
                    res.status(403);
                    res.set("Content-Type", "text/plain");
                    res.send("Cannot delete message");
                    break;
                }
                // Call database to DELETE the specified message using the messageID
                let err = await mongo.deleteMessage(messages, resultMessage)
                if (err) {
                    res.status(500);
                    res.send();
                    return;
                }
                mongo.getChannelByID(channels, resultMessage.channelID).then((result) => {
                    // add to rabbitMQ queue
                    let members = null
                    if (result.channel.private) {
                        members = result.channel.members
                    }
                    let post = new RabbitObject('message-delete', null, null,
                        members, null, resultMessage.id)
                    sendObjectToQueue(queue, post)
                })

                res.set("Content-Type", "text/plain");
                res.send("Message deleted");
                return;
            default:
                res.status(405);
                res.set("Content-Type", "text/plain");
                res.send("Method Not Allowed");
                break;
        }
    })

    function createChannel(req: any, creator: User): Channel {
        let c = req.body;

        c.members.push(creator.ID)
        return new Channel("", c.name, c.description, c.private,
            c.members, c.createdAt, creator, c.editedAt);
    }

    function createMessage(req: any, creator: User): Message {
        let m = req.body;
        return new Message("", req.params.channelID, m.createdAt, m.body,
            creator, m.editedAt);
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