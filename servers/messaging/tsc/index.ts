"use strict";

// to compile run tsc --outDir ../

//require the express and morgan packages
const express = require("express");
const morgan = require("morgan");
import { MongoClient, Db, Collection } from "mongodb";
import * as mongo from "./mongo_handlers";
import { Message } from "./message";
import { Channel } from "./channel";
import assert from "assert";

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
const url = 'mongodb://localhost:27017';

// Database Name
const dbName = 'messaging';

// Create a new MongoClient
const client = new MongoClient(url);
var db: Db;
var messages: Collection;
var channels: Collection;

// Reasoning for refactor: 
// https://mongodb.github.io/node-mongodb-native/driver-articles/mongoclient.html#mongoclient-connection-pooling
// Use connect method to connect to the mongo DB
client.connect(function (err: any) {
    assert.equal(null, err);
    console.log("Connected successfully to server");

    db = client.db(dbName);
    // Create db.channels and db.messages collections in mongo
    // https://mongodb.github.io/node-mongodb-native/api-articles/nodekoarticle1.html#mongo-db-and-collections
    db.createCollection('channels', function (err, collection) {
        channels = collection;
    });
    db.createCollection('messages', function (err, collection) {
        messages = collection;
    });
    var general = new Channel("general", "an open channel for all", false, [], "enter timestamp here", -1, "not yet edited");
    // channel that we always want at startup
    let result = mongo.insertNewChannel(channels, general);
    if (result == null) {
        console.log("failed to create general channel upon opening connection to DB");
        // res.status(500);
    }

    // Start the application after the database connection is ready
    app.listen(port, "", () => {
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
            if (insertResult == null) {
                res.status(500);
            }
            res.set("Content-Type", "application/json");
            res.json(insertResult);
            res.status(201);  //probably cant do this >>> .send("success");
            break;
        default:
            break;
    }
});

// Specific channel handler
app.use("/v1/channels/:channelID", (req: any, res: any, next: any) => {
    // QUERY for the channel based on req.params.channelID
    let resultChannel = mongo.getChannelByID(channels, req.params.channelID);
    if (resultChannel === null) {
        res.status(404);
    }
    switch (req.method) {
        case 'GET':
            if (!isChannelMember(resultChannel, req.Header.Xuser)) {
                res.status(403);
                break;
            }
            // QUERY for last 100 messages here
            let last100Messages = mongo.last100Messages(messages, resultChannel._id);
            if (last100Messages == null) {
                res.status(500);
                break;
            }
            res.set("Content-Type", "application/json");
            // write last 100 messages to the client, encoded in JSON 
            res.json(last100Messages);
            break;

        case 'POST':
            if (!isChannelMember(resultChannel, req.Header.Xuser)) {
                res.status(403);
                break;
            }
            // Create a new message
            // Call database to INSERT a new message to the channel
            let newMessage = createMessage(req);
            let insertedMessage = mongo.insertNewMessage(messages, newMessage);
            if (insertedMessage == null) {
                res.status(500);
            }
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
            let updatedChannel = mongo.updateChannel(channels, resultChannel, req);
            if (updatedChannel == null) {
                res.status(500);
            }
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
            if (result == null) {
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
    let resultChannel = mongo.getChannelByID(channels, req.params.channelID);
    if (resultChannel == null) {
        res.status(404);
    }
    switch (req.method) {
        case 'POST':
            if (!isChannelCreator(resultChannel, req.Header.Xuser)) {
                res.status(403);
                break;
            }
            // Call database to UPDATE the current channel
            let updatedChannel = mongo.addChannelMember(channels, resultChannel, req);
            if (updatedChannel == null) {
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
            updatedChannel = mongo.removeChannelMember(channels, resultChannel, req);
            if (updatedChannel == null) {
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
    let resultMessage = mongo.getMessageByID(messages, req.params.messageID);
    if (resultMessage == null) {
        res.status(404);
    }
    switch (req.method) {
        case 'PATCH':
            if (!isMessageCreator(resultMessage, req.Header.Xuser)) {
                res.status(403);
                break;
            }
            // TODO: Call the database to UPDATE the message in the database using the messageID
            let updatedMessage = mongo.updateMessage(messages, resultMessage, req);
            if (updatedMessage == null) {
                res.status(500);
                break;
            }
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
            if (result == null) {
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