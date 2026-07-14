package fa

import (
	"bytes"
	"context"
	"errors"
	"io"

	fa "github.com/fxamacker/cbor/v2"
	ch "github.com/takanoriyanagitani/go-slice2cbor2hreqs"
)

type Decoder struct{ *fa.Decoder }

func (d Decoder) ToRequests() ch.PostReqsDTOs {
	return func(yield func(ch.HTTPPostRequestBasicPartialDTO, error) bool) {
		var dto ch.HTTPPostRequestBasicPartialDTO

		for {
			err := d.Decoder.Decode(&dto)

			// no more data
			if errors.Is(err, io.EOF) {
				return
			}

			// no error
			if nil == err {
				if !yield(dto, nil) {
					return
				}
				continue
			}

			ute, ok := errors.AsType[*fa.UnmarshalTypeError](err)
			if ok {
				var cborTyp string = ute.CBORType
				// slice may contain trailing zeros
				// TODO: more reliably detect EOF in the slice
				if "positive integer" == cborTyp {
					return
				}
			}

			if !yield(dto, err) {
				return
			}
		}
	}
}

func CborToReqsFa(_ context.Context, raw ch.CborSlice) ch.PostReqsDTOs {
	dec := Decoder{Decoder: fa.NewDecoder(bytes.NewReader(raw))}
	return dec.ToRequests()
}

var CborToReqsBasicPartialFa ch.CborToReqsBasicPartial = CborToReqsFa
