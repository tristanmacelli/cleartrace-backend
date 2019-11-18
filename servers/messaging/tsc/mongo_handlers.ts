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

// insertNewChannel takes in a new Channel and
export function insertNewChannel(channels: Collection, newChannel: Channel) {
    let errString: string = "";
    let idWeWant: any;
    channels.save({
        name: newChannel.name, description: newChannel.description,
        private: newChannel.private, members: newChannel.members,
        createdAt: newChannel.createdAt, creator: newChannel.creator,
        editedAt: newChannel.editedAt
    }).catch(() => {
        errString = "Error inserting new channel";
    });
    channels.find({ name: newChannel.name, createdAt: newChannel.createdAt }).next()
        .then(doc => {
            idWeWant = doc.id
        }).catch(err => {
            idWeWant = ""
        });
    newChannel._id = idWeWant;
    return { newChannel, errString };
}

// insertNewMessage takes in a new Message and
export function insertNewMessage(messages: Collection, newMessage: Message) {
    let errString: string = "";
    let idWeWant: any;
    if (newMessage.channelID == null) {
        errString = "Could not find ID";
        return { newMessage, errString };
    }
    let result = messages.save({
        channelID: newMessage.channelID, createdAt: newMessage.createdAt,
        body: newMessage.body, creator: newMessage.creator,
        editedAt: newMessage.editedAt
    }).catch(() => {
        errString = "Error inserting new message";
    });
    messages.find({ body: newMessage.body, createdAt: newMessage.createdAt }).next()
        .then(doc => {
            idWeWant = doc.id
        }).catch(err => {
            idWeWant = ""
        });
    newMessage._id = idWeWant;
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
        errString = "Error updating message";
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
        errString = "Error updating message";
    });
    return errString;
}

// removeChannelMember takes an existing Channel and removes members using a req (request) object
export function removeChannelMember(channels: Collection, existingChannel: Channel, req: any): string {
    // Remove the specified member from this channel's list of members
    let errString: string = "";
    existingChannel.members.splice(req.body.message.id, 1);
    channels.save({
        name: existingChannel.name, description: existingChannel.description,
        private: existingChannel.private, members: existingChannel.members,
        createdAt: existingChannel.createdAt, creator: existingChannel.creator,
        editedAt: existingChannel.editedAt
    }).catch(() => {
        errString = "Error updating message";
    });
    return errString;
}

export function updateMessage(messages: Collection, existingMessage: Message, req: any) {
    let errString: string = "";
    messages.save({
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
    // We are not allowed to delete the general channel
    let errString: string = "";
    if (existingChannel.creator == -1) {
        return "Error deleting channel";
    }
    channels.remove({ _id: new ObjectID(existingChannel._id) }).catch(() => {
        errString = "Error deleting channel";
    });
    messages.remove({ channelID: existingChannel._id }).catch(() => {
        errString = "Error deleting channel";
    });
    return errString;
}

// deleteMessage does something
export function deleteMessage(messages: Collection, existingMessage: Message): string {
    let errString: string = "";
    messages.remove({ messageID: existingMessage._id }).catch(() => {
        errString = "Error deleting message";
    });
    return errString;
}

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