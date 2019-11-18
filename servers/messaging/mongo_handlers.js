"use strict";
function __export(m) {
    for (var p in m) if (!exports.hasOwnProperty(p)) exports[p] = m[p];
}
Object.defineProperty(exports, "__esModule", { value: true });
// to compile run tsc --outDir ../
var mongodb_1 = require("mongodb");
var channel_1 = require("./channel");
var message_1 = require("./message");
// getAllChannels does something
// TODO: make sure the returned value is a shape that we can actually use
function getAllChannels(channels) {
    // if channels does not yet exist
    var cursor = channels.find();
    if (!cursor.hasNext()) {
        // Throw error
        console.log("No channels collection found");
        return null;
    }
    return cursor.forEach(function (m) { JSON.stringify(m); });
}
exports.getAllChannels = getAllChannels;
// insertNewChannel takes in a new Channel and
function insertNewChannel(channels, newChannel) {
    var errString = "";
    var idWeWant;
    channels.save({
        name: newChannel.name, description: newChannel.description,
        private: newChannel.private, members: newChannel.members,
        createdAt: newChannel.createdAt, creator: newChannel.creator,
        editedAt: newChannel.editedAt
    }).catch(function () {
        errString = "Error inserting new channel";
    });
    channels.find({ name: newChannel.name, createdAt: newChannel.createdAt }).next()
        .then(function (doc) {
        idWeWant = doc.id;
    }).catch(function (err) {
        idWeWant = "";
    });
    newChannel._id = idWeWant;
    return { newChannel: newChannel, errString: errString };
}
exports.insertNewChannel = insertNewChannel;
// insertNewMessage takes in a new Message and
function insertNewMessage(messages, newMessage) {
    var errString = "";
    var idWeWant;
    if (newMessage.channelID == null) {
        errString = "Could not find ID";
        return { newMessage: newMessage, errString: errString };
    }
    var result = messages.save({
        channelID: newMessage.channelID, createdAt: newMessage.createdAt,
        body: newMessage.body, creator: newMessage.creator,
        editedAt: newMessage.editedAt
    }).catch(function () {
        errString = "Error inserting new message";
    });
    messages.find({ body: newMessage.body, createdAt: newMessage.createdAt }).next()
        .then(function (doc) {
        idWeWant = doc.id;
    }).catch(function (err) {
        idWeWant = "";
    });
    newMessage._id = idWeWant;
    return { newMessage: newMessage, errString: errString };
}
exports.insertNewMessage = insertNewMessage;
// updatedChannel updates name and body of an existing Channel using a req (request) object
function updateChannel(channels, existingChannel, req) {
    var errString = "";
    channels.save({
        name: req.body.name, description: req.body.description,
        private: existingChannel.private, members: existingChannel.members,
        createdAt: existingChannel.createdAt, creator: existingChannel.creator,
        editedAt: existingChannel.editedAt
    }).catch(function () {
        errString = "Error updating message";
    });
    existingChannel.name = req.body.name;
    existingChannel.description = req.body.description;
    return { existingChannel: existingChannel, errString: errString };
}
exports.updateChannel = updateChannel;
// addChannelMembers takes an existing Channel and adds members using a req (request) object
function addChannelMember(channels, existingChannel, req) {
    var errString = "";
    existingChannel.members.push(req.body.message.id);
    channels.save({
        name: existingChannel.name, description: existingChannel.description,
        private: existingChannel.private, members: existingChannel.members,
        createdAt: existingChannel.createdAt, creator: existingChannel.creator,
        editedAt: existingChannel.editedAt
    }).catch(function () {
        errString = "Error updating message";
    });
    return errString;
}
exports.addChannelMember = addChannelMember;
// removeChannelMember takes an existing Channel and removes members using a req (request) object
function removeChannelMember(channels, existingChannel, req) {
    // Remove the specified member from this channel's list of members
    var errString = "";
    existingChannel.members.splice(req.body.message.id, 1);
    channels.save({
        name: existingChannel.name, description: existingChannel.description,
        private: existingChannel.private, members: existingChannel.members,
        createdAt: existingChannel.createdAt, creator: existingChannel.creator,
        editedAt: existingChannel.editedAt
    }).catch(function () {
        errString = "Error updating message";
    });
    return errString;
}
exports.removeChannelMember = removeChannelMember;
function updateMessage(messages, existingMessage, req) {
    var errString = "";
    messages.save({
        body: req.body, creator: existingMessage.creator,
        createdAt: existingMessage.createdAt, channelID: existingMessage.channelID,
        editedAt: existingMessage.editedAt
    }).catch(function () {
        errString = "Error updating message";
    });
    existingMessage.body = req.body;
    return { existingMessage: existingMessage, errString: errString };
}
exports.updateMessage = updateMessage;
// deleteChannel does something
function deleteChannel(channels, messages, existingChannel) {
    // We are not allowed to delete the general channel
    var errString = "";
    if (existingChannel.creator == -1) {
        return "Error deleting channel";
    }
    channels.remove({ _id: new mongodb_1.ObjectID(existingChannel._id) }).catch(function () {
        errString = "Error deleting channel";
    });
    messages.remove({ channelID: existingChannel._id }).catch(function () {
        errString = "Error deleting channel";
    });
    return errString;
}
exports.deleteChannel = deleteChannel;
// deleteMessage does something
function deleteMessage(messages, existingMessage) {
    var errString = "";
    messages.remove({ messageID: existingMessage._id }).catch(function () {
        errString = "Error deleting message";
    });
    return errString;
}
exports.deleteMessage = deleteMessage;
function getChannelByID(channels, id) {
    // Since id's are auto-generated and unique we chose to use find instead of findOne() 
    var finalResponse;
    var errString;
    channels.find({ _id: id }).next().then(function (doc) {
        finalResponse = doc;
        errString = "";
    }).catch(function (err) {
        finalResponse = null;
        errString = err;
    });
    var finalChannel;
    if (finalResponse == null) {
        finalChannel = new channel_1.Channel("", "", false, [], "", -1, "");
        return { finalChannel: finalChannel, errString: errString };
    }
    finalChannel = new channel_1.Channel(finalResponse.name, finalResponse.description, finalResponse.private, finalResponse.members, finalResponse.createdAt, finalResponse.Creator, finalResponse.editedAt);
    return { finalChannel: finalChannel, errString: errString };
}
exports.getChannelByID = getChannelByID;
function getMessageByID(messages, id) {
    // Since id's are auto-generated and unique we chose to use find instead of findOne() 
    var finalResponse;
    var errString;
    messages.find({ _id: id }).next().then(function (doc) {
        finalResponse = doc;
        errString = "";
    }).catch(function (err) {
        finalResponse = null;
        errString = err;
    });
    var finalMessage;
    if (finalResponse == null) {
        finalMessage = new message_1.Message("", "", "", "", "");
        return { finalMessage: finalMessage, errString: errString };
    }
    finalMessage = new message_1.Message(finalResponse.channelID, finalResponse.createdAt, finalResponse.body, finalResponse.creator, finalResponse.editedAt);
    return { finalMessage: finalMessage, errString: errString };
}
exports.getMessageByID = getMessageByID;
// last100Messages does something
function last100Messages(messages, id) {
    if (id == null) {
        throw "No id value passed";
    }
    id = id.toString();
    return messages.find({ channelID: id }).sort({ createdAt: -1 }).limit(100);
}
exports.last100Messages = last100Messages;
// last100Messages does something
function last100SpecificMessages(messages, channelID, messageID) {
    if (channelID == null) {
        throw "No id value passed";
    }
    channelID = channelID.toString();
    return messages.find({ channelID: channelID, _id: { $lt: messageID } }).sort({ createdAt: -1 }).limit(100);
}
exports.last100SpecificMessages = last100SpecificMessages;
__export(require("./mongo_handlers"));
