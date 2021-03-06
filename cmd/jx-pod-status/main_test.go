package main

import (
	"os"
	"reflect"
	"testing"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/fake"
)

func TestOptions_findErrors(t *testing.T) {

	objects := getTestPods()

	type fields struct {
		objects   []runtime.Object
		namespace string
	}
	tests := []struct {
		name    string
		fields  fields
		want    []string
		wantErr bool
	}{
		{name: "single_namespace", fields: fields{
			namespace: "foo",
		}, want: []string{"foo-pod is in pod status phase Pending "}, wantErr: false},
		{name: "multi_namespace", fields: fields{
			namespace: "",
		}, want: []string{"bar-pod is in pod status phase Pending ", "foo-pod is in pod status phase Pending "}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			os.Setenv(envVarTargetNamespace, tt.fields.namespace)

			client := fake.NewSimpleClientset(objects...)
			o := Options{
				client: client,
			}
			got, err := o.findErrors()
			if (err != nil) != tt.wantErr {
				t.Errorf("findErrors() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("findErrors() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func getTestPods() []runtime.Object {

	return []runtime.Object{
		&v1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: "foo",
			},
		},
		&v1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "foo-pod",
				Namespace: "foo",
			},
			Status: v1.PodStatus{
				Phase: v1.PodPending,
			},
		},
		&v1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: "bar",
			},
		},
		&v1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "bar-pod",
				Namespace: "bar",
			},
			Status: v1.PodStatus{
				Phase: v1.PodPending,
			},
		},
	}
}
