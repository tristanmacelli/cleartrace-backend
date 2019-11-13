"use strict";



//require the express and morgan packages
const express = require("express");
const morgan = require("morgan");
var http = require('http');

//create a new express application
const app = express();

const addr = process.env.ADDR || ":80";
//split host and port using destructuring
const [host, port] = addr.split(":");

//add JSON request body parsing middleware
app.use(express.json());
//add the request logging middleware
app.use(morgan("dev"));


app.use("/v1/channels", (req, res, next) => {
    switch (req.method) {
        case 'GET':
            res.set("Content-Type", "application/json");

            //write those to the client, encoded in JSON
            res.json(allChannels);
            break;
            
        case 'POST':
                console.log(req.body)
                if (req.body.channel.name == null) {
                    next()
                    //do something about the name property being null
                } 
                var bodyChannel = req.body.channel;
                var insert = new channel(bodyChannel.ID, bodyChannel.Name, bodyChannel.Description, bodyChannel.Private,
                                         bodyChannel.Members, bodyChannel.CreatedAt, bodyChannel.Creator, bodyChannel.EditedAt);
                // Call databse to insert this new channel
                res.set("Content-Type", "application/json");
                res.json(insert);
                res.status(201)  //probably cant do this >>> .send("success");
                break;
        default:
            break;
    }
});


//error handler that will be called if
//any handler earlier in the chain throws
//an exception or passes an error to next()
app.use((err, req, res, next) => {
    //write a stack trace to standard out,
    //which writes to the server's log
    console.error(err.stack)

    //but only report the error message
    //to the client, with a 500 status code
    res.set("Content-Type", "text/plain");
    res.status(500).send(err.message);
});

