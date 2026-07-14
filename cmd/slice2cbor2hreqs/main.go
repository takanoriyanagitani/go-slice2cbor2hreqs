package main

import (
	"context"
	"io"
	"iter"
	"log/slog"
	"net/http"
	"os"

	ch "github.com/takanoriyanagitani/go-slice2cbor2hreqs"
	fa "github.com/takanoriyanagitani/go-slice2cbor2hreqs/cbor/decoder/fa"
)

var c2r ch.CborToReqsBasicPartial = fa.CborToReqsBasicPartialFa

var baseURL ch.RawURL = ch.RawURL(os.Getenv("ENV_BASE_URL"))

func printHreq(hreq *http.Request) {
	slog.Info(
		"http-request",
		"method", hreq.Method,
		"url", hreq.URL.String(),
		"proto", hreq.Proto,
		"len", hreq.ContentLength,
		"host", hreq.Host,
	)
}

func sub(ctx context.Context) error {
	var rdr io.Reader = os.Stdin
	lmtd := &io.LimitedReader{
		R: rdr,
		N: 1048576,
	}
	slice, err := io.ReadAll(lmtd)
	if nil != err {
		return err
	}

	parsedURL, err := baseURL.Parse()
	if nil != err {
		return err
	}

	var dtos ch.PostReqsDTOs = c2r(ctx, slice)
	var partials iter.Seq2[ch.HTTPPostRequestBasicPartial, error] = dtos.
		Convert()
	var basics iter.Seq2[ch.HTTPPostRequestBasic, error] = ch.
		PartialReqs(partials).
		Convert(parsedURL)
	var hreqs iter.Seq2[*http.Request, error] = ch.
		BasicReqs(basics).
		Convert(ctx)

	for hreq, err := range hreqs {
		if nil != err {
			return err
		}

		printHreq(hreq)
	}
	return nil
}

func main() {
	err := sub(context.Background())
	if nil != err {
		slog.Error("error", "detail", err)
	}
}
