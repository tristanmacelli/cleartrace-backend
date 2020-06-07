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
const mc = new MongoClient(url, { useUnifiedTopology: true });

function sleep(seconds: number) {
    let milliseconds = seconds * 1000
    const stop = new Date().getTime() + milliseconds;
    while(new Date().getTime() < stop);       
}

const createConnection = async (): Promise<MongoClient> => {
    let client: MongoClient;
    while (1) {
        try {
            client = await mc.connect();
            break;
        } catch (e) {
            console.log("Cannot connect to mongo: MongoNetworkError: failed to connect to server");
            console.log("Retrying in 1 second");
            sleep(1)
        }
    }
    return client!;
}

const startUp = async () => {
    const client = await createConnection();
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