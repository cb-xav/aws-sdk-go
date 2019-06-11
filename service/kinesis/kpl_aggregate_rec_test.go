package kinesis_test

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	rec "github.com/cb-xav/aws-sdk-go/service/kinesis/records"
	"testing"
)

func generateAggregateRecord() []byte { 

	aggr := &rec.AggregatedRecord{}

	numRecords := 10
	for i := 0; i < numRecords; i++ {
		var partKey uint64
		var hashKey uint64
		partKey = uint64(i)
		hashKey = uint64(i) * uint64(10)
		r := &rec.Record{
			PartitionKeyIndex: &partKey,
			ExplicitHashKeyIndex: &hashKey,
			Data: []byte("OMG SOME TEST DATA"),
			Tags: make([]*rec.Tag, 0),
		}

		aggr.Records = append(aggr.Records, r)
	}


	byteArr, _ := proto.Marshal(aggr)
	return byteArr
}

func TestGettingAggregateRecord(t *testing.T) { 
	// Check to make sure the generation works
	byteArr := generateAggregateRecord()

	// print out the byte version
	fmt.Printf("The Protobuf Data = `%q`", byteArr)
}

func TestSmallLength(t *testing.T) {

	//smallByte := []byte("No")
	
	

}

func TestManipulation(t *testing.T) {
	
}
