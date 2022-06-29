// Package data_search keys index
package data_search

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"

	"hms/gateway/pkg/common"
	"hms/gateway/pkg/crypto/hm"
	"hms/gateway/pkg/docs/model"
	"hms/gateway/pkg/errors"
	"hms/gateway/pkg/indexer"
)

type DataEntry struct {
	GroupId       *uuid.UUID
	ValueSet      map[string]interface{}
	DocStorIdEncr []byte
}

type Element struct {
	ItemType    string
	ElementType string
	NodeId      string
	Name        string
	DataEntries []*DataEntry
}

type Node struct {
	NodeType string
	NodeId   string
	Next     map[string]*Node
	Items    map[string]*Element // nodeId -> Element
}

func newNode(_type, nodeId string) *Node {
	return &Node{
		NodeType: _type,
		NodeId:   nodeId,
		Next:     make(map[string]*Node),
	}
}

type DataSearchIndex struct {
	index indexer.Indexer
}

func New() *DataSearchIndex {
	return &DataSearchIndex{
		index: indexer.Init("data_search"),
	}
}

func (d *DataSearchIndex) Add(key string, value interface{}) error {
	return nil
}

func (n *Node) dump() {
	data, err := json.MarshalIndent(n, "", "    ")
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(string(data))
}

func (i *DataSearchIndex) UpdateIndexWithNewContent(content interface{}, groupAccess *model.GroupAccess, docStorageIdEncrypted []byte) error {
	var (
		key   = (*[32]byte)(groupAccess.Key)
		nonce = groupAccess.Nonce
	)
	node := &Node{}
	err := i.index.GetById("INDEX", node)
	if err != nil {
		if errors.Is(err, errors.IsNotExist) {
			node = newNode("INDEX", "")
		} else {
			return err
		}
	}

	var iterate func(items interface{}, node *Node)

	iterate = func(items interface{}, node *Node) {
		switch items.(type) {
		case []interface{}:
			// ok
		default:
			log.Println("Unexpected type of content item:", items)
			return
		}

		for _, item := range items.([]interface{}) {
			item := item.(map[string]interface{})

			_type := item["_type"].(string)
			itemNodeId := item["archetype_node_id"].(string)

			switch _type {
			case "SECTION":
				iterate(item["items"].([]interface{}), node)
			case "EVALUATION", "OBSERVATION", "INSTRUCTION":
				if node.Next[_type] == nil {
					node.Next[_type] = newNode(_type, itemNodeId)
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
						nodeType.Next[key] = newNode(itemsKeyType, itemsKeyNodeId)
					}
					nodeCurrent := nodeType.Next[key]

					if nodeCurrent.Next[itemsKeyNodeId] == nil {
						nodeCurrent.Next[itemsKeyNodeId] = newNode(itemsKeyType, itemsKeyNodeId)
					}
					nodeCurrent = nodeCurrent.Next[itemsKeyNodeId]

					if itemsKey["items"] != nil {
						if nodeCurrent.Next["items"] == nil {
							nodeCurrent.Next["items"] = newNode("items", "")
						}
						nodeCurrent = nodeCurrent.Next["items"]
						iterate(itemsKey["items"].([]interface{}), nodeCurrent)
					}

					if itemsKey["events"] != nil {
						if nodeCurrent.Next["events"] == nil {
							nodeCurrent.Next["events"] = newNode("events", "")
						}
						nodeCurrent = nodeCurrent.Next["events"]
						iterate(itemsKey["events"].([]interface{}), nodeCurrent)
					}
				}

				if item["activities"] != nil {
					if nodeType.Next["activities"] == nil {
						nodeType.Next["activities"] = newNode("activities", "")
					}
					nodeCurrent := nodeType.Next["activities"]
					iterate(item["activities"].([]interface{}), nodeCurrent)
				}
			case "ACTION":
				if node.Next["ACTION"] == nil {
					node.Next["ACTION"] = newNode("ACTION", itemNodeId)
				}
				nodeCurrent := node.Next["ACTION"]

				if nodeCurrent.Next[itemNodeId] == nil {
					nodeCurrent.Next[itemNodeId] = newNode(_type, itemNodeId)
				}

				for _, key := range []string{"protocol", "description"} {
					if item[key] == nil {
						continue
					}
					itemsKey := item[key].(map[string]interface{})
					itemsKeyType := itemsKey["_type"].(string)
					itemsKeyNodeId := itemsKey["archetype_node_id"].(string)
					if nodeCurrent.Next[key] == nil {
						nodeCurrent.Next[key] = newNode(itemsKeyType, itemsKeyNodeId)
					}
					nodeCurrent = nodeCurrent.Next[key]

					if nodeCurrent.Next[itemsKeyNodeId] == nil {
						nodeCurrent.Next[itemsKeyNodeId] = newNode(itemsKeyType, itemsKeyNodeId)
					}
					nodeCurrent = nodeCurrent.Next[itemsKeyNodeId]

					iterate(itemsKey["items"].([]interface{}), nodeCurrent)
				}
			case "CLUSTER":
				if node.Next["CLUSTER"] == nil {
					node.Next["CLUSTER"] = newNode("CLUSTER", itemNodeId)
				}
				nodeCluster := node.Next["CLUSTER"]

				if nodeCluster.Next[itemNodeId] == nil {
					nodeCluster.Next[itemNodeId] = newNode(_type, itemNodeId)
				}
				iterate(item["items"].([]interface{}), nodeCluster.Next[itemNodeId])
			case "ACTIVITY":
				itemsDescription := item["description"].(map[string]interface{})
				itemsDescriptionType := itemsDescription["_type"].(string)
				itemsDescriptionNodeId := itemsDescription["archetype_node_id"].(string)
				if node.Next["description"] == nil {
					node.Next["description"] = newNode("description", itemsDescriptionNodeId)
				}
				nodeCurrent := node.Next["description"]

				if nodeCurrent.Next[itemsDescriptionNodeId] == nil {
					nodeCurrent.Next[itemsDescriptionNodeId] = newNode(itemsDescriptionType, itemsDescriptionNodeId)
				}
				nodeCurrent = nodeCurrent.Next[itemsDescriptionNodeId]
				iterate(itemsDescription["items"].([]interface{}), nodeCurrent)
			case "POINT_EVENT":
				itemsData := item["data"].(map[string]interface{})
				itemsDataType := itemsData["_type"].(string)
				itemsDataNodeId := itemsData["archetype_node_id"].(string)
				if node.Next["data"] == nil {
					node.Next["data"] = newNode("data", itemsDataNodeId)
				}
				nodeData := node.Next["data"]

				if nodeData.Next[itemsDataNodeId] == nil {
					nodeData.Next[itemsDataNodeId] = newNode(itemsDataType, itemsDataNodeId)
				}
				iterate(itemsData["items"].([]interface{}), nodeData.Next[itemsDataNodeId])
			case "ITEM_TREE":
				iterate(item["items"].([]interface{}), node)
			case "HISTORY":
				iterate(item["events"].([]interface{}), node)
			case "ELEMENT":
				var (
					valueSet  map[string]interface{}
					itemValue = item["value"].(map[string]interface{})
					itemName  = item["name"].(map[string]interface{})
					valueType = itemValue["_type"].(string)
					errors    []error
				)
				switch valueType {
				case "DV_TEXT":
					valueSet = map[string]interface{}{
						"value": hm.EncryptString(itemValue["value"].(string), key, nonce, errors),
					}
				case "DV_CODED_TEXT":
					defCode := itemValue["defining_code"].(map[string]interface{})
					codeString := defCode["code_string"].(string)
					codeString = codeString[2:] // format at0000
					valueSet = map[string]interface{}{
						"value":       hm.EncryptString(itemValue["value"].(string), key, nonce, errors),
						"code_string": hm.EncryptInt(codeString, key, errors),
					}
				case "DV_IDENTIFIER":
					valueSet = map[string]interface{}{
						"id": hm.EncryptString(itemValue["id"].(string), key, nonce, errors),
					}
				case "DV_MULTIMEDIA":
					uri := itemValue["uri"].(map[string]interface{})
					valueSet = map[string]interface{}{
						"uri": hm.EncryptString(uri["value"].(string), key, nonce, errors),
					}
				case "DV_DATE_TIME":
					dateTime, err := time.Parse(common.OPENEHR_TIME_FORMAT, itemValue["value"].(string))
					if err != nil {
						errors = append(errors, err)
						break
					}
					valueSet = map[string]interface{}{
						"value": hm.EncryptInt(dateTime.Unix(), key, errors),
					}
				case "DV_DATE":
					dateTime, err := time.Parse("2006-01-02", itemValue["value"].(string))
					if err != nil {
						errors = append(errors, err)
						break
					}
					valueSet = map[string]interface{}{
						"value": hm.EncryptInt(dateTime.Unix(), key, errors),
					}
				case "DV_TIME":
					dateTime, err := time.Parse("15:04:05.999", itemValue["value"].(string))
					if err != nil {
						errors = append(errors, err)
						break
					}
					valueSet = map[string]interface{}{
						"value": hm.EncryptInt(dateTime.Unix(), key, errors),
					}
				case "DV_QUANTITY":
					valueSet = map[string]interface{}{
						"units": hm.EncryptString(itemValue["units"].(string), key, nonce, errors),
					}

					log.Printf("magnitude type: %T", itemValue["magnitude"])

					switch itemValue["magnitude"].(type) {
					case float64:
						valueSet["magnitude"] = hm.EncryptFloat(itemValue["magnitude"], key, errors)
					default:
						valueSet["magnitude"] = hm.EncryptInt(itemValue["magnitude"], key, errors)
					}
					if fmt.Sprintf("%T", itemValue["precision"]) != "<nil>" {
						valueSet["precision"] = hm.EncryptInt(itemValue["precision"], key, errors)
					}
				case "DV_COUNT":
					valueSet = map[string]interface{}{
						"magnitude": hm.EncryptInt(itemValue["magnitude"], key, errors),
					}
				case "DV_PROPORTION":
					valueSet = map[string]interface{}{
						"numerator":   hm.EncryptFloat(itemValue["numerator"], key, errors),
						"denominator": hm.EncryptFloat(itemValue["denominator"], key, errors),
						"type":        hm.EncryptInt(itemValue["type"], key, errors),
					}
				case "DV_URI":
					switch itemValue["uri"].(type) {
					case string:
						valueSet = map[string]interface{}{
							"uri": hm.EncryptString(itemValue["uri"].(string), key, nonce, errors),
						}
					}
				case "DV_BOOLEAN":
					if fmt.Sprintf("%T", itemValue["value"]) != "bool" {
						errors = append(errors, fmt.Errorf("Incorrect DV_BOOLEAN element %v", itemValue["value"]))
						break
					}
					var value string
					if itemValue["value"].(bool) == true {
						value = "true"
					} else {
						value = "false"
					}

					valueSet = map[string]interface{}{
						"value": hm.EncryptString(value, key, nonce, errors),
					}
				case "DV_DURATION":
					// TODO make comparable duration
					valueSet = map[string]interface{}{
						"value": hm.EncryptString(itemValue["value"].(string), key, nonce, errors),
					}
				}

				if len(errors) > 0 {
					log.Printf("Errors in item %v processing. Errors: %v", item, errors)
					continue
				}

				if node.Items == nil {
					node.Items = make(map[string]*Element)
				}

				element, ok := node.Items[itemNodeId]
				if !ok {
					element = &Element{
						ItemType:    _type,
						ElementType: hex.EncodeToString(hm.EncryptString(valueType, key, nonce, errors)), // TODO make ElementType - []byte
						NodeId:      itemNodeId,
						Name:        hex.EncodeToString(hm.EncryptString(itemName["value"].(string), key, nonce, errors)), // TODO make Name - []byte
						DataEntries: []*DataEntry{},
					}
					node.Items[itemNodeId] = element
				}
				dataEntry := &DataEntry{
					GroupId:       groupAccess.GroupUUID,
					ValueSet:      valueSet,
					DocStorIdEncr: docStorageIdEncrypted,
				}
				element.DataEntries = append(element.DataEntries, dataEntry)
			}
		}
	}

	iterate(content, node)

	node.dump()

	return i.index.Replace("INDEX", node)
}
