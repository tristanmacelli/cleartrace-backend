package sessions

import (
	"encoding/json"

	"time"

	"github.com/go-redis/redis"
)

//RedisStore represents a session.Store backed by redis.
type RedisStore struct {
	//Redis client used to talk to redis server.
	Client *redis.Client
	//Used for key expiry time on redis.
	SessionDuration time.Duration
}

//NewRedisStore constructs a new RedisStore
func NewRedisStore(client *redis.Client, sessionDuration time.Duration) *RedisStore {
	//initialize and return a new RedisStore struct
	// keets track of client and session duration
	return &RedisStore{
		Client:          client,
		SessionDuration: sessionDuration,
	}
}

//Store implementation

//Save saves the provided `sessionState` and associated SessionID to the store.
//The `sessionState` parameter is typically a pointer to a struct containing
//all the data you want to associated with the given SessionID.
func (rs *RedisStore) Save(sid SessionID, sessionState interface{}) error {
	//TODO: marshal the `sessionState` to JSON and save it in the redis database,
	//using `sid.getRedisKey()` for the key.
	//return any errors that occur along the way.

	// marshal the json and resolve errors
	j, err := json.Marshal(&sessionState)
	if nil != err {
		// json can not be marshalled
		return err
	}

	// get redis key
	redisKey := sid.getRedisKey()

	// save state to database
	err = rs.Client.Set(redisKey, j, 0).Err()
	if err != nil {
		return err
	}

	// no error return
	return nil
}

//Get populates `sessionState` with the data previously saved
//for the given SessionID
func (rs *RedisStore) Get(sid SessionID, sessionState interface{}) error {
	//TODO: get the previously-saved session state data from redis,
	//unmarshal it back into the `sessionState` parameter
	//and reset the expiry time, so that it doesn't get deleted until
	//the SessionDuration has elapsed.

	// get redis frmatted key
	redisKey := sid.getRedisKey()

	// get the state information from redis
	marshalledJSON, err := rs.Client.Get(redisKey).Result()
	if err != nil {
		// state does not exist for the key
		return ErrStateNotFound
	}

	// unmarshal the fetched state and put it into session state using pointer
	err = json.Unmarshal([]byte(marshalledJSON), &sessionState)
	if err != nil {
		return err
	}

	// reset the expiry time
	rs.SessionDuration = time.Hour
	err = rs.Client.Set(redisKey, marshalledJSON, rs.SessionDuration).Err()
	if err != nil {
		return err
	}

	return nil
}

//Delete deletes all state data associated with the SessionID from the store.
func (rs *RedisStore) Delete(sid SessionID) error {
	//TODO: delete the data stored in redis for the provided SessionID
	redisKey := sid.getRedisKey()
	// handle errors here
	err := rs.Client.Del(redisKey).Err()
	if err != nil {
		return err
	}
	return nil
}

//getRedisKey() returns the redis key to use for the SessionID
func (sid SessionID) getRedisKey() string {
	//convert the SessionID to a string and add the prefix "sid:" to keep
	//SessionID keys separate from other keys that might end up in this
	//redis instance
	return "sid:" + sid.String()
}
