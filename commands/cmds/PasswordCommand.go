package cmds

import (
	"arismcnc/database"
	"arismcnc/utils"
	"fmt"
	"io"

	"github.com/gliderlabs/ssh"
)

type PasswordCommand struct{}

func (c *PasswordCommand) Name() string {
	return "password"
}

func (c *PasswordCommand) Execute(session ssh.Session, db *database.Database, args []string, output io.Writer) {
	// Ask the user for the current password
	fmt.Fprintf(output, "Enter your old password (leave empty to cancel): ")
	currentPassword, err := utils.ReadLine(session)
	if err != nil {
		fmt.Fprintln(output, "Error reading input:", err)
		return
	}

	// Cancel if the user leaves the input empty
	if currentPassword == "" {
		fmt.Fprintln(output, "\033[37;1mPassword change canceled.")
		return
	}

	// Verify the current password
	isValid, err := db.VerifyPassword(session.User(), currentPassword)
	if err != nil {
		fmt.Fprintln(output, "\033[31;1mError verifying password:", err)
		return
	}

	if !isValid {
		fmt.Fprintln(output, "\033[31;1mInvalid current password. Please try again.")
		return
	}

	// Ask the user for the new password
	fmt.Fprintf(output, "Enter new password: ")
	password1, err := utils.ReadLine(session)
	if err != nil {
		fmt.Fprintln(output, "Error reading input:", err)
		return
	}

	// Ask the user to confirm the new password
	fmt.Fprintf(output, "Confirm new password: ")
	password2, err := utils.ReadLine(session)
	if err != nil {
		fmt.Fprintln(output, "Error reading input:", err)
		return
	}

	// Check if both passwords match
	if password1 != password2 {
		fmt.Fprintln(output, "\033[37;1mPasswords do not match. Please try again.")
		return
	}

	// Change the password in the database
	err = db.ChangePassword(session.User(), password1)
	if err != nil {
		fmt.Fprintln(output, "\033[37;1mError changing password:", err)
		return
	}

	fmt.Fprintln(output, "\033[37;1mPassword changed successfully!")
}

func (c *PasswordCommand) AdminOnly() bool {
	return false
}

func (c *PasswordCommand) Aliases() []string {
	return []string{"password", "passwd"}
}

// Register command in CommandMap
func init() {
	CommandMap["password"] = &PasswordCommand{}
}
