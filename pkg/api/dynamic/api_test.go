package dynamic

import (
	"testing"

	"github.com/fission/fission-workflows/pkg/api/invocation"
	"github.com/fission/fission-workflows/pkg/api/workflow"
	"github.com/fission/fission-workflows/pkg/fes/backend/mem"
	"github.com/fission/fission-workflows/pkg/fnenv"
	"github.com/fission/fission-workflows/pkg/fnenv/mock"
	"github.com/fission/fission-workflows/pkg/types"
	"github.com/stretchr/testify/assert"
)

func TestApi_AddDynamicTask(t *testing.T) {
	api, es := setupApi()

	wfSpec := types.NewWorkflowSpec().
		SetOutput("task1").
		AddTask("task1", types.NewTaskSpec("someFn"))
	wfSpec.ForceId = "forcedId"
	id, err := api.wfApi.Create(wfSpec)
	assert.NoError(t, err)

	wfiSpec := types.NewWorkflowInvocationSpec(id)
	invocationId, err := api.wfiApi.Invoke(wfiSpec)
	assert.NoError(t, err)

	taskSpec := types.NewTaskSpec("someFn")
	err = api.AddDynamicTask(invocationId, "task-parent", taskSpec)
	assert.NoError(t, err)

	s := es.Snapshot()
	assert.Len(t, s, 3)
}

func setupApi() (*Api, *mem.Backend) {
	es := mem.NewBackend()
	resolver := fnenv.NewMetaResolver(map[string]fnenv.RuntimeResolver{
		"mock": mock.NewResolver(),
	})

	wfiApi := invocation.NewApi(es)
	wfApi := workflow.NewApi(es, resolver)

	return NewApi(wfApi, wfiApi), es
}
