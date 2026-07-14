package slice2hreqs

import (
	"bytes"
	"context"
	"iter"
	"net/http"
	"net/url"
)

type PostBody []byte

type HTTPPostRequestBasic struct {
	*url.URL
	http.Header
	PostBody
}

func (b HTTPPostRequestBasic) PathAppended(path string) HTTPPostRequestBasic {
	return HTTPPostRequestBasic{
		URL:      b.URL.JoinPath(path),
		Header:   b.Header,
		PostBody: b.PostBody,
	}
}

func (b HTTPPostRequestBasic) ToRequest(
	ctx context.Context,
) (*http.Request, error) {
	return http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		b.URL.String(),
		bytes.NewReader(b.PostBody),
	)
}

type RawURL string

func (r RawURL) Parse() (*url.URL, error) {
	return url.ParseRequestURI(string(r))
}

type URLPath string

type HTTPPostRequestBasicPartial struct {
	URLPath
	http.Header
	PostBody
}

func (p HTTPPostRequestBasicPartial) ToReq(u *url.URL) HTTPPostRequestBasic {
	return HTTPPostRequestBasic{
		URL:      u,
		Header:   p.Header,
		PostBody: p.PostBody,
	}.
		PathAppended(string(p.URLPath))
}

type HTTPPostRequestBasicPartialDTO struct {
	URLPathRaw  string              `json:"url_path"`
	HeaderRaw   map[string][]string `json:"header"`
	PostBodyRaw []byte              `json:"post_body"`
}

func (d HTTPPostRequestBasicPartialDTO) Convert() HTTPPostRequestBasicPartial {
	return HTTPPostRequestBasicPartial{
		URLPath:  URLPath(d.URLPathRaw),
		Header:   d.HeaderRaw,
		PostBody: d.PostBodyRaw,
	}
}

type CborSlice []byte

type PostReqsDTOs iter.Seq2[HTTPPostRequestBasicPartialDTO, error]

type CborToReqsBasicPartial func(context.Context, CborSlice) PostReqsDTOs

func (d PostReqsDTOs) Convert() iter.Seq2[HTTPPostRequestBasicPartial, error] {
	return func(yield func(HTTPPostRequestBasicPartial, error) bool) {
		for dto, err := range d {
			if nil != err {
				yield(HTTPPostRequestBasicPartial{}, err)
				return
			}

			var converted HTTPPostRequestBasicPartial = dto.Convert()
			if !yield(converted, nil) {
				return
			}
		}
	}
}

type PartialReqs iter.Seq2[HTTPPostRequestBasicPartial, error]

func (p PartialReqs) Convert(
	url *url.URL,
) iter.Seq2[HTTPPostRequestBasic, error] {
	return func(yield func(HTTPPostRequestBasic, error) bool) {
		for partial, err := range p {
			if nil != err {
				yield(HTTPPostRequestBasic{}, err)
				return
			}

			var basic HTTPPostRequestBasic = partial.ToReq(url)
			if !yield(basic, nil) {
				return
			}
		}
	}
}

type BasicReqs iter.Seq2[HTTPPostRequestBasic, error]

func (b BasicReqs) Convert(
	ctx context.Context,
) iter.Seq2[*http.Request, error] {
	return func(yield func(*http.Request, error) bool) {
		for basic, err := range b {
			if nil != err {
				yield(nil, err)
				return
			}

			hreq, e := basic.ToRequest(ctx)
			if !yield(hreq, e) {
				return
			}
		}
	}
}
