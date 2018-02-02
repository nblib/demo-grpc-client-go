package demo_grpc_client_go

import (
	"testing"
	"google.golang.org/grpc"
	"log"
	"golang.org/x/net/context"
	"github.com/nblib/demo-grpc-client-go/clients/hello"
	"github.com/nblib/demo-grpc-client-go/clients/demo"
	"github.com/nblib/demo-grpc-client-go/clients/sample"
	"github.com/nblib/demo-grpc-client-go/clients/common"
	"fmt"
	"io"
	"time"
)

const address string = "localhost:50051"

//测试demo客户端
func TestHello(t *testing.T) {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	client := hello.NewHelloSerivceClient(conn)

	reply, err := client.TickInfo(context.Background(), &hello.HelloRequest{
		Name:    "hewe",
		Age:     30,
		IsAdult: true,
	})

	if err != nil {
		log.Fatalf("could not tick: %v", err)
	}
	log.Printf("tick reply: info==> %s, receive time==> %s", reply.GetInfo(), reply.GetReceiveTime())
}

//测试list和map,oneof
func TestDemo(t *testing.T) {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := demo.NewDemoServiceClient(conn)
	result, err := client.CheckIfBlack(context.Background(), &demo.CheckIps{
		Name: "testListAndMap",
		Ips:  []string{"192.158.22.33", "10.169.2.121", "192.168.23.111"},
	})
	if err != nil {
		log.Fatalf("check error: %v", err)
	}
	results := result.GetResults()
	for k, v := range results {
		log.Printf("ip: %s, isBlack: %t", k, v)
	}

	//测试获取oneof
	info, err := client.GetContectInfo(context.Background(), &demo.Empty{})
	if err != nil {
		log.Fatalf("GetContectInfo error: %v", err)
	}
	log.Printf("get contect info : tel= %v,cell=%v", info.GetTel(), info.GetCell())

}

func TestSamplePost(t *testing.T) {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	client := sample.NewPositionClient(conn)
	postLocationClient, err := client.PostLocation(context.Background())
	if err != nil {
		log.Fatalf("did not postLocationClient: %v", err)
	}
	postLocationClient.Send(&sample.Location{
		Lat: 11.323,
		Lon: 92.24,
	})
	postLocationClient.Send(&sample.Location{
		Lat: 11.323,
		Lon: 92.24,
	})
	postLocationClient.Send(&sample.Location{
		Lat: 11.323,
		Lon: 92.24,
	})

	_, err1 := postLocationClient.CloseAndRecv()
	if err1 != nil {
		log.Fatalf("did not CloseAndRecv: %v", err)
	}

}
func TestSamplePull(t *testing.T) {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	client := sample.NewPositionClient(conn)
	pullLocationClient, err := client.PullLocation(context.Background(), &common.Empty{})
	if err != nil {
		log.Fatalf("did not pullLocationClient: %v", err)
	}

	for {
		location, err := pullLocationClient.Recv()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatalf("did not Recv: %v", err)
		} else {
			fmt.Println(location.GetLat(), location.GetLon())
		}
	}

}
//测试是否自动重连
func TestInterruptConnect(t *testing.T) {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := demo.NewDemoServiceClient(conn)
	for {
		time.Sleep(2 * time.Second)
		result, err := client.CheckIfBlack(context.Background(), &demo.CheckIps{
			Name: "testListAndMap",
			Ips:  []string{"192.158.22.33", "10.169.2.121", "192.168.23.111"},
		})
		if err != nil {
			log.Printf("check error: %v", err)
			continue
		}
		results := result.GetResults()
		for k, v := range results {
			log.Printf("ip: %s, isBlack: %t", k, v)
		}
	}

}
