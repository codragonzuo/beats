package rpcserver


import (
    "context"
    "fmt"
    "strconv"
    "github.com/codragonzuo/beats/collector/commu/sharedlib"
    "github.com/codragonzuo/beats/collector/commu/commu"
    "github.com/codragonzuo/beats/collector/callbeater"
)

type CalculatorHandler struct {
    log map[int]*sharedlib.SharedStruct
    cb  * callbeater.Callbeat
}

func NewCalculatorHandler() *CalculatorHandler {
    //
    return &CalculatorHandler{log: make(map[int]*sharedlib.SharedStruct), cb: callbeater.New()}
}

func (p *CalculatorHandler) Ping(ctx context.Context) (err error) {
    fmt.Print("ping()\n")
    return nil
}

func (p *CalculatorHandler) Add(ctx context.Context, num1 int32, num2 int32) (retval17 int32, err error) {
    fmt.Print("add(", num1, ",", num2, ")\n")
    return num1 + num2, nil
}

func (p *CalculatorHandler) Calculate(ctx context.Context, logid int32, w *commu.Work) (val int32, err error) {
    fmt.Print("calculate(", logid, ", {", w.Op, ",", w.Num1, ",", w.Num2, "})\n")
    switch w.Op {
    case commu.Operation_ADD:
        val = w.Num1 + w.Num2
        break
    case commu.Operation_SUBTRACT:
        val = w.Num1 - w.Num2
        break
    case commu.Operation_MULTIPLY:
        val = w.Num1 * w.Num2
        break
    case commu.Operation_DIVIDE:
        if w.Num2 == 0 {
            ouch := commu.NewInvalidOperation()
            ouch.WhatOp = int32(w.Op)
            ouch.Why = "Cannot divide by 0"
            err = ouch
            return
        }
        val = w.Num1 / w.Num2
        break
    default:
        ouch := commu.NewInvalidOperation()
        ouch.WhatOp = int32(w.Op)
        ouch.Why = "Unknown operation"
        err = ouch
        return
    }
    entry := sharedlib.NewSharedStruct()
    entry.Key = logid
    entry.Value = strconv.Itoa(int(val))
    k := int(logid)
    /*
       oldvalue, exists := p.log[k]
       if exists {
         fmt.Print("Replacing ", oldvalue, " with ", entry, " for key ", k, "\n")
       } else {
         fmt.Print("Adding ", entry, " for key ", k, "\n")
       }
    */
    p.log[k] = entry
    return val, err
}

func  (p *CalculatorHandler) Doconfig(ctx context.Context, commandid int32, operationid int32, jsonconfig string) (val int32, err error) {

    fmt.Printf("operationid=%d\n", operationid)
    fmt.Printf("jsonconfig=%s\n", jsonconfig)
    switch (operationid) {
    case 1: 
        p.cb.Callfilebeat()
        break

    case 2:
        p.cb.Stopfilebeat()
        break

    case 3:
        p.cb.Callpacketbeat()
        break

    case 4:
        p.cb.Stoppacketbeat()
        break

    case 5:
        fmt.Printf("Doconfig commandid=5\n")
        break
    default:
        break
    }
    fmt.Printf("Doconfig called Over\n-----------------------\n")
    return 1,nil
}

func (p *CalculatorHandler) GetStruct(ctx context.Context, key int32) (*sharedlib.SharedStruct, error) {
    fmt.Print("getStruct(", key, ")\n")
    v, _ := p.log[int(key)]
    return v, nil
}

func (p *CalculatorHandler) Zip(ctx context.Context) (err error) {
    fmt.Print("zip()\n")
    return nil
}


