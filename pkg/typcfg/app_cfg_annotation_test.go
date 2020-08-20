package typcfg_test

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/typical-go/typical-go/pkg/execkit"
	"github.com/typical-go/typical-go/pkg/typannot"
	"github.com/typical-go/typical-go/pkg/typgo"
	"github.com/typical-go/typical-rest-server/pkg/typcfg"
)

func TestCfgAnnotation_Annotate(t *testing.T) {
	os.MkdirAll("somepkg1", 0777)
	defer os.RemoveAll("somepkg1")

	unpatch := execkit.Patch([]*execkit.RunExpectation{})
	defer unpatch(t)

	var out strings.Builder
	typcfg.Stdout = &out
	defer func() { typcfg.Stdout = os.Stdout }()

	AppCfgAnnotation := &typcfg.AppCfgAnnotation{}
	c := &typannot.Context{
		Destination: "somepkg1",
		Context: &typgo.Context{
			BuildSys: &typgo.BuildSys{
				Descriptor: &typgo.Descriptor{ProjectName: "some-project"},
			},
		},
		Summary: &typannot.Summary{
			Annots: []*typannot.Annot{
				{
					TagName: "@app-cfg",
					Decl: &typannot.Decl{
						Name:    "SomeSample",
						Package: "mypkg",
						Type: &typannot.StructType{
							Fields: []*typannot.Field{
								{Name: "SomeField1", Type: "string", StructTag: `default:"some-text"`},
								{Name: "SomeField2", Type: "int", StructTag: `default:"9876"`},
							},
						},
					},
				},
			},
		},
	}

	require.NoError(t, AppCfgAnnotation.Annotate(c))

	b, _ := ioutil.ReadFile("somepkg1/app_cfg_annotated.go")
	require.Equal(t, `package somepkg1

// Autogenerated by Typical-Go. DO NOT EDIT.

import (
	"github.com/kelseyhightower/envconfig"
)

func init() { 
	typapp.AppendCtor(
		&typapp.Constructor{Name: "",Fn: LoadSomeSample},
	)
}

// LoadSomeSample load env to new instance of SomeSample
func LoadSomeSample() (*mypkg.SomeSample, error) {
	var cfg mypkg.SomeSample
	if err := envconfig.Process("SOMESAMPLE", &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
`, string(b))

	require.Equal(t, "Generate @app-cfg to somepkg1/app_cfg_annotated.go\n", out.String())

}

func TestCfgAnnotation_Annotate_GenerateDotEnvAndUsageDoc(t *testing.T) {
	unpatch := execkit.Patch([]*execkit.RunExpectation{})
	defer unpatch(t)

	var out strings.Builder
	typcfg.Stdout = &out
	defer func() { typcfg.Stdout = os.Stdout }()

	defer os.Clearenv()

	a := &typcfg.AppCfgAnnotation{
		Target:   "some-target",
		Template: "some-template",
		DotEnv:   ".env33",
		UsageDoc: "some-usage.md",
	}
	c := &typannot.Context{
		Context: &typgo.Context{
			BuildSys: &typgo.BuildSys{
				Descriptor: &typgo.Descriptor{ProjectName: "some-project"},
			},
		},
		Summary: &typannot.Summary{Annots: []*typannot.Annot{
			{
				TagName:  "@app-cfg",
				TagParam: `ctor_name:"ctor1" prefix:"SS"`,
				Decl: &typannot.Decl{
					Name:    "SomeSample",
					Package: "mypkg",
					Type: &typannot.StructType{
						Fields: []*typannot.Field{
							{Name: "SomeField1", Type: "string", StructTag: `default:"some-text"`},
							{Name: "SomeField2", Type: "int", StructTag: `default:"9876"`},
						},
					},
				},
			},
		}},
	}

	require.NoError(t, a.Annotate(c))
	defer os.Remove(a.Target)
	defer os.Remove(a.DotEnv)
	defer os.Remove(a.UsageDoc)

	b, _ := ioutil.ReadFile(a.Target)
	require.Equal(t, `some-template`, string(b))

	b, _ = ioutil.ReadFile(a.DotEnv)
	require.Equal(t, "SS_SOMEFIELD1=some-text\nSS_SOMEFIELD2=9876\n", string(b))
	require.Equal(t, "some-text", os.Getenv("SS_SOMEFIELD1"))
	require.Equal(t, "9876", os.Getenv("SS_SOMEFIELD2"))

	require.Equal(t, "Generate @app-cfg to some-target\nNew keys added in '.env33': SS_SOMEFIELD1 SS_SOMEFIELD2\nGenerate 'some-usage.md'\n", out.String())
}

func TestCfgAnnotation_Annotate_Predefined(t *testing.T) {
	target := "cfg-target"

	unpatch := execkit.Patch([]*execkit.RunExpectation{})
	defer unpatch(t)
	defer os.RemoveAll(target)

	appCfgAnnotation := &typcfg.AppCfgAnnotation{
		TagName:  "@some-tag",
		Template: "some-template",
		Target:   target,
	}
	c := &typannot.Context{
		Context: &typgo.Context{
			BuildSys: &typgo.BuildSys{
				Descriptor: &typgo.Descriptor{ProjectName: "some-project"},
			},
		},
		Summary: &typannot.Summary{
			Annots: []*typannot.Annot{
				{
					TagName: "@some-tag",
					Decl: &typannot.Decl{
						Name:    "SomeSample",
						Package: "mypkg",
						Type:    &typannot.StructType{Fields: []*typannot.Field{}},
					},
				},
			},
		},
	}
	require.NoError(t, appCfgAnnotation.Annotate(c))

	b, _ := ioutil.ReadFile(target)
	require.Equal(t, `some-template`, string(b))
}

func TestCfgAnnotation_Annotate_RemoveTargetWhenNoAnnotation(t *testing.T) {
	target := "target1"
	defer os.Remove(target)
	ioutil.WriteFile(target, []byte("some-content"), 0777)
	c := &typannot.Context{
		Context: &typgo.Context{},
		Summary: &typannot.Summary{},
	}

	AppCfgAnnotation := &typcfg.AppCfgAnnotation{Target: target}
	require.NoError(t, AppCfgAnnotation.Annotate(c))
	_, err := os.Stat(target)
	require.True(t, os.IsNotExist(err))
}
