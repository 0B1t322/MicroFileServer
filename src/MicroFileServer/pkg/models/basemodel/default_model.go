package basemodel

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TODO В будущем сделать возможность не удаления из базы данных а помечания как удаленного и не отдавать результат просто так

type BaseModel struct {
	IDField			`json:",inline" bson:",inline"`
	DateFields		`json:",inline" bson:",inline"`
}
// TODO все же надо поменять ID на string
type IDField struct {
	ID		primitive.ObjectID	`json:"id" bson:"_id,omitempty"`
}

func (f *IDField) PrepareID(id interface{}) (interface{}, error) {
	if idStr, ok := id.(string); ok {
		return primitive.ObjectIDFromHex(idStr)
	}

	// Otherwise id must be ObjectId
	return id, nil
}

// GetID method returns a model's ID
func (f *IDField) GetID() interface{} {
	return f.ID
}

// SetID sets the value of a model's ID field.
func (f *IDField) SetID(id interface{}) {
	f.ID = id.(primitive.ObjectID)
}

type DateFields struct {
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}

// Creating hook is used here to set the `created_at` field
// value when inserting a new model into the database.
func (f *DateFields) Creating() error {
	f.CreatedAt = time.Now().UTC().Round(time.Millisecond)
	return nil
}

// Saving hook is used here to set the `updated_at` field 
// value when creating or updateing a model.
func (f *DateFields) Saving() error {
	f.UpdatedAt = time.Now().UTC().Round(time.Millisecond)
	return nil
}

func (b *BaseModel) Creating() error {
	return b.DateFields.Creating()
}

func (b *BaseModel) Saving() error {
	return b.DateFields.Saving()
}