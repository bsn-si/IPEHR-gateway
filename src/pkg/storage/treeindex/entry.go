package treeindex

import (
	"github.com/bsn-si/IPEHR-gateway/src/pkg/docs/model/base"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/errors"
)

func processCareEntry(node Noder, entry *base.CareEntry) (Noder, error) {
	// todo: add processing for CareEntry struct fields
	return node, nil
}

func processAction(node Noder, act *base.Action) (Noder, error) {
	node, err := processCareEntry(node, &act.CareEntry)
	if err != nil {
		return nil, errors.Wrap(err, "cannot process ACTION.base")
	}

	descriptionNode, err := walk(act.Description)
	if err != nil {
		return nil, errors.Wrap(err, "cannot process ACTION.description")
	}

	node.addAttribute("description", descriptionNode)

	//todo: add processing for ACTION struct fields

	return node, nil
}

func processEvaluation(node Noder, evaluation *base.Evaluation) (Noder, error) {
	node, err := processCareEntry(node, &evaluation.CareEntry)
	if err != nil {
		return nil, errors.Wrap(err, "cannot process EVALUATION.base")
	}

	dataNode, err := walk(evaluation.Data)
	if err != nil {
		return nil, errors.Wrap(err, "cannon process EVALUATION.Data")
	}

	node.addAttribute("data", dataNode)

	//todo: add processing for ACTION struct fields

	return node, nil
}

func processInstruction(node Noder, instr *base.Instruction) (Noder, error) {
	node, err := processCareEntry(node, &instr.CareEntry)
	if err != nil {
		return nil, errors.Wrap(err, "cannot process INSTRUCTION.base")
	}

	//todo: add processing for INSTRUCTION struct fields

	return node, nil
}

func processObservation(node Noder, obs *base.Observation) (Noder, error) {
	node, err := processCareEntry(node, &obs.CareEntry)
	if err != nil {
		return nil, errors.Wrap(err, "cannot process OBSERVATION.base")
	}

	dataNode, err := walk(obs.Data)
	if err != nil {
		return nil, errors.Wrap(err, "cannon process OBSERVATION.Data")
	}

	node.addAttribute("data", dataNode)

	if obs.State != nil {
		stateNode, err := walk(*obs.State)
		if err != nil {
			return nil, errors.Wrap(err, "cannot process OBSERVATION.state")
		}

		node.addAttribute("state", stateNode)
	}

	if obs.Protocol != nil && obs.Protocol.Data != nil {
		protocolNode, err := walk(*obs.Protocol)
		if err != nil {
			return nil, errors.Wrap(err, "cannon process OBSERVATION.Protocol")
		}

		node.addAttribute("protocol", protocolNode)
	}

	return node, nil
}
