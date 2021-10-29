package main_test

import (
	"context"
	"testing"

	"github.com/MicroFileServer/proto"
	"google.golang.org/grpc"
)

const token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6IjEyMyJ9.eyJhdWQiOiJpdGxhYiIsImlzcyI6Imh0dHBzOi8vZGV2LmlkZW50aXR5LnJ0dWl0bGFiLnJ1IiwiaWF0IjoxNTE2MjM5MDIyLCJleHAiOjE1MDU0Njc3NTY4NjksInN1YiI6InVzZXItMSIsInJvbGUiOiJ1c2VyIiwiaXRsYWIiOlsidXNlciJdLCJzY29wZSI6WyJyb2xlcyIsIm9wZW5pZCIsInByb2ZpbGUiLCJpdGxhYi5ldmVudHMiLCJpdGxhYi5yZXBvcnRzIl19.fF1qXeJXQdukk2uTziUkXMcYYHU6CvDpH8TkKKyAP0I"
func TestFunc_GRPC(t *testing.T) {
	conn, err := grpc.Dial("localhost:8082", grpc.WithInsecure(), grpc.WithAuthority(token))
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	defer conn.Close()

	c := proto.NewMicroFileServerClient(conn)

	_, err = c.DeleteFile(context.Background(), &proto.DeleteFileReq{FileId: "asdasd"})
	if err != nil {
		t.Log(err)
	}
}