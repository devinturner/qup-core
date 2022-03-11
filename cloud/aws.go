package cloud

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
)

const (
	DefaultImageId      = "ami-0fb653ca2d3203ac1"
	DefaultInstanceType = "t2.micro"
)

type client struct {
	ec2iface.EC2API
}

type Instance struct {
	Id string
}

// Create an ec2 instance. You can specify options by passing in an arbitrary number of InstanceOption functions.
func (c *client) CreateInstance(opts ...InstanceOption) (*Instance, error) {
	input := constructInput(opts)
	res, err := c.RunInstances(input)
	if err != nil {
		return nil, err
	}

	if l := len(res.Instances); l != 1 {
		return nil, fmt.Errorf("invalid number of instances in reservation: %d", l)
	}

	instance := new(Instance)
	instance.Id = *res.Instances[0].InstanceId

	return instance, nil
}

var WithDefaultImage InstanceOption = WithImageId(DefaultImageId)
var WithDefaultInstanceType InstanceOption = WithInstanceType(DefaultInstanceType)

// construct the input from a slice of InstanceOptions and set defaults.
func constructInput(opts []InstanceOption) *ec2.RunInstancesInput {
	out := &ec2.RunInstancesInput{
		MaxCount: aws.Int64(1),
		MinCount: aws.Int64(1),
	}

	for _, opt := range opts {
		opt(out)
	}

	if out.ImageId == nil {
		WithDefaultImage(out)
	}

	if out.InstanceType == nil {
		WithDefaultInstanceType(out)
	}
	return out
}

type InstanceOption func(*ec2.RunInstancesInput)

// Option to specify an AMI
func WithImageId(imageId string) InstanceOption {
	return func(rii *ec2.RunInstancesInput) {
		rii.ImageId = aws.String(imageId)
	}
}

// Option to specify an instance type
func WithInstanceType(instanceType string) InstanceOption {
	return func(rii *ec2.RunInstancesInput) {
		rii.InstanceType = aws.String(instanceType)
	}
}

// Option to specify a list of security group ids
func WithSecurityGroups(securityGroupIds ...string) InstanceOption {
	sgs := make([]*string, 0)
	for _, sg := range securityGroupIds {
		sgs = append(sgs, aws.String(sg))
	}

	return func(rii *ec2.RunInstancesInput) {
		rii.SecurityGroupIds = sgs
	}
}

// Option to specify a subnet
func WithSubnetId(subnetId string) InstanceOption {
	return func(rii *ec2.RunInstancesInput) {
		rii.SubnetId = aws.String(subnetId)
	}
}

// Option to add custom tags
func WithTags(tags map[string]string) InstanceOption {
	ts := make([]*ec2.Tag, 0)
	for k, v := range tags {
		ts = append(ts, &ec2.Tag{
			Key:   aws.String(k),
			Value: aws.String(v),
		})
	}

	return func(rii *ec2.RunInstancesInput) {
		rii.TagSpecifications = []*ec2.TagSpecification{
			{
				ResourceType: aws.String("instance"),
				Tags:         ts,
			},
		}
	}
}

// TODO: add option to customize region
