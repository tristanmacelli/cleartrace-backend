"use strict";

import { ObjectID, Collection, MongoClient, Cursor } from "mongodb";
import { Channel, isChannelMember } from "./channel";
import { Message } from "./message";
import { User } from "./user";

const mongoURL = 'mongodb://mongodb:27017/mongodb';

// Create a new MongoClient
const mc = new MongoClient(mongoURL, { useUnifiedTopology: true });

export async function createConnection() {
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

// getAllChannels does something
export async function getChannels(channels: Collection, userID: number, search: string) {
    let err: boolean = false;
    let allChannels: Channel[] = []
    // if channels does not yet exist
    let cursor: Cursor<any>
    if (!search) {
        cursor = channels.find();
    } else {
        cursor = channels.find({ name: { $regex: "/^"+ search +"/i" } });
    }
    if (await cursor.hasNext()) {    
        let result = await cursor.toArray()
        for (let i = 0; i < result.length; i++) {
            let channel = new Channel(result[i]._id, result[i].name, result[i].description, result[i].private,
                result[i].members, result[i].createdAt, result[i].creator, result[i].editedAt);
            if (isChannelMember(channel, userID)) {
                allChannels.push(channel)
            }               
        }
    } else {
        err = true
    }
    return { allChannels, err };
}

// insertNewChannel takes in a new Channel and
export async function insertNewChannel(channels: Collection, newChannel: Channel) {
    let err: boolean = false;
    let duplicates: boolean = true;
    let rightNow = new Date()

    let cursor = channels.find({ name: newChannel.name })
    if (!await cursor.hasNext()) {
        duplicates = false;
        newChannel.createdAt = rightNow
        await channels.save({
            name: newChannel.name, description: newChannel.description,
            private: newChannel.private, members: newChannel.members,
            createdAt: newChannel.createdAt, creator: newChannel.creator,
            editedAt: newChannel.editedAt
        }).catch(() => {
            err = true;
        });
        cursor = channels.find({ name: newChannel.name, createdAt: newChannel.createdAt })
        if (await cursor.hasNext()) {
            let doc = await cursor.next()
            newChannel.id = doc._id
        } else {
            err = true;
        }
    }           
    return { newChannel, duplicates, err }
}

// insertNewMessage takes in a new Message and
export async function insertNewMessage(messages: Collection, newMessage: Message) {
    let err: boolean = false;
    newMessage.createdAt = new Date()

    await messages.save({
        channelID: newMessage.channelID, createdAt: newMessage.createdAt,
        body: newMessage.body, creator: newMessage.creator,
        editedAt: newMessage.editedAt
    }).catch(() => {
        err = true;
    });

    let cursor = messages.find({ body: newMessage.body, createdAt: newMessage.createdAt })
    if (await cursor.hasNext()) {
        let doc = await cursor.next()
        newMessage.id = doc._id
    } else {
        newMessage.id = ""
    }

    return { newMessage, err };
}

// updatedChannel updates name and body of an existing Channel using a req (request) object
export async function updateChannel(channels: Collection, existingChannel: Channel, req: any) {
    let err: boolean = false;
    existingChannel.editedAt = new Date()

    await channels.save({
        name: req.body.name, description: req.body.description,
        private: existingChannel.private, members: existingChannel.members,
        createdAt: existingChannel.createdAt, creator: existingChannel.creator,
        editedAt: existingChannel.editedAt
    }).catch(() => {
        err = true;
    });
    existingChannel.name = req.body.name;
    existingChannel.description = req.body.description;
    return { existingChannel, err };
}

// addChannelMembers takes an existing Channel and adds members using a req (request) object
export async function addChannelMember(channels: Collection, existingChannel: Channel, req: any): Promise<boolean> {
    let err: boolean = false;
    existingChannel.members.push(req.body.id);

    await channels.save({
        name: existingChannel.name, description: existingChannel.description,
        private: existingChannel.private, members: existingChannel.members,
        createdAt: existingChannel.createdAt, creator: existingChannel.creator,
        editedAt: existingChannel.editedAt
    }).catch(() => {
        err = true;
    });
    return err;
}

// removeChannelMember takes an existing Channel and removes members using a req (request) object
export async function removeChannelMember(channels: Collection, existingChannel: Channel, req: any): Promise<boolean> {
    // Remove the specified member from this channel's list of members
    let err: boolean = false;
    existingChannel.members.splice(req.body.id, 1);
    await channels.save({
        name: existingChannel.name, description: existingChannel.description,
        private: existingChannel.private, members: existingChannel.members,
        createdAt: existingChannel.createdAt, creator: existingChannel.creator,
        editedAt: existingChannel.editedAt
    }).catch(() => {
        err = true;
    });
    return err;
}

// updateMessage takes an existing Message and a request with updates to apply to the Message's body 
export async function updateMessage(messages: Collection, existingMessage: Message, req: any) {
    let err: boolean = false;
    existingMessage.editedAt = new Date()

    await messages.save({
        body: req.body.body, creator: existingMessage.creator,
        createdAt: existingMessage.createdAt, channelID: existingMessage.channelID,
        editedAt: existingMessage.editedAt
    }).catch(() => {
        err = true
    });

    existingMessage.body = req.body.body;
    existingMessage.id = req.params.messageID;
    return { existingMessage, err };
}

// deleteChannel does something
export async function deleteChannel(channels: Collection, messages: Collection, existingChannel: Channel): Promise<boolean> {
    let err: boolean = false;
    // The general channel never gets deleted
    if (existingChannel.creator.ID == -1) {
        err = true;
    }

    let chanID = new ObjectID(existingChannel.id.toString())
    await channels.remove({ _id: chanID }).catch(() => {
        err = true;
    });
    await messages.remove({ channelID: existingChannel.id }).catch(() => {
        err = true;
    });
    // Try this version which is not deprecated
    // await messages.deleteMany({ channelID: existingChannel.id }).catch(() => {
    //     errString = "Error deleting messages associated with the channel";
    // });
    return err;
}

// deleteMessage does something
export async function deleteMessage(messages: Collection, existingMessage: Message): Promise<boolean> {
    await messages.remove({ messageID: existingMessage.id }).catch(() => {
        return true;
    });
    return false;
}

// getChannelByID does something
export async function getChannelByID(channels: Collection, id: string) {
    // Since id's are auto-generated and unique we chose to use find instead of findOne() 
    let result: any = null;
    let err: boolean = true;
    let mongoID = new ObjectID(id)

    let cursor = channels.find({ _id: mongoID })
    if (await cursor.hasNext()) {
        result = await cursor.next()
        err = false;
    }
    let channel: Channel;
    if (result == null) {
        let emptyUser = new User(-1, "", "", "", "", "")
        let dummyDate = new Date()
        channel = new Channel("", "", "", false, [], dummyDate, emptyUser, dummyDate);
        return { channel, err };
    }
    channel = new Channel(result._id, result.name, result.description, result.private,
        result.members, result.createdAt, result.creator, result.editedAt);
    return { channel, err };
}

// getMessageByID does something
export async function getMessageByID(messages: Collection, id: string) {
    // Since id's are auto-generated and unique we chose to use find instead of findOne() 
    let result: any = null;
    let err: boolean = true;
    let mongoID = new ObjectID(id)

    let cursor = messages.find({ _id: mongoID })
    if (await cursor.hasNext()) {
        result = await cursor.next()
        err = false;
    }

    let message: Message;
    if (result == null) {
        let emptyUser = new User(-1, "", "", "", "", "")
        let dummyDate = new Date()
        message = new Message("", "", dummyDate, "", emptyUser, dummyDate);
        return { message, err };
    }
    message = new Message(result._id, result.channelID, result.createdAt, result.body,
        result.creator, result.editedAt)
    return { message, err }
}

// TODO: Reshape the return value of find to a JSON array of message model objects
// last100Messages does something
export async function last100Messages(messages: Collection, channelID: string, messageID: string) {
    let last100messages: Message[] = []
    let err: boolean = false;

    if (channelID == null) {
        err = true;
        return { last100messages, err };
    }
    channelID = channelID.toString();
    let cursor
    if (!messageID) {
        cursor = messages.find({ channelID: channelID }).sort({ createdAt: -1 }).limit(100);
    } else {
        let objID = new ObjectID(messageID)
        cursor = messages.find({ channelID: channelID, _id: { $lt: objID } }).sort({ createdAt: -1 }).limit(100);
    }

    if (await cursor.hasNext()) {    
        let result = await cursor.toArray()
        for (let i = 0; i < result.length; i++) {
            let message = new Message(result[i]._id, result[i].channelID, result[i].createdAt, result[i].body,
                result[i].creator, result[i].editedAt)
            last100messages.push(message)
        }
    }
    return { last100messages, err };
}

export function sleep(seconds: number) {
    let milliseconds = seconds * 1000
    const stop = new Date().getTime() + milliseconds;
    while(new Date().getTime() < stop);       
}

export * from "./mongo_handlers";