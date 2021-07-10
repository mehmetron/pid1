package userPort

import (
	"fmt"
	"os/exec"
)

func GetOpenPortLib() uint16 {
	//fmt.Println("path to tcp tab ", pathTCPTab)
	//// TCP sockets
	//socks, err := TCP6Socks(NoopFilter)
	//if err != nil {
	//	fmt.Println("86 ", err)
	//	//return err
	//}
	//for _, e := range socks {
	//	fmt.Printf("%v\n", e)
	//}
	//fmt.Println("tcpsocks 92 ", socks)

	// get only listening TCP sockets
	tabs, err := TCP6Socks(func(s *SockTabEntry) bool {
		return s.State == Listen
	})
	if err != nil {
		fmt.Println("98 ", err)
		//return err
	}
	var demoPort uint16
	for _, e := range tabs {
		fmt.Printf("%v\n", e)
		fmt.Printf("%v\n", e.LocalAddr)
		fmt.Printf("%v\n", e.LocalAddr.Port)
		fmt.Printf("%v\n", e.RemoteAddr)
		fmt.Printf("%v\n", e.ino)
		fmt.Printf("%v\n", e.State)
		fmt.Printf("%v\n", e.Process)

		var pid1Port uint16 = 8080
		if e.LocalAddr.Port != pid1Port {
			demoPort = e.LocalAddr.Port
		}
	}
	fmt.Println("tcpsocks 105 ", tabs)

	return demoPort
}

func GetOpenPort() string {
	cmd := exec.Command("netstat", "-peanut")
	grep := exec.Command("grep", "bob")

	// Get ps's stdout and attach it to grep's stdin.
	pipe, _ := cmd.StdoutPipe()
	defer pipe.Close()

	grep.Stdin = pipe

	// Run ps first.
	cmd.Start()

	// Run and get the output of grep.
	res, _ := grep.Output()

	fmt.Println(string(res))
	return string(res)

	//var stdoutBuf, stderrBuf bytes.Buffer
	//
	//cmd.Stdout = &stdoutBuf
	//cmd.Stderr = &stderrBuf
	//
	//cmd.Run()
	//
	//if ctx.Err() == context.DeadlineExceeded {
	//	fmt.Println("Command was killed")
	//	fmt.Fprintf(w, "Command was killed")
	//	return
	//}
	//outStr, errStr := string(stdoutBuf.Bytes()), string(stderrBuf.Bytes())
	//fmt.Printf("\nout:\n%s\nerr:\n%s\n", outStr, errStr)
	//
	//fmt.Fprintf(w, fmt.Sprintf("Output: %s \n Err: %s", outStr, errStr))

}
