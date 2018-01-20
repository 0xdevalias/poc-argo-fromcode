package main

import (
	"encoding/json"
	"log"

	argoCmd "github.com/argoproj/argo/cmd/argo/commands"
	wfv1 "github.com/argoproj/argo/pkg/apis/workflow/v1alpha1"
	"github.com/argoproj/argo/pkg/client/clientset/versioned/typed/workflow/v1alpha1"
	"github.com/argoproj/argo/workflow/common"

	"github.com/ghodss/yaml"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func main() {
	namespace := "poc-argo-fromcode"

	log.Printf("Defining argo workflow..")
	workflowDef := exampleWorkflow()

	log.Printf("Creating argo client..")
	// TODO: Remove the need to rely on argoCmd
	wfClient := argoCmd.InitWorkflowClient(namespace)

	log.Printf("Submitting argo workflow..")
	workflow, err := submitWorkflow(&workflowDef, wfClient)
	if err != nil {
		log.Fatalf("Workflow manifest failed submission: %v", err)
	}

	printWorkflow("yaml", workflow)
}

// Example argo workflow.
// This could also be built up automagically by reading in a YAML/etc file and unmarshalling it
// Ref:
//   https://godoc.org/github.com/argoproj/argo/pkg/apis/workflow/v1alpha1
//   https://godoc.org/github.com/argoproj/argo/pkg/apis/workflow/v1alpha1#Workflow
//   https://godoc.org/github.com/argoproj/argo/pkg/apis/workflow/v1alpha1#WorkflowSpec
//   https://godoc.org/github.com/argoproj/argo/pkg/apis/workflow/v1alpha1#Template
//   https://godoc.org/k8s.io/api/core/v1#Container
func exampleWorkflow() wfv1.Workflow {
	return wfv1.Workflow{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "argoproj.io/v1alpha1",
			Kind:       "Workflow",
		},
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: "poc-argo-fromcode-",
		},
		Spec: wfv1.WorkflowSpec{
			Entrypoint: "helloWorld",
			Templates: []wfv1.Template{
				{
					Name: "helloWorld",
					Container: &corev1.Container{
						Image:   "alpine:latest",
						Command: []string{"echo"},
						Args:    []string{"poc-argo-fromcode world!"},
					},
				},
			},
		},
	}
}

// Validate and submit the argo workflow
// Ref:
//   https://godoc.org/github.com/argoproj/argo/pkg/client/clientset/versioned/typed/workflow/v1alpha1#WorkflowInterface
//   https://github.com/argoproj/argo/blob/master/cmd/argo/commands/submit.go#L90-L132
func submitWorkflow(wf *wfv1.Workflow, wfClient v1alpha1.WorkflowInterface) (*wfv1.Workflow, error) {
	err := common.ValidateWorkflow(wf)
	if err != nil {
		return nil, err
	}

	created, err := wfClient.Create(wf)
	if err != nil {
		return nil, err
	}

	return created, nil
}

// Ref: https://github.com/argoproj/argo/blob/a0b2d78c869f277c20c4cd3ba18b8d2688674e54/cmd/argo/commands/get.go#L53-L68
func printWorkflow(outFmt string, wf *wfv1.Workflow) {
	switch outFmt {
	case "name":
		log.Println(wf.ObjectMeta.Name)
	case "json":
		outBytes, _ := json.MarshalIndent(wf, "", "    ")
		log.Printf("\n%s", string(outBytes))
	case "yaml":
		outBytes, _ := yaml.Marshal(wf)
		log.Print(string(outBytes))
	default:
		log.Fatalf("Unknown output format: %s", outFmt)
	}
}
