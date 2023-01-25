package treeindex

import (
	"fmt"

	"github.com/bsn-si/IPEHR-gateway/src/pkg/docs/model"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/docs/model/base"
	"github.com/vmihailenco/msgpack/v5"

	"github.com/pkg/errors"
)

const (
	ACTION      = "ACTION"
	EVALUATION  = "EVALUATION"
	INSTRUCTION = "INSTRUCTION"
	OBSERVATION = "OBSERVATION"
)

type Tree struct {
	Data map[string]Container
}

func NewTree() *Tree {
	return &Tree{
		Data: map[string]Container{
			ACTION:      {},
			EVALUATION:  {},
			INSTRUCTION: {},
			OBSERVATION: {},
		},
	}
}

func (t *Tree) GetDataSourceByName(name string) (Container, error) {
	c, ok := t.Data[name]
	if !ok {
		return nil, fmt.Errorf("unexpected source type: %v", name) //nolint
	}

	return c, nil
}

func (t *Tree) AddComposition(com model.Composition) error {
	return t.processCompositionContent(com.Content)
}

type Container map[string][]Noder

func (c *Container) DecodeMsgpack(dec *msgpack.Decoder) error {
	tempMap := map[string][]NodeWrapper{}
	if err := dec.Decode(&tempMap); err != nil {
		return errors.Wrap(err, "cannot unmarshal attributes map")
	}

	container := Container{}
	for k, nodes := range tempMap {
		container[k] = make([]Noder, 0, len(nodes))

		for _, n := range nodes {
			container[k] = append(container[k], n.data)
		}
	}

	*c = container
	return nil
}

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
			if err := addObjectIntoCollection(t.Data[ACTION], obj); err != nil {
				return errors.Wrap(err, "cannot process ACTION in section")
			}
		case *base.Evaluation:
			if err := addObjectIntoCollection(t.Data[EVALUATION], obj); err != nil {
				return errors.Wrap(err, "cannot process EVALUATION in section")
			}
		case *base.Instruction:
			if err := addObjectIntoCollection(t.Data[INSTRUCTION], obj); err != nil {
				return errors.Wrap(err, "cannot process INSTRUCTION in section")
			}
		case *base.Observation:
			if err := addObjectIntoCollection(t.Data[OBSERVATION], obj); err != nil {
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

	container[node.GetID()] = append(container[node.GetID()], node)

	return nil
}
