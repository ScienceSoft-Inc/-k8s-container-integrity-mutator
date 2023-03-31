package mutate

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckProcessImage(t *testing.T) {

	tests := []struct {
		name        string
		annotations map[string]string
		expect      struct {
			result bool
			err    error
		}
	}{
		{
			"success one image",
			map[string]string{
				AnnotationProcessImage: "nginx=my-image:latest",
			},
			struct {
				result bool
				err    error
			}{
				true, nil,
			},
		},
		{
			"success few images",
			map[string]string{
				AnnotationProcessImage: "nginx=my-image:latest,mongodb=mongodb:6-nanoserver,postgresql=postgresql:15-alpine3.17",
			},
			struct {
				result bool
				err    error
			}{
				true, nil,
			},
		},
		{
			"success no processes",
			map[string]string{
				AnnotationProcessImage: "",
			},
			struct {
				result bool
				err    error
			}{
				true, nil,
			},
		},
		{
			"failed no annotation",
			map[string]string{},
			struct {
				result bool
				err    error
			}{
				false, nil,
			},
		},
		{
			"invalid annotation format",
			map[string]string{
				AnnotationProcessImage: "myprocess:myimage",
			},
			struct {
				result bool
				err    error
			}{
				false, ErrInvalidProcessImageFormat,
			},
		},
		{
			"invalid annotation format 2",
			map[string]string{
				AnnotationProcessImage: "nginx=nginx:latest,postgresql",
			},
			struct {
				result bool
				err    error
			}{
				false, ErrInvalidProcessImageFormat,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {})
		ok, err := checkProcessImage(test.annotations)
		assert.Equal(t, test.expect.result, ok)
		assert.Equal(t, test.expect.err, err)
	}
}
