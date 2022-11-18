package treeindex

import (
	"fmt"

	"hms/gateway/pkg/docs/model"
	"hms/gateway/pkg/docs/model/base"

	"github.com/pkg/errors"
)

type Tree struct {
	root map[string]*Node

	actions       map[string]Container
	evaluations   map[string]Container
	instructions  map[string]Container
	obeservations map[string]Container
}

func NewTree() *Tree {
	return &Tree{
		root: make(map[string]*Node),

		actions:       make(map[string]Container),
		evaluations:   make(map[string]Container),
		instructions:  make(map[string]Container),
		obeservations: make(map[string]Container),
	}
}

type Container map[string][]*Node

type Node struct {
	ID      string
	Type    base.ItemType
	IsArray bool

	Attributes Attributes
	Value      base.Root
}

type Attributes map[string]map[string]*Node

func (a Attributes) add(name string, node *Node) {
	m, ok := a[name]
	if !ok {
		m = map[string]*Node{}
		a[name] = m
	}

	m[node.ID] = node
}

func NewNode(obj base.Root) *Node {
	l := obj.GetLocatable()

	return &Node{
		ID:         l.ArchetypeNodeID,
		Type:       l.Type,
		IsArray:    false,
		Attributes: map[string]map[string]*Node{},
	}
}

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
	return nil
}

func (t *Tree) processEvaluation(evaluation *base.Evaluation) error {
	return nil
}

func (t *Tree) processInstruction(instruction *base.Instruction) error {
	return nil
}

func (t *Tree) processObservation(observation *base.Observation) error {
	container, ok := t.obeservations[observation.ArchetypeNodeID]
	if !ok {
		container = Container{}
	}

	node, err := walk(observation)
	if err != nil {
		return errors.Wrap(err, "cannot get node for observation")
	}

	container[node.ID] = append(container[node.ID], node)
	t.obeservations[observation.ArchetypeNodeID] = container
	return nil
}
