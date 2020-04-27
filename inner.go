package gentest

import (
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/golang/protobuf/ptypes/wrappers"
	"google.golang.org/genproto/protobuf/field_mask"
	"google.golang.org/protobuf/runtime/protoimpl"
)

type Inner struct {
	fieldMask field_mask.FieldMask
	Foo       string `json:"foo,omitempty"`
	Bar       int64  `json:"bar,omitempty"`
}

type InnerProto struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Uri          string                `protobuf:"bytes,1,opt,name=uri,proto3" json:"uri,omitempty"`
	FriendlyName *wrappers.StringValue `protobuf:"bytes,2,opt,name=friendly_name,json=friendlyName,proto3" json:"friendly_name,omitempty"`
	// The time this host was created.
	CreatedTime *timestamp.Timestamp `protobuf:"bytes,3,opt,name=created_time,json=createdTime,proto3" json:"created_time,omitempty"`
	// The time this host was last updated.
	UpdatedTime *timestamp.Timestamp `protobuf:"bytes,4,opt,name=updated_time,json=updatedTime,proto3" json:"updated_time,omitempty"`
	// Marks the host as disabled.  Default is false.
	Disabled *wrappers.BoolValue `protobuf:"bytes,5,opt,name=disabled,proto3" json:"disabled,omitempty"`
	// This field is required.
	Address *wrappers.StringValue `protobuf:"bytes,6,opt,name=address,proto3" json:"address,omitempty"`
}
