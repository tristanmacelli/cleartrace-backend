package sessions

import (
	"bytes"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	b64 "encoding/base64"
	"errors"
	"fmt"
)

//InvalidSessionID represents an empty, invalid session ID
const InvalidSessionID SessionID = ""

//idLength is the length of the ID portion
const idLength = 32

//signedLength is the full length of the signed session ID
//(ID portion plus signature)
const signedLength = idLength + sha256.Size

//SessionID represents a valid, digitally-signed session ID.
//This is a base64 URL encoded string created from a byte slice
//where the first `idLength` bytes are crytographically random
//bytes representing the unique session ID, and the remaining bytes
//are an HMAC hash of those ID bytes (i.e., a digital signature).
//The byte slice layout is like so:
//+-----------------------------------------------------+
//|...32 crypto random bytes...|HMAC hash of those bytes|
//+-----------------------------------------------------+
type SessionID string

//ErrInvalidID is returned when an invalid session id is passed to ValidateID()
var ErrInvalidID = errors.New("Invalid Session ID")

//NewSessionID creates and returns a new digitally-signed session ID,
//using `signingKey` as the HMAC signing key. An error is returned only
//if there was an error generating random bytes for the session ID
func NewSessionID(signingKey string) (SessionID, error) {
	//TODO: if `signingKey` is zero-length, return InvalidSessionID
	//and an error indicating that it may not be empty
	if signingKey.length == 0 {
		// TODO RETURN AN ERROR HERE
		return InvalidSessionID
	}
	//TODO: Generate a new digitally-signed SessionID by doing the following:
	//- create a byte slice where the first `idLength` of bytes
	//  are cryptographically random bytes for the new session ID,
	//  and the remaining bytes are an HMAC hash of those ID bytes,
	//  using the provided `signingKey` as the HMAC key.
	//- encode that byte slice using base64 URL Encoding and return
	//  the result as a SessionID type

	s_id := make([]byte, idLength)
	_, err := rand.Read(s_id)
	if err != nil {
		fmt.Println("error:", err)
		return InvalidSessionID
	}
	// first half of id
	fmt.Println(s_id)

	key := []byte(signingKey)
	//create a new HMAC hasher
	h := hmac.New(sha256.New, key)
	//write the message into it
	h.Write(s_id)
	//calculate the HMAC signature
	signature := h.Sum(nil)

	// final session id in bytes
	s_id = append(s_id, signature...)

	// encode using Base64 URL encoding
	s_id_encoded := b64.StdEncoding.EncodeToString([]byte(s_id))
	fmt.Println(s_id_encoded)

	return s_id_encoded, nil
}

//ValidateID validates the string in the `id` parameter
//using the `signingKey` as the HMAC signing key
//and returns an error if invalid, or a SessionID if valid
func ValidateID(id string, signingKey string) (SessionID, error) {

	//TODO: validate the `id` parameter using the provided `signingKey`.
	//base64 decode the `id` parameter, HMAC hash the
	//ID portion of the byte slice, and compare that to the
	//HMAC hash stored in the remaining bytes. If they match,
	//return the entire `id` parameter as a SessionID type.
	//If not, return InvalidSessionID and ErrInvalidID.

	// decode the id parameter
	s_id_Dec, _ := b64.StdEncoding.DecodeString(id)
	// get the first half of the slice
	first_half := s_id_Dec[0:idLength]
	fmt.Println(first_half)

	// HMAC hash the first half
	k := []byte(signingKey)
	he := hmac.New(sha256.New, k)
	he.Write(first_half)
	signature := he.Sum(nil)

	// compare this to the second half
	res := bytes.Compare(signature, s_id_Dec[idLength:])

	if res == 0 {
		return id, nil
	} else {
		return InvalidSessionID, ErrInvalidID
	}
}

//String returns a string representation of the sessionID
func (sid SessionID) String() string {
	return string(sid)
}
