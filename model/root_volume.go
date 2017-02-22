package model

import (
	"fmt"
)

type EbsVolume struct {
	Type      string `yaml:"type,omitempty"`
	Device    string `yaml:"device,omitempty"`
	MountPath string `yaml:"mountPath,omitempty"`
	Iops      int    `yaml:"iops,omitempty"`
	Size      int    `yaml:"size,omitempty"`
}

type RootVolume struct {
	Type      string `yaml:"rootVolumeType,omitempty"`
	Iops      int    `yaml:"rootVolumeIOPS,omitempty"`
	Size      int    `yaml:"rootVolumeSize,omitempty"`
}

func (rootVolume RootVolume) ToEbsVolume() (ebsVolume EbsVolume)  {
	ebsVolume.Size      = rootVolume.Size
	ebsVolume.Iops      = rootVolume.Iops
	ebsVolume.Type      = rootVolume.Type
	ebsVolume.Device    = "/dev/xvda"
	ebsVolume.MountPath = "/root"
	return ebsVolume
}

func (ebsVolume EbsVolume) ToRootVolume() (rootVolume RootVolume) {
	rootVolume.Size      = ebsVolume.Size
	rootVolume.Iops      = ebsVolume.Iops
	rootVolume.Type      = ebsVolume.Type
	return rootVolume
}

func NewGp2RootVolume(size int) RootVolume {
	return RootVolume{
		Size: size,
		Iops: 0,
		Type: "gp2",
	}
}

func NewIo1RootVolume(size int, iops int) RootVolume {
	return RootVolume{
		Size: size,
		Iops: iops,
		Type: "io1",
	}
}

func (v EbsVolume) Validate() error {
	if v.Type == "io1" {
		if v.Iops < 100 || v.Iops > 2000 {
			return fmt.Errorf(`invalid rootVolumeIOPS %d in %+v: rootVolumeIOPS must be between 100 and 2000`, v.Iops, v)
		}
	} else {
		if v.Iops != 0 {
			return fmt.Errorf(`invalid rootVolumeIOPS %d for volume type "%s" in %+v": rootVolumeIOPS must be 0 when rootVolumeType is "standard" or "gp1"`, v.Iops, v.Type, v)
		}

		if v.Type != "standard" && v.Type != "gp2" {
			return fmt.Errorf(`invalid rootVolumeType "%s" in %+v: rootVolumeType must be one of "standard", "gp1", "io1"`, v.Type, v)
		}
	}

	//if v.MountPath == "/root" {
	//	if v.Device != "/dev/xvda" {
	//		return fmt.Errorf(`invalid rootVolumeIOPS %d in%v`, v.Iops, v)
	//	}
	//} else {
	//	if !len(v.Device) == 9 || !strings.HasSuffix(v.Device, "/dev/xvd") || !strings.LastIndexAny(v.Device, "fghijklmnopqrstuvwxyz") {
	//		return fmt.Errorf(`invalid rootVolumeType "%s" in %+v: rootVolumeType must be one of "standard", "gp1", "io1"`, v.Type, v)
	//	}
	//}
	return nil
}
