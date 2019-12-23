package main

import (
	"bytes"
	"fmt"
	"log"
	"os"

	flag "github.com/spf13/pflag"
	"golang.org/x/crypto/ssh/terminal"
)

// Goal: Rebuild htpasswd in go :)
// Usage:
// 	htpasswd [-cimBdpsDv] [-C cost] passwordfile username
// 	htpasswd -b[cmBdpsDv] [-C cost] passwordfile username password

// 	htpasswd -n[imBdps] [-C cost] username
// 	htpasswd -nb[mBdps] [-C cost] username password
//  -c  Create a new file.
//  -n  Don't update file; display results on stdout.
//  -b  Use the password from the command line rather than prompting for it.
//  -i  Read password from stdin without verification (for script usage).
//  -m  Force MD5 encryption of the password (default).
//  -2  Force SHA-256 crypt() hash of the password (very secure).
//  -5  Force SHA-512 crypt() hash of the password (very secure).
//  -B  Force bcrypt encryption of the password (very secure).
//  -C  Set the computing time used for the bcrypt algorithm
//      (higher is more secure but slower, default: 5, valid: 4 to 17).
//  -r  Set the number of rounds used for the SHA-256, SHA-512 algorithms
//      (higher is more secure but slower, default: 5000).
//  -d  Force CRYPT encryption of the password (8 chars max, insecure).
//  -s  Force SHA-1 encryption of the password (insecure).
//  -p  Do not encrypt the password (plaintext, insecure).
//  -D  Delete the specified user.
//  -v  Verify password for the specified user.
// On other systems than Windows and NetWare the '-p' flag will probably not work.
// The SHA-1 algorithm does not use a salt and is less secure than the MD5 algorithm
func main() {
	var user string
	var pw string
	var pwfile string
	c := flag.BoolP("c", "c", false, "Create a new file.")
	n := flag.BoolP("n", "n", false, "Don't update file; display results on stdout.")
	b := flag.BoolP("b", "b", false, "Use the password from the command line rather than prompting for it.")
	i := flag.BoolP("i", "i", false, "Read password from stdin without verification (for script usage).")
	m := flag.BoolP("m", "m", false, "Force MD5 encryption of the password (default).")
	sha256 := flag.BoolP("2", "2", false, "Force SHA-256 crypt() hash of the password (very secure).")
	sha512 := flag.BoolP("5", "5", false, "Force SHA-512 crypt() hash of the password (very secure).")
	B := flag.BoolP("B", "B", false, "Force bcrypt encryption of the password (very secure).")
	C := flag.IntP("C", "C", 10, "Set the computing time used for the bcrypt algorithm\n(higher is more secure but slower, default: 10, valid: 4 to 17).")
	r := flag.UintP("r", "r", 5000, "Set the number of rounds used for the SHA-256, SHA-512 algorithms\n(higher is more secure but slower, default: 5000).")
	d := flag.BoolP("d", "d", false, "Force CRYPT encryption of the password (8 chars max, insecure).")
	s := flag.BoolP("s", "s", false, "Force SHA-1 encryption of the password (insecure).")
	p := flag.BoolP("p", "p", false, " Do not encrypt the password (plaintext, insecure).")
	D := flag.BoolP("D", "D", false, "Delete the specified user.")
	v := flag.BoolP("v", "v", false, "Verify password for the specified user.")
	flag.Parse()
	args := flag.Args()
	if *n {
		if *b {
			lenArgs(2, args)
			user = args[0]
			pw = args[1]
			return
		}
		lenArgs(1, args)
		user = args[0]
		getPW()
		return
	}
	if *b {
		lenArgs(3, args)
		pwfile = args[0]
		user = args[1]
		pw = args[2]
	}

	lenArgs(2, args)
	fmt.Println(*c, *n, *b, *i, *m, *sha256, *sha512, *B, *C, *r, *d, *s, *p, *D, *v)
	fmt.Println(pwfile, user, pw)
	fmt.Println(args)
}

func getPW() (pw []byte) {
	fmt.Print("New password: ")
	pw, err := terminal.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println()
	fmt.Print("Re-type new password: ")
	pw2, err := terminal.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println()
	if bytes.Compare(pw, pw2) != 0 {
		fmt.Println("htpwd: password verification error")
		os.Exit(3)
	}
	return
}

func usage() {
	text := `Usage:
	htpasswd [-cimBdpsDv] [-C cost] passwordfile username
	htpasswd -b[cmBdpsDv] [-C cost] passwordfile username password
	`
	fmt.Println(text)
}

func lenArgs(l int, args []string) {
	if len(args) < l {
		usage()
		flag.PrintDefaults()
		os.Exit(2)
	}
}
