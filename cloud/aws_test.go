package cloud

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/stretchr/testify/assert"
)

var (
	testImageId        = "ami-12345"
	testInstanceType   = "t2.testtype"
	testSecurityGroups = []string{
		"sg-12345",
		"sg-67890",
	}
	testSubnetId = "subnet-12345"
	testTagKeys  = []string{
		"key-1",
		"key-2",
	}
	testTagValues = []string{
		"value-1",
		"value-2",
	}
	testTags = map[string]string{
		testTagKeys[0]: testTagValues[0],
		testTagKeys[1]: testTagValues[1],
	}
)

func Test_constructInput(t *testing.T) {
	tests := []struct {
		name string
		opts []InstanceOption
		want *ec2.RunInstancesInput
	}{
		{
			"no constructors returns default",
			nil,
			&ec2.RunInstancesInput{
				MaxCount:     aws.Int64(1),
				MinCount:     aws.Int64(1),
				ImageId:      aws.String(DefaultImageId),
				InstanceType: aws.String(DefaultInstanceType),
			},
		},
		{
			"WithImageId constructor",
			[]InstanceOption{WithImageId(testImageId)},
			&ec2.RunInstancesInput{
				MaxCount:     aws.Int64(1),
				MinCount:     aws.Int64(1),
				ImageId:      aws.String(testImageId),
				InstanceType: aws.String(DefaultInstanceType),
			},
		},
		{
			"WithInstanceType constructor",
			[]InstanceOption{WithInstanceType(testInstanceType)},
			&ec2.RunInstancesInput{
				MaxCount:     aws.Int64(1),
				MinCount:     aws.Int64(1),
				ImageId:      aws.String(DefaultImageId),
				InstanceType: aws.String(testInstanceType),
			},
		},
		{
			"WithSecurityGroups single entry constructor",
			[]InstanceOption{WithSecurityGroups(testSecurityGroups[0])},
			&ec2.RunInstancesInput{
				MaxCount:     aws.Int64(1),
				MinCount:     aws.Int64(1),
				ImageId:      aws.String(DefaultImageId),
				InstanceType: aws.String(DefaultInstanceType),
				SecurityGroupIds: []*string{
					aws.String(testSecurityGroups[0]),
				},
			},
		},
		{
			"WithSecurityGroups multi entry constructor",
			[]InstanceOption{WithSecurityGroups(testSecurityGroups...)},
			&ec2.RunInstancesInput{
				MaxCount:     aws.Int64(1),
				MinCount:     aws.Int64(1),
				ImageId:      aws.String(DefaultImageId),
				InstanceType: aws.String(DefaultInstanceType),
				SecurityGroupIds: []*string{
					aws.String(testSecurityGroups[0]),
					aws.String(testSecurityGroups[1]),
				},
			},
		},
		{
			"WithSubnetId constructor",
			[]InstanceOption{WithSubnetId(testSubnetId)},
			&ec2.RunInstancesInput{
				MaxCount:     aws.Int64(1),
				MinCount:     aws.Int64(1),
				ImageId:      aws.String(DefaultImageId),
				InstanceType: aws.String(DefaultInstanceType),
				SubnetId:     aws.String(testSubnetId),
			},
		},
		{
			"WithTags constructor",
			[]InstanceOption{WithTags(testTags)},
			&ec2.RunInstancesInput{
				MaxCount:     aws.Int64(1),
				MinCount:     aws.Int64(1),
				ImageId:      aws.String(DefaultImageId),
				InstanceType: aws.String(DefaultInstanceType),
				TagSpecifications: []*ec2.TagSpecification{
					{
						ResourceType: aws.String("instance"),
						Tags: []*ec2.Tag{
							{
								Key:   aws.String(testTagKeys[0]),
								Value: aws.String(testTagValues[0]),
							},
							{
								Key:   aws.String(testTagKeys[1]),
								Value: aws.String(testTagValues[1]),
							},
						},
					},
				},
			},
		},
		{
			"WithImageId + WithTags constructor",
			[]InstanceOption{
				WithImageId(testImageId),
				WithTags(testTags),
			},
			&ec2.RunInstancesInput{
				MaxCount:     aws.Int64(1),
				MinCount:     aws.Int64(1),
				ImageId:      aws.String(testImageId),
				InstanceType: aws.String(DefaultInstanceType),
				TagSpecifications: []*ec2.TagSpecification{
					{
						ResourceType: aws.String("instance"),
						Tags: []*ec2.Tag{
							{
								Key:   aws.String(testTagKeys[0]),
								Value: aws.String(testTagValues[0]),
							},
							{
								Key:   aws.String(testTagKeys[1]),
								Value: aws.String(testTagValues[1]),
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := constructInput(tt.opts)
			assert.Equal(t, got, tt.want)
		})
	}
}

type ec2Mock struct {
	ec2iface.EC2API
	Reservation ec2.Reservation
}

type ec2ErrMock struct {
	ec2iface.EC2API
	Reservation ec2.Reservation
}

func (c ec2Mock) RunInstances(*ec2.RunInstancesInput) (*ec2.Reservation, error) {
	return &c.Reservation, nil
}

func (c ec2ErrMock) RunInstances(*ec2.RunInstancesInput) (*ec2.Reservation, error) {
	return nil, fmt.Errorf("mock error")
}

var (
	testSimpleInstanceId = "i-abcdefg"
	testReservations     = map[string]ec2.Reservation{
		"simple": {
			Instances: []*ec2.Instance{
				{
					InstanceId: &testSimpleInstanceId,
				},
			},
		},
	}
)

// TODO: add more test cases
func TestCreateInstance(t *testing.T) {
	tests := []struct {
		name    string
		client  *client
		opts    []InstanceOption
		want    *Instance
		wantErr bool
	}{
		{
			"client error returns an error",
			&client{ec2ErrMock{}},
			nil,
			nil,
			true,
		},
		{
			"empty reservation returns an error",
			&client{ec2Mock{}},
			nil,
			nil,
			true,
		},
		{
			"simple instance with id",
			&client{ec2Mock{Reservation: testReservations["simple"]}},
			nil,
			&Instance{Id: testSimpleInstanceId},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.client.CreateInstance(tt.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("client.CreateInstance() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, got, tt.want)
		})
	}
}
