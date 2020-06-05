"use strict";

export class User {
    ID: number;
    email: string;
    passHash: Uint8Array;
    userName: string;
    firstName: string;
    lastName: string;
    photoURL: string;

    constructor(ID: number, email: string, passHash: Uint8Array, userName: string, firstName: string, lastName: string, photoURL: string) {
        this.ID = ID;
        this.email = email;
        this.passHash = passHash;
        this.userName = userName;
        this.firstName = firstName;
        this.lastName = lastName;
        this.photoURL = photoURL;
    }
}