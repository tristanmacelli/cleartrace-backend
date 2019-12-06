"use strict";

export class User {
    ID: number;
    Email: string;
    PassHash: Uint8Array;
    UserName: string;
    FirstName: string;
    LastName: string;
    PhotoURL: string;

    constructor(ID: number, Email: string, PassHash: Uint8Array, UserName: string, FirstName: string, LastName: string, PhotoURL: string) {
        this.ID = ID;
        this.Email = Email;
        this.PassHash = PassHash;
        this.UserName = UserName;
        this.FirstName = FirstName;
        this.LastName = LastName;
        this.PhotoURL = PhotoURL;
    }
}