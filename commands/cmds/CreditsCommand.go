package cmds

import (
	"arismcnc/database"
	"fmt"
	"io"

	"github.com/gliderlabs/ssh"
)

type CreditsCommand struct{}

func (c *CreditsCommand) Name() string {
	return "credits"
}

func (c *CreditsCommand) Execute(session ssh.Session, db *database.Database, args []string, output io.Writer) {
	fmt.Fprintln(output, "\033[8;30;80t")
	fmt.Fprintln(output, "\033[2J\033[1;1H")
	fmt.Fprintln(output, "")
	fmt.Fprintln(output, "")
	fmt.Fprintln(output, " \x1b[38;5;231m[ \x1b[38;5;93mAll Credits \x1b[38;5;231m·]")
	fmt.Fprintln(output, "  \x1b[38;5;231m• \x1b[38;5;92mMade by \x1b[38;5;231m··················· \x1b[38;5;32mTry999_9")
	fmt.Fprintln(output, "  \x1b[38;5;231m• \x1b[38;5;92mOwner \x1b[38;5;231m····················· \x1b[38;5;32mTry999_9 - Mrbrew - freeze")
	fmt.Fprintln(output, "  \x1b[38;5;231m• \x1b[38;5;92mGalaxia cnc version \x1b[38;5;231m······· \x1b[38;5;32m1.0.0")
	fmt.Fprintln(output, "  \x1b[38;5;231m• \x1b[38;5;92mGalaxia started at \x1b[38;5;231m········ \x1b[38;5;32m1/1/2025")
	fmt.Fprintln(output, "  \x1b[38;5;231m• \x1b[38;5;92mDiscord group \x1b[38;5;231m············· \x1b[38;5;32mdiscord.gg/pyBUFfzp")
	fmt.Fprintln(output, "  \x1b[38;5;231m• \x1b[38;5;92mTelegram group \x1b[38;5;231m··········· \x1b[38;5;32mt.me/galaxia_network")
	fmt.Fprintln(output, "")
}

func (c *CreditsCommand) AdminOnly() bool {
	return false
}

func (c *CreditsCommand) Aliases() []string {
	return []string{"Credits", "credits"}
}

func init() {
	CommandMap["credits"] = &CreditsCommand{}
}
