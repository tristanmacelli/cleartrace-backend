"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
var Message = /** @class */ (function () {
    function Message(ChannelID, CreatedAt, Body, Creator, EditedAt) {
        this._id = "";
        this.channelID = ChannelID;
        this.createdAt = CreatedAt;
        this.body = Body;
        this.creator = Creator;
        this.editedAt = EditedAt;
    }
    return Message;
}());
exports.Message = Message;
