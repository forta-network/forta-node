package messaging

func schemaReg(subject string) (interface{}, bool) {
	switch subject {
	case SubjectAgentsVersionsLatest, SubjectAgentsActionRun, SubjectAgentsActionStop,
		SubjectAgentsStatusRunning, SubjectAgentsStatusStopped:
		return &AgentPayload{}, true
	}
	return nil, false
}
