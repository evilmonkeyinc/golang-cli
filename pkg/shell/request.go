package shell

import "context"

type Request interface {
	Args() []string
	Context() context.Context
	Routes() Routes
	WithContext(context.Context) Request
	WithArgs([]string) Request
	WithRoutes(Routes) Request
}

func newRequest(ctx context.Context, args []string, routes Routes) *request {
	return &request{
		ctx:    ctx,
		args:   args,
		routes: routes,
	}
}

type request struct {
	ctx    context.Context
	args   []string
	routes Routes
}

func (req *request) Args() []string {
	return req.args
}

func (req *request) Context() context.Context {
	return req.ctx
}

func (req *request) Routes() Routes {
	return req.routes
}

func (req *request) WithContext(ctx context.Context) Request {
	return &request{
		ctx:    ctx,
		args:   req.args,
		routes: req.routes,
	}
}

func (req *request) WithArgs(args []string) Request {
	return &request{
		ctx:    req.ctx,
		args:   args,
		routes: req.routes,
	}
}

func (req *request) WithRoutes(routes Routes) Request {
	return &request{
		ctx:    req.ctx,
		args:   req.args,
		routes: routes,
	}
}
