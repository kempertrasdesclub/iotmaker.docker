package iotmakerDocker

import (
	"github.com/docker/go-connections/nat"
)

// Mount nat por list by image config
func (el *DockerSystem) ImageMountNatPortListChangeExposedWithIpAddress(imageId, ipAddress string, currentPortList, changeToPortList []nat.Port) (error, nat.PortMap) {
	var err error
	var portList []nat.Port
	var ret nat.PortMap = make(map[nat.Port][]nat.PortBinding)

	err, portList = el.ImageListExposedPorts(imageId)
	if err != nil {
		return err, nat.PortMap{}
	}

	for _, port := range portList {
		inPort := ""
		for k, currPort := range currentPortList {
			if currPort.Port() == port.Port() && currPort.Proto() == port.Proto() {
				inPort = changeToPortList[k].Port()
				break
			}
		}

		ret[port] = []nat.PortBinding{
			{
				HostPort: inPort,
				HostIP:   ipAddress,
			},
		}
	}

	return err, ret
}
