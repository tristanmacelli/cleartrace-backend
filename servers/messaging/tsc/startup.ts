// import the MongoClient and Db
import { Db, MongoError } from "mongodb";
import * as mongo from "./mongo_handlers";
import { Channel } from "./channel";
import { User } from "./user";

async function startUp() {
    const db = await mongo.createConnection();
    console.log("Connected successfully to userMessageStore");

    // check if any collection exists
    let channels = await db.createCollection('channels')
    let emptyUser = new User(-1, "", "", "", "", "")
    let dummyDate = new Date()
    let general = new Channel("", "general", "an open channel for all", false, [], dummyDate, emptyUser, dummyDate);
    
    let result = await mongo.insertNewChannel(channels, general)
    // check for insertion errors
    if (result.duplicates) {
        console.log("Error inserting new channel, channel name already exists")
    }
    if (result.err) {
        console.log("Error inserting new general channel", result.err)
        process.exit(1)
    }
        
    db.createCollection('messages')
    .catch((err: MongoError) => {
        console.log("Error creating new collection", err)
    })

    console.log("MongoDB start up complete")
    process.exit(0)
}

startUp();