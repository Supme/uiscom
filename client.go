package uiscom

import (
	"context"
	"fmt"
	"github.com/ybbus/jsonrpc/v3"
	"time"
)

type Target string

func (t Target) URL() string {
	return string(t)
}

const (
	TargetComagic = Target("https://dataapi.comagic.ru/v2.0")
	TargetUiscom  = Target("https://dataapi.uiscom.ru/v2.0")
)

type Client struct {
	client      jsonrpc.RPCClient
	AccessToken string
}

func NewWithToken(target Target, token string) *Client {
	c := jsonrpc.NewClientWithOpts(
		target.URL(),
		&jsonrpc.RPCClientOpts{
			DefaultRequestID: int(time.Now().UTC().Unix()),
		})
	client := Client{
		client:      c,
		AccessToken: token,
	}
	return &client
}

func (c Client) call(ctx context.Context, method string, params ...any) (any, error) {
	resp, err := c.client.Call(ctx, method, params...)
	switch e := err.(type) {
	case nil:
	case *jsonrpc.HTTPError:
		return resp, fmt.Errorf("%d %s", e.Code, e.Error())
	default:
		return resp, e
	}
	if resp.Error != nil {
		//return resp, fmt.Errorf("%d %s %#v", resp.Error.Code, resp.Error.Message, resp.Error.Data)
		return resp, fmt.Errorf("%d %s", resp.Error.Code, resp.Error.Message)
	}
	return resp.Result, nil
}

func (c Client) GetAccount(ctx context.Context) (any, error) {
	return c.call(ctx, "get.account", map[string]string{"access_token": c.AccessToken})
}

func (c Client) GetCalls(ctx context.Context, userID int, dateFrom, dateTill time.Time, limit, offset int, fields ...Field) (any, error) {
	params := map[string]any{"access_token": c.AccessToken}
	if userID >= 0 {
		params["user_id"] = userID
	}
	params["date_from"] = TimeToString(dateFrom)
	params["date_till"] = TimeToString(dateTill)
	params["limit"] = limit
	params["offset"] = offset
	if fields != nil {
		params["fields"] = fields
	}
	return c.call(ctx, "get.calls_report", params)
}

func (c Client) GetCallLegs(ctx context.Context, userID int, dateFrom, dateTill time.Time, limit, offset int, fields ...Field) (any, error) {
	params := map[string]any{"access_token": c.AccessToken}
	if userID >= 0 {
		params["user_id"] = userID
	}
	params["date_from"] = TimeToString(dateFrom)
	params["date_till"] = TimeToString(dateTill)
	params["limit"] = limit
	params["offset"] = offset
	if fields != nil {
		params["fields"] = fields
	}
	return c.call(ctx, "get.call_legs_report", params)
}

func (c Client) GetEmployeeStat(ctx context.Context, userID int, dateFrom, dateTill time.Time, limit, offset int, defaultStatuses bool, fields ...Field) (any, error) {
	params := map[string]any{"access_token": c.AccessToken}
	if userID >= 0 {
		params["user_id"] = userID
	}
	params["date_from"] = TimeToString(dateFrom)
	params["date_till"] = TimeToString(dateTill)
	params["limit"] = limit
	params["offset"] = offset
	params["only_default_statuses_in_stats"] = defaultStatuses
	if fields != nil {
		params["fields"] = fields
	}
	return c.call(ctx, "get.employee_stat", params)
}
