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
    editedAt: string;
    constructor(Name: string, Description: string, Private: boolean, Members: number[], CreatedAt: Date, Creator: User, EditedAt: string) {
        this.id = "";
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