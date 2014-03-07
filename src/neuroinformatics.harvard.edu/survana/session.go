package survana

import (
	_ "log"
)

const (
	SESSION_ID         = "SSESSIONID"
	SESSION_COLLECTION = "sessions"
)

//Represents a user's session. Id and _id are kept separate so that in
//the future, Id's can be regenerated on every request.
//Id and Authenticated are aliases for Values['id'] and Values['authenticated']
type Session struct {
    DBO                             `bson:",inline,omitempty" json:"-"`
	Id            string            `bson:"id,omitempty" json:"-"`   //the publicly visible session id
	UserId        string            `bson:"user_id,omitempty" json:"-"` //the user id this session is associated with
	Authenticated bool              `bson:"authenticated" json:"-"`  //whether the user has logged in or not
	Values        map[string]string `bson:"values" json:"-"`         //all other values go here
}

//creates a new Session object with no Id.
func NewSession() *Session {
	return &Session{
        DBO: DBO { Collection: SESSION_COLLECTION },
		Authenticated: false,
		Values:        make(map[string]string, 0),
	}
}

// Loads session info from the database
// returns nil if the session doesn't exist
func FindSession(id string, db Database) (session *Session, err error) {
	session = NewSession()
	err = db.FindId(id, session)

	//if the session doesn't exist, return error
	if err != nil {
		//use nil session to show that it was not found
		if err == ErrNotFound {
			err = nil
		}
		return nil, err
	}

	return
}

// Deletes itself from the database
func (s *Session) Delete(db Database) (err error) {
	return db.Delete(s)
}

// Stores the session in the database
func (s *Session) Save(db Database) (err error) {
	return db.Save(s)
}
