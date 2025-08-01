package cmd

import (
	"strconv"

	"github.com/spf13/cobra"

	"github.com/filebrowser/filebrowser/v2/settings"
	"github.com/filebrowser/filebrowser/v2/users"
)

func init() {
	rulesCmd.AddCommand(rulesRmCommand)
	rulesRmCommand.Flags().Uint("index", 0, "index of rule to remove")
	_ = rulesRmCommand.MarkFlagRequired("index")
}

var rulesRmCommand = &cobra.Command{
	Use:   "rm <index> [index_end]",
	Short: "Remove a global rule or user rule",
	Long: `Remove a global rule or user rule. The provided index
is the same that's printed when you run 'rules ls'. Note
that after each removal/addition, the index of the
commands change. So be careful when removing them after each
other.

You can also specify an optional parameter (index_end) so
you can remove all commands from 'index' to 'index_end',
including 'index_end'.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if err := cobra.RangeArgs(1, 2)(cmd, args); err != nil {
			return err
		}

		for _, arg := range args {
			if _, err := strconv.Atoi(arg); err != nil {
				return err
			}
		}

		return nil
	},
	RunE: python(func(cmd *cobra.Command, args []string, d *pythonData) error {
		i, err := strconv.Atoi(args[0])
		if err != nil {
			return err
		}
		f := i
		if len(args) == 2 {
			f, err = strconv.Atoi(args[1])
			if err != nil {
				return err
			}
		}

		user := func(u *users.User) error {
			u.Rules = append(u.Rules[:i], u.Rules[f+1:]...)
			return d.store.Users.Save(u)
		}

		global := func(s *settings.Settings) error {
			s.Rules = append(s.Rules[:i], s.Rules[f+1:]...)
			return d.store.Settings.Save(s)
		}

		return runRules(d.store, cmd, user, global)
	}, pythonConfig{}),
}
