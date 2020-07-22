package callbeater


import(
	_ "errors"
	"fmt"
        "sync"
        "os/exec"
        "log"
        "time"
        "syscall"
)




func init() {
}

// Input define a snmptrap input
type Callbeat struct {
    sync.Mutex
    started bool
    cmd  * exec.Cmd 
    cmdpkt * exec.Cmd 
}

func New() (* Callbeat)  {

    cmd := exec.Cmd {
        Path:  "/root/beats/filebeat/filebeat",
        Dir:  "/root/beats/filebeat/",
    }

    cmdpkt := exec.Cmd {
        Path:  "/root/beats/packetbeat/packetbeat",
        Dir:  "/root/beats/filebeat/",
    }

    return &Callbeat{
        started: false,
        cmd :  &cmd,
        cmdpkt : &cmdpkt,
        }
}



func (cb * Callbeat) Callfilebeat () {
    if err := cb.cmd.Start(); err != nil { 
        log.Panic(err) 
    } 
    cb.cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true} 

    fmt.Println("Start child process with pid", cb.cmd.Process.Pid) 

}


func (cb* Callbeat)   Stopfilebeat () {

    if cb.cmd.Process == nil {
         return
    }

    go func() {
        if err := cb.cmd.Wait(); err != nil {
            fmt.Printf("Child process %d exit with err: %v\n", cb.cmd.Process.Pid, err)
            fmt.Printf("process state : %v\n", cb.cmd.ProcessState)

            cmd := exec.Cmd {
                Path:  "/root/beats/filebeat/filebeat",
                Dir:  "/root/beats/filebeat/",
            }
            cb.cmd = &cmd
        }
         
    }()


    if err := cb.cmd.Process.Kill(); err != nil {
        log.Fatal("failed to kill process: ", err)
    }

    time.Sleep(3 * time.Second)

    fmt.Printf("stopfilebeat over\n")


}


func (cb * Callbeat) Callpacketbeat() {
    
    if err := cb.cmdpkt.Start(); err != nil {
        log.Panic(err)
    }
    cb.cmdpkt.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

    fmt.Println("Start child process with pid", cb.cmdpkt.Process.Pid)
}


func (cb * Callbeat) Stoppacketbeat() {

    if cb.cmdpkt.Process == nil {
         return
    }

    go func() {
        if err := cb.cmdpkt.Wait(); err != nil {
            fmt.Printf("Child process %d exit with err: %v\n", cb.cmdpkt.Process.Pid, err)
            fmt.Printf("process state : %v\n", cb.cmdpkt.ProcessState)

            cmdpkt := exec.Cmd {
                Path:  "/root/beats/packetbeat/packetbeat",
                Dir:  "/root/beats/packetbeat/",
            }
            cb.cmdpkt = &cmdpkt
        }

    }()



    if err := cb.cmdpkt.Process.Kill(); err != nil {
        log.Fatal("failed to kill process: ", err)
    }    

    time.Sleep(5 * time.Second)

    fmt.Printf("stoppacketbeat over\n")
}


