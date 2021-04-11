package grpc

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
	"time"
	pb "traffic_jam_direction/pkg/grpc/proto/gen"
	"traffic_jam_direction/pkg/logging"
	"traffic_jam_direction/pkg/setting"
)

type TrafficClient struct {
	conn   *grpc.ClientConn
	client pb.TrafficServiceClient
	initialized bool
}

func (c *TrafficClient) Init() error {
	var err error
	if credential ==  nil {
		credential,err  = credentials.NewClientTLSFromFile("conf/client.crt", setting.GrpcSetting.Host)
	}
	if err != nil {
		return err
	}

	c.conn, err = grpc.Dial(setting.GrpcSetting.TrafficPort, grpc.WithTransportCredentials(credential))
	if err != nil {
		return err
	}

	c.client = pb.NewTrafficServiceClient(c.conn)
	c.initialized = true
	return nil
}

func (c *TrafficClient) Destroy() error {
	if !c.initialized {
		logging.Fatal("TrafficClient use before init");
	}
	return c.conn.Close()
}

func (c *TrafficClient) QueryTraffic(longitude, latitude float64) (result map[string]interface{}, err error) {
	result, err = nil, nil
	if !c.initialized {
		logging.Fatal("TrafficClient use before init")
	}
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(5*time.Second))
	defer cancel()

	resp, err := c.client.QueryTraffic(ctx, &pb.TrafficRequest{Latitude: latitude, Longitude: longitude})
	if err != nil {
		statusErr, ok := status.FromError(err)
		if ok && statusErr.Code() == codes.DeadlineExceeded {
			logging.Warn("client.QueryTraffic err: deadline")
		} else {
			logging.WarnF("client.QueryTraffic err: %v", err)
		}
		return
	}
	result = map[string]interface{}{
		"speed_normal": resp.SpeedNormal,
		"speed_hot": resp.SpeedHot,
		"speed_extreme": resp.SpeedExtreme,
	}
	return
}
