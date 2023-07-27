"use strict";

import { User } from "./user";

export class Channel {
  id: string;
  name: string;
  description: string;
  private: boolean;
  members: number[];
  createdAt: Date;
  creator: User;
  editedAt: Date;
  constructor(
    id: string,
    Name: string,
    Description: string,
    Private: boolean = true,
    Members: number[],
    CreatedAt: Date = new Date(),
    Creator: User,
    EditedAt: Date,
  ) {
    this.id = id;
    this.name = Name;
    this.description = Description;
    this.private = Private;
    this.members = Members;
    this.createdAt = CreatedAt;
    this.creator = Creator;
    this.editedAt = EditedAt;
  }

  isChannelCreator = (userID: number): boolean => {
    return this.creator.ID === userID;
  };

  isChannelMember = (userID: number): boolean => {
    if (this.private) {
      for (let i = 0; i < this.members.length; i++) {
        if (this.members[i] === userID) {
          return true;
        }
      }
    } else {
      return true;
    }
    return false;
  };
}
