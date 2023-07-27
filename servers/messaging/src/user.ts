"use strict";

export class User {
  ID: number;
  email: string;
  userName: string;
  firstName: string;
  lastName: string;
  photoURL: string;

  constructor(ID: number, email: string, userName: string, firstName: string, lastName: string, photoURL: string) {
    this.ID = ID;
    this.email = email;
    this.userName = userName;
    this.firstName = firstName;
    this.lastName = lastName;
    this.photoURL = photoURL;
  }
}
