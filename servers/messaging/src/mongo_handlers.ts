"use strict";

import { ObjectId, Collection, MongoClient, Sort, Db, Document, Filter } from "mongodb";
import { Channel } from "./channel";
import { Message } from "./message";

const mongoContainerName = "userMessageStore";
const dbName = "userMessageDB";
const mongoURL = "mongodb://" + mongoContainerName + ":27017/" + dbName;

// Create a new MongoClient
const mc = new MongoClient(mongoURL);

export const createConnection = async (): Promise<Db> => {
  const retryInterval: number = 1;

  const client = await recursiveCreateConnection(retryInterval);
  if (!client) {
    throw new Error("Could not connect to the database. Goodbye.");
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
  } catch (e) {
    console.error(`mongo_handlers createConnection ${e}`);
    console.error("Cannot connect to the database: MongoNetworkError: failed to connect to server");
    console.info(`Retrying in ${retryInterval} second(s)`);
    await sleep(retryInterval);
    recursiveCreateConnection(retryInterval * 2);
  }
};

// getChannels gets some (if a search term is passed) or all of the users channels
export const getChannels = async (
  channels: Collection,
  userID: number,
  search: string,
): Promise<{
  usersChannels: Channel[];
  err: Error | null;
}> => {
  const searchTerm = search ? { name: { $regex: "/^" + search + "/i" } } : {};
  const cursor = channels.find(searchTerm);

  if (!(await cursor.hasNext())) {
    return { usersChannels: [], err: new Error("No channels found") };
  }

  const results = await cursor.toArray();
  const usersChannels = results
    .map((channel) => {
      return new Channel(
        channel._id.toString(),
        channel.name,
        channel.description,
        channel.private,
        channel.members,
        channel.createdAt,
        channel.creator,
        channel.editedAt,
      );
    })
    .filter((channel) => {
      return channel.isChannelMember(userID);
    });
  return { usersChannels, err: null };
};

// insertNewChannel takes in a new Channel and inserts it into the messaging DB
export const insertNewChannel = async (
  channels: Collection,
  newChannel: Channel,
): Promise<{
  newChannel: Channel;
  err: Error | null;
}> => {
  // Duplicate Check
  const filter = { name: newChannel.name };
  if (await channels.find(filter).hasNext()) {
    return {
      newChannel,
      err: new Error("duplicate: a channel with the provided name already exists "),
    };
  }
  // const newGroupHash = hash(newChannel.members.toString());
  // const filter = { hash: newGroupHash };
  // if (await groupMembersHash.find(filter).hasNext()) {
  //   return { newChannel, err: new Error("duplicate: a channel with the provided members already exists ") }
  // }

  const insertDoc = {
    name: newChannel.name,
    description: newChannel.description,
    private: newChannel.private,
    members: newChannel.members,
    createdAt: new Date(),
    creator: newChannel.creator,
    editedAt: newChannel.editedAt,
  };
  const insertResult = await channels
    .insertOne(insertDoc)
    .then((res) => res)
    .catch((reason) => {
      console.error(`mongo_handlers insertNewChannel ${reason}`);
      return new Error(reason);
    });

  if (insertResult instanceof Error) {
    return { newChannel, err: insertResult };
  }
  newChannel.id = insertResult.insertedId.toString();
  // const hashInsertResult = await groupMembersHash.insertOne({ hash: newGroupHash })
  //   .then(() => null)
  //   .catch((reason) => {
  //     console.error(`mongo_handlers insertNewChannel ${reason}`)
  //     return new Error(reason)
  //   })

  return { newChannel, err: null };
};

// insertNewMessage takes in a new Message and inserts it into the messaging DB
export const insertNewMessage = async (
  messages: Collection,
  newMessage: Message,
): Promise<{
  newMessage: Message;
  err: Error | null;
}> => {
  const insertDoc = {
    channelID: newMessage.channelID,
    createdAt: new Date(),
    body: newMessage.body,
    creator: newMessage.creator,
    editedAt: newMessage.editedAt,
  };

  const insertResult = await messages
    .insertOne(insertDoc)
    .then((res) => res)
    .catch((reason) => {
      console.error(`mongo_handlers insertNewMessage ${reason}`);
      return new Error(reason);
    });

  if (insertResult instanceof Error) {
    return { newMessage, err: insertResult };
  }
  newMessage.id = insertResult.insertedId.toString();

  return { newMessage, err: null };
};

// addChannelMembers takes an existing Channel and adds members using a req (request) object
export const addChannelMember = async (channels: Collection, existingChannel: Channel, userId: number) => {
  // Add the specified member to this channel's list of members
  existingChannel.members.push(userId);
  return await updateChannelMembers(channels, existingChannel, existingChannel.members);
};

// removeChannelMember takes an existing Channel and removes members using a req (request) object
export const removeChannelMember = async (channels: Collection, existingChannel: Channel, userId: number) => {
  // Remove the specified member from this channel's list of members
  existingChannel.members.splice(userId, 1);
  return updateChannelMembers(channels, existingChannel, existingChannel.members);
};

const updateChannelMembers = async (channels: Collection, existingChannel: Channel, newMembersList: number[]) => {
  const updatedChannel: Channel = {
    ...existingChannel,
    members: newMembersList,
    editedAt: new Date(),
  };
  const filter = { _id: new ObjectId(updatedChannel.id) };
  const updateDoc = {
    $set: {
      members: updatedChannel.members,
      editedAt: updatedChannel.editedAt,
    },
  };

  const updateError = await channels
    .updateOne(filter, updateDoc)
    .then(() => null)
    .catch((reason) => {
      console.error(`mongo_handlers updateChannelMembers ${reason}`);
      return new Error(reason);
    });
  return { updatedChannel, err: updateError };
};

// updatedChannel updates name and body of an existing Channel using a req (request) object
export const updateChannel = async (
  channels: Collection,
  existingChannel: Channel,
  updates: { name: string; description: string },
) => {
  const updatedChannel: Channel = {
    ...existingChannel,
    editedAt: new Date(),
    name: updates.name,
    description: updates.description,
  };

  const filter = { _id: new ObjectId(updatedChannel.id) };
  const updateDoc = {
    $set: {
      name: updatedChannel.name,
      description: updatedChannel.description,
      editedAt: updatedChannel.editedAt,
    },
  };

  const updateError = await channels
    .updateOne(filter, updateDoc)
    .then(() => null)
    .catch((reason) => {
      console.error(`mongo_handlers updateChannel ${reason}`);
      return new Error(reason);
    });

  return { updatedChannel, err: updateError };
};

// updateMessage takes an existing Message and a request with updates to apply to the Message's body
export const updateMessage = async (messages: Collection, existingMessage: Message, updates: { body: string }) => {
  const updatedMessage: Message = {
    ...existingMessage,
    body: updates.body,
    editedAt: new Date(),
  };
  const filter = { messageID: updatedMessage.id };
  const updateDoc = {
    $set: {
      body: updatedMessage.body,
      editedAt: updatedMessage.editedAt,
    },
  };

  const updateError = await messages
    .updateOne(filter, updateDoc)
    .then(() => null)
    .catch((reason) => {
      console.error(`mongo_handlers updateMessage ${reason}`);
      return new Error(reason);
    });

  return { updatedMessage, err: updateError };
};

// deleteChannel deletes a single channel & its associated messages
export const deleteChannel = async (
  channels: Collection,
  messages: Collection,
  existingChannel: Channel,
): Promise<Error | null> => {
  // The general channel never gets deleted
  if (existingChannel.creator.ID == -1) {
    return new Error("channel not found");
  }

  const channelFilter = { _id: new ObjectId(existingChannel.id) };
  const channelDeletionError = await channels
    .deleteOne(channelFilter)
    .then(() => null)
    .catch((reason) => {
      console.error(`mongo_handlers deleteChannel ${reason}`);
      return new Error(reason);
    });
  if (channelDeletionError) return channelDeletionError;

  const messageFilter = { channelID: existingChannel.id };
  const messagesDeletionError = await messages
    .deleteMany(messageFilter)
    .then(() => null)
    .catch((reason) => {
      console.error(`mongo_handlers deleteChannel ${reason}`);
      return new Error(reason);
    });
  if (messagesDeletionError) return messagesDeletionError;

  // !(void) === true; void is returned from mongo driver functions when a callback const is =  execute=> d
  return null;
};

// deleteMessage deletes a single message
export const deleteMessage = async (messages: Collection, existingMessage: Message): Promise<Error | null> => {
  const filter = { messageID: existingMessage.id };
  const deleteError = await messages
    .deleteOne(filter)
    .then(() => null)
    .catch((reason) => {
      console.error(`mongo_handlers deleteMessage ${reason}`);
      return new Error(reason);
    });

  return deleteError;
};

// getChannelByID returns the channel associated with the provided id value. If there no channel
// is no channel associated with the provided id then an error indicator is returned
export const getChannelByID = async (
  channels: Collection,
  id: string,
): Promise<{ channel?: Channel; err: Error | null }> => {
  const filter: Filter<Document> = { _id: new ObjectId(id) };
  // Since id's are auto-generated and unique we chose to use findOne() instead of find()
  const result = await channels.findOne(filter).catch((reason) => {
    console.error(`mongo_handlers getChannelByID ${reason}`);
    return new Error(reason);
  });

  if (!result) {
    return { err: new Error("Not found") };
  }
  if (result instanceof Error) {
    return { err: result };
  }

  const channel = new Channel(
    result._id.toString(),
    result.name,
    result.description,
    result.private,
    result.members,
    result.createdAt,
    result.creator,
    result.editedAt,
  );
  return { channel, err: null };
};

// getMessageByID returns the message associated with the provided id value. If there no message
// is no message associated with the provided id then an error indicator is returned
export const getMessageByID = async (
  messages: Collection,
  id: string,
): Promise<{ message?: Message; err: Error | null }> => {
  const filter = { _id: new ObjectId(id) };
  // Since id's are auto-generated and unique we chose to use findOne() instead of find()
  const result = await messages.findOne(filter).catch((reason) => {
    console.error(`mongo_handlers getMessageByID ${reason}`);
    return new Error(reason);
  });

  if (!result) {
    return { err: new Error("Not found") };
  }
  if (result instanceof Error) {
    return { err: result };
  }

  const message = new Message(
    result._id.toString(),
    result.channelID,
    result.createdAt,
    result.body,
    result.creator,
    result.editedAt,
  );
  return { message, err: null };
};

// TODO: Reshape the return value of find to a JSON array of message model objects
// last100Messages gets the most recent 100 messages of channel
export const last100Messages = async (
  messages: Collection,
  channelID: string,
  messageID: string,
): Promise<{
  last100messages: Message[];
  err: Error | null;
}> => {
  const sortFilter: Sort = { createdAt: -1 };

  if (channelID == null) {
    return { last100messages: [], err: new Error("Unknown channelID") };
  }

  const findFilter = messageID
    ? { channelID: channelID, _id: { $lt: new ObjectId(messageID) } }
    : { channelID: channelID };
  const cursor = messages.find(findFilter).sort(sortFilter).limit(100);

  if (!(await cursor.hasNext())) {
    return { last100messages: [], err: null };
  }

  const results = await cursor.toArray();
  const last100messages: Message[] = results.map((message) => {
    return new Message(
      message._id.toString(),
      message.channelID,
      message.createdAt,
      message.body,
      message.creator,
      message.editedAt,
    );
  });

  return { last100messages, err: null };
};

export const sleep = async (seconds: number): Promise<void> => {
  return new Promise((resolve) => setTimeout(resolve, seconds * 1000));
};

// hash creates a deterministic integer value
// const hash = (inputString: string) => {
//   if (inputString.length === 0) return 0;
//   let hash = 0;
//   let char;

//   for (let i = 0; i < inputString.length; i++) {
//     char = inputString.charCodeAt(i);
//     hash = ((hash << 5) - hash) + char;
//     hash |= 0; // Convert to 32bit integer
//   }
//   return hash;
// };

export * from "./mongo_handlers";
