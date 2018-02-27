package dynamic

import (
	"github.com/fission/fission-workflows/pkg/api/invocation"
	"github.com/fission/fission-workflows/pkg/api/workflow"
	"github.com/fission/fission-workflows/pkg/fnenv/workflows"
	"github.com/fission/fission-workflows/pkg/types"
	"github.com/fission/fission-workflows/pkg/types/typedvalues"
	"github.com/fission/fission-workflows/pkg/types/validate"
	"github.com/gogo/protobuf/proto"
	"github.com/mitchellh/hashstructure"
	"github.com/sirupsen/logrus"
)

// Api that servers mainly as a function.Runtime wrapper that deals with the higher-level logic workflow-related logic.
type Api struct {
	wfApi  *workflow.Api
	wfiApi *invocation.Api
}

func NewApi(wfApi *workflow.Api, wfiApi *invocation.Api) *Api {
	return &Api{
		wfApi:  wfApi,
		wfiApi: wfiApi,
	}
}

func (ap *Api) AddDynamicTask(invocationId string, parentTaskId string, taskSpec *types.TaskSpec) error {

	// Transform TaskSpec into WorkflowSpec
	wfSpec := &types.WorkflowSpec{
		OutputTask: "main",
		Tasks: map[string]*types.TaskSpec{
			"main": taskSpec,
		},
		Dynamic:    true, // FUTURE: use internal as indicator to cleanup and hide these generated workflows
		ApiVersion: types.WorkflowApiVersion,
	}
	hash, err := hashstructure.Hash(wfSpec, nil)
	if err == nil {
		wfSpec.ForceId = string(hash)
	} else {
		logrus.Errorf("Failed to generate hash of workflow; defaulting to random id: %v", err)
	}
	return ap.addDynamicWorkflow(invocationId, parentTaskId, wfSpec, taskSpec)
}

func (ap *Api) AddDynamicWorkflow(invocationId string, parentTaskId string, workflowSpec *types.WorkflowSpec) error {
	// TODO add inputs to WorkflowSpec
	return ap.addDynamicWorkflow(invocationId, parentTaskId, workflowSpec, &types.TaskSpec{})
}

func (ap *Api) addDynamicWorkflow(invocationId string, parentTaskId string, wfSpec *types.WorkflowSpec,
	stubTask *types.TaskSpec) error {

	// Clean-up WorkflowSpec and submit
	sanitizeWorkflow(wfSpec)
	err := validate.WorkflowSpec(wfSpec)
	if err != nil {
		return err
	}
	wfId, err := ap.wfApi.Create(wfSpec)
	if err != nil && err != workflow.ErrWorkflowAlreadyExists {
		return err
	}

	// Generate Proxy Task
	proxyTaskSpec := proto.Clone(stubTask).(*types.TaskSpec)
	proxyTaskSpec.FunctionRef = wfId
	proxyTaskSpec.Input("_parent", typedvalues.ParseString(invocationId))
	proxyTaskId := parentTaskId + "_child"
	proxyTask := types.NewTask(proxyTaskId, proxyTaskSpec.FunctionRef)
	proxyTask.Spec = proxyTaskSpec
	proxyTask.Status.Status = types.TaskStatus_READY
	proxyTask.Status.FnRef = workflows.CreateFnRef(wfId)

	// Ensure that the only link of the dynamic task is with its parent
	proxyTaskSpec.Requires = map[string]*types.TaskDependencyParameters{
		parentTaskId: {
			Type: types.TaskDependencyParameters_DYNAMIC_OUTPUT,
		},
	}

	err = validate.TaskSpec(proxyTaskSpec)
	if err != nil {
		return err
	}

	// Submit added task to workflow invocation
	return ap.wfiApi.AddTask(invocationId, proxyTask)

}

func sanitizeWorkflow(v *types.WorkflowSpec) {
	if len(v.ApiVersion) == 0 {
		v.ApiVersion = types.WorkflowApiVersion
	}

	// ForceID is not supported for internal workflows
	v.ForceId = ""
}
