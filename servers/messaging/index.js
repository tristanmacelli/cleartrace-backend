"use strict";
// "use strict";
// version 0.1
// to compile run tsc --outDir ../
var __awaiter = (this && this.__awaiter) || function (thisArg, _arguments, P, generator) {
    return new (P || (P = Promise))(function (resolve, reject) {
        function fulfilled(value) { try { step(generator.next(value)); } catch (e) { reject(e); } }
        function rejected(value) { try { step(generator["throw"](value)); } catch (e) { reject(e); } }
        function step(result) { result.done ? resolve(result.value) : new P(function (resolve) { resolve(result.value); }).then(fulfilled, rejected); }
        step((generator = generator.apply(thisArg, _arguments || [])).next());
    });
};
var __generator = (this && this.__generator) || function (thisArg, body) {
    var _ = { label: 0, sent: function() { if (t[0] & 1) throw t[1]; return t[1]; }, trys: [], ops: [] }, f, y, t, g;
    return g = { next: verb(0), "throw": verb(1), "return": verb(2) }, typeof Symbol === "function" && (g[Symbol.iterator] = function() { return this; }), g;
    function verb(n) { return function (v) { return step([n, v]); }; }
    function step(op) {
        if (f) throw new TypeError("Generator is already executing.");
        while (_) try {
            if (f = 1, y && (t = y[op[0] & 2 ? "return" : op[0] ? "throw" : "next"]) && !(t = t.call(y, op[1])).done) return t;
            if (y = 0, t) op = [0, t.value];
            switch (op[0]) {
                case 0: case 1: t = op; break;
                case 4: _.label++; return { value: op[1], done: false };
                case 5: _.label++; y = op[1]; op = [0]; continue;
                case 7: op = _.ops.pop(); _.trys.pop(); continue;
                default:
                    if (!(t = _.trys, t = t.length > 0 && t[t.length - 1]) && (op[0] === 6 || op[0] === 2)) { _ = 0; continue; }
                    if (op[0] === 3 && (!t || (op[1] > t[0] && op[1] < t[3]))) { _.label = op[1]; break; }
                    if (op[0] === 6 && _.label < t[1]) { _.label = t[1]; t = op; break; }
                    if (t && _.label < t[2]) { _.label = t[2]; _.ops.push(op); break; }
                    if (t[2]) _.ops.pop();
                    _.trys.pop(); continue;
            }
            op = body.call(thisArg, _);
        } catch (e) { op = [6, e]; y = 0; } finally { f = t = 0; }
        if (op[0] & 5) throw op[1]; return { value: op[0] ? op[1] : void 0, done: true };
    }
};
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
}
var __importStar = (this && this.__importStar) || function (mod) {
    if (mod && mod.__esModule) return mod;
    var result = {};
    if (mod != null) for (var k in mod) if (Object.hasOwnProperty.call(mod, k)) result[k] = mod[k];
    result["default"] = mod;
    return result;
}
var _this = this;
Object.defineProperty(exports, "__esModule", { value: true });
//require the express and morgan packages
var express_1 = __importDefault(require("express"));
var morgan_1 = __importDefault(require("morgan"));
var mongodb_1 = require("mongodb");
var mongo = __importStar(require("./mongo_handlers"));
var message_1 = require("./message");
var channel_1 = require("./channel");
//create a new express application
var app = express_1.default();
var addr = process.env.ADDR || "80";
//split host and port using destructuring
// const [host, port] = addr.split(":");
// let portNum = parseInt(port);
//add JSON request body parsing middleware
app.use(express_1.default.json());
//add the request logging middleware
app.use(morgan_1.default("dev"));
// Connection URL
var url = 'mongodb://mongodb:27017/mongodb';
// Database Name
var dbName = 'mongodb';
var db;
var messages;
var channels;
// Reasoning for refactor: 
// https://bit.ly/342jCtj
// Create a new MongoClient
var mc = new mongodb_1.MongoClient(url, { useUnifiedTopology: true });
var mongoClient = function () { return __awaiter(_this, void 0, void 0, function () {
    var client, e_1;
    return __generator(this, function (_a) {
        switch (_a.label) {
            case 0:
                _a.trys.push([0, 2, , 3]);
                return [4 /*yield*/, mc.connect()];
            case 1:
                client = _a.sent();
                return [3 /*break*/, 3];
            case 2:
                e_1 = _a.sent();
                console.log("Cannot connect to mongo: MongoNetworkError: failed to connect to server");
                console.log("Restarting Messaging server");
                process.exit(1);
                return [3 /*break*/, 3];
            case 3: return [2 /*return*/, client];
        }
    });
}); };
var checkConnection = function () { return __awaiter(_this, void 0, void 0, function () {
    var client;
    return __generator(this, function (_a) {
        switch (_a.label) {
            case 0: return [4 /*yield*/, mongoClient()];
            case 1:
                client = _a.sent();
                return [2 /*return*/, client];
        }
    });
}); };
var main = function () { return __awaiter(_this, void 0, void 0, function () {
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
        return channel.creator == user.id;
    }
    function isMessageCreator(message, user) {
        return message.creator == user.id;
    }
    var client;
    return __generator(this, function (_a) {
        switch (_a.label) {
            case 0: return [4 /*yield*/, checkConnection()];
            case 1:
                client = _a.sent();
                db = client.db(dbName);
                channels = db.collection("channels");
                messages = db.collection("messages");
                // TODO: We should do a test of the mongo helper methods
                app.listen(+addr, "", function () {
                    //callback is executed once server is listening
                    console.log("server is listening at http://:" + addr + "...");
                    // console.log("port : " + port);
                    // console.log("host : " + host);
                });
                app.use("/v1/channels", function (req, res) {
                    switch (req.method) {
                        case 'GET':
                            res.set("Content-Type", "application/json");
                            // QUERY for all channels here
                            var allChannels = mongo.getAllChannels(channels, res);
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
                app.use("/v1/channels/:channelID", function (req, res) {
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
                app.use("/v1/channels/:channelID/members", function (req, res) {
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
                app.use("/v1/messages/:messageID", function (req, res) {
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
                // error handler that will be called if
                // any handler earlier in the chain throws
                // an exception or passes an error to next()
                app.use(function (err, req, res) {
                    //write a stack trace to standard out,
                    //which writes to the server's log
                    console.error(err.stack);
                    //but only report the error message
                    //to the client, with a 500 status code
                    res.set("Content-Type", "text/plain");
                    res.status(500).send(err.message);
                });
                return [2 /*return*/];
        }
    });
}); };
main();
