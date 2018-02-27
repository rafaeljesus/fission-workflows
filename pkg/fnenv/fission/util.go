package fission

import (
	controller "github.com/fission/fission/controller/client"
	executor "github.com/fission/fission/executor/client"
)

func SetupRuntime(executorAddr string) *FunctionEnv {
	client := executor.MakeClient(executorAddr)
	return NewFunctionEnv(client)
}

func SetupResolver(controllerAddr string) *Resolver {
	controllerClient := controller.MakeClient(controllerAddr)
	return NewResolver(controllerClient)
}

// FUTURE: cleanup function
