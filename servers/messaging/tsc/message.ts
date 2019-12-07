"use strict";

import { User } from "./user";

export class Message {
    id: string;
    channelID: string;
    createdAt: Date;
    body: string;
    creator: User;
    editedAt: Date;
    constructor(ChannelID: string, CreatedAt: Date, Body: string, Creator: User, EditedAt: Date) {
        this.id = "";
        this.channelID = ChannelID;
        this.createdAt = CreatedAt;
        this.body = Body;
        this.creator = Creator;
        this.editedAt = EditedAt;
    }
}

// export default Message;

// to compile run tsc --outDir ../