"use strict";



//require the express and morgan packages
const express = require("express");
const morgan = require("morgan");
var http = require('http');

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
            // TODO: QUERY for all channels here
            // write those to the client, encoded in JSON
            res.json(allChannels);
            break;

        case 'POST':
            console.log(req.body)
            if (req.body.channel.name == null) {
                next()
                //do something about the name property being null
            }
            var insert = createChannel(req)
            // TODO: Call database to INSERT this new channel
            res.set("Content-Type", "application/json");
            res.json(insert);
            res.status(201)  //probably cant do this >>> .send("success");
            break;
        default:
            break;
    }
});

// Specific channel handler
app.use("/v1/channels/{channelID}", (req, res, next) => {
    switch (req.method) {
        case 'GET':
            // Is this necessary if we already have it in JSON in the request?
            var channel = createChannel(req);
            if (!isChannelMember(channel, req.Xuser)) {
                res.status(403)
                break;
            }
            // TODO: QUERY for last 100 messages here
            last100Messages = null // fn call to query database
            res.set("Content-Type", "application/json");
            //write last 100 messages to the client, encoded in JSON 
            res.json(last100Messages);
            break;

        case 'POST':
            // Is this necessary if we already have it in JSON in the request?
            var channel = createChannel(req);
            if (!isChannelMember(channel, req.Xuser)) {
                res.status(403)
                break;
            }
            // Create a new message
            newMessage = createMessage(req)
            // TODO: Call database to INSERT a new message to the channel
            insertedMessage = null // fn call to insert into database new message at specified channel
            res.set("Content-Type", "application/json");
            res.json(insertedMessage);
            res.status(201)  // probably cant do this >>> .send("success");
            break;
        case 'PATCH':
            // Is this necessary if we already have it in JSON in the request?
            var channel = createChannel(req);
            if (!isChannelCreator(channel, req.Xuser)) {
                res.status(403)
                break;
            }
            // TODO: Call database to UPDATE the channel name and/or description
            updatedChannel = null // fn call to UPDATE existing channel in the database
            res.set("Content-Type", "application/json");
            res.json(updatedChannel);
            break;
        case 'DELETE':
            // Is this necessary if we already have it in JSON in the request?
            var channel = createChannel(req);
            if (!isChannelCreator(channel, req.Xuser)) {
                res.status(403)
                break;
            }
            // TODO: Call database to DELETE this channel
            res.set("Content-Type", "text/plain");
            res.send("Channel was successfully deleted")
            break;
        default:
            break;
    }
});

// Adding and removing members from your channel
app.use("/v1/channels/{channelID}/members", (req, res, next) => {
    switch (req.method) {
        case 'POST':
            // Is this necessary if we already have it in JSON in the request?
            var channel = createChannel(req);
            if (!isChannelCreator(channel, req.Xuser)) {
                res.status(403)
                break;
            }
            // TODO: QUERY the database for the specified channel using the channelID
            // TODO: Get the list of members
            // TODO: Add the specified member
            // TODO: Call database to UPDATE the current channel
            updatedChannel = null // fn call to UPDATE existing channel in the database
            res.set("Content-Type", "application/json");
            res.status(201).send(req.user.ID + " was added to your channel")
            break;
        case 'DELETE':
            // Is this necessary if we already have it in JSON in the request?
            var channel = createChannel(req);
            if (!isChannelCreator(channel, req.Xuser)) {
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
app.use("/v1/messages/{messageID}", (req, res, next) => {
    switch (req.method) {
        case 'PATCH':
            // Is this necessary if we already have it in JSON in the request?
            var channel = createMessage(req);
            if (!isMessageCreator(channel, req.Xuser)) {
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
            var channel = createMessage(req);
            if (!isMessageCreator(channel, req.Xuser)) {
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
    var c = req.body.channel;
    return new channel(c.Name, c.Description, c.Private,
        c.Members, c.CreatedAt, c.Creator, c.EditedAt);
}

function createMessage(req) {
    var m = req.body.message;
    return new message(m.ChannelID, m.CreatedAt, m.Body,
        m.Creator, m.EditedAt);
}

function isChannelMember(channel, user) {
    var isMember = false;
    if (channel.Private) {
        for (i = 0; i < channel.Members.length; i++) {
            if (channel.Members[i].ID == user.ID) {
                isMember = true;
                break;
            }
        }
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

