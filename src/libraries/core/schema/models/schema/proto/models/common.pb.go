// Code generated by protoc-gen-go. DO NOT EDIT.
// source: schema/proto/model/common.proto

package model

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	timestamp "github.com/golang/protobuf/ptypes/timestamp"
	_ "github.com/infobloxopen/protoc-gen-gorm/options"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type AccountType int32

const (
	AccountType_RegularUser AccountType = 0
	AccountType_Startup     AccountType = 1
	AccountType_Investor    AccountType = 2
)

var AccountType_name = map[int32]string{
	0: "RegularUser",
	1: "Startup",
	2: "Investor",
}

var AccountType_value = map[string]int32{
	"RegularUser": 0,
	"Startup":     1,
	"Investor":    2,
}

func (x AccountType) String() string {
	return proto.EnumName(AccountType_name, int32(x))
}

func (AccountType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_8bd0ff5c490a2f61, []int{0}
}

type TeamAccountType int32

const (
	TeamAccountType_StartupTeam  TeamAccountType = 0
	TeamAccountType_InvestorTeam TeamAccountType = 1
)

var TeamAccountType_name = map[int32]string{
	0: "StartupTeam",
	1: "InvestorTeam",
}

var TeamAccountType_value = map[string]int32{
	"StartupTeam":  0,
	"InvestorTeam": 1,
}

func (x TeamAccountType) String() string {
	return proto.EnumName(TeamAccountType_name, int32(x))
}

func (TeamAccountType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_8bd0ff5c490a2f61, []int{1}
}

type Address struct {
	Id                   int32                `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Longitude            string               `protobuf:"bytes,2,opt,name=longitude,proto3" json:"longitude,omitempty"`
	Latitude             string               `protobuf:"bytes,3,opt,name=latitude,proto3" json:"latitude,omitempty"`
	City                 string               `protobuf:"bytes,4,opt,name=city,proto3" json:"city,omitempty"`
	State                string               `protobuf:"bytes,5,opt,name=state,proto3" json:"state,omitempty"`
	Country              string               `protobuf:"bytes,6,opt,name=country,proto3" json:"country,omitempty"`
	ZipCode              string               `protobuf:"bytes,7,opt,name=zipCode,proto3" json:"zipCode,omitempty"`
	Street               string               `protobuf:"bytes,8,opt,name=street,proto3" json:"street,omitempty"`
	BuildingNumber       string               `protobuf:"bytes,9,opt,name=building_number,json=buildingNumber,proto3" json:"building_number,omitempty"`
	CreatedAt            *timestamp.Timestamp `protobuf:"bytes,10,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	UpdatedAt            *timestamp.Timestamp `protobuf:"bytes,11,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty"`
	DeletedAt            *timestamp.Timestamp `protobuf:"bytes,12,opt,name=deleted_at,json=deletedAt,proto3" json:"deleted_at,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *Address) Reset()         { *m = Address{} }
func (m *Address) String() string { return proto.CompactTextString(m) }
func (*Address) ProtoMessage()    {}
func (*Address) Descriptor() ([]byte, []int) {
	return fileDescriptor_8bd0ff5c490a2f61, []int{0}
}

func (m *Address) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Address.Unmarshal(m, b)
}
func (m *Address) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Address.Marshal(b, m, deterministic)
}
func (m *Address) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Address.Merge(m, src)
}
func (m *Address) XXX_Size() int {
	return xxx_messageInfo_Address.Size(m)
}
func (m *Address) XXX_DiscardUnknown() {
	xxx_messageInfo_Address.DiscardUnknown(m)
}

var xxx_messageInfo_Address proto.InternalMessageInfo

func (m *Address) GetId() int32 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *Address) GetLongitude() string {
	if m != nil {
		return m.Longitude
	}
	return ""
}

func (m *Address) GetLatitude() string {
	if m != nil {
		return m.Latitude
	}
	return ""
}

func (m *Address) GetCity() string {
	if m != nil {
		return m.City
	}
	return ""
}

func (m *Address) GetState() string {
	if m != nil {
		return m.State
	}
	return ""
}

func (m *Address) GetCountry() string {
	if m != nil {
		return m.Country
	}
	return ""
}

func (m *Address) GetZipCode() string {
	if m != nil {
		return m.ZipCode
	}
	return ""
}

func (m *Address) GetStreet() string {
	if m != nil {
		return m.Street
	}
	return ""
}

func (m *Address) GetBuildingNumber() string {
	if m != nil {
		return m.BuildingNumber
	}
	return ""
}

func (m *Address) GetCreatedAt() *timestamp.Timestamp {
	if m != nil {
		return m.CreatedAt
	}
	return nil
}

func (m *Address) GetUpdatedAt() *timestamp.Timestamp {
	if m != nil {
		return m.UpdatedAt
	}
	return nil
}

func (m *Address) GetDeletedAt() *timestamp.Timestamp {
	if m != nil {
		return m.DeletedAt
	}
	return nil
}

type Education struct {
	Id                   int32                `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	School               string               `protobuf:"bytes,2,opt,name=school,proto3" json:"school,omitempty"`
	Degree               string               `protobuf:"bytes,3,opt,name=degree,proto3" json:"degree,omitempty"`
	FieldOfStudy         string               `protobuf:"bytes,4,opt,name=field_of_study,json=fieldOfStudy,proto3" json:"field_of_study,omitempty"`
	StartDate            *timestamp.Timestamp `protobuf:"bytes,5,opt,name=start_date,json=startDate,proto3" json:"start_date,omitempty"`
	EndDate              *timestamp.Timestamp `protobuf:"bytes,6,opt,name=end_date,json=endDate,proto3" json:"end_date,omitempty"`
	CurrentlyAttending   bool                 `protobuf:"varint,7,opt,name=currently_attending,json=currentlyAttending,proto3" json:"currently_attending,omitempty"`
	Gpa                  float32              `protobuf:"fixed32,8,opt,name=gpa,proto3" json:"gpa,omitempty"`
	Activities           string               `protobuf:"bytes,9,opt,name=activities,proto3" json:"activities,omitempty"`
	Societies            string               `protobuf:"bytes,10,opt,name=societies,proto3" json:"societies,omitempty"`
	Description          string               `protobuf:"bytes,11,opt,name=description,proto3" json:"description,omitempty"`
	CreatedAt            *timestamp.Timestamp `protobuf:"bytes,12,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	UpdatedAt            *timestamp.Timestamp `protobuf:"bytes,13,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty"`
	DeletedAt            *timestamp.Timestamp `protobuf:"bytes,14,opt,name=deleted_at,json=deletedAt,proto3" json:"deleted_at,omitempty"`
	MediaId              *Media               `protobuf:"bytes,15,opt,name=media_id,json=mediaId,proto3" json:"media_id,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *Education) Reset()         { *m = Education{} }
func (m *Education) String() string { return proto.CompactTextString(m) }
func (*Education) ProtoMessage()    {}
func (*Education) Descriptor() ([]byte, []int) {
	return fileDescriptor_8bd0ff5c490a2f61, []int{1}
}

func (m *Education) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Education.Unmarshal(m, b)
}
func (m *Education) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Education.Marshal(b, m, deterministic)
}
func (m *Education) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Education.Merge(m, src)
}
func (m *Education) XXX_Size() int {
	return xxx_messageInfo_Education.Size(m)
}
func (m *Education) XXX_DiscardUnknown() {
	xxx_messageInfo_Education.DiscardUnknown(m)
}

var xxx_messageInfo_Education proto.InternalMessageInfo

func (m *Education) GetId() int32 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *Education) GetSchool() string {
	if m != nil {
		return m.School
	}
	return ""
}

func (m *Education) GetDegree() string {
	if m != nil {
		return m.Degree
	}
	return ""
}

func (m *Education) GetFieldOfStudy() string {
	if m != nil {
		return m.FieldOfStudy
	}
	return ""
}

func (m *Education) GetStartDate() *timestamp.Timestamp {
	if m != nil {
		return m.StartDate
	}
	return nil
}

func (m *Education) GetEndDate() *timestamp.Timestamp {
	if m != nil {
		return m.EndDate
	}
	return nil
}

func (m *Education) GetCurrentlyAttending() bool {
	if m != nil {
		return m.CurrentlyAttending
	}
	return false
}

func (m *Education) GetGpa() float32 {
	if m != nil {
		return m.Gpa
	}
	return 0
}

func (m *Education) GetActivities() string {
	if m != nil {
		return m.Activities
	}
	return ""
}

func (m *Education) GetSocieties() string {
	if m != nil {
		return m.Societies
	}
	return ""
}

func (m *Education) GetDescription() string {
	if m != nil {
		return m.Description
	}
	return ""
}

func (m *Education) GetCreatedAt() *timestamp.Timestamp {
	if m != nil {
		return m.CreatedAt
	}
	return nil
}

func (m *Education) GetUpdatedAt() *timestamp.Timestamp {
	if m != nil {
		return m.UpdatedAt
	}
	return nil
}

func (m *Education) GetDeletedAt() *timestamp.Timestamp {
	if m != nil {
		return m.DeletedAt
	}
	return nil
}

func (m *Education) GetMediaId() *Media {
	if m != nil {
		return m.MediaId
	}
	return nil
}

type Media struct {
	Id                   int32                `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	DocumentLinks        []string             `protobuf:"bytes,2,rep,name=document_links,json=documentLinks,proto3" json:"document_links,omitempty"`
	PhotoLinks           []string             `protobuf:"bytes,3,rep,name=photo_links,json=photoLinks,proto3" json:"photo_links,omitempty"`
	VideoLinks           []string             `protobuf:"bytes,4,rep,name=video_links,json=videoLinks,proto3" json:"video_links,omitempty"`
	PresentationLinks    []string             `protobuf:"bytes,5,rep,name=presentation_links,json=presentationLinks,proto3" json:"presentation_links,omitempty"`
	CreatedAt            *timestamp.Timestamp `protobuf:"bytes,6,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	DeletedAt            *timestamp.Timestamp `protobuf:"bytes,7,opt,name=deleted_at,json=deletedAt,proto3" json:"deleted_at,omitempty"`
	UpdatedAt            *timestamp.Timestamp `protobuf:"bytes,8,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *Media) Reset()         { *m = Media{} }
func (m *Media) String() string { return proto.CompactTextString(m) }
func (*Media) ProtoMessage()    {}
func (*Media) Descriptor() ([]byte, []int) {
	return fileDescriptor_8bd0ff5c490a2f61, []int{2}
}

func (m *Media) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Media.Unmarshal(m, b)
}
func (m *Media) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Media.Marshal(b, m, deterministic)
}
func (m *Media) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Media.Merge(m, src)
}
func (m *Media) XXX_Size() int {
	return xxx_messageInfo_Media.Size(m)
}
func (m *Media) XXX_DiscardUnknown() {
	xxx_messageInfo_Media.DiscardUnknown(m)
}

var xxx_messageInfo_Media proto.InternalMessageInfo

func (m *Media) GetId() int32 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *Media) GetDocumentLinks() []string {
	if m != nil {
		return m.DocumentLinks
	}
	return nil
}

func (m *Media) GetPhotoLinks() []string {
	if m != nil {
		return m.PhotoLinks
	}
	return nil
}

func (m *Media) GetVideoLinks() []string {
	if m != nil {
		return m.VideoLinks
	}
	return nil
}

func (m *Media) GetPresentationLinks() []string {
	if m != nil {
		return m.PresentationLinks
	}
	return nil
}

func (m *Media) GetCreatedAt() *timestamp.Timestamp {
	if m != nil {
		return m.CreatedAt
	}
	return nil
}

func (m *Media) GetDeletedAt() *timestamp.Timestamp {
	if m != nil {
		return m.DeletedAt
	}
	return nil
}

func (m *Media) GetUpdatedAt() *timestamp.Timestamp {
	if m != nil {
		return m.UpdatedAt
	}
	return nil
}

type Subscriptions struct {
	Id                   uint32               `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	SubscriptionName     string               `protobuf:"bytes,2,opt,name=subscription_name,json=subscriptionName,proto3" json:"subscription_name,omitempty"`
	SubscriptionStatus   string               `protobuf:"bytes,3,opt,name=subscription_status,json=subscriptionStatus,proto3" json:"subscription_status,omitempty"`
	StartDate            *timestamp.Timestamp `protobuf:"bytes,4,opt,name=start_date,json=startDate,proto3" json:"start_date,omitempty"`
	EndDate              *timestamp.Timestamp `protobuf:"bytes,5,opt,name=end_date,json=endDate,proto3" json:"end_date,omitempty"`
	AccessType           string               `protobuf:"bytes,6,opt,name=access_type,json=accessType,proto3" json:"access_type,omitempty"`
	IsActive             bool                 `protobuf:"varint,7,opt,name=is_active,json=isActive,proto3" json:"is_active,omitempty"`
	CreatedAt            *timestamp.Timestamp `protobuf:"bytes,8,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	DeletedAt            *timestamp.Timestamp `protobuf:"bytes,9,opt,name=deleted_at,json=deletedAt,proto3" json:"deleted_at,omitempty"`
	UpdatedAt            *timestamp.Timestamp `protobuf:"bytes,10,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *Subscriptions) Reset()         { *m = Subscriptions{} }
func (m *Subscriptions) String() string { return proto.CompactTextString(m) }
func (*Subscriptions) ProtoMessage()    {}
func (*Subscriptions) Descriptor() ([]byte, []int) {
	return fileDescriptor_8bd0ff5c490a2f61, []int{3}
}

func (m *Subscriptions) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Subscriptions.Unmarshal(m, b)
}
func (m *Subscriptions) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Subscriptions.Marshal(b, m, deterministic)
}
func (m *Subscriptions) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Subscriptions.Merge(m, src)
}
func (m *Subscriptions) XXX_Size() int {
	return xxx_messageInfo_Subscriptions.Size(m)
}
func (m *Subscriptions) XXX_DiscardUnknown() {
	xxx_messageInfo_Subscriptions.DiscardUnknown(m)
}

var xxx_messageInfo_Subscriptions proto.InternalMessageInfo

func (m *Subscriptions) GetId() uint32 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *Subscriptions) GetSubscriptionName() string {
	if m != nil {
		return m.SubscriptionName
	}
	return ""
}

func (m *Subscriptions) GetSubscriptionStatus() string {
	if m != nil {
		return m.SubscriptionStatus
	}
	return ""
}

func (m *Subscriptions) GetStartDate() *timestamp.Timestamp {
	if m != nil {
		return m.StartDate
	}
	return nil
}

func (m *Subscriptions) GetEndDate() *timestamp.Timestamp {
	if m != nil {
		return m.EndDate
	}
	return nil
}

func (m *Subscriptions) GetAccessType() string {
	if m != nil {
		return m.AccessType
	}
	return ""
}

func (m *Subscriptions) GetIsActive() bool {
	if m != nil {
		return m.IsActive
	}
	return false
}

func (m *Subscriptions) GetCreatedAt() *timestamp.Timestamp {
	if m != nil {
		return m.CreatedAt
	}
	return nil
}

func (m *Subscriptions) GetDeletedAt() *timestamp.Timestamp {
	if m != nil {
		return m.DeletedAt
	}
	return nil
}

func (m *Subscriptions) GetUpdatedAt() *timestamp.Timestamp {
	if m != nil {
		return m.UpdatedAt
	}
	return nil
}

type SocialMedia struct {
	Id                   int32                `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	GithubUrl            uint32               `protobuf:"varint,2,opt,name=github_url,json=githubUrl,proto3" json:"github_url,omitempty"`
	WebsiteUrl           string               `protobuf:"bytes,3,opt,name=website_url,json=websiteUrl,proto3" json:"website_url,omitempty"`
	FacebookUrl          string               `protobuf:"bytes,4,opt,name=facebook_url,json=facebookUrl,proto3" json:"facebook_url,omitempty"`
	TwitterUrl           string               `protobuf:"bytes,5,opt,name=twitter_url,json=twitterUrl,proto3" json:"twitter_url,omitempty"`
	LinkedUrl            string               `protobuf:"bytes,6,opt,name=linked_url,json=linkedUrl,proto3" json:"linked_url,omitempty"`
	YoutubeUrl           string               `protobuf:"bytes,7,opt,name=youtube_url,json=youtubeUrl,proto3" json:"youtube_url,omitempty"`
	CreatedAt            *timestamp.Timestamp `protobuf:"bytes,8,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	DeletedAt            *timestamp.Timestamp `protobuf:"bytes,9,opt,name=deleted_at,json=deletedAt,proto3" json:"deleted_at,omitempty"`
	UpdatedAt            *timestamp.Timestamp `protobuf:"bytes,10,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *SocialMedia) Reset()         { *m = SocialMedia{} }
func (m *SocialMedia) String() string { return proto.CompactTextString(m) }
func (*SocialMedia) ProtoMessage()    {}
func (*SocialMedia) Descriptor() ([]byte, []int) {
	return fileDescriptor_8bd0ff5c490a2f61, []int{4}
}

func (m *SocialMedia) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SocialMedia.Unmarshal(m, b)
}
func (m *SocialMedia) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SocialMedia.Marshal(b, m, deterministic)
}
func (m *SocialMedia) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SocialMedia.Merge(m, src)
}
func (m *SocialMedia) XXX_Size() int {
	return xxx_messageInfo_SocialMedia.Size(m)
}
func (m *SocialMedia) XXX_DiscardUnknown() {
	xxx_messageInfo_SocialMedia.DiscardUnknown(m)
}

var xxx_messageInfo_SocialMedia proto.InternalMessageInfo

func (m *SocialMedia) GetId() int32 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *SocialMedia) GetGithubUrl() uint32 {
	if m != nil {
		return m.GithubUrl
	}
	return 0
}

func (m *SocialMedia) GetWebsiteUrl() string {
	if m != nil {
		return m.WebsiteUrl
	}
	return ""
}

func (m *SocialMedia) GetFacebookUrl() string {
	if m != nil {
		return m.FacebookUrl
	}
	return ""
}

func (m *SocialMedia) GetTwitterUrl() string {
	if m != nil {
		return m.TwitterUrl
	}
	return ""
}

func (m *SocialMedia) GetLinkedUrl() string {
	if m != nil {
		return m.LinkedUrl
	}
	return ""
}

func (m *SocialMedia) GetYoutubeUrl() string {
	if m != nil {
		return m.YoutubeUrl
	}
	return ""
}

func (m *SocialMedia) GetCreatedAt() *timestamp.Timestamp {
	if m != nil {
		return m.CreatedAt
	}
	return nil
}

func (m *SocialMedia) GetDeletedAt() *timestamp.Timestamp {
	if m != nil {
		return m.DeletedAt
	}
	return nil
}

func (m *SocialMedia) GetUpdatedAt() *timestamp.Timestamp {
	if m != nil {
		return m.UpdatedAt
	}
	return nil
}

type Details struct {
	IPOStatus            string   `protobuf:"bytes,1,opt,name=IPOStatus,proto3" json:"IPOStatus,omitempty"`
	CompanyType          string   `protobuf:"bytes,2,opt,name=CompanyType,proto3" json:"CompanyType,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Details) Reset()         { *m = Details{} }
func (m *Details) String() string { return proto.CompactTextString(m) }
func (*Details) ProtoMessage()    {}
func (*Details) Descriptor() ([]byte, []int) {
	return fileDescriptor_8bd0ff5c490a2f61, []int{5}
}

func (m *Details) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Details.Unmarshal(m, b)
}
func (m *Details) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Details.Marshal(b, m, deterministic)
}
func (m *Details) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Details.Merge(m, src)
}
func (m *Details) XXX_Size() int {
	return xxx_messageInfo_Details.Size(m)
}
func (m *Details) XXX_DiscardUnknown() {
	xxx_messageInfo_Details.DiscardUnknown(m)
}

var xxx_messageInfo_Details proto.InternalMessageInfo

func (m *Details) GetIPOStatus() string {
	if m != nil {
		return m.IPOStatus
	}
	return ""
}

func (m *Details) GetCompanyType() string {
	if m != nil {
		return m.CompanyType
	}
	return ""
}

type Experience struct {
	Id                   int32                `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	CompanyName          string               `protobuf:"bytes,2,opt,name=company_name,json=companyName,proto3" json:"company_name,omitempty"`
	Title                string               `protobuf:"bytes,3,opt,name=title,proto3" json:"title,omitempty"`
	EmploymentType       string               `protobuf:"bytes,4,opt,name=employment_type,json=employmentType,proto3" json:"employment_type,omitempty"`
	Location             string               `protobuf:"bytes,5,opt,name=location,proto3" json:"location,omitempty"`
	StartDate            *timestamp.Timestamp `protobuf:"bytes,6,opt,name=start_date,json=startDate,proto3" json:"start_date,omitempty"`
	EndDate              *timestamp.Timestamp `protobuf:"bytes,7,opt,name=end_date,json=endDate,proto3" json:"end_date,omitempty"`
	IsCurrentJob         bool                 `protobuf:"varint,8,opt,name=is_current_job,json=isCurrentJob,proto3" json:"is_current_job,omitempty"`
	Headline             string               `protobuf:"bytes,9,opt,name=headline,proto3" json:"headline,omitempty"`
	Description          string               `protobuf:"bytes,10,opt,name=description,proto3" json:"description,omitempty"`
	MediaId              *Media               `protobuf:"bytes,11,opt,name=media_id,json=mediaId,proto3" json:"media_id,omitempty"`
	CreatedAt            *timestamp.Timestamp `protobuf:"bytes,12,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	UpdatedAt            *timestamp.Timestamp `protobuf:"bytes,13,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty"`
	DeletedAt            *timestamp.Timestamp `protobuf:"bytes,14,opt,name=deleted_at,json=deletedAt,proto3" json:"deleted_at,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *Experience) Reset()         { *m = Experience{} }
func (m *Experience) String() string { return proto.CompactTextString(m) }
func (*Experience) ProtoMessage()    {}
func (*Experience) Descriptor() ([]byte, []int) {
	return fileDescriptor_8bd0ff5c490a2f61, []int{6}
}

func (m *Experience) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Experience.Unmarshal(m, b)
}
func (m *Experience) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Experience.Marshal(b, m, deterministic)
}
func (m *Experience) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Experience.Merge(m, src)
}
func (m *Experience) XXX_Size() int {
	return xxx_messageInfo_Experience.Size(m)
}
func (m *Experience) XXX_DiscardUnknown() {
	xxx_messageInfo_Experience.DiscardUnknown(m)
}

var xxx_messageInfo_Experience proto.InternalMessageInfo

func (m *Experience) GetId() int32 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *Experience) GetCompanyName() string {
	if m != nil {
		return m.CompanyName
	}
	return ""
}

func (m *Experience) GetTitle() string {
	if m != nil {
		return m.Title
	}
	return ""
}

func (m *Experience) GetEmploymentType() string {
	if m != nil {
		return m.EmploymentType
	}
	return ""
}

func (m *Experience) GetLocation() string {
	if m != nil {
		return m.Location
	}
	return ""
}

func (m *Experience) GetStartDate() *timestamp.Timestamp {
	if m != nil {
		return m.StartDate
	}
	return nil
}

func (m *Experience) GetEndDate() *timestamp.Timestamp {
	if m != nil {
		return m.EndDate
	}
	return nil
}

func (m *Experience) GetIsCurrentJob() bool {
	if m != nil {
		return m.IsCurrentJob
	}
	return false
}

func (m *Experience) GetHeadline() string {
	if m != nil {
		return m.Headline
	}
	return ""
}

func (m *Experience) GetDescription() string {
	if m != nil {
		return m.Description
	}
	return ""
}

func (m *Experience) GetMediaId() *Media {
	if m != nil {
		return m.MediaId
	}
	return nil
}

func (m *Experience) GetCreatedAt() *timestamp.Timestamp {
	if m != nil {
		return m.CreatedAt
	}
	return nil
}

func (m *Experience) GetUpdatedAt() *timestamp.Timestamp {
	if m != nil {
		return m.UpdatedAt
	}
	return nil
}

func (m *Experience) GetDeletedAt() *timestamp.Timestamp {
	if m != nil {
		return m.DeletedAt
	}
	return nil
}

type Investment struct {
	Id                   int32    `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	CompanyName          string   `protobuf:"bytes,2,opt,name=company_name,json=companyName,proto3" json:"company_name,omitempty"`
	Industry             string   `protobuf:"bytes,3,opt,name=industry,proto3" json:"industry,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Investment) Reset()         { *m = Investment{} }
func (m *Investment) String() string { return proto.CompactTextString(m) }
func (*Investment) ProtoMessage()    {}
func (*Investment) Descriptor() ([]byte, []int) {
	return fileDescriptor_8bd0ff5c490a2f61, []int{7}
}

func (m *Investment) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Investment.Unmarshal(m, b)
}
func (m *Investment) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Investment.Marshal(b, m, deterministic)
}
func (m *Investment) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Investment.Merge(m, src)
}
func (m *Investment) XXX_Size() int {
	return xxx_messageInfo_Investment.Size(m)
}
func (m *Investment) XXX_DiscardUnknown() {
	xxx_messageInfo_Investment.DiscardUnknown(m)
}

var xxx_messageInfo_Investment proto.InternalMessageInfo

func (m *Investment) GetId() int32 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *Investment) GetCompanyName() string {
	if m != nil {
		return m.CompanyName
	}
	return ""
}

func (m *Investment) GetIndustry() string {
	if m != nil {
		return m.Industry
	}
	return ""
}

func init() {
	proto.RegisterEnum("AccountType", AccountType_name, AccountType_value)
	proto.RegisterEnum("TeamAccountType", TeamAccountType_name, TeamAccountType_value)
	proto.RegisterType((*Address)(nil), "Address")
	proto.RegisterType((*Education)(nil), "Education")
	proto.RegisterType((*Media)(nil), "Media")
	proto.RegisterType((*Subscriptions)(nil), "Subscriptions")
	proto.RegisterType((*SocialMedia)(nil), "SocialMedia")
	proto.RegisterType((*Details)(nil), "Details")
	proto.RegisterType((*Experience)(nil), "Experience")
	proto.RegisterType((*Investment)(nil), "Investment")
}

func init() {
	proto.RegisterFile("schema/proto/model/common.proto", fileDescriptor_8bd0ff5c490a2f61)
}

var fileDescriptor_8bd0ff5c490a2f61 = []byte{
	// 1139 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xd4, 0x57, 0xdb, 0x6e, 0xdb, 0x46,
	0x13, 0x8e, 0xce, 0xd4, 0x50, 0x96, 0x95, 0xfd, 0x83, 0x1f, 0x8c, 0xdb, 0x34, 0x4a, 0x90, 0xa2,
	0x46, 0x82, 0x48, 0x40, 0x0f, 0x68, 0x93, 0x3b, 0xe5, 0x70, 0x91, 0xa2, 0x4d, 0x5a, 0x3a, 0xb9,
	0xe9, 0x0d, 0x41, 0x72, 0x47, 0xf2, 0x36, 0x24, 0x97, 0xe0, 0x2e, 0x9d, 0xa8, 0xcf, 0xd4, 0x17,
	0xb0, 0x81, 0xbe, 0x4a, 0x9f, 0xa3, 0xe8, 0x45, 0x51, 0xec, 0x81, 0x34, 0xad, 0xa2, 0xb0, 0x2c,
	0xf4, 0x26, 0x77, 0x9a, 0x6f, 0xbe, 0xa1, 0xa8, 0xfd, 0xbe, 0x99, 0x1d, 0xc1, 0x54, 0xc4, 0xc7,
	0x98, 0x86, 0xf3, 0xbc, 0xe0, 0x92, 0xcf, 0x53, 0x4e, 0x31, 0x11, 0xf3, 0x98, 0xa7, 0x29, 0xcf,
	0x66, 0x1a, 0x3b, 0x78, 0xbc, 0x62, 0xf2, 0xb8, 0x8c, 0x66, 0x31, 0x4f, 0xe7, 0x2c, 0x5b, 0xf2,
	0x28, 0xe1, 0xef, 0x79, 0x8e, 0x99, 0x29, 0x89, 0x1f, 0xae, 0x30, 0x7b, 0xb8, 0xe2, 0x45, 0x3a,
	0xe7, 0xb9, 0x64, 0x3c, 0x13, 0x73, 0x15, 0xd8, 0xda, 0xaf, 0x1b, 0xb5, 0x1a, 0x89, 0xca, 0xe5,
	0x5c, 0x14, 0xf1, 0x7c, 0xc5, 0xf9, 0x2a, 0xc1, 0x73, 0x4c, 0xb2, 0x14, 0x85, 0x0c, 0xd3, 0xdc,
	0x14, 0xde, 0xfd, 0xad, 0x03, 0x83, 0x05, 0xa5, 0x05, 0x0a, 0x41, 0xa6, 0xd0, 0x66, 0xd4, 0x6b,
	0x4d, 0x5b, 0x87, 0xbd, 0x27, 0x93, 0xb3, 0xd3, 0x9b, 0x23, 0x00, 0xd2, 0x17, 0x58, 0xb0, 0x30,
	0x39, 0x6c, 0xf9, 0x6d, 0x46, 0xc9, 0xc7, 0x30, 0x4c, 0x78, 0xb6, 0x62, 0xb2, 0xa4, 0xe8, 0xb5,
	0xa7, 0xad, 0xc3, 0xa1, 0x7f, 0x0e, 0x90, 0x03, 0x70, 0x92, 0x50, 0x9a, 0x64, 0x47, 0x27, 0xeb,
	0x98, 0x10, 0xe8, 0xc6, 0x4c, 0xae, 0xbd, 0xae, 0xc6, 0xf5, 0x67, 0x72, 0x03, 0x7a, 0x42, 0x86,
	0x12, 0xbd, 0x9e, 0x06, 0x4d, 0x40, 0x3c, 0x18, 0xc4, 0xbc, 0xcc, 0x64, 0xb1, 0xf6, 0xfa, 0x1a,
	0xaf, 0x42, 0x95, 0xf9, 0x85, 0xe5, 0x4f, 0x39, 0x45, 0x6f, 0x60, 0x32, 0x36, 0x24, 0xff, 0x87,
	0xbe, 0x90, 0x05, 0xa2, 0xf4, 0x1c, 0x9d, 0xb0, 0x11, 0xf9, 0x0c, 0xf6, 0xa3, 0x92, 0x25, 0x94,
	0x65, 0xab, 0x20, 0x2b, 0xd3, 0x08, 0x0b, 0x6f, 0xa8, 0x09, 0xe3, 0x0a, 0x7e, 0xa9, 0x51, 0xf2,
	0x08, 0x20, 0x2e, 0x30, 0x94, 0x48, 0x83, 0x50, 0x7a, 0x30, 0x6d, 0x1d, 0xba, 0x9f, 0x1f, 0xcc,
	0xcc, 0xe1, 0xcd, 0xaa, 0xc3, 0x9b, 0xbd, 0xae, 0x0e, 0xcf, 0x1f, 0x5a, 0xf6, 0x42, 0xaa, 0xd2,
	0x32, 0xa7, 0x55, 0xa9, 0x7b, 0x79, 0xa9, 0x65, 0x9b, 0x52, 0x8a, 0x09, 0xda, 0xd2, 0xd1, 0xe5,
	0xa5, 0x96, 0xbd, 0x90, 0x8f, 0xfb, 0x67, 0xa7, 0x37, 0xdb, 0x4e, 0xeb, 0xee, 0x5f, 0x5d, 0x18,
	0x3e, 0xa7, 0x65, 0x1c, 0x2a, 0x4b, 0x6c, 0xa1, 0xa0, 0x3a, 0xa9, 0xf8, 0x98, 0xf3, 0xc4, 0xca,
	0x67, 0x23, 0x85, 0x53, 0x5c, 0x15, 0x58, 0x29, 0x67, 0x23, 0x72, 0x0f, 0xc6, 0x4b, 0x86, 0x09,
	0x0d, 0xf8, 0x32, 0x10, 0xb2, 0xa4, 0x95, 0x82, 0x23, 0x8d, 0xbe, 0x5a, 0x1e, 0x29, 0x4c, 0xfd,
	0x10, 0x21, 0xc3, 0x42, 0x06, 0xb4, 0x92, 0xf3, 0x92, 0x1f, 0xa2, 0xd9, 0xcf, 0x94, 0xdc, 0x5f,
	0x81, 0x83, 0x19, 0x35, 0x85, 0xfd, 0x4b, 0x0b, 0x07, 0x98, 0x51, 0x5d, 0x36, 0x87, 0xff, 0xc5,
	0x65, 0x51, 0x60, 0x26, 0x93, 0x75, 0x10, 0x4a, 0x89, 0x99, 0x52, 0x53, 0xfb, 0xc2, 0xf1, 0x49,
	0x9d, 0x5a, 0x54, 0x19, 0x32, 0x81, 0xce, 0x2a, 0x0f, 0xb5, 0x3f, 0xda, 0xbe, 0xfa, 0x48, 0x3e,
	0x01, 0x08, 0x63, 0xc9, 0x4e, 0x98, 0x64, 0x28, 0xac, 0x2f, 0x1a, 0x88, 0x32, 0xbb, 0xe0, 0x31,
	0x43, 0x9d, 0x06, 0x63, 0xf6, 0x1a, 0x20, 0x53, 0x70, 0x29, 0x8a, 0xb8, 0x60, 0xba, 0x19, 0xb5,
	0xee, 0x43, 0xbf, 0x09, 0x6d, 0x78, 0x6a, 0xb4, 0xbb, 0xa7, 0xf6, 0x76, 0xf7, 0xd4, 0xf8, 0x0a,
	0x9e, 0x22, 0x0f, 0xc0, 0x49, 0x91, 0xb2, 0x30, 0x60, 0xd4, 0xdb, 0xd7, 0x85, 0xfd, 0xd9, 0xf7,
	0x0a, 0x78, 0xe2, 0x9c, 0x9d, 0xde, 0xec, 0x1e, 0xb4, 0xbf, 0x69, 0xf9, 0x03, 0xcd, 0x78, 0x41,
	0x6b, 0x03, 0xfe, 0xd9, 0x86, 0x9e, 0x26, 0x6d, 0x61, 0xbe, 0x4f, 0x61, 0x4c, 0x79, 0x5c, 0xa6,
	0x98, 0xc9, 0x20, 0x61, 0xd9, 0x5b, 0xe1, 0xb5, 0xa7, 0x9d, 0xc3, 0xa1, 0xbf, 0x57, 0xa1, 0xdf,
	0x29, 0x90, 0xdc, 0x06, 0x37, 0x3f, 0xe6, 0x92, 0x5b, 0x4e, 0x47, 0x73, 0x40, 0x43, 0x35, 0xe1,
	0x84, 0x51, 0xac, 0x08, 0x5d, 0x43, 0xd0, 0x90, 0x21, 0x3c, 0x04, 0x92, 0x17, 0x28, 0x30, 0x93,
	0xba, 0x2f, 0x2c, 0xaf, 0xa7, 0x79, 0xd7, 0x9b, 0x19, 0x43, 0xbf, 0xa8, 0x54, 0xff, 0x8a, 0x4a,
	0x35, 0x8e, 0x7b, 0x70, 0x95, 0xe3, 0xbe, 0x28, 0xb2, 0x73, 0x05, 0x91, 0xeb, 0xc3, 0xff, 0xa3,
	0x03, 0x7b, 0x47, 0x65, 0x54, 0x7b, 0xae, 0x39, 0xc3, 0xf7, 0xfe, 0x45, 0x84, 0x07, 0x70, 0x5d,
	0x34, 0x4a, 0x82, 0x2c, 0x4c, 0xab, 0x59, 0x3e, 0x69, 0x26, 0x5e, 0x86, 0xa9, 0x6e, 0xb3, 0x0b,
	0x64, 0x35, 0xa2, 0x4b, 0x61, 0x67, 0x04, 0x69, 0xa6, 0x8e, 0x74, 0x66, 0x63, 0x12, 0x74, 0x77,
	0x9d, 0x04, 0xbd, 0xed, 0x27, 0xc1, 0x6d, 0x70, 0xc3, 0x38, 0x46, 0x21, 0x02, 0xb9, 0xce, 0xd1,
	0xde, 0x19, 0x60, 0xa0, 0xd7, 0xeb, 0x1c, 0xc9, 0x47, 0x30, 0x64, 0x22, 0xd0, 0x8d, 0x8d, 0x76,
	0x40, 0x38, 0x4c, 0x2c, 0x74, 0xbc, 0x21, 0xbd, 0xb3, 0xbb, 0xf4, 0xc3, 0xdd, 0xa5, 0x87, 0x5d,
	0xa4, 0xff, 0xb5, 0x03, 0xee, 0x11, 0x8f, 0x59, 0x98, 0x6c, 0xdb, 0x7d, 0xb7, 0x00, 0xcc, 0x96,
	0x10, 0x94, 0x85, 0x19, 0xff, 0x7b, 0xfe, 0xd0, 0x20, 0x6f, 0x8a, 0x44, 0x9d, 0xe3, 0x3b, 0x8c,
	0x04, 0x93, 0xa8, 0xf3, 0x46, 0x62, 0xb0, 0x90, 0x22, 0xdc, 0x81, 0xd1, 0x32, 0x8c, 0x31, 0xe2,
	0xfc, 0xad, 0x66, 0x98, 0x8b, 0xc0, 0xad, 0x30, 0xfb, 0x0c, 0xf9, 0x8e, 0x49, 0x89, 0x85, 0x66,
	0x98, 0x7b, 0x1d, 0x2c, 0xa4, 0x08, 0xb7, 0x00, 0x54, 0x2f, 0x22, 0xd5, 0xf9, 0xbe, 0xdd, 0x20,
	0x34, 0x62, 0xeb, 0xd7, 0xbc, 0x94, 0x65, 0x64, 0xde, 0xc1, 0xdc, 0xf2, 0x60, 0x21, 0x45, 0xf8,
	0x60, 0xe5, 0xfa, 0x11, 0x06, 0xcf, 0x50, 0x86, 0x2c, 0xd1, 0xf7, 0xca, 0x8b, 0x1f, 0x5e, 0x99,
	0x7e, 0xd1, 0x82, 0x0d, 0xfd, 0x73, 0x40, 0xdd, 0x2b, 0x4f, 0x79, 0x9a, 0x87, 0xd9, 0x5a, 0x99,
	0xd7, 0x36, 0x66, 0x13, 0xaa, 0x1f, 0xf9, 0x7b, 0x17, 0xe0, 0xf9, 0xfb, 0x1c, 0x0b, 0x86, 0x59,
	0x8c, 0x5b, 0x18, 0xe0, 0x0e, 0x8c, 0x62, 0xf3, 0x9c, 0x66, 0xd3, 0xbb, 0x16, 0xd3, 0xfd, 0x7e,
	0x03, 0x7a, 0x92, 0xc9, 0xa4, 0xda, 0x02, 0x4c, 0xa0, 0xd6, 0x28, 0x4c, 0xf3, 0x84, 0xaf, 0xf5,
	0xe4, 0xd6, 0x6d, 0x66, 0xc4, 0x1f, 0x9f, 0xc3, 0xba, 0xd5, 0xd4, 0x06, 0xc8, 0xcd, 0x2e, 0x62,
	0xc5, 0xaf, 0xe3, 0x8d, 0xc9, 0xd0, 0xdf, 0x75, 0x32, 0x0c, 0xb6, 0x9f, 0x0c, 0xf7, 0x60, 0xcc,
	0x44, 0x60, 0x77, 0x81, 0xe0, 0x67, 0x1e, 0x69, 0xc3, 0x38, 0xfe, 0x88, 0x89, 0xa7, 0x06, 0xfc,
	0x96, 0x47, 0xea, 0x9d, 0x8f, 0x31, 0xa4, 0x09, 0xcb, 0xd0, 0x2e, 0x01, 0x75, 0xbc, 0x79, 0xc9,
	0xc3, 0x3f, 0x2f, 0xf9, 0x3b, 0x8d, 0x3b, 0xd3, 0x6d, 0xde, 0x99, 0xf5, 0x4d, 0xf9, 0xc1, 0xed,
	0x01, 0xb5, 0xc1, 0x4a, 0x80, 0x17, 0xd9, 0x09, 0x0a, 0xa9, 0xf4, 0xfd, 0x6f, 0xfc, 0x75, 0x00,
	0x0e, 0xcb, 0x68, 0x29, 0xd4, 0x76, 0x6f, 0xff, 0x22, 0x54, 0x71, 0xf5, 0xb5, 0xf7, 0x1f, 0x81,
	0xbb, 0x88, 0xf5, 0xce, 0xaf, 0x3d, 0xb5, 0x0f, 0xae, 0x8f, 0xab, 0x32, 0x09, 0x8b, 0x37, 0x02,
	0x8b, 0xc9, 0x35, 0xe2, 0xc2, 0xe0, 0x48, 0x59, 0xa3, 0xcc, 0x27, 0x2d, 0x32, 0x02, 0xc7, 0xbc,
	0x23, 0x2f, 0x26, 0xed, 0xfb, 0x5f, 0xc2, 0xfe, 0x6b, 0x0c, 0xd3, 0x8d, 0x72, 0xcb, 0x56, 0x99,
	0xc9, 0x35, 0x32, 0x81, 0x51, 0x55, 0xa1, 0x91, 0xd6, 0x93, 0xc1, 0x4f, 0x3d, 0xfd, 0x7f, 0x2c,
	0xea, 0xeb, 0x73, 0xf9, 0xe2, 0xef, 0x00, 0x00, 0x00, 0xff, 0xff, 0x23, 0x43, 0x7c, 0x87, 0xac,
	0x0d, 0x00, 0x00,
}
