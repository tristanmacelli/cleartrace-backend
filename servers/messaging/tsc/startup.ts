// import the MongoClient and Db
import { MongoClient, Db, Collection, MongoError } from "mongodb";
import * as mongo from "./mongo_handlers";
import { Channel } from "./channel";
import { User } from "./user";

// Mongo Connection URL
const url = 'mongodb://mongodb:27017/mongodb';

// Database Name
const dbName = 'mongodb';

// Create a new MongoClient
const client = new MongoClient(url);

const startUp = async () => {
    client.connect(function (err: any) {
        if (err) {
            console.log("Error connecting to Mongodb: ", err);
        } else {
            console.log("Connected successfully to Mongodb");
        }
        var db: Db = client.db(dbName);

        // check if any collection exists
        db.createCollection('channels', async function (err: MongoError, collection: Collection<any>) {
            if (err) {
                console.log("Error creating new collection: ", err);
            }
            // create general channel (we always want this at startup)
            let channels: Collection = collection;
            let emptyUser = new User(-1, "", "", "", "", "")
            let dummyDate = new Date()
            let general = new Channel("", "general", "an open channel for all", false, [], dummyDate, emptyUser, dummyDate);
            await mongo.insertNewChannel(channels, general).then(result => {
                // check for insertion errors
                if (result.err) {
                    console.log("Failed to create new general channel upon opening connection to DB");
                }
            })
            
        });

        db.createCollection('messages');
        // client.close();
    });
}
startUp();