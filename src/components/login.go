package command

import (
	"strconv"

	"github.com/racker/passport-experiments/passport/login"
	"github.com/spf13/cobra"
)

func (p *Passport) GetLoginCommand() *cobra.Command {
	loginCommand := &cobra.Command{
		Use:  "login",
		Long: "TODO",
		RunE: p.LoginRun,
	}
	return loginCommand
}

func (p *Passport) LoginRun(cmd *cobra.Command, args []string) error {
	identity, err := login.Run()
	if err != nil {
		return err
	}

	racker, err := strconv.ParseBool(cmd.Flag("racker").Value.String())
	if err != nil {
		return err
	}

	return p.Config.UpdateConfigWithIdentity(*identity, racker)
}
