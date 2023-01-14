package domain

type Admin struct {
	UserID   string  `json:"userId" bson:"userId"`
	Login    string  `json:"login" bson:"login"`
	Password string  `json:"password" bson:"password"`
	Session  Session `json:"session,omitempty" bson:"session,omitempty"`
}
