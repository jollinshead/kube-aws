package model

import (
	"fmt"
	"strings"
	"regexp"
)

type EbsVolume struct {
	Type      string `yaml:"type,omitempty"`
	Iops      int    `yaml:"iops,omitempty"`
	SizeGb    int    `yaml:"sizeGb,omitempty"`
	Device    string `yaml:"device,omitempty"`
	MountPath string `yaml:"mountPath,omitempty"`
}

func (v EbsVolume) SystemdMountName() string {
	return strings.Join(strings.Split(strings.TrimLeft(v.MountPath, "/"), "/"), "-")
}

func (v EbsVolume) Validate() error {
	if v.Type == "io1" {
		if v.Iops < 100 || v.Iops > 2000 {
			return fmt.Errorf(`invalid iops "%d" in %+v: iops must be between "100" and "2000"`, v.Iops, v)
		}
	} else {
		if v.Iops != 0 {
			return fmt.Errorf(`invalid iops "%d" for volume type "%s" in %+v": iops must be "0" when type is "standard" or "gp1"`, v.Iops, v.Type, v)
		}

		if v.Type != "standard" && v.Type != "gp2" {
			return fmt.Errorf(`invalid type "%s" in %+v: type must be one of "standard", "gp1", "io1"`, v.Type, v)
		}
	}

	if v.SizeGb <= 0 {
		return fmt.Errorf(`invalid sizeGb "%d" in %+v: sizeGb must be greater than "0"`, v.SizeGb, v)
	}

	if v.MountPath == "/root" {
		return fmt.Errorf(`invalid mountPath "%s" in%v: mountPath cannot be "/root"`, v.MountPath, v)
	} else if regexp.MustCompile("^[a-zA-Z0-9/]*$").MatchString(v.MountPath) != true || strings.Contains(v.MountPath, "//") == true {
		return fmt.Errorf(`invalid mountPath "%s" in%v`, v.MountPath, v)
	}

	if strings.Compare(v.Device, "/dev/xvdf") == -1 || strings.Compare(v.Device, "/dev/xvdz") == 1 {
		return fmt.Errorf(`invalid device "%s" in%v: device must be a value from "/dev/xvdf" to "/dev/xvdz"`, v.Device, v)
	}

	return nil
}
