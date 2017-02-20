package config

import (
	cfg "github.com/coreos/kube-aws/core/controlplane/config"
	"github.com/coreos/kube-aws/model"
)

type WorkerNodePoolConfig struct {
	model.NodePoolConfig `yaml:",inline"`
}

func NewWorkerNodePoolConfig() WorkerNodePoolConfig {
	return WorkerNodePoolConfig{
		NodePoolConfig: model.NewDefaultNodePoolConfig(),
	}

}

func (c WorkerNodePoolConfig) ValidateInputs() error {
	return nil
}

func (c WorkerNodePoolConfig) Validate() error {
	if err := c.NodePoolConfig.Valid(); err != nil {
		return err
	}
	return nil
}

func (c WorkerNodePoolConfig) WithDefaultsFrom(main cfg.DefaultWorkerSettings) WorkerNodePoolConfig {
	if c.EbsVolume.EbsVolumeType == "" {
		c.EbsVolume.EbsVolumeType = main.WorkerEbsVolumeType
	}

	if c.EbsVolume.EbsVolumeIOPS == 0 && c.EbsVolume.EbsVolumeType == "io1" {
		c.EbsVolume.EbsVolumeIOPS = main.WorkerEbsVolumeIOPS
	}

	if c.SpotFleet.EbsVolumeType == "" {
		c.SpotFleet.EbsVolumeType = c.EbsVolume.EbsVolumeType
	}

	if c.EbsVolumeSize == 0 {
		c.EbsVolumeSize = main.WorkerEbsVolumeSize
	}

	if c.Tenancy == "" {
		c.Tenancy = main.WorkerTenancy
	}

	if c.InstanceType == "" {
		c.InstanceType = main.WorkerInstanceType
	}

	if c.CreateTimeout == "" {
		c.CreateTimeout = main.WorkerCreateTimeout
	}

	return c
}
