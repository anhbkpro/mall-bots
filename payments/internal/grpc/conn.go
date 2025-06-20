package grpc

import (
	"context"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func Dial(ctx context.Context, address string) (conn *grpc.ClientConn, err error) {
	conn, err = grpc.DialContext(ctx, address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return
	}

	defer func() {
		if err != nil {
			if err = conn.Close(); err != nil {
				log.Printf("failed to close connection: %v", err)
			}
			return
		}

		go func() {
			<-ctx.Done()
			if err = conn.Close(); err != nil {
				log.Printf("failed to close connection: %v", err)
			}
		}()
	}()

	return
}
