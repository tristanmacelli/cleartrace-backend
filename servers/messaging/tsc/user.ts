"use strict";

export class User {
    id: number;
    email: string;
    passHash: Uint8Array;
    userName: string;
    firstName: string;
    lastName: string;
    photoURL: string;

    constructor(id: number, email: string, passHash: Uint8Array, userName: string, firstName: string, lastName: string, photoURL: string) {
        this.id = id;
        this.email = email;
        this.passHash = passHash;
        this.userName = userName;
        this.firstName = firstName;
        this.lastName = lastName;
        this.photoURL = photoURL;
    }
}