"use strict";

import { User } from "./user";

export class Message {
  id: string;
  channelID: string;
  createdAt: Date;
  body: string;
  creator: User;
  editedAt: Date;
  constructor(id: string, ChannelID: string, CreatedAt: Date, Body: string, Creator: User, EditedAt: Date) {
    this.id = id;
    this.channelID = ChannelID;
    this.createdAt = CreatedAt;
    this.body = Body;
    this.creator = Creator;
    this.editedAt = EditedAt;
  }

  isMessageCreator = (userID: number): boolean => {
    return this.creator.ID === userID;
  };
}
