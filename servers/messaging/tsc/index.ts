"use strict";
// version 0.1

//require the express and morgan packages
import express from "express";
import morgan from "morgan";
import { Collection, Db } from "mongodb";
import * as mongo from "./mongo_handlers";
import { Message, isMessageCreator } from "./message";
import { Channel, isChannelCreator, isChannelMember } from "./channel";
import { RabbitObject, sendObjectToQueue, createMQChannel, createMQConnection } from "./rabbit";
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

var db: Db;
var messages: Collection;
var channels: Collection;

const main = async () => {
    const db = await mongo.createConnection();
    
    // These collections save interactions so that users may log in later
    // and review old conversations/pick up where things left off
    var channels = db.collection("channels");
    var messages = db.collection("messages");

    // This queue sends all the pertinent messages back to the users through the gateway
    // via a websockets connection
    let mqClient = await createMQConnection();
    let mqChannel = await createMQChannel(mqClient);

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

    app.use("/v1/channels/:channelID/members", async (req: any, res: any) => {
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
        let result = await mongo.getChannelByID(channels, req.params.channelID)
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
                let addErr = await mongo.addChannelMember(channels, resultChannel, req)
                if (addErr) {
                    res.status(500);
                    res.send()
                    return;
                }
                res.set("Content-Type", "text/plain");
                res.status(201)
                res.send(user.ID + " was added to your channel");
                break;
            case 'DELETE':
                if (!isChannelCreator(resultChannel, user.ID)) {
                    res.status(403)
                    res.set("Content-Type", "text/plain");
                    res.send("Cannot delete members")
                    break;
                }
                // database to UPDATE the current channel members
                let removeErr = await mongo.removeChannelMember(channels, resultChannel, req)
                if (removeErr) {
                    res.status(500);
                    res.send()
                    return;
                }
                res.set("Content-Type", "text/plain");
                res.status(201)
                res.send(user.ID + " was removed from your channel");
                break;
            default:
                res.status(405);
                res.set("Content-Type", "text/plain");
                res.send("Method Not Allowed")
                break;
        }
    });

    // Specific channel handler
    app.use("/v1/channels/:channelID", async (req: any, res: any) => {

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
        let result = await mongo.getChannelByID(channels, req.params.channelID)
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
                let result
                if (!req.params.before) {
                    result = await mongo.last100Messages(messages, resultChannel.id, "")
                } else {
                    result = await mongo.last100Messages(messages, resultChannel.id, req.params.before)
                }
                if (result.err) {
                    res.status(500);
                    res.send()
                    return;
                }
                res.set("Content-Type", "application/json");
                res.json(result.last100messages);
                res.send()
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
                let insertResult = await mongo.insertNewMessage(messages, newMessage)
                if (insertResult.err) {
                    res.status(500);
                    res.send()
                    return;
                }
                let insertedMessage = insertResult.newMessage;
                // add to rabbitMQ queue
                let postMembers = (resultChannel.private ? resultChannel.members : null)

                let post = new RabbitObject('message-new', null, insertedMessage,
                    postMembers, null, null)
                // sendObjectToQueue(queue, post)
                sendObjectToQueue(mqChannel, post)

                res.status(201);
                res.set("Content-Type", "application/json");
                res.json(insertedMessage);
                res.send()
                break;
            case 'PATCH':
                if (!isChannelCreator(resultChannel, user.ID)) {
                    res.status(403);
                    res.set("Content-Type", "text/plain");
                    res.send("Cannot amend channel")
                    break;
                }
                // Call database to UPDATE the channel name and/or description
                let updateResult = await mongo.updateChannel(channels, resultChannel, req)
                if (updateResult.err) {
                    res.status(500);
                    res.send()
                    return;
                }
                let updatedChannel = updateResult.existingChannel;
                // add to rabbitMQ queue
                let patchMembers = (updatedChannel.private ? updatedChannel.members : null)

                let PatchObj = new RabbitObject('channel-update', updatedChannel, null,
                    patchMembers, null, null)
                // sendObjectToQueue(queue, PatchObj)
                sendObjectToQueue(mqChannel, PatchObj)
                
                res.set("Content-Type", "application/json");
                res.json(updatedChannel);
                res.send()
                break;
            case 'DELETE':
                if (!isChannelCreator(resultChannel, user.ID)) {
                    res.status(403);
                    res.set("Content-Type", "text/plain");
                    res.send("You cannot delete this channel")
                    break;
                }
                // Call database to DELETE this channel
                let err = await mongo.deleteChannel(channels, messages, resultChannel)
                if (err) {
                    res.status(500);
                    res.send()
                    return;
                }
                // add to rabbitMQ queue
                let deleteMembers = (resultChannel.private ? resultChannel.members : null)

                let obj = new RabbitObject('channel-delete', null, null, deleteMembers,
                    resultChannel.id, null)
                // sendObjectToQueue(queue, obj)
                sendObjectToQueue(mqChannel, obj)

                res.set("Content-Type", "text/plain");
                res.send("Channel was successfully deleted");
                break;
            default:
                res.status(405);
                res.set("Content-Type", "text/plain");
                res.send("Method Not Allowed")
                break;
        }
    });


    app.use("/v1/channels", async (req: any, res: any) => {
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
                let getResult
                if (!req.params.startsWith) {
                    getResult = await mongo.getChannels(channels, user.ID, "")
                } else {
                    getResult = await mongo.getChannels(channels, user.ID, req.params.startsWith)
                }
                if (getResult.err) {
                    res.status(500);
                    res.send()
                    return;
                }
                res.set("Content-Type", "application/json");
                res.json(getResult.allChannels);
                res.send();
                break
            case 'POST':
                if (req.body.name == null) {
                    res.status(500);
                    res.send();
                    break;
                    //do something about the name property being null
                }
                let newChannel = createChannel(req, user);

                let insertResult = await mongo.insertNewChannel(channels, newChannel)
                if (insertResult.duplicates) {
                    res.status(400);
                    res.send();
                    return;
                }

                if (insertResult.err) {
                    res.status(500);
                    res.send();
                    return;
                }
                let insertChannel = insertResult.newChannel;
                // add to rabbitMQ queue
                let members = (insertChannel.private ? insertChannel.members : null)

                let post = new RabbitObject('channel-new', insertChannel, null,
                    members, null, null)
                // sendObjectToQueue(queue, post)
                sendObjectToQueue(mqChannel, post)

                res.status(201);
                res.set("Content-Type", "application/json");
                res.json(insertChannel);
                res.send();
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

                mongo.getChannelByID(channels, updatedMessage.channelID).then((result) => {
                    // add to rabbitMQ queue
                    let members = (result.channel.private ? result.channel.members : null)

                    let post = new RabbitObject('message-update', null, updatedMessage,
                        members, null, null)
                    // sendObjectToQueue(queue, post)
                    sendObjectToQueue(mqChannel, post)
                })
                res.set("Content-Type", "application/json");
                res.json(updatedMessage);
                res.send();
                break;
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
                    let members = (result.channel.private ? result.channel.members : null)

                    let post = new RabbitObject('message-delete', null, null,
                        members, null, resultMessage.id)
                    // sendObjectToQueue(queue, post)
                    sendObjectToQueue(mqChannel, post)
                })

                res.set("Content-Type", "text/plain");
                res.send("Message deleted");
                break;
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