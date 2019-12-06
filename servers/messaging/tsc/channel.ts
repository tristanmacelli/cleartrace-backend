"use strict";

import { User } from "./user";

export class Channel {
    _id: string;
    name: string;
    description: string;
    private: boolean;
    members: string[];
    createdAt: string;
    creator: User;
    editedAt: string;
    constructor(Name: string, Description: string, Private: boolean, Members: string[], CreatedAt: string, Creator: User, EditedAt: string) {
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

// to compile run tsc --outDir ../