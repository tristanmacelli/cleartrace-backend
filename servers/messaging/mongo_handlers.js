"use strict";

const MongoClient = require('mongodb').MongoClient;
const assert = require('assert');

// Connection URL
const url = 'mongodb://localhost:27017';

// Database Name
const dbName = 'messaging';

// Create a new MongoClient
const client = new MongoClient(url);
const db = openConnection();

function openConnection() {
    const db;
    // Use connect method to connect to the Server
    client.connect(function (err) {
        assert.equal(null, err);
        console.log("Connected successfully to server");

        db = client.db(dbName);
    });
    return db;
}

function getAllChannels() {
    // if channels does not yet exist
    cursor = db.channels.find()
    if (!cursor.hasNext()) {
        // Throw error
        console.log("No channels collection found");
    }
    return cursor.forEach(printjson);
}

function insertNewChannel(newChannel) {
    result = db.channels.save({
        name: newChannel.Name, description: newChannel.Description,
        private: newChannel.Private, members: newChannel.Members,
        createdAt: newChannel.CreatedAt, creator: newChannel.Creator,
        editedAt: newChannel.EditedAt
    });
    if (result.hasWriteError()) {
        return;
    }
    newChannel._id = result._id;
    return newChannel;
}

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

function updatedChannel(existingChannel, req) {
    result = db.channels.save({
        name: req.body.name, description: req.body.description,
        private: existingChannel.Private, members: existingChannel.Members,
        createdAt: existingChannel.CreatedAt, creator: existingChannel.Creator,
        editedAt: existingChannel.EditedAt
    });
    if (result.hasWriteError()) {
        return;
    }
    existingChannel.name = result.name;
    existingChannel.description = result.description;
    return newChannel;
}

// Please make sure that these lines will delete all messages for the specified channelID
function deleteChannel(existingChannel) {
    db.channels.remove({ _id: ObjectId(existingChannel._id) })
    result = db.messages.remove({ channelID: existingChannel._id })
    if (result.hasWriteError()) {
        return null;
    }
    return result;
}

// 
function queryByChannelID(id) {
    if (id == null) {
        return null;
    }
    return db.channels.find(_id = id);
}

function last100Messages(id) {
    if (id == null) {
        return null;
    }
    id = toString(id);
    return db.messages.find(channelID = id).sort({ createdAt = -1 }).limit(100);
}

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
    deleteChannel,
    queryByChannelID,
    last100Messages,
    closeConnection
}