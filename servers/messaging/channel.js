"use strict";

class Channel {
    constructor() {
        this.ID = null;
        this.Name = null;
        this.Description =  null;
        this.Private = null;
        this.Members = null;
        this.CreatedAt = null;
        this.Creator = null;
        this.EditedAt = null;
    }

    constructor(ID, Name, Description, Private, Members, CreatedAt, Creator, EditedAt) {
        this.ID = ID;
        this.Name = Name;
        this.Description =  Description;
        this.Private = Private;
        this.Members = Members;
        this.CreatedAt = CreatedAt;
        this.Creator = Creator;
        this.EditedAt = EditedAt;
    }    
  }