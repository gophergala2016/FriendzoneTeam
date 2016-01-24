package ssh

import (
    "bytes"
    "fmt"
    "log"
    "golang.org/x/crypto/ssh"
    "sync"
)

// Conecta al SSH y ejecuta el comando especificado
func Conekta(user, password, host,cmd string)(string, error){
    var wg sync.WaitGroup
    result := make(chan string)
    hosts := make([]string, 1)
    hosts[0] = host
    sshConfig := &ssh.ClientConfig{
        User: user,
        Auth: []ssh.AuthMethod{
            ssh.Password(password),
        },
    }
    wg.Add(len(hosts))
    for _ , hostname := range hosts {
        go func(hostname string) {
            defer wg.Done()
            result <- executeCmd(cmd,  hostname, sshConfig)
        }(hostname)
    }
    
    for i := 0; i < len(hosts); i++ {
        select {
        case res := <-result:
            fmt.Print(res)
            return res, nil
        }
    }
    
    wg.Wait()
    return <-result, nil
}

func executeCmd(cmd, hostname string, config *ssh.ClientConfig) string {
    conn, err := ssh.Dial("tcp", hostname+":22", config)
    if err != nil {
        log.Fatalf("Failed to dial: %s", err)
    }
    session, err := conn.NewSession()
    if err != nil {
        log.Fatalf("Failed to create session: %s", err)
    }
    defer session.Close()

    var stdoutBuf bytes.Buffer
    session.Stdout = &stdoutBuf
    session.Run(cmd)

    return hostname + ": " + stdoutBuf.String()
}