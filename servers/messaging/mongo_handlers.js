"use strict";

const MongoClient = require('mongodb').MongoClient;
const assert = require('assert');

// Connection URL
const url = 'mongodb://localhost:27017';

// Database Name
const dbName = 'messaging';

// Create a new MongoClient
const client = new MongoClient(url);

// openConnection does something
function openConnection() {
    const db;
    // Use connect method to connect to the Server
    client.connect(function (err) {
        assert.equal(null, err);
        console.log("Connected successfully to server");

        db = client.db(dbName);
    });
    general = new Channel("general", "an open channel for all", false,
        [], "enter timestamp here", "-1", "not yet edited");
    // channel that we always want at startup
    result = insertNewChannel(general);
    if (result == null) {
        console.log("failed to create general channel upon opening connection to DB");
        // res.status(500);
    }
    return db;
}

const db = openConnection();

// getAllChannels does something
function getAllChannels() {
    // if channels does not yet exist
    cursor = db.channels.find();
    if (!cursor.hasNext()) {
        // Throw error
        console.log("No channels collection found");
        return null;
    }
    return cursor.forEach(printjson);
}

// insertNewChannel does something
function insertNewChannel(newChannel) {
    result = db.channels.save({
        name: newChannel.Name, description: newChannel.Description,
        private: newChannel.Private, members: newChannel.Members,
        createdAt: newChannel.CreatedAt, creator: newChannel.Creator,
        editedAt: newChannel.EditedAt
    });
    if (result.hasWriteError()) {
        return null;
    }
    newChannel._id = result._id;
    return newChannel;
}

// insertNewMessage does something
function insertNewMessage(newMessage) {
    if (newMessage.ChannelID == null) {
        return null;
    }
    result = db.messages.save({
        channelID: newMessage.ChannelID, createdAt: newMessage.CreatedAt,
        body: newMessage.Body, creator: newMessage.Creator,
        editedAt: newMessage.EditedAt
    });
    if (result.hasWriteError()) {
        return null;
    }
    newMessage._id = result._id;
    return newMessage;
}

// updatedChannel updates name and body of channel
function updatedChannel(existingChannel, req) {
    result = db.channels.save({
        name: req.body.name, description: req.body.description,
        private: existingChannel.Private, members: existingChannel.Members,
        createdAt: existingChannel.CreatedAt, creator: existingChannel.Creator,
        editedAt: existingChannel.EditedAt
    });
    if (result.hasWriteError()) {
        return null;
    }
    existingChannel.name = result.name;
    existingChannel.description = result.description;
    return newChannel;
}

function addChannelMembers(existingChannel, req) {
    existingChannel.members.push(req.body.message.id);
    result = db.channels.save({
        name: existingChannel.name, description: existingChannel.description,
        private: existingChannel.private, members: existingChannel.members,
        createdAt: existingChannel.createdAt, creator: existingChannel.creator,
        editedAt: existingChannel.editedAt
    });
    if (result.hasWriteError()) {
        return null;
    }
    // Add the specified member
    existingChannel.members = newMembers;
    return existingChannel;
}

function removeChannelMembers(existingChannel, req) {
    // Remove the specified member from this channel's list of members
    existingChannel.members.splice(req.body.message.id, 1);
    result = db.channels.save({
        name: existingChannel.name, description: existingChannel.description,
        private: existingChannel.private, members: existingChannel.members,
        createdAt: existingChannel.createdAt, creator: existingChannel.creator,
        editedAt: existingChannel.editedAt
    });
    if (result.hasWriteError()) {
        return null;
    }
    return existingChannel;
}

function updateMessage(existingMessage, req) {
    result = db.messages.save({
        body: req.body, creator: existingChannel.Creator,
        createdAt: existingChannel.CreatedAt, channelID: existingMessage.channelID,
        editedAt: existingChannel.EditedAt
    });
    if (result.hasWriteError()) {
        return null;
    }
    existingMessage.body = result.body;
    return existingMessage;
}

// deleteChannel does something
function deleteChannel(existingChannel) {
    // We are not allowed to delete the general channel
    if (existingChannel.Creator == -1) {
        return null;
    }
    db.channels.remove({ _id: ObjectId(existingChannel._id) });
    result = db.messages.remove({ channelID: existingChannel._id });
    if (result.hasWriteError()) {
        return null;
    }
    return result;
}

function deleteMessage(existingMessage) {
    result = db.messages.remove({ messageID: existingMessage._id });
    if (result.hasWriteError()) {
        return null;
    }
    return result;
}

// queryByChannelID does something
function getChannelByID(id) {
    if (id == null) {
        return null;
    }
    return db.channels.find(_id = id);
}

function getMessageByID(id) {
    if (id == null) {
        return null;
    }
    return db.messages.find(_id = id);
}

// last100Messages does something
function last100Messages(id) {
    if (id == null) {
        return null;
    }
    id = toString(id);
    return db.messages.find(channelID = id).sort({ createdAt = -1 }).limit(100);
}

// closeConnection does something
function closeConnection() {
    client.close();
}

//export the public functions
module.exports = {
    openConnection,
    getAllChannels,
    insertNewChannel,
    insertNewMessage,
    updatedChannel,
    addChannelMembers,
    removeChannelMembers,
    updateMessage,
    deleteChannel,
    deleteMessage,
    getChannelByID,
    getMessageByID,
    last100Messages,
    closeConnection
}