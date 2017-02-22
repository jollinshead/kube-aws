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
	if c.RootVolume.Type == "" {
		c.RootVolume.Type = main.WorkerRootVolumeType
	}

	if c.RootVolume.Iops == 0 && c.RootVolume.Type == "io1" {
		c.RootVolume.Iops = main.WorkerRootVolumeIOPS
	}

	if c.SpotFleet.RootVolumeType == "" {
		c.SpotFleet.RootVolumeType = c.RootVolume.Type
	}

	if c.RootVolume.Size == 0 {
		c.RootVolume.Size = main.WorkerRootVolumeSize
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
