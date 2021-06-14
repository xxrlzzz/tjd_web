package grpc

import (
	"google.golang.org/grpc/credentials"
	"traffic_jam_direction/pkg/setting"
)

func SetUp() error {
	var err error
	credential, err = credentials.NewClientTLSFromFile("conf/client.crt", setting.GrpcSetting.Host)
	if err != nil {
		return err
	}
	return nil
}
