// Copyright (c) 2018, Sylabs Inc. All rights reserved.
// This software is licensed under a 3-clause BSD license. Please consult the
// LICENSE file distributed with the sources of this project regarding your
// rights to use or distribute this software.

package client

import (
	"time"

	"github.com/globalsign/mgo/bson"
)

// LibraryModels lists names of valid models in the database
var LibraryModels = []string{"Entity", "Collection", "Container", "Image", "Blob"}

// ModelManager - Generic interface for models which must have a bson ObjectID
type ModelManager interface {
	GetID() bson.ObjectId
}

// BaseModel - has an ID, soft deletion marker, and Audit struct
type BaseModel struct {
	ModelManager `bson:",omitempty" json:",omitempty"`
	Deleted      bool      `bson:"deleted" json:"deleted"`
	CreatedBy    string    `bson:"createdBy" json:"createdBy"`
	CreatedAt    time.Time `bson:"createdAt" json:"createdAt"`
	UpdatedBy    string    `bson:"updatedBy,omitempty" json:"updatedBy,omitempty"`
	UpdatedAt    time.Time `bson:"updatedAt,omitempty" json:"updatedAt,omitempty"`
	DeletedBy    string    `bson:"deletedBy,omitempty" json:"deletedBy,omitempty"`
	DeletedAt    time.Time `bson:"deletedAt,omitempty" json:"deletedAt,omitempty"`
}

// IsDeleted - Convenience method to check soft deletion state if working with
// an interface
func (m BaseModel) IsDeleted() bool {
	return m.Deleted
}

// GetCreated - Convenience method to get creation stamps if working with an
// interface
func (m BaseModel) GetCreated() (auditUser string, auditTime time.Time) {
	return m.CreatedBy, m.CreatedAt
}

// GetUpdated - Convenience method to get update stamps if working with an
// interface
func (m BaseModel) GetUpdated() (auditUser string, auditTime time.Time) {
	return m.UpdatedBy, m.UpdatedAt
}

// GetDeleted - Convenience method to get deletino stamps if working with an
// interface
func (m BaseModel) GetDeleted() (auditUser string, auditTime time.Time) {
	return m.DeletedBy, m.DeletedAt
}

// Check BaseModel implements ModelManager at compile time
var _ ModelManager = (*BaseModel)(nil)

// Entity - Top level entry in the library, contains collections of images
// for a user or group
type Entity struct {
	BaseModel
	ID          bson.ObjectId   `bson:"_id" json:"id"`
	Name        string          `bson:"name" json:"name"`
	Description string          `bson:"description" json:"description"`
	Collections []bson.ObjectId `bson:"collections" json:"collections"`
}

// GetID - Convenience method to get model ID if working with an interface
func (e Entity) GetID() bson.ObjectId {
	return e.ID
}

// Collection - Second level in the library, holds a collection of containers
type Collection struct {
	BaseModel
	ID          bson.ObjectId   `bson:"_id" json:"id"`
	Name        string          `bson:"name" json:"name"`
	Description string          `bson:"description" json:"description"`
	Entity      bson.ObjectId   `bson:"entity" json:"entity"`
	Containers  []bson.ObjectId `bson:"containers" json:"containers"`
}

// GetID - Convenience method to get model ID if working with an interface
func (c Collection) GetID() bson.ObjectId {
	return c.ID
}

// Container - Third level of library. Inside a collection, holds images for
// a particular container
type Container struct {
	BaseModel
	ID          bson.ObjectId            `bson:"_id" json:"id"`
	Name        string                   `bson:"name" json:"name"`
	Description string                   `bson:"description" json:"description"`
	Collection  bson.ObjectId            `bson:"collection" json:"collection"`
	Images      []bson.ObjectId          `bson:"images" json:"images"`
	ImageTags   map[string]bson.ObjectId `bson:"imageTags" json:"imageTags"`
}

// GetID - Convenience method to get model ID if working with an interface
func (c Container) GetID() bson.ObjectId {
	return c.ID
}

// Image - Represents a Singularity image held by the library for a particular
// Container
type Image struct {
	BaseModel
	ID          bson.ObjectId `bson:"_id" json:"id"`
	Hash        string        `bson:"hash" json:"hash"`
	Description string        `bson:"description" json:"description"`
	Container   bson.ObjectId `bson:"container" json:"container"`
	Blob        bson.ObjectId `bson:"blob,omitempty" json:"blob,omitempty"`
	Size        int64         `bson:"size" json:"size"`
	Uploaded    bool          `bson:"uploaded" json:"uploaded"`
}

// GetID - Convenience method to get model ID if working with an interface
func (img Image) GetID() bson.ObjectId {
	return img.ID
}

// Blob - Binary data object (e.g. container image file) stored in a Backend
// Uses object store bucket/key semantics
type Blob struct {
	BaseModel
	ID          bson.ObjectId `bson:"_id" json:"id"`
	Bucket      string        `bson:"bucket" json:"bucket"`
	Key         string        `bson:"key" json:"key"`
	Size        int64         `bson:"size" json:"size"`
	ContentHash string        `bson:"contentHash" json:"contentHash"`
	Status      string        `bson:"status" json:"status"`
}

// GetID - Convenience method to get model ID if working with an interface
func (b Blob) GetID() bson.ObjectId {
	return b.ID
}

// ImageTag - A single mapping from a string to bson ID. Not stored in the DB
// but used by API calls setting tags
type ImageTag struct {
	Tag     string
	ImageID bson.ObjectId
}

// TagMap - A map of tags to imageIDs for a container
type TagMap map[string]bson.ObjectId
