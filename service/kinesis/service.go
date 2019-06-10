// Code generated by private/model/cli/gen-api/main.go. DO NOT EDIT.

package kinesis

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/client/metadata"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/signer/v4"
	"github.com/aws/aws-sdk-go/private/protocol/jsonrpc"
)

// Kinesis provides the API operation methods for making requests to
// Amazon Kinesis. See this package's package overview docs
// for details on the service.
//
// Kinesis methods are safe to use concurrently. It is not safe to
// modify mutate any of the struct's properties though.
type Kinesis struct {
	*client.Client
}

// Used for custom client initialization logic
var initClient func(*client.Client)

// Used for custom request initialization logic
var initRequest func(*request.Request)

// Magic File Header for a KPL Aggregated Record
var KplMagicHeader = fmt.Sprintf("%q", []byte("\xf3\x89\x9a\xc2"))

// Service information constants
const (
	ServiceName = "kinesis"   // Name of service.
	EndpointsID = ServiceName // ID to lookup a service endpoint with.
	ServiceID   = "Kinesis"   // ServiceID is a unique identifer of a specific service.
	Md5Buffer = 15
)

// New creates a new instance of the Kinesis client with a session.
// If additional configuration is needed for the client instance use the optional
// aws.Config parameter to add your extra config.
//
// Example:
//     // Create a Kinesis client from just a session.
//     svc := kinesis.New(mySession)
//
//     // Create a Kinesis client with additional configuration
//     svc := kinesis.New(mySession, aws.NewConfig().WithRegion("us-west-2"))
func New(p client.ConfigProvider, cfgs ...*aws.Config) *Kinesis {
	c := p.ClientConfig(EndpointsID, cfgs...)
	return newClient(*c.Config, c.Handlers, c.Endpoint, c.SigningRegion, c.SigningName)
}

// newClient creates, initializes and returns a new service client instance.
func newClient(cfg aws.Config, handlers request.Handlers, endpoint, signingRegion, signingName string) *Kinesis {
	svc := &Kinesis{
		Client: client.New(
			cfg,
			metadata.ClientInfo{
				ServiceName:   ServiceName,
				ServiceID:     ServiceID,
				SigningName:   signingName,
				SigningRegion: signingRegion,
				Endpoint:      endpoint,
				APIVersion:    "2013-12-02",
				JSONVersion:   "1.1",
				TargetPrefix:  "Kinesis_20131202",
			},
			handlers,
		),
	}

	// Handlers
	svc.Handlers.Sign.PushBackNamed(v4.SignRequestHandler)
	svc.Handlers.Build.PushBackNamed(jsonrpc.BuildHandler)
	svc.Handlers.Unmarshal.PushBackNamed(jsonrpc.UnmarshalHandler)
	svc.Handlers.UnmarshalMeta.PushBackNamed(jsonrpc.UnmarshalMetaHandler)
	svc.Handlers.UnmarshalError.PushBackNamed(jsonrpc.UnmarshalErrorHandler)

	svc.Handlers.UnmarshalStream.PushBackNamed(jsonrpc.UnmarshalHandler)

	// Run custom client initialization if present
	if initClient != nil {
		initClient(svc.Client)
	}

	return svc
}

// newRequest creates a new request for a Kinesis operation and runs any
// custom request initialization.
func (c *Kinesis) newRequest(op *request.Operation, params, data interface{}) *request.Request {
	req := c.NewRequest(op, params, data)

	// Run custom request initialization if present
	if initRequest != nil {
		initRequest(req)
	}

	return req
}

// getProtoRecords takes an array of Kinesis records and expands any Protobuf
// records within that array, returning an array of all records
func getProtoRecords(records []*Record) ([]*Record, error) {
  allRecords := make([]*Record, 0)
	for _, record := range records {
		header := fmt.Sprintf("%q", record.Data[:4])
		if header == KplMagicHeader {
			protoRecords, err := expandProtoRecord(record)
			if err != nil {
				return nil, err
			}
			allRecords = append(allRecords, protoRecords...)
		} else {
			allRecords = append(allRecords, record)
		}
	}
	return allRecords, nil
}

// expandProtoRecord expands a Protobuf record into an array of the records
// contained by the given Protobuf record
func expandProtoRecord(record *Record) ([]*Record, error) {
	expandedRecords := make([]*Record, 0)
	msg := record.Data[4:len(record.Data)-1-Md5Buffer]
	aggRecord := &AggregatedRecord{}
	err := proto.Unmarshal(msg, aggRecord)

	if err != nil {
		return expandedRecords, err
	}

	for _, aggrec := range aggRecord.Records {
		r := &Record{
			ApproximateArrivalTimestamp: record.ApproximateArrivalTimestamp,
			Data: aggrec.Data,
			EncryptionType: record.EncryptionType,
			PartitionKey: record.PartitionKey,
			SequenceNumber: record.SequenceNumber,
		}
		expandedRecords = append(expandedRecords, r)
	}

	return expandedRecords, nil
}
