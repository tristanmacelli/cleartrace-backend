"use strict";

//require the express and morgan packages
const express = require("express");
const morgan = require("morgan");
// let http = require('http');
const mongo = require('./mongo_handlers.js');

//create a new express application
const app = express();

const addr = process.env.ADDR || ":80";
//split host and port using destructuring
const [host, port] = addr.split(":");

//add JSON request body parsing middleware
app.use(express.json());
//add the request logging middleware
app.use(morgan("dev"));

// All channel handler
app.use("/v1/channels", (req, res, next) => {
    switch (req.method) {
        case 'GET':
            res.set("Content-Type", "application/json");
            // QUERY for all channels here
            allChannels = mongo.getAllChannels();
            if (allChannels == null) {
                res.status(500);
            }
            // write those to the client, encoded in JSON
            res.json(allChannels);
            break;

        case 'POST':
            console.log(req.body)
            if (req.body.channel.name == null) {
                next()
                //do something about the name property being null
            }
            let insert = createChannel(req);
            // Call database to INSERT this new channel
            insertResult = mongo.insertNewChannel(insert);
            if (insertResult == null) {
                res.status(500);
            }
            res.set("Content-Type", "application/json");
            res.json(insertResult);
            res.status(201)  //probably cant do this >>> .send("success");
            break;
        default:
            break;
    }
});

// Specific channel handler
app.use("/v1/channels/:channelID", (req, res, next) => {
    // TODO: QUERY for the channel based on req.params.channelID
    resultChannel = mongo.queryByChannelID(req.params.channelID)
    if (resultChannel == null) {
        res.status(500);
    }
    switch (req.method) {
        case 'GET':
            if (!isChannelMember(resultChannel, req.Header.Xuser)) {
                res.status(403);
                break;
            }
            // QUERY for last 100 messages here
            last100Messages = mongo.last100Messages(resultChannel._id);
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
            newMessage = createMessage(req);
            // Call database to INSERT a new message to the channel
            insertedMessage = mongo.insertNewMessage(newMessage);
            if (insertedMessage == null) {
                res.status(500);
            }
            res.set("Content-Type", "application/json");
            res.json(insertedMessage);
            res.status(201)  // probably cant do this >>> .send("success");
            break;
        case 'PATCH':
            if (!isChannelCreator(resultChannel, req.Header.Xuser)) {
                res.status(403)
                break;
            }
            // Call database to UPDATE the channel name and/or description
            updatedChannel = mongo.updatedChannel(resultChannel, req);
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
            result = mongo.deleteChannel(resultChannel);
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
app.use("/v1/channels/:channelID/members", (req, res, next) => {
    switch (req.method) {
        case 'POST':
            // Is this necessary if we already have it in JSON in the request?
            let channel = createChannel(req);
            if (!isChannelCreator(channel, req.Header.Xuser)) {
                res.status(403);
                break;
            }
            // TODO: QUERY the database for the specified channel using the channelID
            // TODO: Get the list of members
            // TODO: Add the specified member
            // TODO: Call database to UPDATE the current channel
            updatedChannel = null // fn call to UPDATE existing channel in the database
            res.set("Content-Type", "application/json");
            res.status(201).send(req.user.ID + " was added to your channel");
            break;
        case 'DELETE':
            // Is this necessary if we already have it in JSON in the request?
            // TODO: QUERY for the channel based on req.params.channelID
            let channel = createChannel(req);
            if (!isChannelCreator(channel, req.Header['X-user'])) {
                res.status(403)
                break;
            }
            // TODO: QUERY the database for the specified channel
            // TODO: Get the list of members
            // TODO: Remove the specified member from this channel's list of members
            // TODO: Call database to UPDATE the current channel
            res.set("Content-Type", "text/plain");
            res.status(200).send(req.user.ID + " was removed from your channel")
            break;
        default:
            break;
    }
});

// message handler
app.use("/v1/messages/:messageID", (req, res, next) => {
    switch (req.method) {
        case 'PATCH':
            // Is this necessary if we already have it in JSON in the request?
            let channel = createMessage(req);
            if (!isMessageCreator(channel, req.Header.Xuser)) {
                res.status(403)
                break;
            }
            // TODO: QUERY the database for the message using the messageID 
            // TODO: Update the message body
            // TODO: Call the database to UPDATE the message in the database using the messageID
            updatedMessage = null // fn call to UPDATE existing message in the database
            res.set("Content-Type", "application/json");
            res.json(updatedMessage);
            break;
        case 'DELETE':
            // Is this necessary if we already have it in JSON in the request?
            let channel = createMessage(req);
            if (!isMessageCreator(channel, req.Header.Xuser)) {
                res.status(403)
                break;
            }
            // TODO: Call database to DELETE the specified message using the messageID
            res.set("Content-Type", "text/plain");
            res.send("Message deleted")
            break;
        default:
            break;
    }
});

function createChannel(req) {
    let c = req.body.channel;
    return new channel(c.Name, c.Description, c.Private,
        c.Members, c.CreatedAt, c.Creator, c.EditedAt);
}

function createMessage(req) {
    let m = req.body.message;
    return new message(req.params.ChannelID, m.CreatedAt, m.Body,
        m.Creator, m.EditedAt);
}

function isChannelMember(channel, user) {
    let isMember = false;
    if (channel.Private) {
        for (i = 0; i < channel.Members.length; i++) {
            if (channel.Members[i].ID == user.ID) {
                isMember = true;
                break;
            }
        }
    } else {
        isMember = true;
    }
    return isMember;
}

function isChannelCreator(channel, user) {
    return channel.Creator.ID == user.ID
}

function isMessageCreator(message, user) {
    return message.Creator.ID == user.ID
}

//error handler that will be called if
//any handler earlier in the chain throws
//an exception or passes an error to next()
app.use((err, req, res, next) => {
    //write a stack trace to standard out,
    //which writes to the server's log
    console.error(err.stack)

    //but only report the error message
    //to the client, with a 500 status code
    res.set("Content-Type", "text/plain");
    res.status(500).send(err.message);
});

