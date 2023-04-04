package mutate

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
)

func TestSidecarConfig_ConfigFromAnnotations(t *testing.T) {
	assert := assert.New(t)
	type fields struct {
		Containers []corev1.Container
	}
	type args struct {
		annotations map[string]string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []string
	}{
		{
			name: "config from annotation append",
			fields: fields{
				Containers: []corev1.Container{
					{
						Args: []string{"--testpar1=1", "--testpar2=2"},
					},
				},
			},
			args: args{
				annotations: map[string]string{
					AnnotationIntegrityMonitorInject:                   "true",
					AnnotationProcessImage:                             "nginx=nginx:latest",
					fmt.Sprintf("nginx.%s", AnnotationMonitoringPaths): "/proc,/dir1, dir2",
					fmt.Sprintf("redis.%s", AnnotationMonitoringPaths): "/proc, /dir1,dir2",
					"unneeded-annotation":                              "value",
				},
			},
			want: []string{"--testpar1=1", "--testpar2=2", "--monitoring-options=nginx=/proc,/dir1,dir2 redis=/proc,/dir1,dir2", "--process-image=nginx=nginx:latest"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sc := &SidecarConfig{
				Containers: tt.fields.Containers,
			}
			sc.ConfigFromAnnotations(tt.args.annotations)
			assert.ElementsMatchf(tt.want, sc.Containers[0].Args, fmt.Sprintf("slices %v and %v is not equal", tt.want, sc.Containers[0].Args))
		})
	}
}
