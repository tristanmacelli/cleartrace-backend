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
    constructor(id: string, Name: string, Description: string, Private: boolean, Members: number[], CreatedAt: Date, Creator: User, EditedAt: Date) {
        this.id = id
        this.name = Name;
        this.description = Description;
        this.private = Private;
        this.members = Members;
        this.createdAt = CreatedAt;
        this.creator = Creator;
        this.editedAt = EditedAt;
    }
}

export function isChannelCreator(channel: Channel, userID: number): boolean {
    return channel.creator.ID === userID;
}

export function isChannelMember(channel: Channel, userID: number): boolean {
    if (channel.private) {
        for (let i = 0; i < channel.members.length; i++) {
            if (channel.members[i] === userID) {
                return true;
            }
        }
    } else {
        return true;
    }
    return false;
}

export function initializeDummyChannel(): Channel {
  let emptyUser = new User(-1, "", "", "", "", "")
  let dummyDate = new Date()
  let dummyChannel = new Channel("", "", "", false, [], dummyDate, emptyUser, dummyDate);
  return dummyChannel;
}

// export default Channel;

// to compile run tsc --outDir ../