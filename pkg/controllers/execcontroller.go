// Copyright 2018 Oracle and/or its affiliates. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package controller

import (
	apps "k8s.io/api/apps/v1beta1"
	"statefulguardian/pkg/apis/statefulguardian/v1alpha1"
	"context"
	"sync"
	"strconv"
	"regexp"
	"strings"
	"bytes"
	"github.com/golang/glog"
	"github.com/pkg/errors"
	utilexec "k8s.io/utils/exec"
)

var errorRegex = regexp.MustCompile(`Traceback.*\n(?:  (.*)\n){1,}(?P<type>[\w\.]+)\: (?P<message>.*)`)

const (
   KUBECTL_PATH="/kubectl"
   KUBECONFIG_PATH="/etc/kubeconfig/config"
)

type ExecController struct {
	SG *v1alpha1.Statefulguardian
	r *runner
}

// NewController creates a new SGController.
func NewExecController( sg *v1alpha1.Statefulguardian ) *ExecController {
	ec := ExecController {
		SG: sg,
		r: New(utilexec.New())}

	return &ec
}

// On cluster ready
func (e *ExecController) ClusterInit(ctx context.Context) error {
	_, err:=e.RunClusterOps(ctx, &e.SG.Spec.App, e.SG.Spec.Lifecycle.Init.Cluster)

	return err
}

func (e *ExecController) PodInit(ctx context.Context, pod string) error {
	_, err:=e.RunPodOps(ctx, &e.SG.Spec.App, e.SG.Spec.Lifecycle.Init.Pod, pod)

	return err
}

// On cluster quit
func (e *ExecController) ClusterQuit(ctx context.Context) error {
	_, err:=e.RunClusterOps(ctx, &e.SG.Spec.App, e.SG.Spec.Lifecycle.Quit.Cluster)

	return err
}

func (e *ExecController) PodQuit(ctx context.Context, pod string) error {
	_, err:=e.RunPodOps(ctx, &e.SG.Spec.App, e.SG.Spec.Lifecycle.Quit.Pod, pod)

	return err
}

func (e *ExecController) RunPodOps(ctx context.Context,ss *apps.StatefulSet, ops []v1alpha1.Operation, pod string) (string, error) {
	ns:=ss.ObjectMeta.Namespace
	for _, op:= range(ops) {
		container := op.Container
		cmds:=op.Cmd
		glog.V(6).Infof("Running operation: %v in pod %v container %v", op.Name, pod, container)

		output, err := e.r.run(ctx, ns, pod, container, cmds)
		if err != nil {
			return "Run operation "+op.Name+" error, OUTPUT:"+string(output), err
		}
	}

	return "", nil

}

func (e *ExecController) RunClusterOps(ctx context.Context,ss *apps.StatefulSet, ops []v1alpha1.Operation) (string, error) {
	glog.Infof("Running cluster operations")
	replis:= int(*ss.Spec.Replicas)
	ns:=ss.ObjectMeta.Namespace
	ssname:=ss.ObjectMeta.Name

	// if not set, chose the first container as default container
	for _, op:= range(ops) {
		hosts := strings.Fields(op.Container)
		pod:=hosts[0]
		cmds:=op.Cmd
		if pod=="__ALL__" {
			for i:=0; i<replis; i++ {
				rpod := ssname+"-"+strconv.Itoa(i)
				container := ""
				if len(hosts)==2 {
					container=hosts[1]
				}
	        		glog.Infof("Running operation: %v in pod %v container %v", op.Name, rpod, container)
				output, err := e.r.run(ctx, ns, rpod, container, cmds)
				if err != nil {
	        		        glog.Infof("Running operation error: cmd: %v output: %v", cmds, string(output))
					return "Run operation "+op.Name+" error, OUTPUT:"+string(output), err
				}
			}
		} else {
			container := ""
			if len(hosts)==2 {
				container=hosts[1]
			}
	        	glog.Infof("Running operation: %v in pod %v container %v", op.Name, pod, container)
			output, err := e.r.run(ctx, ns, pod, container, cmds)
			if err != nil {
				return "Run operation "+op.Name+" error, OUTPUT:"+string(output), err
			}
		}
	}

	return "", nil
}

type runner struct {
	mu   sync.Mutex
	exec utilexec.Interface
}

func New(exec utilexec.Interface) *runner{
	return &runner{exec: exec}
}

func (r *runner) run(ctx context.Context, ns string, pod string, container string, cmds []string) ([]byte, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	args := []string{"--kubeconfig", KUBECONFIG_PATH, "exec", "-n", ns, pod}

	if len(container)>0 {
		args = append(args, "-c", container)
	}
	args = append(args, "--")
	args = append(args, cmds...)

	cmd := r.exec.CommandContext(ctx, KUBECTL_PATH, args...)

	cmd.SetStdout(stdout)
	cmd.SetStderr(stderr)

	glog.Infof("Running command: %v %v", KUBECONFIG_PATH, args)
	err := cmd.Run()
	glog.Infof("    stdout: %s\n    stderr: %s\n    err: %s", stdout, stderr, err)
	if err != nil {
		underlying := NewErrorFromStderr(stderr.String())
		if underlying != nil {
			return nil, errors.WithStack(underlying)
		}
	}

	return stdout.Bytes(), err
}

type Error struct {
	error
	Type    string
	Message string
}


// one is present.
func NewErrorFromStderr(stderr string) error {
	matches := errorRegex.FindAllStringSubmatch(stderr, -1)
	if len(matches) == 0 {
		return nil
	}
	result := make(map[string]string)
	for i, name := range errorRegex.SubexpNames() {
		if i != 0 && name != "" {
			result[name] = matches[len(matches)-1][i]
		}
	}
	return &Error{
		Type:    result["type"],
		Message: result["message"],
	}
}
