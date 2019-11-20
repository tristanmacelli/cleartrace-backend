// import the MongoClient and Db
import { MongoClient, Db, Collection} from "mongodb";
import { Channel } from "./channel";
import * as mongo from "./mongo_handlers";

// Mongo Connection URL
// const url = 'mongo://mongodb:27017';
const url = 'mongodb://mongodb:27017/mongodb';


// Database Name
const dbName = 'mongodb';
var channels: Collection;


// Create a new MongoClient
const client = new MongoClient(url);

client.connect(async function (err: any) {
    console.log("Connected successfully to Mongodb");
    if (err) {
        console.log("the error is ", err);
    }
    var db: Db = client.db(dbName);

    // check if any collection exists
    db.createCollection('channels', function(e,a){
        if (e) {
            console.log("error is ", e);
        }
        channels = a;
        // create general channel
        var general = new Channel("general", "an open channel for all", false, [], "enter timestamp here", -1, "not yet edited");
        // channel that we always want at startup
        let result = mongo.insertNewChannel(channels, general);
        // check for error while inserting
        if (result.errString.length > 0) {
            console.log("failed to create general channel upon opening connection to DB");
            // res.status(500);
        }
    });

    db.createCollection('messages');
    client.close();
});


// Reasoning for refactor: 
// https://bit.ly/342jCtj
// Use connect method to connect to the mongo DB
// MongoClient.connect(url, function (err: any, mc:MongoClient) {
//     console.log("Connected successfully to server");
//     console.log("there is an error:", err.err);
//     if (err) throw err;

//     // db = client.db(dbName);
//     var db = mc.db(dbName);
//     // Create channels and messages collections in mongo
//     // https://bit.ly/2r57au6
//     db.createCollection('channels');
//     db.createCollection('messages');
//     // var general = new Channel("general", "an open channel for all", false, [], "enter timestamp here", -1, "not yet edited");
//     // channel that we always want at startup
//     // let result  = mongo.insertNewChannel(channels, general);

//     // db.channels.save({
//     //     name: newChannel.name, description: newChannel.description,
//     //     private: newChannel.private, members: newChannel.members,
//     //     createdAt: newChannel.createdAt, creator: newChannel.creator,
//     //     editedAt: newChannel.editedAt
//     // }).catch(() => {
//     //     errString = "Error inserting new channel";
//     // });

//     // if (result.errString.length > 0) {
//     //     console.log("failed to create general channel upon opening connection to DB");
//     //     // res.status(500);
//     // }
//     mc.close();
// });