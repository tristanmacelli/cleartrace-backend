"use strict";

export class Channel {
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
    _id : string;
    name :string;
    description : string;
    private : boolean;
    members : number[];
    createdAt : string;
    creator : number;
    editedAt : string;
    constructor(Name:string, Description:string, Private:boolean, Members:number[], CreatedAt:string, Creator:number, EditedAt:string) {
        this._id = "";
        this.name = Name;
        this.description = Description;
        this.private = Private;
        this.members = Members;
        this.createdAt = CreatedAt;
        this.creator = Creator;
        this.editedAt = EditedAt;
    }
}

// export default Channel;