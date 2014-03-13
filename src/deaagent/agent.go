package deaagent

import (
	"github.com/cloudfoundry/gosteno"
	"github.com/cloudfoundry/loggregatorlib/emitter"
	"github.com/howeyc/fsnotify"
	"io/ioutil"
	"path"
	"time"
)

type agent struct {
	InstancesJsonFilePath string
	logger                *gosteno.Logger
	knownInstancesChan    chan<- func(map[string]*Task)
}

func NewAgent(instancesJsonFilePath string, logger *gosteno.Logger) *agent {
	knownInstancesChan := atomicCacheOperator()
	return &agent{instancesJsonFilePath, logger, knownInstancesChan}
}

func (agent *agent) Start(emitter emitter.Emitter) {
	go agent.pollInstancesJson(emitter)
}

func (agent *agent) processTasks(currentTasks map[string]Task, emitter emitter.Emitter) func(knownTasks map[string]*Task) {
	return func(knownTasks map[string]*Task) {
		agent.logger.Debug("Reading tasks data after event on instances.json")
		agent.logger.Debugf("Current known tasks are %v", knownTasks)
		for taskIdentifier, _ := range knownTasks {
			_, present := currentTasks[taskIdentifier]
			if present {
				continue
			}
			knownTasks[taskIdentifier].stopListening()
			delete(knownTasks, taskIdentifier)
			agent.logger.Debugf("Removing stale task %v", taskIdentifier)
		}

		for _, task := range currentTasks {
			identifier := task.Identifier()
			_, present := knownTasks[identifier]
			if present {
				continue
			}
			agent.logger.Debugf("Adding new task %s", task.Identifier())
			knownTasks[identifier] = &task

			go func() {
				defer func() {
					agent.knownInstancesChan <- removeFromCache(identifier)
				}()
				knownTasks[identifier].startListening(emitter, agent.logger)
			}()
		}
	}
}

func (agent *agent) pollInstancesJson(emitter emitter.Emitter) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		panic(err)
	}

	for {
		time.Sleep(100 * time.Millisecond)
		err := watcher.Watch(path.Dir(agent.InstancesJsonFilePath))
		if err != nil {
			agent.logger.Warnf("Reading failed, retrying. %s\n", err)
			continue
		}
		break
	}

	agent.logger.Info("Read initial tasks data")
	agent.readInstancesJson(emitter)

	for {
		select {
		case ev := <-watcher.Event:
			agent.logger.Debugf("Got Event: %v\n", ev)
			if ev.IsDelete() {
				agent.knownInstancesChan <- resetCache
			} else {
				agent.readInstancesJson(emitter)
			}
		case err := <-watcher.Error:
			agent.logger.Warnf("Received error from file system notification: %s\n", err)
		}
	}
}

func (agent *agent) readInstancesJson(emitter emitter.Emitter) {
	json, err := ioutil.ReadFile(agent.InstancesJsonFilePath)
	if err != nil {
		agent.logger.Warnf("Reading failed, retrying. %s\n", err)
		return
	}

	currentTasks, err := readTasks(json)
	if err != nil {
		agent.logger.Warnf("Failed parsing json %s: %v Trying again...\n", err, string(json))
		return
	}

	agent.knownInstancesChan <- agent.processTasks(currentTasks, emitter)
}

func removeFromCache(taskId string) func(knownTasks map[string]*Task) {
	return func(knownTasks map[string]*Task) {
		delete(knownTasks, taskId)
	}
}

func resetCache(knownTasks map[string]*Task) {
	for _, task := range knownTasks {
		task.stopListening()
	}
	knownTasks = make(map[string]*Task)
}

func atomicCacheOperator() chan<- func(map[string]*Task) {
	operations := make(chan func(map[string]*Task))
	go func() {
		knownTasks := make(map[string]*Task)
		for operation := range operations {
			operation(knownTasks)
		}
	}()
	return operations
}
