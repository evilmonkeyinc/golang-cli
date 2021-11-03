package test

import (
	"context"

	"github.com/evilmonkeyinc/golang-cli/pkg/shell"
)

type TestRequest struct {
	ctx    context.Context
	args   []string
	routes shell.Routes
}

func (req *TestRequest) Args() []string {
	return req.args
}

func (req *TestRequest) Context() context.Context {
	return req.ctx
}

func (req *TestRequest) Routes() shell.Routes {
	return req.routes
}

func (req *TestRequest) WithContext(ctx context.Context) shell.Request {
	return &TestRequest{
		ctx:    ctx,
		args:   req.args,
		routes: req.routes,
	}
}

func (req *TestRequest) WithArgs(args []string) shell.Request {
	return &TestRequest{
		ctx:    req.ctx,
		args:   args,
		routes: req.routes,
	}
}

func (req *TestRequest) WithRoutes(routes shell.Routes) shell.Request {
	return &TestRequest{
		ctx:    req.ctx,
		args:   req.args,
		routes: routes,
	}
}
