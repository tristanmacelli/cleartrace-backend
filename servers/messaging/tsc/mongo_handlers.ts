"use strict";

import { ObjectID, Collection, MongoClient, Cursor, Db } from "mongodb";
import { Channel, isChannelMember, initializeDummyChannel } from "./channel";
import { Message, initializeDummyMessage } from "./message";

const mongoContainerName = 'userMessageStore'
const dbName = 'userMessageDB';
const mongoURL = 'mongodb://' + mongoContainerName + ':27017/' + dbName;

// Create a new MongoClient
const mc = new MongoClient(mongoURL, { useUnifiedTopology: true });

export async function createConnection(): Promise<Db> {
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
    return client!.db(dbName);
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

// insertNewChannel takes in a new Channel and inserts it into the messaging DB
export async function insertNewChannel(channels: Collection, newChannel: Channel) {
    let err: boolean = false;
    let duplicates: boolean = true;
    newChannel.createdAt = new Date()
    const filter = { name: newChannel.name }
    const insertDoc = {
        name: newChannel.name, 
        description: newChannel.description,
        private: newChannel.private, 
        members: newChannel.members,
        createdAt: newChannel.createdAt, 
        creator: newChannel.creator,
        editedAt: newChannel.editedAt
    }

    let cursor = channels.find(filter)
    if (!await cursor.hasNext()) {
        duplicates = false;
        await channels.insertOne(insertDoc)
        .catch(() => {
            err = true;
        });
        cursor = channels.find(filter)
        if (await cursor.hasNext()) {
            let doc = await cursor.next()
            newChannel.id = doc._id
        } else {
            err = true;
        }
    }           
    return { newChannel, duplicates, err }
}

// insertNewMessage takes in a new Message and inserts it into the messaging DB
export async function insertNewMessage(messages: Collection, newMessage: Message) {
    let err: boolean = false;
    newMessage.createdAt = new Date()
    const insertDoc = {
        channelID: newMessage.channelID, 
        createdAt: newMessage.createdAt,
        body: newMessage.body, 
        creator: newMessage.creator,
        editedAt: newMessage.editedAt
    }

    await messages.insertOne(insertDoc)
    .catch(() => {
        err = true;
    });
    const filterNew = { body: newMessage.body, createdAt: newMessage.createdAt }
    let cursor = messages.find(filterNew)
    if (await cursor.hasNext()) {
        let doc = await cursor.next()
        newMessage.id = doc._id
    } else {
        newMessage.id = ""
    }

    return { newMessage, err };
}

// addChannelMembers takes an existing Channel and adds members using a req (request) object
export async function addChannelMember(channels: Collection, existingChannel: Channel, req: any): Promise<boolean> {
  existingChannel.members.push(req.body.id);
  return await updateChannelMembers(channels, existingChannel.id, existingChannel.members)
}

// removeChannelMember takes an existing Channel and removes members using a req (request) object
export async function removeChannelMember(channels: Collection, existingChannel: Channel, req: any): Promise<boolean> {
    // Remove the specified member from this channel's list of members
    existingChannel.members.splice(req.body.id, 1);
    return updateChannelMembers(channels, existingChannel.id, existingChannel.members)
}

async function updateChannelMembers(channels: Collection, channelID: string, channelMembers: number[]): Promise<boolean> {
  let err: boolean = false;
  let newEditedAt = new Date()
  let channelIDObj = new ObjectID(channelID)
  const filter = { _id: channelIDObj };
  const updateDoc = {
    $set: { 
      members: channelMembers, 
      editedAt: newEditedAt
    },
  };

  await channels.updateOne(filter, updateDoc)
  .catch(() => {
    err = true;
  })
  return err;
}

// updatedChannel updates name and body of an existing Channel using a req (request) object
export async function updateChannel(channels: Collection, existingChannel: Channel, req: any) {
  let err: boolean = false;
  existingChannel.editedAt = new Date()
  let channelID = new ObjectID(existingChannel.id.toString())
  const filter = { _id: channelID };
  const updateDoc = {
    $set: { 
      name: req.body.name, 
      description: req.body.description,
      editedAt: existingChannel.editedAt
    },
  };

  await channels.updateOne(filter, updateDoc)
  .catch(() => {
    err = true;
  })
  existingChannel.name = req.body.name;
  existingChannel.description = req.body.description;
  return { existingChannel, err };
}

// updateMessage takes an existing Message and a request with updates to apply to the Message's body 
export async function updateMessage(messages: Collection, existingMessage: Message, req: any) {
    let err: boolean = false;
    existingMessage.editedAt = new Date()
    const filter = { messageID: existingMessage.id };
    const updateDoc = {
      $set: { 
        body: req.body.body,
        editedAt: existingMessage.editedAt
      },
    };

    await messages.updateOne(filter, updateDoc)
    .catch(() => {
      err = true;
    })

    existingMessage.body = req.body.body;
    existingMessage.id = req.params.messageID;
    return { existingMessage, err };
}

// deleteChannel deletes a single channel & its associated messages
export async function deleteChannel(channels: Collection, messages: Collection, existingChannel: Channel): Promise<boolean> {
    let err: boolean = false;
    // The general channel never gets deleted
    if (existingChannel.creator.ID == -1) {
      return true;
    }

    let chanID = new ObjectID(existingChannel.id.toString())
    const channelFilter = { _id: chanID }
    await channels.deleteOne(channelFilter)
    .catch(() => {
      err = true;
    });
    const messageFilter = { channelID: existingChannel.id.toString() }
    await messages.deleteMany(messageFilter)
    .catch(() => {
        err = true;
    });
    return err;
}

// deleteMessage deletes a single message
export async function deleteMessage(messages: Collection, existingMessage: Message): Promise<boolean> {
  const filter = { messageID: existingMessage.id }  
  await messages.deleteOne(filter)
    .catch(() => {
        return true;
    });
    return false;
}

// getChannelByID returns the channel associated with the provided id value. If there no channel 
// is no channel associated with the provided id then an error indicator is returned
export async function getChannelByID(channels: Collection, id: string) {
    // Since id's are auto-generated and unique we chose to use find instead of findOne() 
    let result: any = null;
    let err: boolean = true;
    let mongoID = new ObjectID(id)
    const filter = { _id: mongoID };

    let cursor = channels.find(filter)
    if (await cursor.hasNext()) {
        result = await cursor.next()
        err = false;
    }
    let channel: Channel;
    if (result == null) {
        channel = initializeDummyChannel();
        return { channel, err };
    }
    channel = new Channel(result._id, result.name, result.description, result.private,
        result.members, result.createdAt, result.creator, result.editedAt);
    return { channel, err };
}

// getMessageByID returns the message associated with the provided id value. If there no message 
// is no message associated with the provided id then an error indicator is returned
export async function getMessageByID(messages: Collection, id: string) {
    // Since id's are auto-generated and unique we chose to use find instead of findOne() 
    let result: any = null;
    let err: boolean = true;
    let mongoID = new ObjectID(id)
    const filter = { _id: mongoID };

    let cursor = messages.find(filter)
    if (await cursor.hasNext()) {
        result = await cursor.next()
        err = false;
    }

    let message: Message;
    if (result == null) {
        message = initializeDummyMessage();
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
    const sortFilter = { createdAt: -1 };

    if (channelID == null) {
        err = true;
        return { last100messages, err };
    }
    channelID = channelID.toString();
    let cursor
    if (!messageID) {
        const findFilter = { channelID: channelID };
        cursor = messages.find(findFilter).sort(sortFilter).limit(100);
    } else {
        let objID = new ObjectID(messageID)
        const findFilter = { channelID: channelID, _id: { $lt: objID } };
        cursor = messages.find(findFilter).sort(sortFilter).limit(100);
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