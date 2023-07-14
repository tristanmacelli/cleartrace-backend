"use strict";

import { ObjectId, Collection, MongoClient, Sort, Db, Document, Filter } from "mongodb";
import { Channel, isChannelMember } from "./channel";
import { Message } from "./message";

const mongoContainerName = 'userMessageStore'
const dbName = 'userMessageDB';
const mongoURL = 'mongodb://' + mongoContainerName + ':27017/' + dbName;

// Create a new MongoClient
const mc = new MongoClient(mongoURL);

export const createConnection = async (): Promise<Db> => {
    let retryInterval: number = 1;

    const client = await recursiveCreateConnection(retryInterval);
    if (!client) {
      throw new Error("Could not connect to the database. Goodbye.")
    }
    return client.db(dbName);
};

const recursiveCreateConnection = async (retryInterval: number): Promise<void | MongoClient> => {
  // 7 retries over the course of ~2 minutes
  if (retryInterval > 64) {
    return;
  }
  
  try {
    const client = await mc.connect();
    return client;
  } catch(e) {
    console.log(`mongo_handlers.ts createConnection ${e}`);
    console.log("Cannot connect to the database: MongoNetworkError: failed to connect to server");
    console.log(`Retrying in ${retryInterval} second(s)`);
    sleep(retryInterval);  
    recursiveCreateConnection(retryInterval * 2)
  }
};

// getChannels gets some (if a search term is passed) or all of the users channels 
export const getChannels = async (channels: Collection, userID: number, search: string) => {
    const searchTerm = search ? { name: { $regex: "/^"+ search +"/i" } } : {}; 
    const cursor = channels.find(searchTerm);

    if (!(await cursor.hasNext())) {    
        return { usersChannels: [], err: true }
    }

    const results = await cursor.toArray()
    const allChannels = results.map((channel) => {
        return new Channel(channel._id.toString(), channel.name, channel.description, channel.private,
            channel.members, channel.createdAt, channel.creator, channel.editedAt);
    });
    const usersChannels = allChannels.filter((channel) => {
        return isChannelMember(channel, userID)
    })
    return { usersChannels, err: false };
};

// insertNewChannel takes in a new Channel and inserts it into the messaging DB
export const insertNewChannel = async (channels: Collection, newChannel: Channel) => {
    // Duplicate Check
    const filter = { name: newChannel.name }
    if (await channels.find(filter).hasNext()) {
        return { newChannel, hasDuplicates: true, err: true }
    }

    const insertDoc = {
        name: newChannel.name, 
        description: newChannel.description,
        private: newChannel.private, 
        members: newChannel.members,
        createdAt: new Date(), 
        creator: newChannel.creator,
        editedAt: newChannel.editedAt
    }
    const insertResult = await channels.insertOne(insertDoc)
      .catch((reason) => {
        console.log(`mongo_handlers.ts insertNewChannel ${reason}`)
      });

    if (!insertResult) {
        return { newChannel, hasDuplicates: false, err: true }
    }
    newChannel.id = insertResult.insertedId.toString()

    return { newChannel, hasDuplicates: false, err: false }
};

// insertNewMessage takes in a new Message and inserts it into the messaging DB
export const insertNewMessage = async (messages: Collection, newMessage: Message) => {
    const insertDoc = {
        channelID: newMessage.channelID, 
        createdAt: new Date(),
        body: newMessage.body, 
        creator: newMessage.creator,
        editedAt: newMessage.editedAt
    }

    const insertResult = await messages.insertOne(insertDoc)
      .catch((reason) => {
        console.log(`mongo_handlers.ts insertNewMessage ${reason}`)
      });

    if (!insertResult) {
        return { newMessage, err: true }
    }
    newMessage.id = insertResult.insertedId.toString()

    return { newMessage, err: false };
};

// addChannelMembers takes an existing Channel and adds members using a req (request) object
export const addChannelMember = async (channels: Collection, existingChannel: Channel, userId: number): Promise<boolean> => {
  // Add the specified member to this channel's list of members
  existingChannel.members.push(userId);
  return await updateChannelMembers(channels, existingChannel.id, existingChannel.members)
};

// removeChannelMember takes an existing Channel and removes members using a req (request) object
export const removeChannelMember = async (channels: Collection, existingChannel: Channel, userId: number): Promise<boolean> => {
    // Remove the specified member from this channel's list of members
    existingChannel.members.splice(userId, 1);
    return updateChannelMembers(channels, existingChannel.id, existingChannel.members)
};
 const updateChannelMembers = async (channels: Collection, channelID: string, channelMembers: number[]): Promise<boolean> => {
  const newEditedAt = new Date()
  const filter = { _id: new ObjectId(channelID) };
  const updateDoc = {
    $set: { 
      members: channelMembers, 
      editedAt: newEditedAt
    },
  };

  const updateResult = await channels.updateOne(filter, updateDoc)
    .catch((reason) => {
      console.log(`mongo_handlers.ts updateChannelMembers ${reason}`)
    })
  return !updateResult;
};

// updatedChannel updates name and body of an existing Channel using a req (request) object
export const updateChannel = async (
    channels: Collection,
    existingChannel: Channel,
    updates: { name: string, description: string }
) => {
  const updatedChannel: Channel = {
    ...existingChannel,
    editedAt: new Date(),
    name: updates.name,
    description: updates.description
  }

  const filter = { _id: new ObjectId(updatedChannel.id) };
  const updateDoc = {
    $set: { 
      name: updatedChannel.name, 
      description: updatedChannel.description,
      editedAt: updatedChannel.editedAt
    },
  };

  const updateResult = await channels.updateOne(filter, updateDoc)
    .catch((reason) => {
      console.log(`mongo_handlers.ts updateChannel() ${reason}`)
    })

  return { updatedChannel, err: !updateResult };
};

// updateMessage takes an existing Message and a request with updates to apply to the Message's body 
export const updateMessage = async (messages: Collection, existingMessage: Message, updates: { body: string }) => {
    const updatedMessage: Message = {
        ...existingMessage,
        body: updates.body,
        editedAt: new Date()
    }
    const filter = { messageID: updatedMessage.id };
    const updateDoc = {
      $set: { 
        body: updatedMessage.body,
        editedAt: updatedMessage.editedAt
      },
    };

    const updateResult = await messages.updateOne(filter, updateDoc)
      .catch((reason) => {
        console.log(`mongo_handlers.ts updateMessage() ${reason}`)
      })

    return { updatedMessage, err: !updateResult };
};

// deleteChannel deletes a single channel & its associated messages
export const deleteChannel = async (channels: Collection, messages: Collection, existingChannel: Channel): Promise<boolean> => {
    // The general channel never gets deleted
    if (existingChannel.creator.ID == -1) {
      return true;
    }

    const channelFilter = { _id: new ObjectId(existingChannel.id) }
    const channelDeletionResult = await channels.deleteOne(channelFilter)
      .catch((reason) => {
        console.log(`mongo_handlers.ts deleteChannel() ${reason}`)
      });

    const messageFilter = { channelID: existingChannel.id }
    const messagesDeletionResult = await messages.deleteMany(messageFilter)
      .catch((reason) => {
        console.log(`mongo_handlers.ts deleteChannel() ${reason}`)
      });

    // !(void) === true; void is returned from mongo driver functions when a callback const is =  execute=> d
    const err: boolean = !channelDeletionResult && !messagesDeletionResult 
    return err;
};

// deleteMessage deletes a single message
export const deleteMessage = async (messages: Collection, existingMessage: Message): Promise<boolean> => {
  const filter = { messageID: existingMessage.id }
  const deleteResult = await messages.deleteOne(filter)
    .catch((reason) => {
      console.log(`mongo_handlers.ts deleteMessage() ${reason}`)
    });

 return !deleteResult;
};

// getChannelByID returns the channel associated with the provided id value. If there no channel 
// is no channel associated with the provided id then an error indicator is returned
export const getChannelByID = async (channels: Collection, id: string): Promise<{ channel?: Channel, err: boolean }> => {
    const filter: Filter<Document> = { _id: new ObjectId(id) };
    // Since id's are auto-generated and unique we chose to use findOne() instead of find()
    const result = await channels.findOne(filter)
      .catch((reason) => {
        console.log(`mongo_handlers.ts getChannelByID ${reason}`)
      })

    if (!result) {
        return { err: true };
    }

    const channel = new Channel(result._id.toString(), result.name, result.description, result.private,
        result.members, result.createdAt, result.creator, result.editedAt);
    return { channel, err: false };
};

// getMessageByID returns the message associated with the provided id value. If there no message 
// is no message associated with the provided id then an error indicator is returned
export const getMessageByID = async (messages: Collection, id: string): Promise<{ message?: Message, err: boolean }>  => {
    const filter = { _id: new ObjectId(id) };
    // Since id's are auto-generated and unique we chose to use findOne() instead of find()
    const result = await messages.findOne(filter)
      .catch((reason) => {
        console.log(`mongo_handlers.ts getMessageByID ${reason}`)
      })

    if (!result) {
        return { err: true };
    }

    const message = new Message(result._id.toString(), result.channelID, result.createdAt, result.body,
        result.creator, result.editedAt)
    return { message, err: false }
};

// TODO: Reshape the return value of find to a JSON array of message model objects
// last100Messages gets the most recent 100 messages of channel
export const last100Messages = async (messages: Collection, channelID: string, messageID: string) => {
    const sortFilter: Sort = { createdAt: -1 };

    if (channelID == null) {
        return { last100messages: [], err: true };
    }

    const findFilter = messageID ? { channelID: channelID, _id: { $lt: new ObjectId(messageID) } } : { channelID: channelID };
    const cursor = messages.find(findFilter).sort(sortFilter).limit(100);

    if (!await cursor.hasNext()) {
        return { last100messages: [], err: false };
    }

    const results = await cursor.toArray()
    const last100messages: Message[] = results.map((message) => {
        return new Message(message._id.toString(), message.channelID, message.createdAt, message.body,
        message.creator, message.editedAt)
    });

    return { last100messages, err: false };
};

export const sleep = async (seconds: number): Promise<void> => {
  return new Promise((resolve) => setTimeout(resolve, seconds * 1000));
};

export * from "./mongo_handlers";