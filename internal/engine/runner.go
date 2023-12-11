package engine

import (
	"fmt"
	"sync"

	"github.com/bondyra/swamp/internal/common"
	"github.com/bondyra/swamp/internal/reader"
	"github.com/google/uuid"
)

type ExecutionRunner func(*ExecutionPlan) *ExecutionResult

func ParallelExecutionRunner(ep *ExecutionPlan) *ExecutionResult {
	workerGroup := sync.WaitGroup{}
	taskGroup := sync.WaitGroup{}
	taskQueue := make(chan *ExecutionTask)
	taskResultQueue := make(chan *ExecutionTaskResult)
	outputQueue := make(chan *ExecutionResult)
	for i := 0; i < 10; i++ {
		workerGroup.Add(1)
		go worker(taskQueue, taskResultQueue, &workerGroup)
	}
	scheduleTasks(taskQueue, &taskGroup, rootTask(ep))
	go resultHandler(ep, taskQueue, taskResultQueue, outputQueue, &taskGroup)

	taskGroup.Wait()
	close(taskQueue)
	workerGroup.Wait()
	close(taskResultQueue)

	return <-outputQueue
}

func rootTask(ep *ExecutionPlan) []*ExecutionTask {
	return []*ExecutionTask{
		{
			Id:             "ROOT",
			Type:           ep.root.Type.Type,
			Profiles:       ep.root.Profiles,
			Attrs:          ep.root.Attrs,
			BaseConditions: ep.root.Conditions,
			Reader:         ep.root.Reader,
			executionNode:  ep.root,
		},
	}
}

func scheduleTasks(taskInQueue chan *ExecutionTask, taskGroup *sync.WaitGroup, tasks []*ExecutionTask) {
	for _, task := range tasks {
		taskGroup.Add(1)
		taskInQueue <- task
	}
}

func worker(taskQueue chan *ExecutionTask, taskResultQueue chan *ExecutionTaskResult, workerGroup *sync.WaitGroup) {
	defer workerGroup.Done()
	for {
		task, open := <-taskQueue
		if !open {
			return
		}
		items, err := task.Reader.GetItems(task.Type, task.Profiles, task.Attrs, append(task.BaseConditions, task.LinkConditions...))
		taskResultQueue <- &ExecutionTaskResult{
			Id:    task.Id,
			Task:  task,
			Items: items,
			Err:   err,
		}
	}
}

func resultHandler(ep *ExecutionPlan, taskQueue chan *ExecutionTask, resultQueue chan *ExecutionTaskResult, outputQueue chan *ExecutionResult, taskGroup *sync.WaitGroup) {
	result := newExecutionResult()
	for {
		taskResult, open := <-resultQueue
		if !open {
			outputQueue <- result
			return
		}
		newTasks := processTaskResult(ep, taskResult, result)
		scheduleTasks(taskQueue, taskGroup, newTasks)
		taskGroup.Done()
	}
}

func processTaskResult(ep *ExecutionPlan, taskResult *ExecutionTaskResult, result *ExecutionResult) []*ExecutionTask {
	var linkErr error
	var links []Link
	links, linkErr = ep.GetLinks(taskResult.Task.executionNode)
	resultGroup := &ResultGroup{
		Type:           taskResult.Task.executionNode.Type,
		LinkConditions: taskResult.Task.LinkConditions,
		Items:          make([]*ResultItem, len(taskResult.Items)),
		LinkError:      linkErr,
	}
	if taskResult.Task.parentResultItemId == "" {
		result.LinkItemToGroup(taskResult.Task.parentResultItemId, taskResult.Id)
	}
	newTasks := []*ExecutionTask{}
	for i, item := range taskResult.Items {
		resultItem := &ResultItem{
			Id:         fmt.Sprintf("%s-%s", taskResult.Id, item.Data.Identifier),
			Item:       item,
			QueryError: taskResult.Err,
			GroupId:    taskResult.Id,
		}
		for _, link := range links {
			newTaskId := uuid.New().String()
			linkAttrValue, linkAttrFound := (*item.Data.Properties)[link.SourceAttr]
			linkConditions := []reader.Condition{{
				Attr:  link.TargetAttr,
				Op:    common.EqualsTo,
				Value: linkAttrValue,
			}}
			if !linkAttrFound {
				resultItem.LinkErrors = append(
					resultItem.LinkErrors,
					fmt.Errorf("cannot link item: attribute \"%s\" not found", link.SourceAttr),
				)
			} else {
				newTask := &ExecutionTask{
					Id:             newTaskId,
					Type:           link.TargetNode.Type.Type,
					Profiles:       link.TargetNode.Profiles,
					Attrs:          link.TargetNode.Attrs,
					BaseConditions: link.TargetNode.Conditions,
					LinkConditions: linkConditions,
					Reader:         link.TargetNode.Reader,

					executionNode: link.TargetNode,
				}
				newTasks = append(newTasks, newTask)
			}
			resultGroup.Items[i] = resultItem
		}
	}
	return newTasks
}
