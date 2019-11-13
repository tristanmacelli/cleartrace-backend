"use strict";

class Message {
    constructor() {
        this.ID = null;
        this.ChannelID = null;
        this.CreatedAt =  null;
        this.Body = null;
        this.Creator = null;
        this.EditedAt = null;
    }

    constructor(ID, ChannelID, CreatedAt, Body, Creator, EditedAt) {
        this.ID = ID;
        this.ChannelID = ChannelID;
        this.CreatedAt =  CreatedAt;
        this.Body = Body;
        this.Creator = Creator;
        this.EditedAt = EditedAt;
    }    
  }