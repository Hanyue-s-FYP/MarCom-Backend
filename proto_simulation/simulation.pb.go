// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v5.27.3
// source: simulation.proto

package proto_simulation

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type AgentAttribute struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Key   string `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Value string `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *AgentAttribute) Reset() {
	*x = AgentAttribute{}
	if protoimpl.UnsafeEnabled {
		mi := &file_simulation_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AgentAttribute) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AgentAttribute) ProtoMessage() {}

func (x *AgentAttribute) ProtoReflect() protoreflect.Message {
	mi := &file_simulation_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AgentAttribute.ProtoReflect.Descriptor instead.
func (*AgentAttribute) Descriptor() ([]byte, []int) {
	return file_simulation_proto_rawDescGZIP(), []int{0}
}

func (x *AgentAttribute) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

func (x *AgentAttribute) GetValue() string {
	if x != nil {
		return x.Value
	}
	return ""
}

type Agent struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id    int32             `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Name  string            `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Desc  string            `protobuf:"bytes,3,opt,name=desc,proto3" json:"desc,omitempty"`
	Attrs []*AgentAttribute `protobuf:"bytes,4,rep,name=attrs,proto3" json:"attrs,omitempty"`
}

func (x *Agent) Reset() {
	*x = Agent{}
	if protoimpl.UnsafeEnabled {
		mi := &file_simulation_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Agent) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Agent) ProtoMessage() {}

func (x *Agent) ProtoReflect() protoreflect.Message {
	mi := &file_simulation_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Agent.ProtoReflect.Descriptor instead.
func (*Agent) Descriptor() ([]byte, []int) {
	return file_simulation_proto_rawDescGZIP(), []int{1}
}

func (x *Agent) GetId() int32 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *Agent) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Agent) GetDesc() string {
	if x != nil {
		return x.Desc
	}
	return ""
}

func (x *Agent) GetAttrs() []*AgentAttribute {
	if x != nil {
		return x.Attrs
	}
	return nil
}

type Product struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id    int32   `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Name  string  `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Desc  string  `protobuf:"bytes,3,opt,name=desc,proto3" json:"desc,omitempty"`
	Price float32 `protobuf:"fixed32,4,opt,name=price,proto3" json:"price,omitempty"`
	Cost  float32 `protobuf:"fixed32,5,opt,name=cost,proto3" json:"cost,omitempty"`
}

func (x *Product) Reset() {
	*x = Product{}
	if protoimpl.UnsafeEnabled {
		mi := &file_simulation_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Product) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Product) ProtoMessage() {}

func (x *Product) ProtoReflect() protoreflect.Message {
	mi := &file_simulation_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Product.ProtoReflect.Descriptor instead.
func (*Product) Descriptor() ([]byte, []int) {
	return file_simulation_proto_rawDescGZIP(), []int{2}
}

func (x *Product) GetId() int32 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *Product) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Product) GetDesc() string {
	if x != nil {
		return x.Desc
	}
	return ""
}

func (x *Product) GetPrice() float32 {
	if x != nil {
		return x.Price
	}
	return 0
}

func (x *Product) GetCost() float32 {
	if x != nil {
		return x.Cost
	}
	return 0
}

type SimulationRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id          int32      `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	EnvDesc     string     `protobuf:"bytes,2,opt,name=env_desc,json=envDesc,proto3" json:"env_desc,omitempty"`
	Agents      []*Agent   `protobuf:"bytes,3,rep,name=agents,proto3" json:"agents,omitempty"`
	Products    []*Product `protobuf:"bytes,4,rep,name=products,proto3" json:"products,omitempty"`
	TotalCycles int32      `protobuf:"varint,5,opt,name=total_cycles,json=totalCycles,proto3" json:"total_cycles,omitempty"`
}

func (x *SimulationRequest) Reset() {
	*x = SimulationRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_simulation_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SimulationRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SimulationRequest) ProtoMessage() {}

func (x *SimulationRequest) ProtoReflect() protoreflect.Message {
	mi := &file_simulation_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SimulationRequest.ProtoReflect.Descriptor instead.
func (*SimulationRequest) Descriptor() ([]byte, []int) {
	return file_simulation_proto_rawDescGZIP(), []int{3}
}

func (x *SimulationRequest) GetId() int32 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *SimulationRequest) GetEnvDesc() string {
	if x != nil {
		return x.EnvDesc
	}
	return ""
}

func (x *SimulationRequest) GetAgents() []*Agent {
	if x != nil {
		return x.Agents
	}
	return nil
}

func (x *SimulationRequest) GetProducts() []*Product {
	if x != nil {
		return x.Products
	}
	return nil
}

func (x *SimulationRequest) GetTotalCycles() int32 {
	if x != nil {
		return x.TotalCycles
	}
	return 0
}

type SimulationResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Message string `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
}

func (x *SimulationResponse) Reset() {
	*x = SimulationResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_simulation_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SimulationResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SimulationResponse) ProtoMessage() {}

func (x *SimulationResponse) ProtoReflect() protoreflect.Message {
	mi := &file_simulation_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SimulationResponse.ProtoReflect.Descriptor instead.
func (*SimulationResponse) Descriptor() ([]byte, []int) {
	return file_simulation_proto_rawDescGZIP(), []int{4}
}

func (x *SimulationResponse) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

type PauseRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SimulationId int32 `protobuf:"varint,1,opt,name=simulation_id,json=simulationId,proto3" json:"simulation_id,omitempty"`
}

func (x *PauseRequest) Reset() {
	*x = PauseRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_simulation_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PauseRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PauseRequest) ProtoMessage() {}

func (x *PauseRequest) ProtoReflect() protoreflect.Message {
	mi := &file_simulation_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PauseRequest.ProtoReflect.Descriptor instead.
func (*PauseRequest) Descriptor() ([]byte, []int) {
	return file_simulation_proto_rawDescGZIP(), []int{5}
}

func (x *PauseRequest) GetSimulationId() int32 {
	if x != nil {
		return x.SimulationId
	}
	return 0
}

type PauseResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Message string `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
}

func (x *PauseResponse) Reset() {
	*x = PauseResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_simulation_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PauseResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PauseResponse) ProtoMessage() {}

func (x *PauseResponse) ProtoReflect() protoreflect.Message {
	mi := &file_simulation_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PauseResponse.ProtoReflect.Descriptor instead.
func (*PauseResponse) Descriptor() ([]byte, []int) {
	return file_simulation_proto_rawDescGZIP(), []int{6}
}

func (x *PauseResponse) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

type StreamRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SimulationId int32 `protobuf:"varint,1,opt,name=simulation_id,json=simulationId,proto3" json:"simulation_id,omitempty"`
}

func (x *StreamRequest) Reset() {
	*x = StreamRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_simulation_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StreamRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StreamRequest) ProtoMessage() {}

func (x *StreamRequest) ProtoReflect() protoreflect.Message {
	mi := &file_simulation_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StreamRequest.ProtoReflect.Descriptor instead.
func (*StreamRequest) Descriptor() ([]byte, []int) {
	return file_simulation_proto_rawDescGZIP(), []int{7}
}

func (x *StreamRequest) GetSimulationId() int32 {
	if x != nil {
		return x.SimulationId
	}
	return 0
}

type SimulationUpdate struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	AgentId      int32  `protobuf:"varint,1,opt,name=agent_id,json=agentId,proto3" json:"agent_id,omitempty"`
	Action       string `protobuf:"bytes,2,opt,name=action,proto3" json:"action,omitempty"`
	Content      string `protobuf:"bytes,3,opt,name=content,proto3" json:"content,omitempty"`                                // see db.py in MarCom-SimulationCore line 35
	Cycle        int32  `protobuf:"varint,4,opt,name=cycle,proto3" json:"cycle,omitempty"`                                   // which cycle is that (if is new cycle then need insert)
	SimulationId int32  `protobuf:"varint,5,opt,name=simulation_id,json=simulationId,proto3" json:"simulation_id,omitempty"` // to associate the cycle with the simulation
}

func (x *SimulationUpdate) Reset() {
	*x = SimulationUpdate{}
	if protoimpl.UnsafeEnabled {
		mi := &file_simulation_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SimulationUpdate) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SimulationUpdate) ProtoMessage() {}

func (x *SimulationUpdate) ProtoReflect() protoreflect.Message {
	mi := &file_simulation_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SimulationUpdate.ProtoReflect.Descriptor instead.
func (*SimulationUpdate) Descriptor() ([]byte, []int) {
	return file_simulation_proto_rawDescGZIP(), []int{8}
}

func (x *SimulationUpdate) GetAgentId() int32 {
	if x != nil {
		return x.AgentId
	}
	return 0
}

func (x *SimulationUpdate) GetAction() string {
	if x != nil {
		return x.Action
	}
	return ""
}

func (x *SimulationUpdate) GetContent() string {
	if x != nil {
		return x.Content
	}
	return ""
}

func (x *SimulationUpdate) GetCycle() int32 {
	if x != nil {
		return x.Cycle
	}
	return 0
}

func (x *SimulationUpdate) GetSimulationId() int32 {
	if x != nil {
		return x.SimulationId
	}
	return 0
}

var File_simulation_proto protoreflect.FileDescriptor

var file_simulation_proto_rawDesc = []byte{
	0x0a, 0x10, 0x73, 0x69, 0x6d, 0x75, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x0a, 0x73, 0x69, 0x6d, 0x75, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x22, 0x38,
	0x0a, 0x0e, 0x41, 0x67, 0x65, 0x6e, 0x74, 0x41, 0x74, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x65,
	0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b,
	0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x22, 0x71, 0x0a, 0x05, 0x41, 0x67, 0x65, 0x6e,
	0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x02, 0x69,
	0x64, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x64, 0x65, 0x73, 0x63, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x04, 0x64, 0x65, 0x73, 0x63, 0x12, 0x30, 0x0a, 0x05, 0x61, 0x74, 0x74,
	0x72, 0x73, 0x18, 0x04, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x73, 0x69, 0x6d, 0x75, 0x6c,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x41, 0x67, 0x65, 0x6e, 0x74, 0x41, 0x74, 0x74, 0x72, 0x69,
	0x62, 0x75, 0x74, 0x65, 0x52, 0x05, 0x61, 0x74, 0x74, 0x72, 0x73, 0x22, 0x6b, 0x0a, 0x07, 0x50,
	0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x05, 0x52, 0x02, 0x69, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x64, 0x65,
	0x73, 0x63, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x64, 0x65, 0x73, 0x63, 0x12, 0x14,
	0x0a, 0x05, 0x70, 0x72, 0x69, 0x63, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x02, 0x52, 0x05, 0x70,
	0x72, 0x69, 0x63, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x63, 0x6f, 0x73, 0x74, 0x18, 0x05, 0x20, 0x01,
	0x28, 0x02, 0x52, 0x04, 0x63, 0x6f, 0x73, 0x74, 0x22, 0xbd, 0x01, 0x0a, 0x11, 0x53, 0x69, 0x6d,
	0x75, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x0e,
	0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x02, 0x69, 0x64, 0x12, 0x19,
	0x0a, 0x08, 0x65, 0x6e, 0x76, 0x5f, 0x64, 0x65, 0x73, 0x63, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x07, 0x65, 0x6e, 0x76, 0x44, 0x65, 0x73, 0x63, 0x12, 0x29, 0x0a, 0x06, 0x61, 0x67, 0x65,
	0x6e, 0x74, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x73, 0x69, 0x6d, 0x75,
	0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x41, 0x67, 0x65, 0x6e, 0x74, 0x52, 0x06, 0x61, 0x67,
	0x65, 0x6e, 0x74, 0x73, 0x12, 0x2f, 0x0a, 0x08, 0x70, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x73,
	0x18, 0x04, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x13, 0x2e, 0x73, 0x69, 0x6d, 0x75, 0x6c, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x2e, 0x50, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x52, 0x08, 0x70, 0x72, 0x6f,
	0x64, 0x75, 0x63, 0x74, 0x73, 0x12, 0x21, 0x0a, 0x0c, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x5f, 0x63,
	0x79, 0x63, 0x6c, 0x65, 0x73, 0x18, 0x05, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0b, 0x74, 0x6f, 0x74,
	0x61, 0x6c, 0x43, 0x79, 0x63, 0x6c, 0x65, 0x73, 0x22, 0x2e, 0x0a, 0x12, 0x53, 0x69, 0x6d, 0x75,
	0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x18,
	0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x22, 0x33, 0x0a, 0x0c, 0x50, 0x61, 0x75, 0x73,
	0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x23, 0x0a, 0x0d, 0x73, 0x69, 0x6d, 0x75,
	0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52,
	0x0c, 0x73, 0x69, 0x6d, 0x75, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x64, 0x22, 0x29, 0x0a,
	0x0d, 0x50, 0x61, 0x75, 0x73, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x18,
	0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x22, 0x34, 0x0a, 0x0d, 0x53, 0x74, 0x72, 0x65,
	0x61, 0x6d, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x23, 0x0a, 0x0d, 0x73, 0x69, 0x6d,
	0x75, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05,
	0x52, 0x0c, 0x73, 0x69, 0x6d, 0x75, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x64, 0x22, 0x9a,
	0x01, 0x0a, 0x10, 0x53, 0x69, 0x6d, 0x75, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x55, 0x70, 0x64,
	0x61, 0x74, 0x65, 0x12, 0x19, 0x0a, 0x08, 0x61, 0x67, 0x65, 0x6e, 0x74, 0x5f, 0x69, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x07, 0x61, 0x67, 0x65, 0x6e, 0x74, 0x49, 0x64, 0x12, 0x16,
	0x0a, 0x06, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06,
	0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x18, 0x0a, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e,
	0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74,
	0x12, 0x14, 0x0a, 0x05, 0x63, 0x79, 0x63, 0x6c, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x05, 0x52,
	0x05, 0x63, 0x79, 0x63, 0x6c, 0x65, 0x12, 0x23, 0x0a, 0x0d, 0x73, 0x69, 0x6d, 0x75, 0x6c, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x69, 0x64, 0x18, 0x05, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0c, 0x73,
	0x69, 0x6d, 0x75, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x64, 0x32, 0xf9, 0x01, 0x0a, 0x11,
	0x53, 0x69, 0x6d, 0x75, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x12, 0x50, 0x0a, 0x0f, 0x53, 0x74, 0x61, 0x72, 0x74, 0x53, 0x69, 0x6d, 0x75, 0x6c, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x12, 0x1d, 0x2e, 0x73, 0x69, 0x6d, 0x75, 0x6c, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x2e, 0x53, 0x69, 0x6d, 0x75, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x1e, 0x2e, 0x73, 0x69, 0x6d, 0x75, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x2e, 0x53, 0x69, 0x6d, 0x75, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x46, 0x0a, 0x0f, 0x50, 0x61, 0x75, 0x73, 0x65, 0x53, 0x69, 0x6d, 0x75,
	0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x18, 0x2e, 0x73, 0x69, 0x6d, 0x75, 0x6c, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x2e, 0x50, 0x61, 0x75, 0x73, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x19, 0x2e, 0x73, 0x69, 0x6d, 0x75, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x50, 0x61,
	0x75, 0x73, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x4a, 0x0a, 0x0d, 0x53,
	0x74, 0x72, 0x65, 0x61, 0x6d, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x73, 0x12, 0x19, 0x2e, 0x73,
	0x69, 0x6d, 0x75, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1c, 0x2e, 0x73, 0x69, 0x6d, 0x75, 0x6c, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x53, 0x69, 0x6d, 0x75, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x55,
	0x70, 0x64, 0x61, 0x74, 0x65, 0x30, 0x01, 0x42, 0x14, 0x5a, 0x12, 0x2e, 0x2f, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x5f, 0x73, 0x69, 0x6d, 0x75, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_simulation_proto_rawDescOnce sync.Once
	file_simulation_proto_rawDescData = file_simulation_proto_rawDesc
)

func file_simulation_proto_rawDescGZIP() []byte {
	file_simulation_proto_rawDescOnce.Do(func() {
		file_simulation_proto_rawDescData = protoimpl.X.CompressGZIP(file_simulation_proto_rawDescData)
	})
	return file_simulation_proto_rawDescData
}

var file_simulation_proto_msgTypes = make([]protoimpl.MessageInfo, 9)
var file_simulation_proto_goTypes = []any{
	(*AgentAttribute)(nil),     // 0: simulation.AgentAttribute
	(*Agent)(nil),              // 1: simulation.Agent
	(*Product)(nil),            // 2: simulation.Product
	(*SimulationRequest)(nil),  // 3: simulation.SimulationRequest
	(*SimulationResponse)(nil), // 4: simulation.SimulationResponse
	(*PauseRequest)(nil),       // 5: simulation.PauseRequest
	(*PauseResponse)(nil),      // 6: simulation.PauseResponse
	(*StreamRequest)(nil),      // 7: simulation.StreamRequest
	(*SimulationUpdate)(nil),   // 8: simulation.SimulationUpdate
}
var file_simulation_proto_depIdxs = []int32{
	0, // 0: simulation.Agent.attrs:type_name -> simulation.AgentAttribute
	1, // 1: simulation.SimulationRequest.agents:type_name -> simulation.Agent
	2, // 2: simulation.SimulationRequest.products:type_name -> simulation.Product
	3, // 3: simulation.SimulationService.StartSimulation:input_type -> simulation.SimulationRequest
	5, // 4: simulation.SimulationService.PauseSimulation:input_type -> simulation.PauseRequest
	7, // 5: simulation.SimulationService.StreamUpdates:input_type -> simulation.StreamRequest
	4, // 6: simulation.SimulationService.StartSimulation:output_type -> simulation.SimulationResponse
	6, // 7: simulation.SimulationService.PauseSimulation:output_type -> simulation.PauseResponse
	8, // 8: simulation.SimulationService.StreamUpdates:output_type -> simulation.SimulationUpdate
	6, // [6:9] is the sub-list for method output_type
	3, // [3:6] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_simulation_proto_init() }
func file_simulation_proto_init() {
	if File_simulation_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_simulation_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*AgentAttribute); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_simulation_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*Agent); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_simulation_proto_msgTypes[2].Exporter = func(v any, i int) any {
			switch v := v.(*Product); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_simulation_proto_msgTypes[3].Exporter = func(v any, i int) any {
			switch v := v.(*SimulationRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_simulation_proto_msgTypes[4].Exporter = func(v any, i int) any {
			switch v := v.(*SimulationResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_simulation_proto_msgTypes[5].Exporter = func(v any, i int) any {
			switch v := v.(*PauseRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_simulation_proto_msgTypes[6].Exporter = func(v any, i int) any {
			switch v := v.(*PauseResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_simulation_proto_msgTypes[7].Exporter = func(v any, i int) any {
			switch v := v.(*StreamRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_simulation_proto_msgTypes[8].Exporter = func(v any, i int) any {
			switch v := v.(*SimulationUpdate); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_simulation_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   9,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_simulation_proto_goTypes,
		DependencyIndexes: file_simulation_proto_depIdxs,
		MessageInfos:      file_simulation_proto_msgTypes,
	}.Build()
	File_simulation_proto = out.File
	file_simulation_proto_rawDesc = nil
	file_simulation_proto_goTypes = nil
	file_simulation_proto_depIdxs = nil
}
