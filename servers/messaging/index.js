"use strict";
// "use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
var __importStar = (this && this.__importStar) || function (mod) {
    if (mod && mod.__esModule) return mod;
    var result = {};
    if (mod != null) for (var k in mod) if (Object.hasOwnProperty.call(mod, k)) result[k] = mod[k];
    result["default"] = mod;
    return result;
};
Object.defineProperty(exports, "__esModule", { value: true });
// to compile run tsc --outDir ../
//require the express and morgan packages
var express_1 = __importDefault(require("express"));
var morgan_1 = __importDefault(require("morgan"));
var mongodb_1 = require("mongodb");
var mongo = __importStar(require("./mongo_handlers"));
var message_1 = require("./message");
var channel_1 = require("./channel");
//create a new express application
var app = express_1.default();
var addr = process.env.ADDR || ":80";
//split host and port using destructuring
var _a = addr.split(":"), host = _a[0], port = _a[1];
//add JSON request body parsing middleware
app.use(express_1.default.json());
//add the request logging middleware
app.use(morgan_1.default("dev"));
// Connection URL
var url = 'mongo://mongodb:27017/mongodb';
// Database Name
var dbName = 'mongodb';
// Create a new MongoClient
// const client = new MongoClient(url);
// var db: Db;
var messages;
var channels;
//new Server("mongo://mongodb", 27017)
// var mc = new MongoClient("mongo://mongodb:27017", {native_parser:true})
// mc.
// Reasoning for refactor: 
// https://mongodb.github.io/node-mongodb-native/driver-articles/mongoclient.html#mongoclient-connection-pooling
// Use connect method to connect to the mongo DB
// client.connect(function (err: any) {
mongodb_1.MongoClient.connect(url, function (err, client) {
    console.log("Connected successfully to server");
    var db = client.db(dbName);
    // db = client.db(dbName);
    // check if any collection exists
    db.collections()
        .then(function (doc) {
        console.log(doc);
    }).catch(function (err) {
        console.log(err);
    });
    // Start the application after the database connection is ready
    app.listen(+port, "", function () {
        //callback is executed once server is listening
        console.log("server is listening at http://:" + port + "...");
        console.log("port : " + port);
        console.log("host : " + host);
    });
});
// All channel handler
// No errors here :)
app.use("/v1/channels", function (req, res, next) {
    switch (req.method) {
        case 'GET':
            res.set("Content-Type", "application/json");
            // QUERY for all channels here
            var allChannels = mongo.getAllChannels(channels);
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
            var newChannel = createChannel(req);
            var insertResult = mongo.insertNewChannel(channels, newChannel);
            if (insertResult.errString.length > 0) {
                res.status(500);
            }
            var insertChannel = insertResult.newChannel;
            res.set("Content-Type", "application/json");
            res.json(insertChannel);
            res.status(201); //probably cant do this >>> .send("success");
            break;
        default:
            break;
    }
});
// Specific channel handler
app.use("/v1/channels/:channelID", function (req, res, next) {
    // QUERY for the channel based on req.params.channelID
    if (req.params.channelID == null) {
        res.status(404);
        return;
    }
    var result = mongo.getChannelByID(channels, req.params.channelID);
    if (result.errString.length() > 0) {
        res.status(500);
        return;
    }
    var resultChannel = result.finalChannel;
    switch (req.method) {
        case 'GET':
            if (!isChannelMember(resultChannel, req.Header.Xuser)) {
                res.status(403);
                break;
            }
            var returnedMessages = void 0;
            // QUERY for last 100 messages here
            if (req.params.before != null) {
                returnedMessages = mongo.last100SpecificMessages(messages, resultChannel._id, req.params.before);
                if (returnedMessages == null) {
                    res.status(500);
                    break;
                }
            }
            else {
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
            var newMessage = createMessage(req);
            var insertedResult = mongo.insertNewMessage(messages, newMessage);
            if (insertedResult.errString.length > 0) {
                res.status(500);
            }
            var insertedMessage = insertedResult.newMessage;
            res.set("Content-Type", "application/json");
            res.json(insertedMessage);
            res.status(201); // probably cant do this >>> .send("success");
            break;
        case 'PATCH':
            if (!isChannelCreator(resultChannel, req.Header.Xuser)) {
                res.status(403);
                break;
            }
            // Call database to UPDATE the channel name and/or description
            var updateResult = mongo.updateChannel(channels, resultChannel, req);
            if (updateResult.errString.length > 0) {
                res.status(500);
            }
            var updatedChannel = updateResult.existingChannel;
            res.set("Content-Type", "application/json");
            res.json(updatedChannel);
            break;
        case 'DELETE':
            if (!isChannelCreator(resultChannel, req.Header.Xuser)) {
                res.status(403);
                break;
            }
            // Call database to DELETE this channel
            var result_1 = mongo.deleteChannel(channels, messages, resultChannel);
            if (result_1.length > 0) {
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
app.use("/v1/channels/:channelID/members", function (req, res, next) {
    // QUERY for the channel based on req.params.channelID
    if (req.params.channelID == null) {
        res.status(404);
        return;
    }
    var result = mongo.getChannelByID(channels, req.params.channelID);
    if (result.errString.length() > 0) {
        res.status(500);
        return;
    }
    var resultChannel = result.finalChannel;
    switch (req.method) {
        case 'POST':
            if (!isChannelCreator(resultChannel, req.Header.Xuser)) {
                res.status(403);
                break;
            }
            // Call database to UPDATE the current channel
            var addResult = mongo.addChannelMember(channels, resultChannel, req);
            if (addResult.length > 0) {
                res.status(500);
                break;
            }
            res.set("Content-Type", "application/json");
            res.status(201).send(req.user.ID + " was added to your channel");
            break;
        case 'DELETE':
            if (!isChannelCreator(resultChannel, req.Header['X-user'])) {
                res.status(403);
                break;
            }
            // database to UPDATE the current channel members
            var errResult = mongo.removeChannelMember(channels, resultChannel, req);
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
app.use("/v1/messages/:messageID", function (req, res, next) {
    if (req.params.messageID == null) {
        res.status(404);
        return;
    }
    var result = mongo.getMessageByID(messages, req.params.messageID);
    if (result.errString.length() > 0) {
        res.status(500);
        return;
    }
    var resultMessage = result.finalMessage;
    switch (req.method) {
        case 'PATCH':
            if (!isMessageCreator(resultMessage, req.Header.Xuser)) {
                res.status(403);
                break;
            }
            // TODO: Call the database to UPDATE the message in the database using the messageID
            var updatedResult = mongo.updateMessage(messages, resultMessage, req);
            if (updatedResult.errString.length > 0) {
                res.status(500);
                break;
            }
            var updatedMessage = updatedResult.existingMessage;
            res.set("Content-Type", "application/json");
            res.json(updatedMessage);
            break;
        case 'DELETE':
            if (!isMessageCreator(resultMessage, req.Header.Xuser)) {
                res.status(403);
                break;
            }
            // Call database to DELETE the specified message using the messageID
            // Call database to DELETE this channel
            var result_2 = mongo.deleteMessage(messages, resultMessage);
            if (result_2.length > 0) {
                res.status(500);
            }
            res.set("Content-Type", "text/plain");
            res.send("Message deleted");
            break;
        default:
            break;
    }
});
function createChannel(req) {
    var c = req.body.channel;
    return new channel_1.Channel(c.name, c.description, c.private, c.members, c.createdAt, c.creator, c.editedAt);
}
function createMessage(req) {
    var m = req.body.message;
    return new message_1.Message(req.params.ChannelID, m.createdAt, m.body, m.creator, m.editedAt);
}
function isChannelMember(channel, user) {
    var isMember = false;
    if (channel.private) {
        for (var i = 0; i < channel.members.length; i++) {
            if (channel.members[i] == user.ID) {
                isMember = true;
                break;
            }
        }
    }
    else {
        isMember = true;
    }
    return isMember;
}
function isChannelCreator(channel, user) {
    return channel.creator == user._id;
}
function isMessageCreator(message, user) {
    return message.creator == user._id;
}
//error handler that will be called if
//any handler earlier in the chain throws
//an exception or passes an error to next()
app.use(function (err, req, res, next) {
    //write a stack trace to standard out,
    //which writes to the server's log
    console.error(err.stack);
    //but only report the error message
    //to the client, with a 500 status code
    res.set("Content-Type", "text/plain");
    res.status(500).send(err.message);
});
