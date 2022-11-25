package treeindex

import (
	"fmt"

	"hms/gateway/pkg/docs/model"
	"hms/gateway/pkg/docs/model/base"

	"github.com/pkg/errors"
)

type Tree struct {
	actions       Container
	evaluations   Container
	instructions  Container
	obeservations Container
}

func NewTree() *Tree {
	return &Tree{
		actions:       Container{},
		evaluations:   Container{},
		instructions:  Container{},
		obeservations: Container{},
	}
}
func (t *Tree) AddComposition(com model.Composition) error {
	return t.processCompositionContent(com.Content)
}

type Container map[string][]noder

func (c Container) Len() int {
	count := 0
	for _, v := range c {
		count += len(v)
	}

	return count
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
			if err := addObjectIntoCollection(t.actions, obj); err != nil {
				return errors.Wrap(err, "cannot process ACTION in section")
			}
		case *base.Evaluation:
			if err := addObjectIntoCollection(t.evaluations, obj); err != nil {
				return errors.Wrap(err, "cannot process EVALUATION in section")
			}
		case *base.Instruction:
			if err := addObjectIntoCollection(t.instructions, obj); err != nil {
				return errors.Wrap(err, "cannot process INSTRUCTION in section")
			}
		case *base.Observation:
			if err := addObjectIntoCollection(t.obeservations, obj); err != nil {
				return errors.Wrap(err, "cannot process OBSERVATION in section")
			}
		default:
			return fmt.Errorf("unexpected node type in SECTION handler: %T", obj) //nolint
		}
	}

	return nil
}

func addObjectIntoCollection(container Container, obj base.Root) error {
	node, err := walk(obj)
	if err != nil {
		return errors.Wrap(err, "cannot get node for collection")
	}

	container[node.getID()] = append(container[node.getID()], node)

	return nil
}
