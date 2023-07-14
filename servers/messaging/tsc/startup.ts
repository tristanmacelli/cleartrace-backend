// import the MongoClient and Db
import { MongoError } from "mongodb";
import * as mongo from "./mongo_handlers";
import { Channel } from "./channel";
import { User } from "./user";

async function startUp() {
    const db = await mongo.createConnection();
    console.log("Connected successfully to userMessageStore");

    // check if any collection exists
    const channels = await db.createCollection('channels')
    const emptyUser = new User(-1, "", "", "", "", "")
    const dummyDate = new Date()
    const general = new Channel("", "General", "an open channel for all", false, [], dummyDate, emptyUser, dummyDate);
    
    const { hasDuplicates, err } = await mongo.insertNewChannel(channels, general)
    // check for insertion errors
    if (hasDuplicates) {
        console.log("Error inserting new channel, channel name already exists")
    }
    if (err) {
        console.log("Error inserting new general channel", err)
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