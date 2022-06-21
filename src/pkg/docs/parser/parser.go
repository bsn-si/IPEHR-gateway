package parser

import (
	"encoding/json"
	"fmt"

	"hms/gateway/pkg/docs/model"
)

func ParseDocument(inDocument []byte) (doc model.EHR, err error) {
	err = json.Unmarshal(inDocument, &doc)
	return
}

type DataEntry struct {
	Type   string
	NodeId string
	Name   string
	Value  map[string]interface{}
}

type Node struct {
	Type   string
	NodeId string
	Next   map[string]*Node
	Items  *[]DataEntry
}

func NewNode(_type, nodeId string) *Node {
	return &Node{
		Type:   _type,
		NodeId: nodeId,
		Next:   make(map[string]*Node),
	}
}

func (x *Node) Dump(level int) {
	/*
		log.Println(strings.Repeat("\t", level), x._type, x.nodeId)
		if x.items != nil {
			for _, item := range *x.items {
				log.Println(strings.Repeat("\t", level), item)
			}
		}
		for k, v := range x.next {
			//if level == 0 {
			log.Println(strings.Repeat("\t", level), k, ":")
			//}
			v.Dump(level + 1)
		}
	*/
	data, err := json.MarshalIndent(x, "", "    ")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(data))

}

func ParseComposition(inComposition []byte) (composition model.Composition, err error) {
	err = json.Unmarshal(inComposition, &composition)
	if err != nil {
		return
	}

	//log.Println("archetype_node_id", composition.ArchetypeNodeId)

	content := composition.Content

	x := &Node{
		Next: make(map[string]*Node),
	}
	iterate(content, x)

	x.Dump(0)

	return
}

func iterate(items interface{}, node *Node) {
	for _, item := range items.([]interface{}) {
		item := item.(map[string]interface{})

		_type := item["_type"].(string)
		nodeId := item["archetype_node_id"].(string)

		switch _type {
		case "SECTION":
			iterate(item["items"].([]interface{}), node)
		case "EVALUATION", "OBSERVATION", "INSTRUCTION":
			if node.Next[_type] == nil {
				node.Next[_type] = NewNode(_type, nodeId)
			}
			nodeType := node.Next[_type]
			for _, key := range []string{"data", "protocol"} {
				if item[key] == nil {
					continue
				}

				itemsKey := item[key].(map[string]interface{})
				itemsKeyType := itemsKey["_type"].(string)
				itemsKeyNodeId := itemsKey["archetype_node_id"].(string)

				if nodeType.Next[key] == nil {
					nodeType.Next[key] = NewNode(itemsKeyType, itemsKeyNodeId)
				}
				nodeCurrent := nodeType.Next[key]

				if nodeCurrent.Next[itemsKeyNodeId] == nil {
					nodeCurrent.Next[itemsKeyNodeId] = NewNode(itemsKeyType, itemsKeyNodeId)
				}
				nodeCurrent = nodeCurrent.Next[itemsKeyNodeId]

				if itemsKey["items"] != nil {
					if nodeCurrent.Next["items"] == nil {
						nodeCurrent.Next["items"] = NewNode("items", "")
					}
					nodeCurrent = nodeCurrent.Next["items"]
					iterate(itemsKey["items"].([]interface{}), nodeCurrent)
				}

				if itemsKey["events"] != nil {
					if nodeCurrent.Next["events"] == nil {
						nodeCurrent.Next["events"] = NewNode("events", "")
					}
					nodeCurrent = nodeCurrent.Next["events"]
					iterate(itemsKey["events"].([]interface{}), nodeCurrent)
				}
			}

			if item["activities"] != nil {
				if nodeType.Next["activities"] == nil {
					nodeType.Next["activities"] = NewNode("activities", "")
				}
				nodeCurrent := nodeType.Next["activities"]
				iterate(item["activities"].([]interface{}), nodeCurrent)
			}
		case "ACTION":
			if node.Next["ACTION"] == nil {
				node.Next["ACTION"] = NewNode("ACTION", nodeId)
			}
			nodeCurrent := node.Next["ACTION"]

			if nodeCurrent.Next[nodeId] == nil {
				nodeCurrent.Next[nodeId] = NewNode(_type, nodeId)
			}

			for _, key := range []string{"protocol", "description"} {
				if item[key] == nil {
					continue
				}
				itemsKey := item[key].(map[string]interface{})
				itemsKeyType := itemsKey["_type"].(string)
				itemsKeyNodeId := itemsKey["archetype_node_id"].(string)
				if nodeCurrent.Next[key] == nil {
					nodeCurrent.Next[key] = NewNode(itemsKeyType, itemsKeyNodeId)
				}
				nodeCurrent = nodeCurrent.Next[key]

				if nodeCurrent.Next[itemsKeyNodeId] == nil {
					nodeCurrent.Next[itemsKeyNodeId] = NewNode(itemsKeyType, itemsKeyNodeId)
				}
				nodeCurrent = nodeCurrent.Next[itemsKeyNodeId]

				iterate(itemsKey["items"].([]interface{}), nodeCurrent)
			}
		case "CLUSTER":
			if node.Next["CLUSTER"] == nil {
				node.Next["CLUSTER"] = NewNode("CLUSTER", nodeId)
			}
			nodeCluster := node.Next["CLUSTER"]

			if nodeCluster.Next[nodeId] == nil {
				nodeCluster.Next[nodeId] = NewNode(_type, nodeId)
			}
			iterate(item["items"].([]interface{}), nodeCluster.Next[nodeId])
		case "ACTIVITY":
			itemsDescription := item["description"].(map[string]interface{})
			itemsDescriptionType := itemsDescription["_type"].(string)
			itemsDescriptionNodeId := itemsDescription["archetype_node_id"].(string)
			if node.Next["description"] == nil {
				node.Next["description"] = NewNode("description", itemsDescriptionNodeId)
			}
			nodeCurrent := node.Next["description"]

			if nodeCurrent.Next[itemsDescriptionNodeId] == nil {
				nodeCurrent.Next[itemsDescriptionNodeId] = NewNode(itemsDescriptionType, itemsDescriptionNodeId)
			}
			nodeCurrent = nodeCurrent.Next[itemsDescriptionNodeId]
			iterate(itemsDescription["items"].([]interface{}), nodeCurrent)
		case "POINT_EVENT":
			itemsData := item["data"].(map[string]interface{})
			itemsDataType := itemsData["_type"].(string)
			itemsDataNodeId := itemsData["archetype_node_id"].(string)
			if node.Next["data"] == nil {
				node.Next["data"] = NewNode("data", itemsDataNodeId)
			}
			nodeData := node.Next["data"]

			if nodeData.Next[itemsDataNodeId] == nil {
				nodeData.Next[itemsDataNodeId] = NewNode(itemsDataType, itemsDataNodeId)
			}
			iterate(itemsData["items"].([]interface{}), nodeData.Next[itemsDataNodeId])
		case "ITEM_TREE":
			iterate(item["items"].([]interface{}), node)
		case "HISTORY":
			iterate(item["events"].([]interface{}), node)
		case "ELEMENT":
			if node.Items == nil {
				node.Items = &[]DataEntry{}
			}
			itemValue := item["value"].(map[string]interface{})
			itemName := item["name"].(map[string]interface{})
			valueType := itemValue["_type"].(string)
			var value map[string]interface{}
			switch valueType {
			case "DV_TEXT":
				value = map[string]interface{}{
					"value": itemValue["value"],
				}
			case "DV_CODED_TEXT":
				defCode := itemValue["defining_code"].(map[string]interface{})
				value = map[string]interface{}{
					"value":       itemValue["value"],
					"code_string": defCode["code_string"],
				}
			case "DV_IDENTIFIER":
				value = map[string]interface{}{
					"id": itemValue["id"],
				}
			case "DV_MULTIMEDIA":
				value = map[string]interface{}{
					"uri": itemValue["uri"],
				}
			case "DV_DATE_TIME", "DV_DATE", "DV_TIME":
				value = map[string]interface{}{
					"value": itemValue["value"],
				}
			case "DV_QUANTITY":
				value = map[string]interface{}{
					"magnitude": itemValue["magnitude"],
					"units":     itemValue["units"],
					"precision": itemValue["precision"],
				}
			case "DV_COUNT":
				value = map[string]interface{}{
					"magnitude": itemValue["magnitude"],
				}
			case "DV_PROPORTION":
				value = map[string]interface{}{
					"numerator":   itemValue["numerator"],
					"denominator": itemValue["denominator"],
					"type":        itemValue["type"],
				}
			case "DV_URI":
				value = map[string]interface{}{
					"uri": itemValue["uri"],
				}
			case "DV_BOOLEAN":
				value = map[string]interface{}{
					"value": itemValue["value"],
				}
			case "DV_DURATION":
				value = map[string]interface{}{
					"value": itemValue["value"],
				}
			}
			//log.Println(valueType)
			value["_type"] = valueType
			*node.Items = append(*node.Items, DataEntry{
				Type:   _type,
				NodeId: nodeId,
				Name:   itemName["value"].(string),
				Value:  value,
			})
		}
	}
}

/*
func iterate(items interface{}, level int) {
	for _, item := range items.([]interface{}) {
		item := item.(map[string]interface{})
		log.Println(strings.Repeat("\t", level), item["_type"], item["archetype_node_id"])
		switch item["_type"] {
		case "SECTION":
			iterate(item["items"].([]interface{}), level+1)
		case "ACTION":
			if item["protocol"] != nil {
				protocol := item["protocol"].(map[string]interface{})
				log.Println(strings.Repeat("\t", level), "PROTOCOL", protocol["archetype_node_id"])
				if protocol["_type"] == "ITEM_TREE" {
					iterate(protocol["items"].([]interface{}), level+1)
				}
			}
		case "EVALUATION", "OBSERVATION", "POINT_EVENT":
			for _, key := range []string{"data", "protocol"} {
				if item[key] == nil {
					continue
				}
				log.Println(strings.Repeat("\t", level), key, item[key].(map[string]interface{})["archetype_node_id"])
				key := item[key].(map[string]interface{})
				switch key["_type"] {
				case "ITEM_TREE":
					iterate(key["items"].([]interface{}), level+1)
				case "HISTORY":
					iterate(key["events"].([]interface{}), level+1)
				}
			}
		case "INSTRUCTION":
			if item["activities"] != nil {
				iterate(item["activities"].([]interface{}), level+1)
			}
			if item["protocol"] != nil {
				key := item["protocol"].(map[string]interface{})
				if key["_type"] == "ITEM_TREE" {
					log.Println(strings.Repeat("\t", level), key["_type"], key["archetype_node_id"])
					iterate(key["items"].([]interface{}), level+1)
				}
			}
		case "CLUSTER":
			iterate(item["items"].([]interface{}), level+1)
		case "ITEM_TREE":
			iterate(item["items"].([]interface{}), level+1)
		case "ELEMENT":
			name := item["name"].(map[string]interface{})
			value := item["value"].(map[string]interface{})
			switch value["_type"] {
			case "DV_TEXT":
				log.Println(strings.Repeat("\t", level+1), name["value"], ":", value["value"])
			case "DV_CODED_TEXT":
				definingCode := value["defining_code"].(map[string]interface{})
				log.Println(strings.Repeat("\t", level+1), name["value"], ":", value["value"], "defining_code/code_string", definingCode["code_string"])
			case "DV_IDENTIFIER":
				log.Println(strings.Repeat("\t", level+1), name["value"], ":", value["id"])
			case "DV_MULTIMEDIA":
				uri := value["uri"].(map[string]interface{})
				log.Println(strings.Repeat("\t", level+1), name["value"], ":", uri["value"])
			case "DV_QUANTITY":
				log.Println(strings.Repeat("\t", level+1), name["value"], ":", value["magnitude"], value["units"])
			case "DV_DATE_TIME":
				log.Println(strings.Repeat("\t", level+1), name["value"], ":", value["value"])
			}
		}
	}
}
*/

/*
func iterate(items interface{}, level int, indexTree map[string]interface{}) {
	var currentLevel = indexTree

	switch items.(type) {
	case []interface{}:
		for _, item := range items.([]interface{}) {
			//log.Println(strings.Repeat("\t", level), item["_type"], item["archetype_node_id"])
			item := item.(map[string]interface{})
			_type := item["_type"].(string)
			nodeId := item["archetype_node_id"].(string)
			switch _type {
			case "SECTION":
				iterate(item["items"], level+1, indexTree)
			case "EVALUATION", "OBSERVATION", "INSTRUCTION":
				if _, ok := indexTree[_type]; !ok {
					currentLevel = make(map[string]interface{})
					indexTree[_type] = currentLevel
				} else {
					currentLevel = indexTree[_type]
				}

				if _, ok := currentLevel[nodeId]; !ok {
					currentLevel = make(map[string]interface{})
					currentLevel[nodeId]
				}

				if item["data"] != nil {
					next1[nodeId] = next2
					next2 := make(map[string]interface{})
					next1["data"] = next2
					itemData := item["data"].(map[string]interface{})
					nodeId := itemData["archetype_node_id"].(string)
					next3 := make(map[string]interface{})
					next2[nodeId] = next3
					if itemData["_type"].(string) == "ITEM_TREE" {
						iterate(itemData["items"], level+1, next3)
					} else if itemData["_type"].(string) == "HISTORY" {
						iterate(itemData["events"], level+1, next3)

					}
				}

				if item["protocol"] != nil {
					itemProtocol := item["protocol"].(map[string]interface{})
					nodeId := itemProtocol["archetype_node_id"].(string)
					next2 := make(map[string]interface{})
					next1[nodeId] = next2
					if itemProtocol["_type"].(string) == "ITEM_TREE" {
						iterate(itemProtocol["items"], level+1, next2)
					} else if itemProtocol["_type"].(string) == "HISTORY" {
						iterate(itemProtocol["events"], level+1, next2)
					}

				}

				/*
						for _, key := range []string{"data", "protocol"} {
							if item[key] == nil {
								continue
							}
							nextLevel[key] = make(map[string]interface{})

							x := item[key]
							//log.Println(strings.Repeat("\t", level), key, item[key].(map[string]interface{})["archetype_node_id"])
							switch key["_type"] {
							case "ITEM_TREE":
								iterate(key["items"].([]interface{}), level+1, nextLevel)
							case "HISTORY":
								iterate(key["events"].([]interface{}), level+1, nextLevel)
							}
						}

					for _, key := range []string{"activities"} {
						if item[key] == nil {
							continue
						}
						iterate(item[key].([]interface{}), level+1, nextLevel)
					}
			case "POINT_EVENT":
				for _, key := range []string{"data", "protocol"} {
					if item[key] == nil {
						continue
					}
					//log.Println(strings.Repeat("\t", level), key, item[key].(map[string]interface{})["archetype_node_id"])
					key := item[key].(map[string]interface{})
					switch key["_type"] {
					case "ITEM_TREE":
						iterate(key["items"].([]interface{}), level+1, next1)
					case "HISTORY":
						iterate(key["events"].([]interface{}), level+1, next1)
					}
				}
			case "ACTION":
				if item["protocol"] != nil {
					protocol := item["protocol"].(map[string]interface{})
					//log.Println(strings.Repeat("\t", level), "PROTOCOL", protocol["archetype_node_id"])
					if protocol["_type"] == "ITEM_TREE" {
						iterate(protocol["items"].([]interface{}), level+1, next1)
					}
				}
			case "CLUSTER":
				iterate(item["items"].([]interface{}), level+1, next1)
			case "ITEM_TREE":
				iterate(item["items"].([]interface{}), level+1, next1)
			case "ELEMENT":
				name := item["name"].(map[string]interface{})
				value := item["value"].(map[string]interface{})
				_ = name
				_ = value
				switch value["_type"] {
				case "DV_TEXT":
					//log.Println(strings.Repeat("\t", level+1), name["value"], ":", value["value"])
				case "DV_CODED_TEXT":
					definingCode := value["defining_code"].(map[string]interface{})
					_ = definingCode
					//log.Println(strings.Repeat("\t", level+1), name["value"], ":", value["value"], "defining_code/code_string", definingCode["code_string"])
				case "DV_IDENTIFIER":
					//log.Println(strings.Repeat("\t", level+1), name["value"], ":", value["id"])
				case "DV_MULTIMEDIA":
					uri := value["uri"].(map[string]interface{})
					_ = uri
					//log.Println(strings.Repeat("\t", level+1), name["value"], ":", uri["value"])
				case "DV_QUANTITY":
					//log.Println(strings.Repeat("\t", level+1), name["value"], ":", value["magnitude"], value["units"])
				case "DV_DATE_TIME":
					//log.Println(strings.Repeat("\t", level+1), name["value"], ":", value["value"])
				}
			}
		}
	case map[string]interface{}:
	}

}
*/
