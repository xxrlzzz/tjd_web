package grpc

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
	"io"
	"time"
	"traffic_jam_direction/pkg/base"
	pb "traffic_jam_direction/pkg/grpc/proto/gen"
	"traffic_jam_direction/pkg/logging"
	"traffic_jam_direction/pkg/setting"
)

// cache for credential file
var credential credentials.TransportCredentials = nil

type NavigationClient struct {
	Conn   *grpc.ClientConn
	Client pb.NavigationServiceClient
}

func toPoint(p base.Point) *pb.Point {
	return &pb.Point{Longitude: p.Longitude, Latitude: p.Latitude}
}
func toStep(paths []string) []map[string]interface{} {
	result := make([]map[string]interface{}, 0)
	for idx, path := range paths {
		if idx == 0 {
			continue
		}
		result = append(result, map[string]interface{}{
			"path": fmt.Sprintf("%s;%s", paths[idx-1], path),
		})
	}
	return result
}

func (client *NavigationClient) Init() error {
	var err error

	client.Conn, err = grpc.Dial(setting.GrpcSetting.Host+setting.GrpcSetting.NavigationPort, grpc.WithTransportCredentials(credential))
	if err != nil {
		return err
	}

	client.Client = pb.NewNavigationServiceClient(client.Conn)
	return nil
}

func (client *NavigationClient) Destroy() error {
	return client.Conn.Close()
}

func (client *NavigationClient) Navigation(start, end base.Point) (result map[string]interface{}, err error) {
	// necessary or not?
	result, err = nil, nil

	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(15*time.Second))
	defer cancel()

	stream, err := client.Client.Navigation(ctx, &pb.NavigationRequest{Start: toPoint(start), End: toPoint(end)})
	if err != nil {
		statusErr, ok := status.FromError(err)
		if ok && statusErr.Code() == codes.DeadlineExceeded {
			logging.Warn("client.Navigation err: deadline")
		} else {
			logging.WarnF("client.Navigation err: %v", err)
		}
		return
	}

	paths := make([]string, 0)
	for {
		nxt, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			logging.WarnF("%v.Navigation(_) = _, %v", client.Client, err)
			return result, err
		}

		path := fmt.Sprintf("%f,%f", nxt.To.Longitude, nxt.To.Latitude)
		paths = append(paths, path)
	}
	//paths = baidu.ConvertBaiduCoordinate(paths)
	result = map[string]interface{}{
		"routes": []map[string]interface{}{
			{"steps": toStep(paths)},
		},
	}
	//fmt.Printf("%#v", result)
	return
}
