"use strict";

// to compile run tsc --outDir ../

import { ObjectID, Collection } from "mongodb";
import { Channel } from "./channel";
import { Message } from "./message";

// getAllChannels does something
// TODO: make sure the returned value is a shape that we can actually use
export function getAllChannels(channels: Collection) {
    // if channels does not yet exist
    let cursor = channels.find();
    if (!cursor.hasNext()) {
        // Throw error
        console.log("No channels collection found");
        return null;
    }
    return cursor.forEach(function (m: any) { JSON.stringify(m) });
}

// TODO: for each function that is returning result as a raw Promise<WriteOpResult>
//       we need to do a transformation of that data into an actual Channel or Message object

// insertNewChannel takes in a new Channel and
export function insertNewChannel(channels: Collection, newChannel: Channel): Channel | null {
    let result = channels.save({
        name: newChannel.name, description: newChannel.description,
        private: newChannel.private, members: newChannel.members,
        createdAt: newChannel.createdAt, creator: newChannel.creator,
        editedAt: newChannel.editedAt
    });
    if (result) {
        return null;
    }
    newChannel._id = result._id;
    return newChannel;
}

// insertNewMessage takes in a new Message and
export function insertNewMessage(messages: Collection, newMessage: Message): Message | null {
    if (newMessage.channelID == null) {
        return null;
    }
    let result = messages.save({
        channelID: newMessage.channelID, createdAt: newMessage.createdAt,
        body: newMessage.body, creator: newMessage.creator,
        editedAt: newMessage.editedAt
    });
    if (result.hasWriteError()) {
        return null;
    }
    newMessage._id = result._id;
    return newMessage;
}

// updatedChannel updates name and body of an existing Channel using a req (request) object
export function updateChannel(channels: Collection, existingChannel: Channel, req: any): Channel | null {
    let result = channels.save({
        name: req.body.name, description: req.body.description,
        private: existingChannel.private, members: existingChannel.members,
        createdAt: existingChannel.createdAt, creator: existingChannel.creator,
        editedAt: existingChannel.editedAt
    });
    if (result.hasWriteError()) {
        return null;
    }
    existingChannel.name = req.body.name;
    existingChannel.description = req.body.description;
    return existingChannel;
}

// addChannelMembers takes an existing Channel and adds members using a req (request) object
export function addChannelMember(channels: Collection, existingChannel: Channel, req: any): Channel | null {
    existingChannel.members.push(req.body.message.id);
    let result = channels.save({
        name: existingChannel.name, description: existingChannel.description,
        private: existingChannel.private, members: existingChannel.members,
        createdAt: existingChannel.createdAt, creator: existingChannel.creator,
        editedAt: existingChannel.editedAt
    });
    if (result.hasWriteError()) {
        return null;
    }
    return existingChannel;
}

// addChannelMembers takes an existing Channel and removes members using a req (request) object
export function removeChannelMember(channels: Collection, existingChannel: Channel, req: any): Channel | null {
    // Remove the specified member from this channel's list of members
    existingChannel.members.splice(req.body.message.id, 1);
    let result = channels.save({
        name: existingChannel.name, description: existingChannel.description,
        private: existingChannel.private, members: existingChannel.members,
        createdAt: existingChannel.createdAt, creator: existingChannel.creator,
        editedAt: existingChannel.editedAt
    });
    if (result.hasWriteError()) {
        return null;
    }
    return existingChannel;
}

export function updateMessage(messages: Collection, existingMessage: Message, req: any) {
    let result = messages.save({
        body: req.body, creator: existingMessage.creator,
        createdAt: existingMessage.createdAt, channelID: existingMessage.channelID,
        editedAt: existingMessage.editedAt
    });
    if (result.hasWriteError()) {
        return null;
    }
    existingMessage.body = req.body;
    return existingMessage;
}

// deleteChannel does something
export function deleteChannel(channels: Collection, messages: Collection, existingChannel: Channel) {
    // We are not allowed to delete the general channel
    if (existingChannel.creator == -1) {
        return null;
    }
    channels.remove({ _id: new ObjectID(existingChannel._id) });
    let result = messages.remove({ channelID: existingChannel._id });
    if (result.hasWriteError()) {
        return null;
    }
    return result;
}

// deleteMessage does something
export function deleteMessage(messages: Collection, existingMessage: Message) {
    let result = messages.remove({ messageID: existingMessage._id });
    if (result.hasWriteError()) {
        return null;
    }
    return result;
}

// getByChannelID does something
// TODO: process cursor returned by .find() and return a new Channel object
// export async function getChannelByID(channels: Collection, id: string, res: any) {
//     // Since id's are auto-generated and unique we chose to use find instead of findOne()
//     let channelCursor = await channels.find({ _id: id }, function (err, res) {

//     });//.limit(1);//.toArray().then((result) => result);
//     // .findOne({ _id: id });
//     // if (channelObj) {
//     //     new Channel(channelObj.name, channelObj.description, channelObj.private,
//     //         channelObj.members, channelObj.createdAt, channelObj.creator, channelObj.editedAt)
//     //     return;
//     // }
//     let channelJSON: string = "";
//     while (await channelCursor.hasNext()) {
//         let channelObj = await channelCursor.next();
//         channelJSON = JSON.stringify(channelObj);
//     }
//     return channelJSON;
// }

export function getChannelByID(channels: Collection, id: string) {
    // Since id's are auto-generated and unique we chose to use find instead of findOne() 
    let finalResponse: any;
    let errString: any;

    channels.find({ _id: id }).next().then(doc => {
        finalResponse = doc
        errString = ""
    }).catch(err => {
        finalResponse = null
        errString = err
    });
    let finalChannel: Channel;
    if (finalResponse == null) {
        finalChannel = new Channel("", "", false, [], "", -1, "");
        return { finalChannel, errString };
    }
    finalChannel = new Channel(finalResponse.name, finalResponse.description, finalResponse.private,
        finalResponse.members, finalResponse.createdAt, finalResponse.Creator, finalResponse.editedAt);
    return { finalChannel, errString };
}

// TODO: process cursor returned by .find() and return a new Message object
export function getMessageByID(messages: Collection, id: string) {

    // Since id's are auto-generated and unique we chose to use find instead of findOne() 
    let finalResponse: any;
    let errString: any;

    messages.find({ _id: id }).next().then(doc => {
        finalResponse = doc;
        errString = "";
    }).catch(err => {
        finalResponse = null;
        errString = err;
    })

    let finalMessage: Message;
    if (finalResponse == null) {
        finalMessage = new Message("", "", "", "", "");
        return { finalMessage, errString };
    }
    finalMessage = new Message(finalResponse.channelID, finalResponse.createdAt, finalResponse.body,
        finalResponse.creator, finalResponse.editedAt)
    return { finalMessage, errString }
}

// let messageObj: Channel[] = await messages.find({ _id: id }).limit(1).toArray();
// if (messageObj.length == 0) {
//     return null;
// }
// return messageObj[0];

// .findOne({ _id: id });
// if (channelObj) {
//     new Channel(channelObj.name, channelObj.description, channelObj.private,
//         channelObj.members, channelObj.createdAt, channelObj.creator, channelObj.editedAt)
//     return;
// }
// find({ _id: id }).forEach(element => {
//     print(element)
//     channelObj = element;
//     newChannel = new Channel(channelObj.name, channelObj.description, channelObj.private,
//         channelObj.members, channelObj.createdAt, channelObj.creator, channelObj.editedAt);
// });

// last100Messages does something
export function last100Messages(messages: Collection, id: string) {
    if (id == null) {
        throw "No id value passed";
    }
    id = id.toString();
    return messages.find({ channelID: id }).sort({ createdAt: -1 }).limit(100);
}

// last100Messages does something
export function last100SpecificMessages(messages: Collection, channelID: string, messageID: string) {
    if (channelID == null) {
        throw "No id value passed";
    }
    channelID = channelID.toString();
    return messages.find({ channelID: channelID, _id: { $lt: messageID } }).sort({ createdAt: -1 }).limit(100);
}

export * from "./mongo_handlers";