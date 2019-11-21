// import the MongoClient and Db
import { MongoClient, Db, Collection, MongoError } from "mongodb";
import * as mongo from "./mongo_handlers";
import { Channel } from "./channel";

// Mongo Connection URL
const url = 'mongodb://mongodb:27017/mongodb';

// Database Name
const dbName = 'mongodb';

// Create a new MongoClient
const client = new MongoClient(url);

client.connect(function (err: any) {
    if (err) {
        console.log("Error connecting to Mongodb: ", err);
    } else {
        console.log("Connected successfully to Mongodb");
    }
    var db: Db = client.db(dbName);

    // check if any collection exists
    db.createCollection('channels', function (err: MongoError, collection: Collection<any>) {
        if (err) {
            console.log("Error creating new collection: ", err);
        }
        // create general channel (we always want this at startup)
        let channels: Collection = collection;
        let general = new Channel("general", "an open channel for all", false, [], "enter timestamp here", -1, "not yet edited");
        let result = mongo.insertNewChannel(channels, general);
        // check for insertion errors
        if (result.errString.length > 0) {
            console.log("Failed to create new general channel upon opening connection to DB");
        }
    });

    db.createCollection('messages');
    // client.close();
});