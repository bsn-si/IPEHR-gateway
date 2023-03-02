package driver

import (
	"database/sql/driver"

	aqlprocessor "github.com/bsn-si/IPEHR-gateway/src/pkg/aql/processor"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/errors"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/storage/treeindex"
)

type executer struct {
	query  *aqlprocessor.Query
	params map[string]driver.Value

	index *treeindex.EHRIndex
}

func (exec *executer) run() (*Rows, error) {
	// handle FROM block
	dataSources, err := exec.findSources()
	if err != nil {
		return nil, errors.Wrap(err, "cannot find data sources")
	}

	// handle WHERE block
	dataSources, err = exec.filterSources(dataSources)
	if err != nil {
		return nil, errors.Wrap(err, "cannot filter data sources")
	}

	// handle SELECT block
	rows, err := exec.queryData(dataSources)
	if err != nil {
		return nil, errors.Wrap(err, "cannot query rows from data sources")
	}

	rows, err = exec.orderRows(rows)
	if err != nil {
		return nil, errors.Wrap(err, "cannot order rows")
	}

	return exec.limitRows(rows), nil
}

func (exec *executer) filterSources(rows dataRows) (dataRows, error) {
	if exec.query.Where == nil {
		return rows, nil
	}

	return processWhere(exec.query.Where, rows)
}

func (exec *executer) orderRows(rows *Rows) (*Rows, error) {
	// handle ORDER block
	//TODO: add order logic
	return rows, nil
}

func (exec *executer) limitRows(rows *Rows) *Rows {
	if exec.query.Limit == nil {
		return rows
	}

	limit := exec.query.Limit.Limit
	offset := exec.query.Limit.Offset

	if offset >= 0 {
		if offset > len(rows.rows) {
			offset = len(rows.rows)
		}

		rows.rows = rows.rows[offset:]
	}

	if limit < len(rows.rows) {
		rows.rows = rows.rows[:limit]
	}

	return rows
}

func getValueForPath(path *aqlprocessor.ObjectPath, node treeindex.Noder) (any, bool) {
	index := 0
	queue := []treeindex.Noder{node}

	for len(queue) > 0 {
		if index >= len(path.Paths) {
			return nil, false
		}

		path := path.Paths[index]

		node := queue[0]
		queue = queue[1:]

		switch node := node.(type) {
		case *treeindex.ObjectNode, *treeindex.EHRNode, *treeindex.CompositionNode, *treeindex.EventContextNode:
			{
				nextNode := node.TryGetChild(path.Identifier)
				if nextNode == nil {
					continue
				}

				switch nextNode.GetNodeType() {
				case treeindex.ObjectNodeType:
					if path.PathPredicate != nil && path.PathPredicate.Type == aqlprocessor.NodePathPredicate {
						if np := path.PathPredicate.NodePredicate; np.AtCode != nil && nextNode.GetID() == np.AtCode.ToString() {
							index++
						}
					}
				case treeindex.DataValueNodeType:
					index++
				}

				queue = append(queue, nextNode)
			}
		case *treeindex.SliceNode:
			if path.PathPredicate != nil && path.PathPredicate.Type == aqlprocessor.NodePathPredicate {
				np := path.PathPredicate.NodePredicate

				if np.AtCode != nil {
					nextNode := node.TryGetChild(np.AtCode.ToString())
					if nextNode != nil {
						queue = append(queue, nextNode)
						index++
					}
				}
			}
		case *treeindex.DataValueNode:
			if valueNode := node.TryGetChild(path.Identifier); valueNode != nil {
				queue = append(queue, valueNode)
			}
		case *treeindex.ValueNode:
			return node.GetData(), true
		}
	}

	return nil, false
}
