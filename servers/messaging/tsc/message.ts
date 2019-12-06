"use strict";

import { User } from "./user";

export class Message {
    _id: string;
    channelID: string;
    createdAt: string;
    body: string;
    creator: User;
    editedAt: string;
    constructor(ChannelID: string, CreatedAt: string, Body: string, Creator: User, EditedAt: string) {
        this._id = "";
        this.channelID = ChannelID;
        this.createdAt = CreatedAt;
        this.body = Body;
        this.creator = Creator;
        this.editedAt = EditedAt;
    }
}

// export default Message;

// to compile run tsc --outDir ../