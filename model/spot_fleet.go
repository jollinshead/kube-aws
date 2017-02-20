package model

import (
	"fmt"
	"strconv"
)

// UnitEbsVolumeSize/IOPS are used for spot fleets instead of WorkerEbsVolumeSize/IOPS,
// so that we can make them clearer that they are not default size/iops for each worker node but "size/iops per unit"
// as their names suggest
type SpotFleet struct {
	TargetCapacity       int                   `yaml:"targetCapacity,omitempty"`
	SpotPrice            string                `yaml:"spotPrice,omitempty"`
	IAMFleetRoleARN      string                `yaml:"iamFleetRoleArn,omitempty"`
	EbsVolumeType       string                `yaml:"ebsVolumeType"`
	UnitEbsVolumeSize   int                   `yaml:"unitEbsVolumeSize"`
	UnitEbsVolumeIOPS   int                   `yaml:"unitEbsVolumeIOPS"`
	LaunchSpecifications []LaunchSpecification `yaml:"launchSpecifications,omitempty"`
}

func (f SpotFleet) Enabled() bool {
	return f.TargetCapacity > 0
}

func (c SpotFleet) Valid() error {
	for i, spec := range c.LaunchSpecifications {
		if err := spec.Valid(); err != nil {
			return fmt.Errorf("invalid launchSpecification at index %d: %v", i, err)
		}
	}
	return nil
}

func (f SpotFleet) WithDefaults() SpotFleet {
	defaults := newDefaultSpotFleet()

	if f.SpotPrice == "" {
		f.SpotPrice = defaults.SpotPrice
	}

	if f.UnitEbsVolumeSize == 0 {
		f.UnitEbsVolumeSize = defaults.UnitEbsVolumeSize
	}

	if f.UnitEbsVolumeIOPS == 0 {
		f.UnitEbsVolumeIOPS = defaults.UnitEbsVolumeIOPS
	}

	if f.EbsVolumeType == "" {
		f.EbsVolumeType = defaults.EbsVolumeType
	}

	if len(f.LaunchSpecifications) == 0 {
		f.LaunchSpecifications = defaults.LaunchSpecifications
	}

	launchSpecs := []LaunchSpecification{}
	for _, spec := range f.LaunchSpecifications {
		if spec.SpotPrice == "" {
			p, err := strconv.ParseFloat(f.SpotPrice, 64)
			if err != nil {
				panic(fmt.Errorf(`failed to parse float from spotPrice "%s" in %+v: %v`, f.SpotPrice, f, err))
			}
			spec.SpotPrice = strconv.FormatFloat(p*float64(spec.WeightedCapacity), 'f', -1, 64)
		}
		if spec.EbsVolumeType == "" {
			spec.EbsVolumeType = f.EbsVolumeType
		}
		if spec.EbsVolumeSize == 0 {
			spec.EbsVolumeSize = f.UnitEbsVolumeSize * spec.WeightedCapacity
		}
		if spec.EbsVolumeType == "io1" && spec.EbsVolumeIOPS == 0 {
			spec.EbsVolumeIOPS = f.UnitEbsVolumeIOPS * spec.WeightedCapacity
		}
		launchSpecs = append(launchSpecs, spec)
	}
	f.LaunchSpecifications = launchSpecs

	return f
}

func (f SpotFleet) IAMFleetRoleRef() string {
	if f.IAMFleetRoleARN == "" {
		return `{"Fn::Join":["", [ "arn:aws:iam::", {"Ref":"AWS::AccountId"}, ":role/aws-ec2-spot-fleet-role" ]]}`
	} else {
		return f.IAMFleetRoleARN
	}
}
