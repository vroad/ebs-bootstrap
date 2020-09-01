package main

import (
	"flag"
	"math"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

var (
	useEBS               bool
	ebsVolumeID          string
	mountPoint           string
	blockDevice          string
	awsRegion            string
	fileSystemFormatType string
	maxAttempts          int
)

func init() {
	flag.StringVar(&awsRegion, "aws-region", "eu-west-1", "AWS region this instance is on")
	flag.StringVar(&ebsVolumeID, "ebs-volume-id", "", "EBS volume to attach to this node")
	flag.StringVar(&mountPoint, "mount-point", "/data", "EBS volume mount point")
	flag.StringVar(&blockDevice, "block-device", "/dev/xvdf", "Block device to attach as")
	flag.StringVar(&fileSystemFormatType, "filesystem-type", "ext4", "Linux filesystem format type")
	flag.BoolVar(&useEBS, "use-ebs", true, "Use EBS instead of instance store")
	flag.IntVar(&maxAttempts, "max-attempts", math.MaxInt64, "Retry attachment until it fails n times")
	flag.Parse()
}

func main() {
	// Initialize AWS session
	awsSession := session.Must(session.NewSession())

	// Create ec2 and metadata svc clients with specified region
	ec2SVC := ec2.New(awsSession, aws.NewConfig().WithRegion(awsRegion))
	metadataSVC := ec2metadata.New(awsSession, aws.NewConfig().WithRegion(awsRegion))

	// obtain current AZ, required for finding volume
	availabilityZone, err := metadataSVC.GetMetadata("placement/availability-zone")
	if err != nil {
		panic(err)
	}

	if useEBS {
		volume, err := volumeFromID(ec2SVC, ebsVolumeID, availabilityZone)
		if err != nil {
			panic(err)
		}

		instanceID, err := metadataSVC.GetMetadata("instance-id")
		if err != nil {
			panic(err)
		}

		err = attachVolume(ec2SVC, instanceID, volume, maxAttempts)
		if err != nil {
			panic(err)
		}
	}

	if err := ensureVolumeInited(blockDevice, fileSystemFormatType); err != nil {
		panic(err)
	}

	if err := ensureVolumeMounted(blockDevice, mountPoint); err != nil {
		panic(err)
	}
}
