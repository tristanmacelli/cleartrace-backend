"use strict";

// to compile run tsc --outDir ../

import { ObjectID, Collection } from "mongodb";
import { Channel } from "./channel";
import { Message } from "./message";
import { is } from "bluebird";
import { User } from "./user";

// getAllChannels does something
// TODO: make sure the returned value is a shape that we can actually use
export async function getAllChannels(channels: Collection, res: any) {
    let resultJSON: string = "bbb"
    let successMessage: string = ""
    // if channels does not yet exist
    let cursor = channels.find();
    await cursor.hasNext().then(async () => {
        await cursor.toArray().then(async (result) => {
            successMessage = "Found channels";
            console.log(successMessage);

            let channelsArray:string[] = []
            for (let i = 0 ; i < result.length; i++) {
                channelsArray.push(JSON.stringify(result[i]))
                // channelsString += ","
            }
            resultJSON = JSON.stringify(channelsArray)
        })         
    })
        return {resultJSON, successMessage};
    }

const createChannel = async (channels: Collection, newChannel: Channel, errString: string) => {
    let rightNow = new Date()
    console.log("the date we are trying to specify is")
    console.log(rightNow)
    try {
        await channels.find({ name: newChannel.name}).next()
        .then(async (doc) => {
            if (doc == null) {
                newChannel.createdAt = rightNow
                errString = ""
                console.log("NOT a duplicate channel")
                errString
                await channels.save({
                    name: newChannel.name, description: newChannel.description,
                    private: newChannel.private, members: newChannel.members,
                    createdAt: newChannel.createdAt, creator: newChannel.creator,
                    editedAt: newChannel.editedAt 
                }).then(async () => {
                    await channels.find({ name: newChannel.name, createdAt: newChannel.createdAt }).next()
                    .then(doc => {
                        console.log("THe doc in the second find is")
                        console.log(doc)
                        newChannel.id = doc._id
                    })
                }).catch(() => {
                    errString = "Error inserting new channel";
                });
                
            }
        })
    } catch(e) {
        console.log(e)
    }
    console.log("The new Channel that we send is")
    console.log(newChannel)
    return {newChannel, errString}
}

// insertNewChannel takes in a new Channel and
export const insertNewChannel = async (channels: Collection, newChannel: Channel) =>{
    let errString: string = "duplicate";

    return await createChannel(channels,newChannel,errString)
}

// insertNewMessage takes in a new Message and
export function insertNewMessage(messages: Collection, newMessage: Message) {
    let errString: string = "";
    let autoAssignedID: any;
    messages.save({
        channelID: newMessage.channelID, createdAt: newMessage.createdAt,
        body: newMessage.body, creator: newMessage.creator,
        editedAt: newMessage.editedAt
    }).catch(() => {
        errString = "Error inserting new message";
    });
    messages.find({ body: newMessage.body, createdAt: newMessage.createdAt }).next()
        .then(doc => {
            autoAssignedID = doc._id
        }).catch(() => {
            autoAssignedID = ""
        });
    newMessage.id = autoAssignedID;
    return { newMessage, errString };
}

// updatedChannel updates name and body of an existing Channel using a req (request) object
export function updateChannel(channels: Collection, existingChannel: Channel, req: any) {
    let errString: string = "";
    channels.save({
        name: req.body.name, description: req.body.description,
        private: existingChannel.private, members: existingChannel.members,
        createdAt: existingChannel.createdAt, creator: existingChannel.creator,
        editedAt: existingChannel.editedAt
    }).catch(() => {
        errString = "Error updating channel";
    });
    existingChannel.name = req.body.name;
    existingChannel.description = req.body.description;
    return { existingChannel, errString };
}

// addChannelMembers takes an existing Channel and adds members using a req (request) object
export function addChannelMember(channels: Collection, existingChannel: Channel, req: any): string {
    let errString: string = "";
    existingChannel.members.push(req.body.message.id);
    channels.save({
        name: existingChannel.name, description: existingChannel.description,
        private: existingChannel.private, members: existingChannel.members,
        createdAt: existingChannel.createdAt, creator: existingChannel.creator,
        editedAt: existingChannel.editedAt
    }).catch(() => {
        errString = "Error adding new members to channel";
    });
    return errString;
}

// removeChannelMember takes an existing Channel and removes members using a req (request) object
export async function removeChannelMember(channels: Collection, existingChannel: Channel, req: any): Promise<string> {
    // Remove the specified member from this channel's list of members
    let errString: string = "";
    existingChannel.members.splice(req.body.message.id, 1);
    await channels.save({
        name: existingChannel.name, description: existingChannel.description,
        private: existingChannel.private, members: existingChannel.members,
        createdAt: existingChannel.createdAt, creator: existingChannel.creator,
        editedAt: existingChannel.editedAt
    }).catch(() => {
        errString = "Error removing member from channel";
    });
    return errString;
}

// updateMessage takes an existing Message and a request with updates to apply to the Message's body 
export async function updateMessage(messages: Collection, existingMessage: Message, req: any) {
    let errString: string = "";
    await messages.save({
        body: req.body, creator: existingMessage.creator,
        createdAt: existingMessage.createdAt, channelID: existingMessage.channelID,
        editedAt: existingMessage.editedAt
    }).catch(() => {
        errString = "Error updating message";
    });
    existingMessage.body = req.body;
    return { existingMessage, errString };
}

// deleteChannel does something
export function deleteChannel(channels: Collection, messages: Collection, existingChannel: Channel): string {
    let errString: string = "";
    // We are not allowed to delete the general channel
    if (existingChannel.creator.id == -1) {
        return "Error deleting channel";
    }
    channels.remove({ id: new ObjectID(existingChannel.id) }).catch(() => {
        errString = "Error deleting channel";
    });
    messages.remove({ channelID: existingChannel.id }).catch(() => {
        errString = "Error deleting messages associated with the channel";
    });
    return errString;
}

// deleteMessage does something
export async function deleteMessage(messages: Collection, existingMessage: Message): Promise<string> {
    let errString: string = "";
    await messages.remove({ messageID: existingMessage.id }).catch(() => {
        errString = "Error deleting message";
    });
    return errString;
}

// getChannelByID does something
export function getChannelByID(channels: Collection, id: string) {
    // Since id's are auto-generated and unique we chose to use find instead of findOne() 
    let finalResponse: any;
    let errString: any;

    channels.find({ id: id }).next().then(doc => {
        finalResponse = doc
        errString = ""
    }).catch(() => {
        finalResponse = null
        errString = "Error finding a channel by id"
    });
    let finalChannel: Channel;
    if (finalResponse == null) {
        let emptyUser = new User(-1, "", new Uint8Array(100), "", "", "", "")
        let dummyDate = new Date()
        finalChannel = new Channel("", "", false, [], dummyDate, emptyUser, "");
        return { finalChannel, errString };
    }
    finalChannel = new Channel(finalResponse.name, finalResponse.description, finalResponse.private,
        finalResponse.members, finalResponse.createdAt, finalResponse.Creator, finalResponse.editedAt);
    return { finalChannel, errString };
}

// getMessageByID does something
export function getMessageByID(messages: Collection, id: string) {
    // Since id's are auto-generated and unique we chose to use find instead of findOne() 
    let finalResponse: any;
    let errString: any;

    messages.find({ id: id }).next().then(doc => {
        finalResponse = doc;
        errString = "";
    }).catch(() => {
        finalResponse = null;
        errString = "Error finding a message by id";
    })

    let finalMessage: Message;
    if (finalResponse == null) {
        let emptyUser = new User(-1, "", new Uint8Array(100), "", "", "", "")
        finalMessage = new Message("", "", "", emptyUser, "");
        return { finalMessage, errString };
    }
    finalMessage = new Message(finalResponse.channelID, finalResponse.createdAt, finalResponse.body,
        finalResponse.creator, finalResponse.editedAt)
    return { finalMessage, errString }
}

// TODO: Reshape the return value of find to a JSON array of message model objects
// last100Messages does something
export function last100Messages(messages: Collection, id: string, res: any) {
    if (id == null) {
        console.log("No id value passed");
        return null;
    }
    id = id.toString();
    let cursor = messages.find({ channelID: id }).sort({ createdAt: -1 }).limit(100);
    // TODO: make sure the returned value is a shape that we can actually use
    cursor.toArray(function (err, result) {
        if (err) {
            console.log("Error getting messages");
            return null;
        } else {
            let successMessage = "Found messages";
            console.log(successMessage);
            res.send(JSON.stringify(result));
            return successMessage;
        }
    })
}

// TODO: Reshape the return value of find to a JSON array of message model objects
// last100Messages does something
export function last100SpecificMessages(messages: Collection, channelID: string, messageID: string, res: any) {
    if (channelID == null) {
        console.log("No id value passed");
    }
    channelID = channelID.toString();
    let cursor = messages.find({ channelID: channelID, id: { $lt: messageID } }).sort({ createdAt: -1 }).limit(100);
    // TODO: make sure the returned value is a shape that we can actually use
    cursor.toArray(function (err, result) {
        if (err) {
            console.log("Error getting messages");
        } else {
            let successMessage = "Found specific messages";
            console.log(successMessage);
            res.send(JSON.stringify(result));
            return successMessage;
        }
    })
}

export * from "./mongo_handlers";


// save({
//     "_id": "5dd720def3df9b13a39876e7",
//     "name": "saurav",
//     "description": "an open channel for all",
//     "private": false,
//     "members": [1],
//     "createdAt": "enter timestamp here",
//     "creator": -1,
//     "editedAt": "not yet edited"
// })