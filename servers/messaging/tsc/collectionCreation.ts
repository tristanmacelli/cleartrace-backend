// import the MongoClient and Db
import { MongoClient, Db } from "mongodb";

// Mongo Connection URL
const url = 'mongo://mongodb:27017';

// Database Name
const dbName = 'mongodb';

// Create a new MongoClient
const client = new MongoClient(url);

client.connect(function (err: any) {
    console.log("Connected successfully to server");
    console.log("the error is ", err);
    var db: Db = client.db(dbName);

    // check if any collection exists
    db.createCollection('channels');
    db.createCollection('messages', function () {
        client.close();
    });
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