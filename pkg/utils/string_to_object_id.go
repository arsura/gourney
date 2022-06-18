package util

import "go.mongodb.org/mongo-driver/bson/primitive"

func StringToObjectId(str string) (*primitive.ObjectID, error) {
	oid, err := primitive.ObjectIDFromHex(str)
	if err != nil {
		return nil, err
	}
	return &oid, nil
}
