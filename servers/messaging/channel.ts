"use strict";

class Channel {
    // constructor() {
    //     this._id = null;
    //     this.Name = null;
    //     this.Description = null;
    //     this.Private = null;
    //     this.Members = null;
    //     this.CreatedAt = null;
    //     this.Creator = null;
    //     this.EditedAt = null;
    // }
    _id = null;
    Name :string;
    Description : string;
    Private : boolean;
    Members : number[];
    CreatedAt : string;
    Creator : number;
    EditedAt : string;
    constructor(Name, Description, Private, Members, CreatedAt, Creator, EditedAt) {
        this._id = null;
        this.Name = Name;
        this.Description = Description;
        this.Private = Private;
        this.Members = Members;
        this.CreatedAt = CreatedAt;
        this.Creator = Creator;
        this.EditedAt = EditedAt;
    }
}

module.exports=Channel