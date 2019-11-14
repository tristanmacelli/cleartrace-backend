"use strict";

const MongoClient = require('mongodb').MongoClient;
const assert = require('assert');

// Connection URL
const url = 'mongodb://localhost:27017';

// Database Name
const dbName = 'messaging';

// Create a new MongoClient
const client = new MongoClient(url);

function openConnection() {
    // Use connect method to connect to the Server
    client.connect(function (err) {
        assert.equal(null, err);
        console.log("Connected successfully to server");

        const db = client.db(dbName);

        getAllChannels(db, function () { });
    });
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



function closeConnection() {
    client.close();
}

//export the public functions
module.exports = {
    openConnection,
    getAllChannels,
    closeConnection
}