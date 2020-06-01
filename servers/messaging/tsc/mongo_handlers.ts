"use strict";

// to compile run tsc --outDir ../

import { ObjectID, Collection } from "mongodb";
import { Channel } from "./channel";
import { Message } from "./message";
import { User } from "./user";

// getAllChannels does something
// TODO: make sure the returned value is a shape that we can actually use
export async function getAllChannels(channels: Collection) {
    let allChannelsJSON: string = ""
    let err: boolean = false;
    // if channels does not yet exist
    let cursor = channels.find();
    await cursor.hasNext().then(async () => {
        await cursor.toArray().then((result) => {
            let channelsArray: string[] = []
            for (let i = 0; i < result.length; i++) {
                channelsArray.push(JSON.stringify(result[i]))
            }
            allChannelsJSON = JSON.stringify(channelsArray)
        }).catch(() => {
            err = true;
        })
    }).catch(() => {
        err = true;
    })
    return { allChannelsJSON, err };
}

const createChannel = async (channels: Collection, newChannel: Channel, duplicates: boolean, err: boolean) => {
    let rightNow = new Date()
    try {
        await channels.find({ name: newChannel.name }).next()
            .then(async (doc) => {
                if (doc == null) {
                    newChannel.createdAt = rightNow
                    duplicates = false;
                    await channels.save({
                        name: newChannel.name, description: newChannel.description,
                        private: newChannel.private, members: newChannel.members,
                        createdAt: newChannel.createdAt, creator: newChannel.creator,
                        editedAt: newChannel.editedAt
                    }).then(async () => {
                        await channels.find({ name: newChannel.name, createdAt: newChannel.createdAt }).next()
                            .then(doc => {
                                newChannel.id = doc._id
                            })
                    }).catch(() => {
                        err = true;
                    });

                }
            })
    } catch (e) {
        console.log(e)
    }
    return { newChannel, duplicates, err }
}

// insertNewChannel takes in a new Channel and
export const insertNewChannel = async (channels: Collection, newChannel: Channel) => {
    let err: boolean = false;
    let duplicates: boolean = true;

    return await createChannel(channels, newChannel, duplicates, err)
}

// insertNewMessage takes in a new Message and
export async function insertNewMessage(messages: Collection, newMessage: Message) {
    let err: boolean = false;
    newMessage.createdAt = new Date()

    await messages.save({
        channelID: newMessage.channelID, createdAt: newMessage.createdAt,
        body: newMessage.body, creator: newMessage.creator,
        editedAt: newMessage.editedAt
    }).then(async () => {
        await messages.find({ body: newMessage.body, createdAt: newMessage.createdAt }).next()
            .then(doc => {
                newMessage.id = doc._id
            }).catch(() => {
                newMessage.id = ""
            });
    }).catch(() => {
        err = true;
    });

    return { newMessage, err };
}

// updatedChannel updates name and body of an existing Channel using a req (request) object
export async function updateChannel(channels: Collection, existingChannel: Channel, req: any) {
    let err: boolean = false;

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
    // Removed .message after body in the statement below
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
    // Removed .message after body in the statement below
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

    await messages.save({
        body: req.body, creator: existingMessage.creator,
        createdAt: existingMessage.createdAt, channelID: existingMessage.channelID,
        editedAt: existingMessage.editedAt
    }).catch(() => {
        err = true
    });

    existingMessage.body = req.body;
    return { existingMessage, err };
}

// deleteChannel does something
export async function deleteChannel(channels: Collection, messages: Collection, existingChannel: Channel): Promise<boolean> {
    let err: boolean = false;
    // We are not allowed to delete the general channel
    if (existingChannel.creator.id == -1) {
        err = true;
    }
    // CHANNEL ID DOES NOT EXIST
    let channelID = new ObjectID(existingChannel.id)
    await channels.remove({ _id: channelID }).catch(() => {
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
    let findResult: any;
    let err: boolean = false;
    let mongoID = new ObjectID(id)
    await channels.find({ _id: mongoID }).next().then(doc => {
        findResult = doc;
        err = false;
    }).catch(() => {
        findResult = null
        err = true;
    });
    let channel: Channel;
    if (findResult == null) {
        let emptyUser = new User(-1, "", new Uint8Array(100), "", "", "", "")
        let dummyDate = new Date()
        channel = new Channel("", "", false, [], dummyDate, emptyUser, "");
        return { channel, err };
    }
    channel = new Channel(findResult.name, findResult.description, findResult.private,
        findResult.members, findResult.createdAt, findResult.creator, findResult.editedAt);
    return { channel, err };
}

// getMessageByID does something
export async function getMessageByID(messages: Collection, id: string) {
    // Since id's are auto-generated and unique we chose to use find instead of findOne() 
    let findResult: any;
    let err: boolean = false;
    let mongoID = new ObjectID(id)
    await messages.find({ _id: mongoID }).next().then(doc => {
        findResult = doc;
        err = false;
    }).catch(() => {
        findResult = null;
        err = true;
    })

    let message: Message;
    if (findResult == null) {
        let emptyUser = new User(-1, "", new Uint8Array(100), "", "", "", "")
        let dummyDate = new Date()
        message = new Message("", dummyDate, "", emptyUser, dummyDate);
        return { message, err };
    }
    message = new Message(findResult.channelID, findResult.createdAt, findResult.body,
        findResult.creator, findResult.editedAt)
    return { message, err }
}

// TODO: Reshape the return value of find to a JSON array of message model objects
// last100Messages does something
export async function last100Messages(messages: Collection, channelID: string) {
    let last100messages: string[] = []
    let err: boolean = false;

    if (channelID == null) {
        err = true;
        return { last100messages, err };
    }
    channelID = channelID.toString();
    let cursor = messages.find({ channelID: channelID }).sort({ createdAt: -1 }).limit(100);

    await cursor.hasNext().then(async () => {
        await cursor.toArray().then((result) => {
            for (let i = 0; i < result.length; i++) {
                last100messages.push(JSON.stringify(result[i]))
            }
        }).catch(() => {
            err = true;
        })
    }).catch(() => {
        err = true;
    })
    return { last100messages, err };
}

// TODO: Reshape the return value of find to a JSON array of message model objects
// last100Messages does something
export async function last100SpecificMessages(messages: Collection, channelID: string, messageID: string) {
    let last100messages: string[] = []
    let err: boolean = false;

    if (channelID == null) {
        err = true;
        return { last100messages, err };
    }
    channelID = channelID.toString();
    let objID = new ObjectID(messageID)
    let cursor = messages.find({ channelID: channelID, _id: { $lt: objID } }).sort({ createdAt: -1 }).limit(100);

    await cursor.hasNext().then(async () => {
        await cursor.toArray().then((result) => {
            for (let i = 0; i < result.length; i++) {
                last100messages.push(JSON.stringify(result[i]))
            }
        }).catch(() => {
            err = true;
        })
    }).catch(() => {
        err = true;
    })
    return { last100messages, err };
}

export * from "./mongo_handlers";