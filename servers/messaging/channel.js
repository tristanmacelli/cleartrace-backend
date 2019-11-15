"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
var Channel = /** @class */ (function () {
    function Channel(Name, Description, Private, Members, CreatedAt, Creator, EditedAt) {
        this._id = "";
        this.name = Name;
        this.description = Description;
        this.private = Private;
        this.members = Members;
        this.createdAt = CreatedAt;
        this.creator = Creator;
        this.editedAt = EditedAt;
    }
    return Channel;
}());
exports.Channel = Channel;
