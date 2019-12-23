package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"

	flag "github.com/spf13/pflag"
	"golang.org/x/crypto/ssh/terminal"
)

// Data holds all command line flags and additional arguments.
type Data struct {
	c      *bool
	n      *bool
	b      *bool
	i      *bool
	m      *bool
	sha256 *bool
	sha512 *bool
	B      *bool
	C      *int
	r      *uint
	d      *bool
	s      *bool
	p      *bool
	D      *bool
	v      *bool
	user   string
	pw     []byte
	pwfile string
	args   []string
}

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
	var data Data
	var user string
	var pw []byte
	var pwfile string
	data.c = flag.BoolP("c", "c", false, "Create a new file.")
	data.n = flag.BoolP("n", "n", false, "Don't update file; display results on stdout.")
	data.b = flag.BoolP("b", "b", false, "Use the password from the command line rather than prompting for it.")
	data.i = flag.BoolP("i", "i", false, "Read password from stdin without verification (for script usage).")
	data.m = flag.BoolP("m", "m", false, "Force MD5 encryption of the password (default).")
	data.sha256 = flag.BoolP("2", "2", false, "Force SHA-256 crypt() hash of the password (very secure).")
	data.sha512 = flag.BoolP("5", "5", false, "Force SHA-512 crypt() hash of the password (very secure).")
	data.B = flag.BoolP("B", "B", false, "Force bcrypt encryption of the password (very secure).")
	data.C = flag.IntP("C", "C", 10, "Set the computing time used for the bcrypt algorithm\n(higher is more secure but slower, default: 10, valid: 4 to 17).")
	data.r = flag.UintP("r", "r", 5000, "Set the number of rounds used for the SHA-256, SHA-512 algorithms\n(higher is more secure but slower, default: 5000).")
	data.d = flag.BoolP("d", "d", false, "Force CRYPT encryption of the password (8 chars max, insecure).")
	data.s = flag.BoolP("s", "s", false, "Force SHA-1 encryption of the password (insecure).")
	data.p = flag.BoolP("p", "p", false, "Do not encrypt the password (plaintext, insecure).")
	data.D = flag.BoolP("D", "D", false, "Delete the specified user.")
	data.v = flag.BoolP("v", "v", false, "Verify password for the specified user.")
	flag.Parse()
	data.args = flag.Args()
	switch {
	case *data.v:
		lenArgs(2, args)
		data.pwfile = args[0]
		data.user = args[1]
		switch {
		case *data.v:
			data.pw = []byte(args[2])
		case *data.v:
			data.pw = readStdin()
		default:
			data.pw = getPW()
		}
		// TODO: verify password for specified user in specified passwordfile
	case *data.D:
		lenArgs(2, args)
		pwfile = args[0]
		user = args[1]
		// TODO: insert delete user function
	case *data.n:
		if *data.c {
			usage()
		}
		if *data.b {
			lenArgs(2, args)
			user = args[0]
			pw = []byte(args[1])
			return
		}
		lenArgs(1, args)
		user = args[0]
		if *data.i {
			pw = readStdin()
			return
		}
		pw = getPW()
		return
	case *data.b:
		lenArgs(3, args)
		pwfile = args[0]
		user = args[1]
		pw = []byte(args[2])
	case *data.i:
		pw = readStdin()
	}

	lenArgs(2, args)
	//fmt.Println(*c, *n, *b, *i, *m, *sha256, *sha512, *B, *C, *r, *d, *s, *p, *D, *v)
	fmt.Println(pwfile, user, pw)
	fmt.Println(args)
}

func (data Data) calcHash() {

}

func readStdin() (pw []byte) {
	reader := bufio.NewReader(os.Stdin)
	pw, err := reader.ReadBytes('\n')
	if err != nil {
		log.Fatalln(err)
	}
	return
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
	flag.PrintDefaults()
	os.Exit(2)
}

func lenArgs(l int, args []string) {
	if len(args) < l {
		usage()
	}
}
