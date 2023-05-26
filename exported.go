package ServiceCore

import (
	"errors"
	"fmt"
	"go/token"
	"log"
	"reflect"
	"sync"
)

func ExportList(rcvr any) []string {
	server := new(Server)

	server.register(rcvr, "", false)

	exportList := []string{}

	server.serviceMap.Range(func(key, value interface{}) bool {
		service := value.(*Service)
		name := key.(string)

		for k, _ := range service.Method {
			value := k
			exportList = append(exportList, fmt.Sprintf("%s.%s", name, value))
		}

		return true
	})
	return exportList
}

// ------------------ copied from net/rpc/server.go ------------------
// golang doesn't export those functions, so we have to copy them here to be able to find all the exported methods for the rpc registration
type Service struct {
	Name   string                 // name of service
	Rcvr   reflect.Value          // receiver of methods for the service
	Typ    reflect.Type           // type of the receiver
	Method map[string]*MethodType // registered methods
}

var typeOfError = reflect.TypeOf((*error)(nil)).Elem()

type MethodType struct {
	Method    reflect.Method
	ArgType   reflect.Type
	ReplyType reflect.Type
}

// Server represents an RPC Server.
type Server struct {
	serviceMap sync.Map // map[string]*service
}

func isExportedOrBuiltinType(t reflect.Type) bool {
	for t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	// PkgPath will be non-empty even for an exported type,
	// so we need to check the type name as well.
	return token.IsExported(t.Name()) || t.PkgPath() == ""
}

func (server *Server) register(rcvr any, name string, useName bool) error {
	s := new(Service)
	s.Typ = reflect.TypeOf(rcvr)
	s.Rcvr = reflect.ValueOf(rcvr)
	sname := name
	if !useName {
		sname = reflect.Indirect(s.Rcvr).Type().Name()
	}
	if sname == "" {
		s := "rpc.Register: no service name for type " + s.Typ.String()
		log.Print(s)
		return errors.New(s)
	}
	if !useName && !token.IsExported(sname) {
		s := "rpc.Register: type " + sname + " is not exported"
		log.Print(s)
		return errors.New(s)
	}
	s.Name = sname

	// Install the methods
	s.Method = suitableMethods(s.Typ)

	if len(s.Method) == 0 {
		str := ""

		// To help the user, see if a pointer receiver would work.
		method := suitableMethods(reflect.PointerTo(s.Typ))
		if len(method) != 0 {
			str = "rpc.Register: type " + sname + " has no exported methods of suitable type (hint: pass a pointer to value of that type)"
		} else {
			str = "rpc.Register: type " + sname + " has no exported methods of suitable type"
		}
		log.Print(str)
		return errors.New(str)
	}

	if _, dup := server.serviceMap.LoadOrStore(sname, s); dup {
		return errors.New("rpc: service already defined: " + sname)
	}
	return nil
}

// suitableMethods returns suitable Rpc methods of typ. It will log
// errors if logErr is true.
func suitableMethods(typ reflect.Type) map[string]*MethodType {
	methods := make(map[string]*MethodType)
	for m := 0; m < typ.NumMethod(); m++ {
		method := typ.Method(m)
		mtype := method.Type
		mname := method.Name
		// Method must be exported.
		if !method.IsExported() {
			continue
		}
		// Method needs three ins: receiver, *args, *reply.
		if mtype.NumIn() != 3 {
			log.Printf("rpc.Register: method %q has %d input parameters; needs exactly three\n", mname, mtype.NumIn())
			continue
		}
		// First arg need not be a pointer.
		argType := mtype.In(1)
		if !isExportedOrBuiltinType(argType) {
			log.Printf("rpc.Register: argument type of method %q is not exported: %q\n", mname, argType)
			continue
		}
		// Second arg must be a pointer.
		replyType := mtype.In(2)
		if replyType.Kind() != reflect.Pointer {
			log.Printf("rpc.Register: reply type of method %q is not a pointer: %q\n", mname, replyType)
			continue
		}
		// Reply type must be exported.
		if !isExportedOrBuiltinType(replyType) {
			log.Printf("rpc.Register: reply type of method %q is not exported: %q\n", mname, replyType)
			continue
		}
		// Method needs one out.
		if mtype.NumOut() != 1 {
			log.Printf("rpc.Register: method %q has %d output parameters; needs exactly one\n", mname, mtype.NumOut())
			continue
		}
		// The return type of the method must be error.
		if returnType := mtype.Out(0); returnType != typeOfError {
			log.Printf("rpc.Register: return type of method %q is %q, must be error\n", mname, returnType)
			continue
		}
		methods[mname] = &MethodType{Method: method, ArgType: argType, ReplyType: replyType}
	}
	return methods
}
