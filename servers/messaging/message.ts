"use strict";

export class Message {
    // constructor() {
    //     this._id = null;
    //     this.ChannelID = null;
    //     this.CreatedAt = null;
    //     this.Body = null;
    //     this.Creator = null;
    //     this.EditedAt = null;
    // }

    _id : string;
    channelID :string;
    createdAt :string;
    body :string;
    creator : string;
    editedAt : string;
    constructor(ChannelID:string, CreatedAt:string, Body:string, Creator:string, EditedAt:string) {
        this._id = "";
        this.channelID = ChannelID;
        this.createdAt = CreatedAt;
        this.body = Body;
        this.creator = Creator;
        this.editedAt = EditedAt;
    }
}

// export default Message;