"use strict";
// version 0.1

import express, { Request, Response } from "express";
import { IncomingHttpHeaders } from 'http';
import morgan from "morgan";
import * as mongo from "./mongo_handlers";
import { Message, isMessageCreator } from "./message";
import { Channel, isChannelCreator, isChannelMember } from "./channel";
import { sendObjectToQueue, createMQChannel, createMQConnection, ChannelTransaction, MessageTransaction } from "./rabbit";
import { User } from "./user";

// Creates a new express application
const app = express();
const addr = process.env.ADDR || "80";

app.use(express.json()); // Adds JSON request body parsing middleware
app.use(morgan("dev")); // Adds the request logging middleware

const main = async () => {
    const db = await mongo.createConnection();
    
    // These collections save interactions so that users may log in later
    // and review old conversations/pick up where things left off
    const channels = db.collection("channels");
    const messages = db.collection("messages");

    // This queue sends all the pertinent messages back to the users through the gateway
    // via a websockets connection
    const mqClient = await createMQConnection();
    const mqChannel = await createMQChannel(mqClient);

    app.listen(+addr, "", () => {
        // This callback is executed once server is listening
        console.log(`server is listening at http://:${addr}...`);
    });

    const isAuthenticated = (req: Request<any>) => {
        return req.headers['x-user'] != null
    }

    // Authentication check
    app.use((req: Request<any>, res: Response, next) => {
        if (!isAuthenticated(req)) {
            res.status(401)
            res.send()
            // If you pass anything to the next() function (except the string 'route' or 'router'), Express regards the
            // current request as being an error and will skip any remaining non-error handling routing and middleware
            // functions. (source: https://expressjs.com/en/guide/writing-middleware.html)
            // next(new Error("401 Unauthorized"))
        } else {
            next()
        }
    });

    // Channel Members Middleware (add & remove)
    app.use("/v1/channels/:channelID/members", async (req: Request<any>, res: Response) => {
        const { headers, params, method, body } = req;
        // We can safely assume that x-user is populated sinced we ran an authentication check
        const user = getXUser(headers);

        // QUERY for the channel based on params.channelID
        if (params.channelID == null) {
            res.status(404);
            res.send()
            return;
        }
        const { channel, err } = await mongo.getChannelByID(channels, params.channelID)
        if (err || !channel) {
            res.status(500);
            res.send()
            return;
        }
        switch (method) {
            // Case statement brackets to scope variables to case 
            // (source: https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Statements/switch)
            case 'POST': {
                if (!isChannelCreator(channel, user.ID)) {
                    res.status(403);
                    res.setHeader("Content-Type", "text/plain");
                    res.send("Cannot change members")
                    break;
                }
                // Call database to UPDATE the current channel
                const err = await mongo.addChannelMember(channels, channel, body.id)
                if (err) {
                    res.status(500);
                    res.send()
                    return;
                }
                res.status(201)
                res.setHeader("Content-Type", "text/plain");
                res.send(user.ID + " was added to your channel");
                break;
            }
            case 'DELETE': {
                if (!isChannelCreator(channel, user.ID)) {
                    res.status(403)
                    res.setHeader("Content-Type", "text/plain");
                    res.send("Cannot delete members")
                    break;
                }
                // database to UPDATE the current channel members
                const err = await mongo.removeChannelMember(channels, channel, body.id)
                if (err) {
                    res.status(500);
                    res.send()
                    return;
                }
                res.status(201)
                res.setHeader("Content-Type", "text/plain");
                res.send(user.ID + " was removed from your channel");
                break;
            }
            default:
                res.status(405);
                res.setHeader("Content-Type", "text/plain");
                res.send("Method Not Allowed")
                break;
        }
    });

    // Specific Channel Middleware (get & add messages, update & delete channels)
    app.use("/v1/channels/:channelID", async (req: Request<any>, res: Response) => {
        const { headers, params, method, body } = req;
        // We can safely assume that x-user is populated sinced we ran an authentication check
        const user = getXUser(headers);

        // QUERY for the channel based on params.channelID
        if (params.channelID == null) {
            res.status(404);
            res.send()
            return;
        }
        const { channel, err } = await mongo.getChannelByID(channels, params.channelID)
        if (err || !channel) {
            res.status(500);
            res.send()
            return;
        }
        switch (method) {
            case 'GET': {
                if (!isChannelMember(channel, user.ID)) {
                    res.status(403);
                    res.setHeader("Content-Type", "text/plain");
                    res.send("Cannot get messages")
                    break;
                }
                // QUERY for last 100 messages here
                const messageID = params.before || "";
                const { last100messages, err } = await mongo.last100Messages(messages, channel.id, messageID)

                if (err) {
                    res.status(500);
                    res.send()
                    return;
                }
                res.setHeader("Content-Type", "application/json");
                res.json(last100messages);
                break;
            }
            case 'POST': {
                if (!isChannelMember(channel, user.ID)) {
                    res.status(403);
                    res.setHeader("Content-Type", "text/plain");
                    res.send("Cannot post message")
                    break;
                }
                // Create a new message
                // Call database to INSERT a new message to the channel
                const message = createMessage(req, user);
                const { newMessage, err } = await mongo.insertNewMessage(messages, message)
                if (err) {
                    res.status(500);
                    res.send()
                    return;
                }

                // add to rabbitMQ queue
                const messageCreatedTransaction: MessageTransaction = {
                    entity: newMessage,
                    type: "message-new",
                    userIDs: channel.members,
                }
                sendObjectToQueue(mqChannel, messageCreatedTransaction)

                res.status(201);
                res.setHeader("Content-Type", "application/json");
                res.json(newMessage);
                break;
            }
            case 'PATCH': {
                if (!isChannelCreator(channel, user.ID)) {
                    res.status(403);
                    res.setHeader("Content-Type", "text/plain");
                    res.send("Cannot amend channel")
                    break;
                }
                // Call database to UPDATE the channel name and/or description
                const { updatedChannel, err } = await mongo.updateChannel(channels, channel, body)
                if (err) {
                    res.status(500);
                    res.send()
                    return;
                }
                
                // add to rabbitMQ queue
                const channelUpdatedTransaction: ChannelTransaction = {
                    entity: updatedChannel,
                    type: "channel-update",
                    userIDs: updatedChannel.members,
                }
                sendObjectToQueue(mqChannel, channelUpdatedTransaction)
                
                res.setHeader("Content-Type", "application/json");
                res.json(updatedChannel);
                break;
            }
            case 'DELETE': {
                if (!isChannelCreator(channel, user.ID)) {
                    res.status(403);
                    res.setHeader("Content-Type", "text/plain");
                    res.send("You cannot delete this channel")
                    break;
                }
                // Call database to DELETE this channel
                const err = await mongo.deleteChannel(channels, messages, channel)
                if (err) {
                    res.status(500);
                    res.send()
                    return;
                }
                
                // add to rabbitMQ queue
                const channelDeletedTransaction: ChannelTransaction = {
                    type: "channel-delete",
                    userIDs: channel.members,
                    channelID: channel.id,
                }
                sendObjectToQueue(mqChannel, channelDeletedTransaction)

                res.setHeader("Content-Type", "text/plain");
                res.send("Channel was successfully deleted");
                break;
            }
            default:
                res.status(405);
                res.setHeader("Content-Type", "text/plain");
                res.send("Method Not Allowed")
                break;
        }
    });

    // Channel Middleware (get & add channels)
    app.use("/v1/channels", async (req: Request<any>, res: Response) => {
        const { headers, params, method, body } = req;
        // We can safely assume that x-user is populated sinced we ran an authentication check
        const user = getXUser(headers);

        switch (method) {
            case 'GET': {
                // QUERY for all channels here
                const searchTerm = params.startsWith || "";
                const { usersChannels, err } = await mongo.getChannels(channels, user.ID, searchTerm)
                
                if (err) {
                    res.status(500);
                    res.send()
                    return;
                }
                res.setHeader("Content-Type", "application/json");
                res.json(usersChannels);
                break;
            }
            case 'POST': {
                if (body.name == null) {
                    res.status(500);
                    res.send();
                    break;
                    //do something about the name property being null
                }
                const channel = createChannel(req, user);

                const { newChannel, hasDuplicates, err } = await mongo.insertNewChannel(channels, channel)
                if (hasDuplicates) {
                    res.status(400);
                    res.send();
                    return;
                }

                if (err) {
                    res.status(500);
                    res.send();
                    return;
                }

                // add to rabbitMQ queue
                const channelCreatedTransaction: ChannelTransaction = {
                    entity: newChannel,
                    type: "channel-new",
                    userIDs: newChannel.members,
                }
                sendObjectToQueue(mqChannel, channelCreatedTransaction)

                res.status(201);
                res.setHeader("Content-Type", "application/json");
                res.json(newChannel);
                break;
            }
            default:
                res.status(405);
                res.setHeader("Content-Type", "text/plain");
                res.send("Method Not Allowed")
                break;
        }
    });

    //  // Message Middleware (update or delete a message)
    app.use("/v1/messages/:messageID", async (req: Request<any>, res: Response) => {
        const { headers, params, method, body } = req;
        // We can safely assume that x-user is populated sinced we ran an authentication check
        const user = getXUser(headers);

        if (params.messageID == null) {
            res.status(404);
            res.send()
            return;
        }
        const { message, err } = await mongo.getMessageByID(messages, params.messageID)
        if (err || !message) {
            res.status(500);
            res.send()
            return;
        }
        switch (method) {
            case 'PATCH': {
                if (!isMessageCreator(message, user.ID)) {
                    res.status(403);
                    res.setHeader("Content-Type", "text/plain");
                    res.send("Cannot update message");
                    break;
                }
                // TODO: Call the database to UPDATE the message in the database using the messageID
                const { updatedMessage, err } = await mongo.updateMessage(messages, message, body)
                if (err) {
                    res.status(500);
                    res.send();
                    return;
                }
                const { channel } = await mongo.getChannelByID(channels, updatedMessage.channelID)
                if (!channel) {
                    res.status(500);
                    res.send();
                    return;                    
                }

                // add to rabbitMQ queue
                const messageUpdatedTransaction: MessageTransaction = {
                    entity: updatedMessage,
                    type: "message-update",
                    userIDs: channel.members,
                }
                sendObjectToQueue(mqChannel, messageUpdatedTransaction)
                res.setHeader("Content-Type", "application/json");
                res.json(updatedMessage);
                break;
            }
            case 'DELETE': {
                if (!isMessageCreator(message, user.ID)) {
                    res.status(403);
                    res.setHeader("Content-Type", "text/plain");
                    res.send("Cannot delete message");
                    break;
                }
                // Call database to DELETE the specified message using the messageID
                const err = await mongo.deleteMessage(messages, message)
                if (err) {
                    res.status(500);
                    res.send();
                    return;
                }
                const { channel } = await mongo.getChannelByID(channels, message.channelID);
                if (!channel) {
                    res.status(500);
                    res.send();
                    return;                    
                }
                
                // add to rabbitMQ queue
                const messageDeletedTransaction: MessageTransaction = {
                    type: "message-delete",
                    userIDs: channel.members,
                    messageID: message.id,
                    channelID: message.channelID
                }
                sendObjectToQueue(mqChannel, messageDeletedTransaction)

                res.setHeader("Content-Type", "text/plain");
                res.send("Message deleted");
                break;
            }
            default:
                res.status(405);
                res.setHeader("Content-Type", "text/plain");
                res.send("Method Not Allowed");
                break;
        }
    })

    const getXUser = (headers: IncomingHttpHeaders) => {
        if (typeof headers['x-user'] === "string") {
            return JSON.parse(headers['x-user']);
        }
        return JSON.parse(headers['x-user']![0]);
    }

    const createChannel = (req: Request<any>, creator: User): Channel => {
        const { body } = req;
        const c = body;

        c.members.push(creator.ID)
        return new Channel("", c.name, c.description, c.private,
            c.members, c.createdAt, creator, c.editedAt);
    }

    const createMessage = (req: Request<any>, creator: User): Message => {
        const { body, params } = req;
        const m = body;
        return new Message("", params.channelID, m.createdAt, m.body,
            creator, m.editedAt);
    }

    // This error handler will be called if any handler earlier in the chain throws
    // an exception or passes an error to next()
    app.use((err: any, _req: Request<any>, res: Response) => {
        // Write a stack trace to standard out which writes to the server's log
        console.error(`index.ts:444 ${err.stack}`)

        // Report the error message to the client with a 500 status code
        res.status(500);
        res.setHeader("Content-Type", "text/plain");
        res.send(err.message);
    });

};

main();