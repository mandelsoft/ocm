package validate

import (
	"encoding/json"

	"github.com/mandelsoft/goutils/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"ocm.software/ocm/api/credentials"
	"ocm.software/ocm/api/ocm/plugin/descriptor"
	"ocm.software/ocm/api/ocm/plugin/ppi"
	"ocm.software/ocm/api/utils/runtime"
)

const Name = "validate"

func New(p ppi.Plugin) *cobra.Command {
	opts := Options{}

	cmd := &cobra.Command{
		Use:   Name + " <spec>",
		Short: "validate input specification",
		Long: `
This command accepts an input specification as argument. It is used to
validate the specification and to provide some metadata for the given
specification.

This metadata has to be provided as JSON string on *stdout* and has the 
following fields: 

- **<code>mediaType</code>** *string*

  The media type of the artifact described by the specification. It may be part
  of the specification or implicitly determined by the input type.

- **<code>description</code>** *string*

  A short textual description of the described location.

- **<code>hint</code>** *string*

  A name hint of the described location used to reconstruct a useful
  name for local blobs uploaded to a dedicated repository technology.

- **<code>consumerId</code>** *map[string]string*

  The consumer id used to determine optional credentials for the
  underlying repository. If specified, at least the <code>type</code> field must be set.
`,
		Args: cobra.RangeArgs(1, 2),
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return opts.Complete(args)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return Command(p, cmd, &opts)
		},
	}
	opts.AddFlags(cmd.Flags())
	return cmd
}

type Options struct {
	Specification json.RawMessage
	Dir           string
}

func (o *Options) AddFlags(fs *pflag.FlagSet) {
}

func (o *Options) Complete(args []string) error {
	creds := 0
	if len(args) > 1 {
		o.Dir = args[creds]
		creds++
	} else {
		o.Dir = "."
	}
	if err := runtime.DefaultYAMLEncoding.Unmarshal([]byte(args[creds]), &o.Specification); err != nil {
		return errors.Wrapf(err, "invalid access specification")
	}
	return nil
}

type Result struct {
	MediaType  string                       `json:"mediaType"`
	Short      string                       `json:"description"`
	Hint       string                       `json:"hint"`
	ConsumerId credentials.ConsumerIdentity `json:"consumerId"`
}

func Command(p ppi.Plugin, cmd *cobra.Command, opts *Options) error {
	spec, err := p.DecodeInputSpecification(opts.Specification)
	if err != nil {
		return errors.Wrapf(err, "input specification")
	}

	m := p.GetInputType(spec.GetType())
	if m == nil {
		return errors.ErrUnknown(descriptor.KIND_INPUTTYPE, spec.GetType())
	}
	info, err := m.ValidateSpecification(p, opts.Dir, spec)
	if err != nil {
		return err
	}
	result := Result{MediaType: info.MediaType, ConsumerId: info.ConsumerId, Hint: info.Hint, Short: info.Short}
	data, err := json.Marshal(result)
	if err != nil {
		return err
	}
	cmd.Printf("%s\n", string(data))
	return nil
}