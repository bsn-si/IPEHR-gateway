package treeindex

import (
	"fmt"

	"hms/gateway/pkg/docs/model"
	"hms/gateway/pkg/docs/model/base"

	"github.com/pkg/errors"
)

type Tree struct {
	root map[string]noder

	actions       map[string]Container
	evaluations   map[string]Container
	instructions  map[string]Container
	obeservations map[string]Container
}

func NewTree() *Tree {
	return &Tree{
		root: make(map[string]noder),

		actions:       make(map[string]Container),
		evaluations:   make(map[string]Container),
		instructions:  make(map[string]Container),
		obeservations: make(map[string]Container),
	}
}

type Container map[string][]noder

func (t *Tree) AddComposition(com model.Composition) error {
	return t.processCompositionContent(com.Content)
}

func (t *Tree) processCompositionContent(objects []base.Root) error {
	for _, obj := range objects {
		switch obj := obj.(type) {
		case *base.Section:
			if err := t.processSection(obj); err != nil {
				return errors.Wrap(err, "cannot process section object")
			}
		default:
			return fmt.Errorf("unexpected node type in COMPOSITION.Content handler: %T", obj) //nolint
		}
	}

	return nil
}

func (t *Tree) processSection(section *base.Section) error {
	for _, item := range section.Items {
		switch obj := item.(type) {
		case *base.Action:
			if err := t.processAction(obj); err != nil {
				return errors.Wrap(err, "cannot process ACTION in section")
			}
		case *base.Evaluation:
			if err := t.processEvaluation(obj); err != nil {
				return errors.Wrap(err, "cannot process EVALUATION in section")
			}
		case *base.Instruction:
			if err := t.processInstruction(obj); err != nil {
				return errors.Wrap(err, "cannot process INSTRUCTION in section")
			}
		case *base.Observation:
			if err := t.processObservation(obj); err != nil {
				return errors.Wrap(err, "cannot process OBSERVATION in section")
			}
		default:
			return fmt.Errorf("unexpected node type in SECTION handler: %T", obj) //nolint
		}
	}

	return nil
}

func (t *Tree) processAction(action *base.Action) error {
	container, ok := t.actions[action.ArchetypeNodeID]
	if !ok {
		container = Container{}
	}

	node, err := walk(action)
	if err != nil {
		return errors.Wrap(err, "cannot get node for ACTION")
	}

	container[node.getID()] = append(container[node.getID()], node)
	t.actions[action.ArchetypeNodeID] = container
	return nil
}

func (t *Tree) processEvaluation(evaluation *base.Evaluation) error {
	container, ok := t.evaluations[evaluation.ArchetypeNodeID]
	if !ok {
		container = Container{}
	}

	node, err := walk(evaluation)
	if err != nil {
		return errors.Wrap(err, "cannot get node for EVALUATION")
	}

	container[node.getID()] = append(container[node.getID()], node)
	t.evaluations[evaluation.ArchetypeNodeID] = container

	return nil
}

func (t *Tree) processInstruction(instruction *base.Instruction) error {
	container, ok := t.instructions[instruction.ArchetypeNodeID]
	if !ok {
		container = Container{}
	}

	node, err := walk(instruction)
	if err != nil {
		return errors.Wrap(err, "cannot get node for INSTRUCTION")
	}

	container[node.getID()] = append(container[node.getID()], node)
	t.instructions[instruction.ArchetypeNodeID] = container

	return nil
}

func (t *Tree) processObservation(observation *base.Observation) error {
	if err := addObjectIntoCollection(t.obeservations, observation); err != nil {
		return errors.Wrap(err, "cannot add OBSERVATION object")
	}
	return nil
}

func addObjectIntoCollection(collection map[string]Container, obj base.Root) error {
	container, ok := collection[obj.GetArchetypeNodeID()]
	if !ok {
		container = Container{}
	}

	node, err := walk(obj)
	if err != nil {
		return errors.Wrap(err, "cannot get node for collection")
	}

	container[node.getID()] = append(container[node.getID()], node)
	collection[obj.GetArchetypeNodeID()] = container

	return nil
}
