// github.com/ProtonMail/gopenpgp/v3
// go-git

package main

import (
	"fmt"
)

func usage() {
	fmt.Println("" +
		"Usage: cl [Options] <command> \n" +
		"Password manager\n" +
		"\nOptions:" +
		"-h,--help\tShow this help message\n" +
		"-v,--version\t\tShow version\n" +
		"-q,--quiet\t\tDon't commit changes\n" +
		"-m,--commit-message\tSpecify commit message\n" +
		"\nCommands:\n" +
		"init\t\t\tInitialize a new git repository in the current working directory(Similar to git init)\n" +
		"list\t\t\tList all passwords in the vault\n" +
		"add <name>\t\tAdd a new password\n" +
		"remove <name>\t\tRemove a password\n" +
		"show <name>\t\tShow content of password")
}

func main() {
	usage()
}
