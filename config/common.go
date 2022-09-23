package config

import (
	"github.com/kadaan/papertrail/lib/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"time"
)

const (
	startKey        = "start"
	endKey          = "end"
	groupIdKey      = "groupId"
	systemIdKey     = "systemId"
	defaultGroupId  = 0
	defaultSystemId = 0
)

var (
	defaultEnd   = time.Now().UTC()
	defaultStart = defaultEnd.Add(-6 * time.Hour)
)

func NewFlagBuilder(cmd *cobra.Command) FlagBuilder {
	return &flagBuilder{
		cmd: cmd,
	}
}

type Flag interface {
	Required() Flag
}

type FileFlag interface {
	Flag
	Extensions(extensions ...string) FileFlag
}

type compositeFlag struct {
	flags []Flag
}

func (f *compositeFlag) Required() Flag {
	for _, c := range f.flags {
		_ = c.Required()
	}
	return f
}

type flag struct {
	builder *flagBuilder
	flag    *pflag.Flag
}

func (f *flag) Required() Flag {
	_ = f.builder.cmd.MarkFlagRequired(f.flag.Name)
	return f
}

func (f *flag) Extensions(extensions ...string) FileFlag {
	_ = f.builder.cmd.MarkFlagFilename(f.flag.Name, extensions...)
	return f
}

type FlagBuilder interface {
	TimeRange(startDest *time.Time, endDest *time.Time, usage string) Flag
	GroupID(dest *uint, usage string) Flag
	SystemID(dest *uint, usage string) Flag
}

type flagBuilder struct {
	cmd *cobra.Command
}

func (fb *flagBuilder) newFlag(name string, creator func(flagSet *pflag.FlagSet)) *flag {
	creator(fb.cmd.Flags())
	f := fb.cmd.Flags().Lookup(name)
	_ = viper.BindPFlag(name, f)
	return &flag{
		builder: fb,
		flag:    f,
	}
}

func (fb *flagBuilder) addValidation(validation func(cmd *cobra.Command, args []string) error) {
	if fb.cmd.PreRunE != nil {
		existingValidation := fb.cmd.PreRunE
		fb.cmd.PreRunE = func(cmd *cobra.Command, args []string) error {
			if err := validation(cmd, args); err != nil {
				return err
			}
			return existingValidation(cmd, args)
		}
	} else {
		fb.cmd.PreRunE = validation
	}
}

func (fb *flagBuilder) TimeRange(startDest *time.Time, endDest *time.Time, usage string) Flag {
	startFlag := fb.Time(startDest, startKey, defaultStart, usage+" from")
	endFlag := fb.Time(endDest, endKey, defaultEnd, usage+" to")
	fb.addValidation(func(cmd *cobra.Command, args []string) error {
		if !(*startDest).Before(*endDest) {
			return errors.New("start time is not before end time")
		}
		return nil
	})
	return &compositeFlag{
		flags: []Flag{startFlag, endFlag},
	}
}

func (fb *flagBuilder) StartTime(dest *time.Time, usage string) Flag {
	return fb.Time(dest, startKey, defaultStart, usage)
}

func (fb *flagBuilder) EndTime(dest *time.Time, usage string) Flag {
	return fb.Time(dest, endKey, defaultEnd, usage)
}

func (fb *flagBuilder) Time(dest *time.Time, name string, defaultValue time.Time, usage string) Flag {
	return fb.newFlag(name, func(flagSet *pflag.FlagSet) {
		flagSet.Var(NewTimeValue(dest, defaultValue), name, usage)
	})
}

func (fb *flagBuilder) GroupID(dest *uint, usage string) Flag {
	return fb.newFlag(groupIdKey, func(flagSet *pflag.FlagSet) {
		flagSet.UintVar(dest, groupIdKey, defaultGroupId, usage)
	})
}

func (fb *flagBuilder) SystemID(dest *uint, usage string) Flag {
	return fb.newFlag(systemIdKey, func(flagSet *pflag.FlagSet) {
		flagSet.UintVar(dest, systemIdKey, defaultSystemId, usage)
	})
}
