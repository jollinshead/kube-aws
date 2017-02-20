package model

import "fmt"

type EbsVolume struct {
	EbsVolumeType string `yaml:"ebsVolumeType,omitempty"`
	EbsVolumeIOPS int    `yaml:"ebsVolumeIOPS,omitempty"`
	EbsVolumeSize int    `yaml:"ebsVolumeSize,omitempty"`
}

func NewGp2EbsVolume(size int) EbsVolume {
	return EbsVolume{
		EbsVolumeSize: size,
		EbsVolumeIOPS: 0,
		EbsVolumeType: "gp2",
	}
}

func NewIo1EbsVolume(size int, iops int) EbsVolume {
	return EbsVolume{
		EbsVolumeSize: size,
		EbsVolumeIOPS: iops,
		EbsVolumeType: "io1",
	}
}

func (v EbsVolume) Validate() error {
	if v.EbsVolumeType == "io1" {
		if v.EbsVolumeIOPS < 100 || v.EbsVolumeIOPS > 2000 {
			return fmt.Errorf(`invalid ebsVolumeIOPS %d in %+v: ebsVolumeIOPS must be between 100 and 2000`, v.EbsVolumeIOPS, v)
		}
	} else {
		if v.EbsVolumeIOPS != 0 {
			return fmt.Errorf(`invalid ebsVolumeIOPS %d for volume type "%s" in %+v": ebsVolumeIOPS must be 0 when ebsVolumeType is "standard" or "gp1"`, v.EbsVolumeIOPS, v.EbsVolumeType, v)
		}

		if v.EbsVolumeType != "standard" && v.EbsVolumeType != "gp2" {
			return fmt.Errorf(`invalid ebsVolumeType "%s" in %+v: ebsVolumeType must be one of "standard", "gp1", "io1"`, v.EbsVolumeType, v)
		}
	}
	return nil
}
