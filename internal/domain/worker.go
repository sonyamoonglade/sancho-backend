package domain

type Worker struct {
	UserID   string  `json:"userId" bson:"userId"`
	Login    string  `json:"login" bson:"login"`
	Role     Role    `json:"role" bson:"role"`
	Password string  `json:"password" bson:"password"`
	Session  Session `json:"session,omitempty" bson:"session,omitempty"`
	Name     string  `json:"name" bson:"name"`
}
