// "use strict";

// to compile run tsc --outDir ../

//require the express and morgan packages
import express from "express";
import morgan from "morgan";
import { MongoClient, Db, Collection, Server } from "mongodb";
import * as mongo from "./mongo_handlers";
import { Message } from "./message";
import { Channel } from "./channel";

//create a new express application
const app = express();

const addr = process.env.ADDR || ":80";
//split host and port using destructuring
const [host, port] = addr.split(":");

//add JSON request body parsing middleware
app.use(express.json());
//add the request logging middleware
app.use(morgan("dev"));

// Connection URL
const url = 'mongo://mongodb:27017/mongodb';

// Database Name
const dbName = 'mongodb';

// Create a new MongoClient
// const client = new MongoClient(url);
// var db: Db;
var messages: any;
var channels: Collection;
//new Server("mongo://mongodb", 27017)
// var mc = new MongoClient("mongo://mongodb:27017", {native_parser:true})
// mc.
// Reasoning for refactor: 
// https://mongodb.github.io/node-mongodb-native/driver-articles/mongoclient.html#mongoclient-connection-pooling
// Use connect method to connect to the mongo DB
// client.connect(function (err: any) {
MongoClient.connect(url, function (err: any, client:MongoClient) {
    console.log("Connected successfully to server");
    const db = client.db(dbName);
    // db = client.db(dbName);

    // check if any collection exists
    db.collections()
    .then(doc => {
        console.log(doc)
    }).catch(err => {
        console.log(err)
    });

    // Start the application after the database connection is ready
    app.listen(+port, "", () => {
        //callback is executed once server is listening
        console.log(`server is listening at http://:${port}...`);
        console.log("port : " + port);
        console.log("host : " + host);
    });
});

// All channel handler
// No errors here :)
app.use("/v1/channels", (req: any, res: any, next: any) => {
    switch (req.method) {
        case 'GET':
            res.set("Content-Type", "application/json");
            // QUERY for all channels here
            let allChannels = mongo.getAllChannels(channels);
            if (allChannels == null) {
                res.status(500);
            }
            // write those to the client, encoded in JSON
            res.json(allChannels);
            break;

        case 'POST':
            console.log(req.body);
            if (req.body.channel.name == null) {
                res.status(500);
                //do something about the name property being null
            }
            // Call database to INSERT this new channel
            let newChannel = createChannel(req);
            let insertResult = mongo.insertNewChannel(channels, newChannel);
            if (insertResult.errString.length > 0) {
                res.status(500);
            }
            let insertChannel = insertResult.newChannel;
            res.set("Content-Type", "application/json");
            res.json(insertChannel);
            res.status(201);  //probably cant do this >>> .send("success");
            break;
        default:
            break;
    }
});

// Specific channel handler
app.use("/v1/channels/:channelID", (req: any, res: any, next: any) => {
    // QUERY for the channel based on req.params.channelID
    if (req.params.channelID == null) {
        res.status(404);
        return;
    }
    let result = mongo.getChannelByID(channels, req.params.channelID);
    if (result.errString.length() > 0) {
        res.status(500);
        return;
    }
    let resultChannel = result.finalChannel;
    switch (req.method) {
        case 'GET':
            if (!isChannelMember(resultChannel, req.Header.Xuser)) {
                res.status(403);
                break;
            }
            let returnedMessages;
            // QUERY for last 100 messages here
            if (req.params.before != null) {
                returnedMessages = mongo.last100SpecificMessages(messages, resultChannel._id, req.params.before);
                if (returnedMessages == null) {
                    res.status(500);
                    break;
                }
            } else {
                returnedMessages = mongo.last100Messages(messages, resultChannel._id);
                if (returnedMessages == null) {
                    res.status(500);
                    break;
                }
            }
            res.set("Content-Type", "application/json");
            // write last 100 messages to the client, encoded in JSON 
            res.json(returnedMessages);
            break;

        case 'POST':
            if (!isChannelMember(resultChannel, req.Header.Xuser)) {
                res.status(403);
                break;
            }
            // Create a new message
            // Call database to INSERT a new message to the channel
            let newMessage = createMessage(req);
            let insertedResult = mongo.insertNewMessage(messages, newMessage);
            if (insertedResult.errString.length > 0) {
                res.status(500);
            }
            let insertedMessage = insertedResult.newMessage;
            res.set("Content-Type", "application/json");
            res.json(insertedMessage);
            res.status(201);  // probably cant do this >>> .send("success");
            break;
        case 'PATCH':
            if (!isChannelCreator(resultChannel, req.Header.Xuser)) {
                res.status(403);
                break;
            }
            // Call database to UPDATE the channel name and/or description
            let updateResult = mongo.updateChannel(channels, resultChannel, req);
            if (updateResult.errString.length > 0) {
                res.status(500);
            }
            let updatedChannel = updateResult.existingChannel;
            res.set("Content-Type", "application/json");
            res.json(updatedChannel);
            break;
        case 'DELETE':
            if (!isChannelCreator(resultChannel, req.Header.Xuser)) {
                res.status(403);
                break;
            }
            // Call database to DELETE this channel
            let result = mongo.deleteChannel(channels, messages, resultChannel);
            if (result.length > 0) {
                res.status(500);
            }
            res.set("Content-Type", "text/plain");
            res.send("Channel was successfully deleted");
            break;
        default:
            break;
    }
});

// Adding and removing members from your channel
app.use("/v1/channels/:channelID/members", (req: any, res: any, next: any) => {
    // QUERY for the channel based on req.params.channelID
    if (req.params.channelID == null) {
        res.status(404);
        return;
    }
    let result = mongo.getChannelByID(channels, req.params.channelID);
    if (result.errString.length() > 0) {
        res.status(500);
        return;
    }
    let resultChannel = result.finalChannel;
    switch (req.method) {
        case 'POST':
            if (!isChannelCreator(resultChannel, req.Header.Xuser)) {
                res.status(403);
                break;
            }
            // Call database to UPDATE the current channel
            let addResult = mongo.addChannelMember(channels, resultChannel, req);
            if (addResult.length > 0) {
                res.status(500);
                break;
            }
            res.set("Content-Type", "application/json");
            res.status(201).send(req.user.ID + " was added to your channel");
            break;
        case 'DELETE':
            if (!isChannelCreator(resultChannel, req.Header['X-user'])) {
                res.status(403)
                break;
            }
            // database to UPDATE the current channel members
            let errResult = mongo.removeChannelMember(channels, resultChannel, req);
            if (errResult.length > 0) {
                res.status(500);
                break;
            }
            res.set("Content-Type", "text/plain");
            res.status(200).send(req.user.ID + " was removed from your channel");
            break;
        default:
            break;
    }
});

// Editing the body of or deleting a message
app.use("/v1/messages/:messageID", (req: any, res: any, next: any) => {
    if (req.params.messageID == null) {
        res.status(404);
        return;
    }
    let result = mongo.getMessageByID(messages, req.params.messageID);
    if (result.errString.length() > 0) {
        res.status(500);
        return;
    }
    let resultMessage = result.finalMessage;
    switch (req.method) {
        case 'PATCH':
            if (!isMessageCreator(resultMessage, req.Header.Xuser)) {
                res.status(403);
                break;
            }
            // TODO: Call the database to UPDATE the message in the database using the messageID
            let updatedResult = mongo.updateMessage(messages, resultMessage, req);
            if (updatedResult.errString.length > 0) {
                res.status(500);
                break;
            }
            let updatedMessage = updatedResult.existingMessage;
            res.set("Content-Type", "application/json");
            res.json(updatedMessage);
            break;
        case 'DELETE':
            if (!isMessageCreator(resultMessage, req.Header.Xuser)) {
                res.status(403)
                break;
            }
            // Call database to DELETE the specified message using the messageID
            // Call database to DELETE this channel
            let result = mongo.deleteMessage(messages, resultMessage);
            if (result.length > 0) {
                res.status(500);
            }
            res.set("Content-Type", "text/plain");
            res.send("Message deleted");
            break;
        default:
            break;
    }
});

function createChannel(req: any): Channel {
    let c = req.body.channel;
    return new Channel(c.name, c.description, c.private,
        c.members, c.createdAt, c.creator, c.editedAt);
}

function createMessage(req: any): Message {
    let m = req.body.message;
    return new Message(req.params.ChannelID, m.createdAt, m.body,
        m.creator, m.editedAt);
}

function isChannelMember(channel: Channel, user: any): boolean {
    let isMember = false;
    if (channel.private) {
        for (let i = 0; i < channel.members.length; i++) {
            if (channel.members[i] == user.ID) {
                isMember = true;
                break;
            }
        }
    } else {
        isMember = true;
    }
    return isMember;
}

function isChannelCreator(channel: Channel, user: any): boolean {
    return channel.creator == user._id;
}

function isMessageCreator(message: Message, user: any): boolean {
    return message.creator == user._id;
}

//error handler that will be called if
//any handler earlier in the chain throws
//an exception or passes an error to next()
app.use((err: any, req: any, res: any, next: any) => {
    //write a stack trace to standard out,
    //which writes to the server's log
    console.error(err.stack)

    //but only report the error message
    //to the client, with a 500 status code
    res.set("Content-Type", "text/plain");
    res.status(500).send(err.message);
});