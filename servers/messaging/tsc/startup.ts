// import the MongoClient and Db
import { Db, MongoError } from "mongodb";
import * as mongo from "./mongo_handlers";
import { Channel } from "./channel";
import { User } from "./user";

// Database Name
const dbName = 'mongodb';

async function startUp() {
    const client = await mongo.createConnection();
    console.log("Connected successfully to Mongodb");
    var db: Db = client.db(dbName);

    // check if any collection exists
    await db.createCollection('channels')
    .then(async (channels) => {
        let emptyUser = new User(-1, "", "", "", "", "")
        let dummyDate = new Date()
        let general = new Channel("", "general", "an open channel for all", false, [], dummyDate, emptyUser, dummyDate);
        
        await mongo.insertNewChannel(channels, general)
        .then(result => {
            // check for insertion errors
            if (result.err) {
                console.log("Failed to create new general channel upon opening connection to DB");
            }
        })
        .catch((err: MongoError) => {
            console.log("Error inserting new channel", err)
        })
    })
    .catch((err: MongoError) => {
        console.log("Error creating new collection", err)
    })

    db.createCollection('messages')
    .catch((err: MongoError) => {
        console.log("Error creating new collection", err)
    })

    console.log("MongoDB start up complete")
    process.exit(0)
}

startUp();