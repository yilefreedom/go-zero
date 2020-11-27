// Code generated by goctl. DO NOT EDIT!
// Source: add.proto

//go:generate mockgen -destination ./adder_mock.go -package adder -source $GOFILE

package adder

import (
	"context"

	add "bookstore/rpc/add/internal/pb"

	"github.com/yileCJW/go-zero/zrpc"
)

type (
	AddReq  = add.AddReq
	AddResp = add.AddResp

	Adder interface {
		Add(ctx context.Context, in *AddReq) (*AddResp, error)
	}

	defaultAdder struct {
		cli zrpc.Client
	}
)

func NewAdder(cli zrpc.Client) Adder {
	return &defaultAdder{
		cli: cli,
	}
}

func (m *defaultAdder) Add(ctx context.Context, in *AddReq) (*AddResp, error) {
	adder := add.NewAdderClient(m.cli.Conn())
	return adder.Add(ctx, in)
}
