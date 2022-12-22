package main

import (
	"context"
	"fmt"
	"gclient/err"
	"log"
	"strconv"
	"strings"
	"time"

	strutils "github.com/randyardiansyah25/libpkg/util/str"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func CastError(er error) error {
	if st, ok := status.FromError(er); ok {
		switch st.Code() {
		case codes.Unavailable:
			return err.ErrRefused
		case codes.DeadlineExceeded:
			return err.ErrTimeout
		default:
			return er
		}
	} else {
		return er
	}
}

func RequestListUser() (er error) {
	deadline := time.Now().Add(10 * time.Second)
	ctx, cancelFunc := context.WithDeadline(context.Background(), deadline)

	defer func() {
		cancelFunc()
	}()

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	conn, er := grpc.Dial("localhost:2200", opts...)
	if er != nil {
		return
	}

	defer func() {
		conn.Close()
	}()

	client := NewUserHandlerClient(conn)
	if er != nil {
		return
	}

	log.Println("Request list user")
	list, er := client.ListUser(ctx, &emptypb.Empty{})
	if er != nil {
		er = CastError(er)
		return
	}

	var table = []string{
		"No  ID        Name                   Password        Gender       \n",
		"==================================================================\n",
		//   123412345678901234567890123456789012312345678901234561234567890123
		//       x         x                      x               x
	}

	for i, item := range list.List {
		row := []string{
			strutils.RightPad(strconv.Itoa(i+1), 4, " "),
			strutils.RightPad(item.Id, 10, " "),
			strutils.RightPad(item.Name, 23, " "),
			strutils.RightPad(item.Password, 16, " "),
			strutils.RightPad(item.Gender.String(), 13, " "),
			"\n",
		}
		table = append(table, strings.Join(row, ""))
	}

	fmt.Println(table)

	return nil
}

func main() {
	er := RequestListUser()
	if er != nil {
		fmt.Println(er.Error())
	}
}
