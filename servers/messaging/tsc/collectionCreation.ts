//require the express and morgan packages
import { MongoClient, Db, Collection } from "mongodb";

// Connection URL
const url = 'mongo://mongodb:27017';

// Database Name
const dbName = 'mongodb';

// Create a new MongoClient
// const client = new MongoClient(url,  { useUnifiedTopology: true });
var db: Db;

const client = new MongoClient(url);

client.connect(function (err: any) {
    console.log("Connected successfully to server");
    console.log("the error is ", err)
    db = client.db(dbName);

    // check if any collection exists
    db.createCollection('channels', function (err, collection) {

    });
    db.createCollection('messages', function (err, collection) {
        client.close()
    });
});


// Reasoning for refactor: 
// https://mongodb.github.io/node-mongodb-native/driver-articles/mongoclient.html#mongoclient-connection-pooling
// Use connect method to connect to the mongo DB
// MongoClient.connect(url, function (err: any, d:MongoClient) {
//     console.log("Connected successfully to server");
//     console.log("there is a fucking error:", err.err)
//     if (err) throw err;

//     // db = client.db(dbName);
//     var db = d.db(dbName);
//     // Create db.channels and db.messages collections in mongo
//     // https://mongodb.github.io/node-mongodb-native/api-articles/nodekoarticle1.html#mongo-db-and-collections
//     db.createCollection('channels', function (err, collection) {

//     });
//     db.createCollection('messages', function (err, collection) {
//         d.close()
//     });
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
// });