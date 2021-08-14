package reddit

import (
	"fmt"
	"strings"

	"github.com/valyala/fastjson"
)

type ResponseHandler func(*fastjson.Value) interface{}

type Error struct {
	Message    string `json:"message"`
	Code       int    `json:"error"`
	StatusCode int
}

func (err *Error) Error() string {
	return fmt.Sprintf("%s (%d)", err.Message, err.Code)
}

func NewError(val *fastjson.Value, status int) *Error {
	err := &Error{}

	err.Message = string(val.GetStringBytes("message"))
	err.Code = val.GetInt("error")

	return err
}

type RefreshTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func NewRefreshTokenResponse(val *fastjson.Value) interface{} {
	rtr := &RefreshTokenResponse{}

	rtr.AccessToken = string(val.GetStringBytes("access_token"))
	rtr.RefreshToken = string(val.GetStringBytes("refresh_token"))

	return rtr
}

type MeResponse struct {
	ID   string `json:"id"`
	Name string
}

func (mr *MeResponse) NormalizedUsername() string {
	return strings.ToLower(mr.Name)
}

func NewMeResponse(val *fastjson.Value) interface{} {
	mr := &MeResponse{}

	mr.ID = string(val.GetStringBytes("id"))
	mr.Name = string(val.GetStringBytes("name"))

	return mr
}

type Thing struct {
	Kind        string  `json:"kind"`
	ID          string  `json:"id"`
	Type        string  `json:"type"`
	Author      string  `json:"author"`
	Subject     string  `json:"subject"`
	Body        string  `json:"body"`
	CreatedAt   float64 `json:"created_utc"`
	Context     string  `json:"context"`
	ParentID    string  `json:"parent_id"`
	LinkTitle   string  `json:"link_title"`
	Destination string  `json:"dest"`
	Subreddit   string  `json:"subreddit"`
}

func (t *Thing) FullName() string {
	return fmt.Sprintf("%s_%s", t.Kind, t.ID)
}

func NewThing(val *fastjson.Value) *Thing {
	t := &Thing{}

	t.Kind = string(val.GetStringBytes("kind"))

	data := val.Get("data")

	t.ID = string(data.GetStringBytes("id"))
	t.Type = string(data.GetStringBytes("type"))
	t.Author = string(data.GetStringBytes("author"))
	t.Subject = string(data.GetStringBytes("subject"))
	t.Body = string(data.GetStringBytes("body"))
	t.CreatedAt = data.GetFloat64("created_utc")
	t.Context = string(data.GetStringBytes("context"))
	t.ParentID = string(data.GetStringBytes("parent_id"))
	t.LinkTitle = string(data.GetStringBytes("link_title"))
	t.Destination = string(data.GetStringBytes("dest"))
	t.Subreddit = string(data.GetStringBytes("subreddit"))

	return t
}

type ListingResponse struct {
	Count    int
	Children []*Thing
	After    string
	Before   string
}

func NewListingResponse(val *fastjson.Value) interface{} {
	lr := &ListingResponse{}

	data := val.Get("data")
	children := data.GetArray("children")

	lr.After = string(data.GetStringBytes("after"))
	lr.Before = string(data.GetStringBytes("before"))
	lr.Count = len(children)

	if lr.Count == 0 {
		return lr
	}

	lr.Children = make([]*Thing, lr.Count)
	for i := 0; i < lr.Count; i++ {
		t := NewThing(children[i])
		lr.Children[i] = t
	}

	return lr
}

var EmptyListingResponse = &ListingResponse{}
