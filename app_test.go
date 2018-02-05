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
	//建立一个连接,WithInsecure表示不使用验证,不然在没有验证的情况下会报错
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	//确保关闭,正式使用时,作为全局对象,不需要关闭
	defer conn.Close()
	//根据连接生成一个客户端,
	client := hello.NewHelloSerivceClient(conn)
	//通过客户端调用方法,会调用服务端的方法.传入请求的message,最后服务端会返回响应message
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
	// 请求的message中包含list
	result, err := client.CheckIfBlack(context.Background(), &demo.CheckIps{
		Name: "testListAndMap",
		Ips:  []string{"192.158.22.33", "10.169.2.121", "192.168.23.111"},
	})
	if err != nil {
		log.Fatalf("check error: %v", err)
	}
	//返回message包含map类型数据
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

//测试发送流数据,也就是一次调用,可以发送多条message
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

//返回获取多条message
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
//当请求完成后,服务端关闭,这是再一次请求会报错,然后再把服务端打开,请求会自动重新连接获取
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
