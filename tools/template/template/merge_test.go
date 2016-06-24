package template

import (
	"reflect"
	"testing"

	"github.com/openshift/origin/pkg/template/api"
	kapi "k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/runtime"
)

func BaseTemplate() *api.Template {
	return &api.Template{
		ObjectMeta: kapi.ObjectMeta{
			Name: "fh-core",
			Annotations: map[string]string{
				"templateVersion": "0.0.0",
				"description":     "RHMAP - Core template",
			},
		},
		Parameters: []api.Parameter{
			{
				Name:  "BASE_TEMPLATE",
				Value: "yes",
			},
		},
		Objects: []runtime.Object{
			Service("base-1"),
			Service("base-2"),
		},
	}
}

func ComponentTemplate(id string) *api.Template {
	name := "component-" + id
	return &api.Template{
		ObjectMeta: kapi.ObjectMeta{
			Name: name,
			Annotations: map[string]string{
				"ignored": "yes",
			},
		},
		Parameters: []api.Parameter{
			{
				Name:  "PARAM-" + name,
				Value: name,
			},
		},
		Objects: []runtime.Object{
			Service(name),
		},
	}
}

func Service(name string) *kapi.Service {
	return &kapi.Service{
		ObjectMeta: kapi.ObjectMeta{
			Name: name,
			Labels: map[string]string{
				"name": name,
			},
		},
		Spec: kapi.ServiceSpec{
			Selector: map[string]string{
				"name": name,
			},
			Ports: []kapi.ServicePort{
				{
					Port: 8080,
				},
			},
		},
	}
}

type MergeTest struct {
	what string          // test description
	in   []*api.Template // input templates
	want *api.Template   // output template
}

var mergeTests = []MergeTest{
	{
		what: "degenerate case, empty input",
		in:   []*api.Template{},
		want: nil,
	},
	{
		what: "simple template is not changed",
		in: []*api.Template{
			BaseTemplate(),
		},
		want: BaseTemplate(),
	},
	{
		what: "add object and parameter to template",
		in: []*api.Template{
			BaseTemplate(),
			ComponentTemplate("1"),
		},
		want: &api.Template{
			ObjectMeta: BaseTemplate().ObjectMeta,
			Parameters: append(BaseTemplate().Parameters,
				ComponentTemplate("1").Parameters...),
			Objects: append(BaseTemplate().Objects,
				ComponentTemplate("1").Objects...),
		},
	},
	{
		what: "add object and parameter from multiple components",
		in: []*api.Template{
			BaseTemplate(),
			ComponentTemplate("1"),
			ComponentTemplate("2"),
		},
		want: &api.Template{
			ObjectMeta: BaseTemplate().ObjectMeta,
			Parameters: append(BaseTemplate().Parameters,
				append(ComponentTemplate("1").Parameters,
					ComponentTemplate("2").Parameters...)...),
			Objects: append(BaseTemplate().Objects,
				append(ComponentTemplate("1").Objects,
					ComponentTemplate("2").Objects...)...),
		},
	},
	{
		what: "object and param deduplication",
		in: []*api.Template{
			BaseTemplate(),
			ComponentTemplate("1"),
			ComponentTemplate("1"),
		},
		want: &api.Template{
			ObjectMeta: BaseTemplate().ObjectMeta,
			Parameters: append(BaseTemplate().Parameters,
				ComponentTemplate("1").Parameters...),
			Objects: append(BaseTemplate().Objects,
				ComponentTemplate("1").Objects...),
		},
	},
}

func TestMerge(t *testing.T) {
	for _, tt := range mergeTests {
		got := Merge(tt.in...)
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%s:\ngot:\n%#v,\nwant:\n%#v", tt.what, got, tt.want)
		}
	}
}
