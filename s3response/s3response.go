// Copyright 2023 Versity Software
// This file is licensed under the Apache License, Version 2.0
// (the "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package s3response

import (
	"encoding/xml"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/versity/versitygw/s3err"
)

const (
	iso8601TimeFormat         = "2006-01-02T15:04:05.000Z"
	iso8601TimeFormatExtended = "2006-01-02T15:04:05.000000Z"
	iso8601TimeFormatWithTZ   = "2006-01-02T15:04:05-0700"
)

type PutObjectOutput struct {
	ETag           string
	VersionID      string
	ChecksumCRC32  *string
	ChecksumCRC32C *string
	ChecksumSHA1   *string
	ChecksumSHA256 *string
}

// Part describes part metadata.
type Part struct {
	PartNumber     int
	LastModified   time.Time
	ETag           string
	Size           int64
	ChecksumCRC32  *string
	ChecksumCRC32C *string
	ChecksumSHA1   *string
	ChecksumSHA256 *string
}

func (p Part) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	type Alias Part
	aux := &struct {
		LastModified string `xml:"LastModified"`
		*Alias
	}{
		Alias: (*Alias)(&p),
	}

	aux.LastModified = p.LastModified.UTC().Format(iso8601TimeFormat)

	return e.EncodeElement(aux, start)
}

// ListPartsResponse - s3 api list parts response.
type ListPartsResult struct {
	XMLName xml.Name `xml:"http://s3.amazonaws.com/doc/2006-03-01/ ListPartsResult" json:"-"`

	Bucket            string
	Key               string
	UploadID          string `xml:"UploadId"`
	ChecksumAlgorithm types.ChecksumAlgorithm

	Initiator Initiator
	Owner     Owner

	// The class of storage used to store the object.
	StorageClass types.StorageClass

	PartNumberMarker     int
	NextPartNumberMarker int
	MaxParts             int
	IsTruncated          bool

	// List of parts.
	Parts []Part `xml:"Part"`
}

type ObjectAttributes string

const (
	ObjectAttributesEtag         ObjectAttributes = "ETag"
	ObjectAttributesChecksum     ObjectAttributes = "Checksum"
	ObjectAttributesObjectParts  ObjectAttributes = "ObjectParts"
	ObjectAttributesStorageClass ObjectAttributes = "StorageClass"
	ObjectAttributesObjectSize   ObjectAttributes = "ObjectSize"
)

func (o ObjectAttributes) IsValid() bool {
	return o == ObjectAttributesChecksum ||
		o == ObjectAttributesEtag ||
		o == ObjectAttributesObjectParts ||
		o == ObjectAttributesObjectSize ||
		o == ObjectAttributesStorageClass
}

type GetObjectAttributesResponse struct {
	ETag         *string
	ObjectSize   *int64
	StorageClass types.StorageClass `xml:",omitempty"`
	ObjectParts  *ObjectParts
	Checksum     *types.Checksum

	// Not included in the response body
	VersionId    *string
	LastModified *time.Time
	DeleteMarker *bool
}

type ObjectParts struct {
	PartNumberMarker     int
	NextPartNumberMarker int
	MaxParts             int
	IsTruncated          bool
	Parts                []types.ObjectPart `xml:"Part"`
}

// ListMultipartUploadsResponse - s3 api list multipart uploads response.
type ListMultipartUploadsResult struct {
	XMLName xml.Name `xml:"http://s3.amazonaws.com/doc/2006-03-01/ ListMultipartUploadsResult" json:"-"`

	Bucket             string
	KeyMarker          string
	UploadIDMarker     string `xml:"UploadIdMarker"`
	NextKeyMarker      string
	NextUploadIDMarker string `xml:"NextUploadIdMarker"`
	Delimiter          string
	Prefix             string
	EncodingType       string `xml:"EncodingType,omitempty"`
	MaxUploads         int
	IsTruncated        bool

	// List of pending uploads.
	Uploads []Upload `xml:"Upload"`

	// Delimed common prefixes.
	CommonPrefixes []CommonPrefix
}

type ListObjectsResult struct {
	XMLName        xml.Name `xml:"http://s3.amazonaws.com/doc/2006-03-01/ ListBucketResult" json:"-"`
	Name           *string
	Prefix         *string
	Marker         *string
	NextMarker     *string
	MaxKeys        *int32
	Delimiter      *string
	IsTruncated    *bool
	Contents       []Object
	CommonPrefixes []types.CommonPrefix
	EncodingType   types.EncodingType
}

type ListObjectsV2Result struct {
	XMLName               xml.Name `xml:"http://s3.amazonaws.com/doc/2006-03-01/ ListBucketResult" json:"-"`
	Name                  *string
	Prefix                *string
	StartAfter            *string
	ContinuationToken     *string
	NextContinuationToken *string
	KeyCount              *int32
	MaxKeys               *int32
	Delimiter             *string
	IsTruncated           *bool
	Contents              []Object
	CommonPrefixes        []types.CommonPrefix
	EncodingType          types.EncodingType
}

type Object struct {
	ChecksumAlgorithm []types.ChecksumAlgorithm
	ETag              *string
	Key               *string
	LastModified      *time.Time
	Owner             *types.Owner
	RestoreStatus     *types.RestoreStatus
	Size              *int64
	StorageClass      types.ObjectStorageClass
}

func (o Object) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	type Alias Object
	aux := &struct {
		LastModified *string `xml:"LastModified,omitempty"`
		*Alias
	}{
		Alias: (*Alias)(&o),
	}

	if o.LastModified != nil {
		formattedTime := o.LastModified.UTC().Format(iso8601TimeFormat)
		aux.LastModified = &formattedTime
	}

	return e.EncodeElement(aux, start)
}

// Upload describes in progress multipart upload
type Upload struct {
	Key               string
	UploadID          string `xml:"UploadId"`
	Initiator         Initiator
	Owner             Owner
	StorageClass      types.StorageClass
	Initiated         time.Time
	ChecksumAlgorithm types.ChecksumAlgorithm
}

func (u Upload) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	type Alias Upload
	aux := &struct {
		Initiated string `xml:"Initiated"`
		*Alias
	}{
		Alias: (*Alias)(&u),
	}

	aux.Initiated = u.Initiated.UTC().Format(iso8601TimeFormat)

	return e.EncodeElement(aux, start)
}

// CommonPrefix ListObjectsResponse common prefixes (directory abstraction)
type CommonPrefix struct {
	Prefix string
}

// Initiator same fields as Owner
type Initiator Owner

// Owner bucket ownership
type Owner struct {
	ID          string
	DisplayName string
}

type Tag struct {
	Key   string `xml:"Key"`
	Value string `xml:"Value"`
}

type TagSet struct {
	Tags []Tag `xml:"Tag"`
}

type Tagging struct {
	XMLName xml.Name `xml:"http://s3.amazonaws.com/doc/2006-03-01/ Tagging" json:"-"`
	TagSet  TagSet   `xml:"TagSet"`
}

type TaggingInput struct {
	TagSet TagSet `xml:"TagSet"`
}

type DeleteObjects struct {
	Objects []types.ObjectIdentifier `xml:"Object"`
}

type DeleteResult struct {
	Deleted []types.DeletedObject
	Error   []types.Error
}
type SelectObjectContentPayload struct {
	Expression          *string
	ExpressionType      types.ExpressionType
	RequestProgress     *types.RequestProgress
	InputSerialization  *types.InputSerialization
	OutputSerialization *types.OutputSerialization
	ScanRange           *types.ScanRange
}

type SelectObjectContentResult struct {
	Records  *types.RecordsEvent
	Stats    *types.StatsEvent
	Progress *types.ProgressEvent
	Cont     *types.ContinuationEvent
	End      *types.EndEvent
}

type Bucket struct {
	Name  string `json:"name"`
	Owner string `json:"owner"`
}

type ListBucketsInput struct {
	Owner             string
	IsAdmin           bool
	ContinuationToken string
	Prefix            string
	MaxBuckets        int32
}

type ListAllMyBucketsResult struct {
	XMLName           xml.Name `xml:"http://s3.amazonaws.com/doc/2006-03-01/ ListAllMyBucketsResult" json:"-"`
	Owner             CanonicalUser
	Buckets           ListAllMyBucketsList
	ContinuationToken string `xml:"ContinuationToken,omitempty"`
	Prefix            string `xml:"Prefix,omitempty"`
}

type ListAllMyBucketsEntry struct {
	Name         string
	CreationDate time.Time
}

func (r ListAllMyBucketsEntry) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	type Alias ListAllMyBucketsEntry
	aux := &struct {
		CreationDate string `xml:"CreationDate"`
		*Alias
	}{
		Alias: (*Alias)(&r),
	}

	aux.CreationDate = r.CreationDate.UTC().Format(iso8601TimeFormat)

	return e.EncodeElement(aux, start)
}

type ListAllMyBucketsList struct {
	Bucket []ListAllMyBucketsEntry
}

type CanonicalUser struct {
	ID          string
	DisplayName string
}

type CopyObjectResult struct {
	XMLName             xml.Name `xml:"http://s3.amazonaws.com/doc/2006-03-01/ CopyObjectResult" json:"-"`
	LastModified        time.Time
	ETag                string
	CopySourceVersionId string `xml:"-"`
}

type CopyPartResult struct {
	XMLName        xml.Name `xml:"http://s3.amazonaws.com/doc/2006-03-01/ CopyPartResult" json:"-"`
	LastModified   time.Time
	ETag           *string
	ChecksumCRC32  *string
	ChecksumCRC32C *string
	ChecksumSHA1   *string
	ChecksumSHA256 *string

	// not included in the body
	CopySourceVersionId string `xml:"-"`
}

func (r CopyObjectResult) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	type Alias CopyObjectResult
	aux := &struct {
		LastModified string `xml:"LastModified"`
		*Alias
	}{
		Alias: (*Alias)(&r),
	}

	aux.LastModified = r.LastModified.UTC().Format(iso8601TimeFormat)

	return e.EncodeElement(aux, start)
}

type AccessControlPolicy struct {
	XMLName           xml.Name `xml:"http://s3.amazonaws.com/doc/2006-03-01/ AccessControlPolicy" json:"-"`
	Owner             CanonicalUser
	AccessControlList AccessControlList
}

type AccessControlList struct {
	Grant []Grant
}

type Grant struct {
	Grantee    Grantee
	Permission string
}

// Set the following to encode correctly:
//
//	Grantee: s3response.Grantee{
//		Xsi:         "http://www.w3.org/2001/XMLSchema-instance",
//		Type:        "CanonicalUser",
//	},
type Grantee struct {
	XMLName     xml.Name `xml:"Grantee"`
	Xsi         string   `xml:"xmlns:xsi,attr,omitempty"`
	Type        string   `xml:"xsi:type,attr,omitempty"`
	ID          string
	DisplayName string
}

type OwnershipControls struct {
	Rules []types.OwnershipControlsRule `xml:"Rule"`
}

type InitiateMultipartUploadResult struct {
	XMLName  xml.Name `xml:"http://s3.amazonaws.com/doc/2006-03-01/ InitiateMultipartUploadResult" json:"-"`
	Bucket   string
	Key      string
	UploadId string
}

type ListVersionsResult struct {
	XMLName             xml.Name `xml:"http://s3.amazonaws.com/doc/2006-03-01/ ListVersionsResult" json:"-"`
	CommonPrefixes      []types.CommonPrefix
	DeleteMarkers       []types.DeleteMarkerEntry `xml:"DeleteMarker"`
	Delimiter           *string
	EncodingType        types.EncodingType
	IsTruncated         *bool
	KeyMarker           *string
	MaxKeys             *int32
	Name                *string
	NextKeyMarker       *string
	NextVersionIdMarker *string
	Prefix              *string
	VersionIdMarker     *string
	Versions            []types.ObjectVersion `xml:"Version"`
}

type GetBucketVersioningOutput struct {
	XMLName   xml.Name `xml:"http://s3.amazonaws.com/doc/2006-03-01/ VersioningConfiguration" json:"-"`
	MFADelete *types.MFADeleteStatus
	Status    *types.BucketVersioningStatus
}

type PutObjectRetentionInput struct {
	XMLName         xml.Name `xml:"Retention"`
	Mode            types.ObjectLockRetentionMode
	RetainUntilDate AmzDate
}

type AmzDate struct {
	time.Time
}

// Parses the date from xml string and validates for predefined date formats
func (d *AmzDate) UnmarshalXML(e *xml.Decoder, startElement xml.StartElement) error {
	var dateStr string
	err := e.DecodeElement(&dateStr, &startElement)
	if err != nil {
		return err
	}

	retDate, err := d.ISO8601Parse(dateStr)
	if err != nil {
		return s3err.GetAPIError(s3err.ErrInvalidRequest)
	}

	*d = AmzDate{retDate}
	return nil
}

// Encodes expiration date if it is non-zero
// Encodes empty string if it's zero
func (d AmzDate) MarshalXML(e *xml.Encoder, startElement xml.StartElement) error {
	if d.IsZero() {
		return nil
	}
	return e.EncodeElement(d.UTC().Format(iso8601TimeFormat), startElement)
}

// Parses ISO8601 date string to time.Time by
// validating different time layouts
func (AmzDate) ISO8601Parse(date string) (t time.Time, err error) {
	for _, layout := range []string{
		iso8601TimeFormat,
		iso8601TimeFormatExtended,
		iso8601TimeFormatWithTZ,
		time.RFC3339,
	} {
		t, err = time.Parse(layout, date)
		if err == nil {
			return t, nil
		}
	}

	return t, err
}

// Admin api response types
type ListBucketsResult struct {
	Buckets []Bucket
}

type Checksum struct {
	Algorithms []types.ChecksumAlgorithm

	CRC32  *string
	CRC32C *string
	SHA1   *string
	SHA256 *string
}
