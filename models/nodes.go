package models

import "time"

const NODE_COLLECTION = "nodes"

type Parameter struct {
	Min   float64 `json:"min" bson:"min"`
	Max   float64 `json:"max" bson:"max"`
	Label string  `json:"label" bson:"label"`
}

type NodeMetadata struct {
	Location    string `json:"location" bson:"location"`
	Sublocation string `json:"subLocation" bson:"subLocation"`
	MachineName string `json:"machineName" bson:"machineName"`
}

type Node struct {
	Uid string `json:"uid" bson:"uid"`

	Metadata NodeMetadata `json:"metadata" bson:"metadata"`
	User     string       `json:"user" bson:"user"`

	IsArchived bool `json:"isArchived" bson:"isArchived"`

	CreatedOn  int64  `json:"createdOn" bson:"createdOn"`
	CreatedBy  string `json:"createdBy" bson:"createdBy"`
	ModifiedOn int64  `json:"modifiedOn" bson:"modifiedOn"`
	ModifiedBy string `json:"modifiedBy" bson:"modifiedBy"`

	Parameters []Parameter `json:"parameters" bson:"parameters"`
}

func NewNode(uid, user string) *Node {
	now := time.Now().Unix()
	return &Node{
		Uid:        uid,
		CreatedOn:  now,
		CreatedBy:  user,
		ModifiedOn: now,
		ModifiedBy: user,
	}
}

func (n *Node) SetMetadata(location, sublocation, machineName string) *Node {
	n.Metadata = NodeMetadata{
		Location:    location,
		Sublocation: sublocation,
		MachineName: machineName,
	}
	return n
}

func (n *Node) Archived(b bool) *Node {
	n.IsArchived = b
	return n
}

func (n *Node) SetParameters(p ...Parameter) *Node {
	n.Parameters = append([]Parameter{}, p...)
	return n
}

func (n *Node) Save() error {
	return nil
}
