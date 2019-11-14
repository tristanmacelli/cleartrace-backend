"use strict";

class Message {
    constructor() {
        this._id = null;
        this.ChannelID = null;
        this.CreatedAt = null;
        this.Body = null;
        this.Creator = null;
        this.EditedAt = null;
    }

    constructor(ChannelID, CreatedAt, Body, Creator, EditedAt) {
        this._id = null;
        this.ChannelID = ChannelID;
        this.CreatedAt = CreatedAt;
        this.Body = Body;
        this.Creator = Creator;
        this.EditedAt = EditedAt;
    }
}