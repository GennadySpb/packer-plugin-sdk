package template

import (
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/go-multierror"
)

// Template represents the parsed template that is used to configure
// Packer builds.
type Template struct {
	Description string
	MinVersion  string

	Variables      map[string]*Variable
	Builders       map[string]*Builder
	Provisioners   []*Provisioner
	PostProcessors [][]*PostProcessor
	Push           *Push
}

// Builder represents a builder configured in the template
type Builder struct {
	Name   string
	Type   string
	Config map[string]interface{}
}

// PostProcessor represents a post-processor within the template.
type PostProcessor struct {
	OnlyExcept `mapstructure:",squash"`

	Type              string
	KeepInputArtifact bool `mapstructure:"keep_input_artifact"`
	Config            map[string]interface{}
}

// Provisioner represents a provisioner within the template.
type Provisioner struct {
	OnlyExcept `mapstructure:",squash"`

	Type        string
	Config      map[string]interface{}
	Override    map[string]interface{}
	PauseBefore time.Duration `mapstructure:"pause_before"`
}

// Push represents the configuration for pushing the template to Atlas.
type Push struct {
	Name    string
	Address string
	BaseDir string `mapstructure:"base_dir"`
	Include []string
	Exclude []string
	Token   string
	VCS     bool
}

// Variable represents a variable within the template
type Variable struct {
	Default  string
	Required bool
}

// OnlyExcept is a struct that is meant to be embedded that contains the
// logic required for "only" and "except" meta-parameters.
type OnlyExcept struct {
	Only   []string
	Except []string
}

//-------------------------------------------------------------------
// Functions
//-------------------------------------------------------------------

// Validate does some basic validation of the template on top of the
// validation that occurs while parsing. If possible, we try to defer
// validation to here. The validation errors that occur during parsing
// are the minimal necessary to make sure parsing builds a reasonable
// Template structure.
func (t *Template) Validate() error {
	var err error

	// At least one builder must be defined
	if len(t.Builders) == 0 {
		err = multierror.Append(err, errors.New(
			"at least one builder must be defined"))
	}

	// Verify that the provisioner overrides target builders that exist
	for i, p := range t.Provisioners {
		for name, _ := range p.Override {
			if _, ok := t.Builders[name]; !ok {
				err = multierror.Append(err, fmt.Errorf(
					"provisioner %d: override '%s' doesn't exist",
					i+1, name))
			}
		}
	}

	return err
}

//-------------------------------------------------------------------
// GoStringer
//-------------------------------------------------------------------

func (b *Builder) GoString() string {
	return fmt.Sprintf("*%#v", *b)
}

func (p *Provisioner) GoString() string {
	return fmt.Sprintf("*%#v", *p)
}

func (p *PostProcessor) GoString() string {
	return fmt.Sprintf("*%#v", *p)
}

func (v *Variable) GoString() string {
	return fmt.Sprintf("*%#v", *v)
}