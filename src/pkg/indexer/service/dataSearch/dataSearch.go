// Package dataSearch keys index
package dataSearch

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"

	"hms/gateway/pkg/common"
	"hms/gateway/pkg/crypto/hm"
	"hms/gateway/pkg/docs/model"
	"hms/gateway/pkg/errors"
	"hms/gateway/pkg/indexer"
	log "hms/gateway/pkg/logging"
)

type DataEntry struct {
	GroupID       *uuid.UUID
	ValueSet      map[string]interface{}
	DocStorIDEncr []byte
}

type Element struct {
	ItemType    string
	ElementType string
	NodeID      string
	Name        string
	DataEntries []*DataEntry
}

type Node struct {
	NodeType string
	NodeID   string
	Next     map[string]*Node
	Items    map[string]*Element // nodeID -> Element
}

func newNode(_type, nodeID string) *Node {
	return &Node{
		NodeType: _type,
		NodeID:   nodeID,
		Next:     make(map[string]*Node),
	}
}

type Index struct {
	index indexer.Indexer
}

func New() *Index {
	return &Index{
		index: indexer.Init("data_search"),
	}
}

func (i *Index) Add(key string, value interface{}) error {
	return nil
}

// nolint
func (n *Node) dump() {
	data, err := json.MarshalIndent(n, "", "    ")
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(string(data))
}

func (i *Index) UpdateIndexWithNewContent(content interface{}, groupAccess *model.GroupAccess, docStorageIDEncrypted []byte) error {
	var (
		key   = (*[32]byte)(groupAccess.Key)
		nonce = groupAccess.Nonce
		node  = &Node{}
	)

	if err := i.index.GetByID("INDEX", node); err != nil {
		if errors.Is(err, errors.ErrIsNotExist) {
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
			itemNodeID := item["archetype_node_id"].(string)

			switch _type {
			case "SECTION":
				iterate(item["items"].([]interface{}), node)
			case "EVALUATION", "OBSERVATION", "INSTRUCTION":
				if node.Next[_type] == nil {
					node.Next[_type] = newNode(_type, itemNodeID)
				}

				nodeType := node.Next[_type]

				for _, key := range []string{"data", "protocol"} {
					if item[key] == nil {
						continue
					}

					itemsKey := item[key].(map[string]interface{})
					itemsKeyType := itemsKey["_type"].(string)
					itemsKeyNodeID := itemsKey["archetype_node_id"].(string)

					if nodeType.Next[key] == nil {
						nodeType.Next[key] = newNode(itemsKeyType, itemsKeyNodeID)
					}

					nodeCurrent := nodeType.Next[key]
					if nodeCurrent.Next[itemsKeyNodeID] == nil {
						nodeCurrent.Next[itemsKeyNodeID] = newNode(itemsKeyType, itemsKeyNodeID)
					}

					nodeCurrent = nodeCurrent.Next[itemsKeyNodeID]

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
					node.Next["ACTION"] = newNode("ACTION", itemNodeID)
				}

				nodeCurrent := node.Next["ACTION"]

				if nodeCurrent.Next[itemNodeID] == nil {
					nodeCurrent.Next[itemNodeID] = newNode(_type, itemNodeID)
				}

				for _, key := range []string{"protocol", "description"} {
					if item[key] == nil {
						continue
					}

					itemsKey := item[key].(map[string]interface{})

					itemsKeyType := itemsKey["_type"].(string)

					itemsKeyNodeID := itemsKey["archetype_node_id"].(string)

					if nodeCurrent.Next[key] == nil {
						nodeCurrent.Next[key] = newNode(itemsKeyType, itemsKeyNodeID)
					}

					nodeCurrent = nodeCurrent.Next[key]
					if nodeCurrent.Next[itemsKeyNodeID] == nil {
						nodeCurrent.Next[itemsKeyNodeID] = newNode(itemsKeyType, itemsKeyNodeID)
					}

					nodeCurrent = nodeCurrent.Next[itemsKeyNodeID]

					iterate(itemsKey["items"].([]interface{}), nodeCurrent)
				}
			case "CLUSTER":
				if node.Next["CLUSTER"] == nil {
					node.Next["CLUSTER"] = newNode("CLUSTER", itemNodeID)
				}

				nodeCluster := node.Next["CLUSTER"]
				if nodeCluster.Next[itemNodeID] == nil {
					nodeCluster.Next[itemNodeID] = newNode(_type, itemNodeID)
				}

				iterate(item["items"].([]interface{}), nodeCluster.Next[itemNodeID])
			case "ACTIVITY":
				itemsDescription := item["description"].(map[string]interface{})

				itemsDescriptionType := itemsDescription["_type"].(string)

				itemsDescriptionNodeID := itemsDescription["archetype_node_id"].(string)

				if node.Next["description"] == nil {
					node.Next["description"] = newNode("description", itemsDescriptionNodeID)
				}

				nodeCurrent := node.Next["description"]
				if nodeCurrent.Next[itemsDescriptionNodeID] == nil {
					nodeCurrent.Next[itemsDescriptionNodeID] = newNode(itemsDescriptionType, itemsDescriptionNodeID)
				}

				nodeCurrent = nodeCurrent.Next[itemsDescriptionNodeID]

				iterate(itemsDescription["items"].([]interface{}), nodeCurrent)
			case "POINT_EVENT":
				itemsData := item["data"].(map[string]interface{})

				itemsDataType := itemsData["_type"].(string)

				itemsDataNodeID := itemsData["archetype_node_id"].(string)

				if node.Next["data"] == nil {
					node.Next["data"] = newNode("data", itemsDataNodeID)
				}

				nodeData := node.Next["data"]
				if nodeData.Next[itemsDataNodeID] == nil {
					nodeData.Next[itemsDataNodeID] = newNode(itemsDataType, itemsDataNodeID)
				}

				iterate(itemsData["items"].([]interface{}), nodeData.Next[itemsDataNodeID])
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
					err       error
				)

				switch valueType {
				case "DV_TEXT":
					switch value := itemValue["value"].(type) {
					case string:
						valueSet = map[string]interface{}{
							"value": hm.EncryptString(value, key, nonce),
						}
					default:
						err = fmt.Errorf("%w: DV_TEXT value element %v", errors.ErrIncorrectFormat, value)
					}
				case "DV_CODED_TEXT":
					switch defCode := itemValue["defining_code"].(type) {
					case map[string]interface{}:
						switch codeString := defCode["code_string"].(type) {
						case string:
							codeString = codeString[2:] // format at0000

							var codeStringInt int64

							codeStringInt, err = strconv.ParseInt(codeString, 10, 64)
							if err != nil {
								err = fmt.Errorf("%w: DV_CODED_TEXT defining_code.code_string element %v", errors.ErrIncorrectFormat, codeString)
								break
							}

							valueSet = map[string]interface{}{
								"code_string": hm.EncryptInt(codeStringInt, key),
							}

							switch value := itemValue["value"].(type) {
							case string:
								valueSet["value"] = hm.EncryptString(value, key, nonce)
							default:
								err = fmt.Errorf("%w: DV_CODED_TEXT value element %v", errors.ErrIncorrectFormat, value)
							}
						default:
							err = fmt.Errorf("%w: DV_CODED_TEXT code_string element %v", errors.ErrIncorrectFormat, codeString)
						}
					default:
						err = fmt.Errorf("%w: DV_CODED_TEXT defining_code element %v", errors.ErrIncorrectFormat, defCode)
					}
				case "DV_IDENTIFIER":
					switch id := itemValue["id"].(type) {
					case string:
						valueSet = map[string]interface{}{
							"id": hm.EncryptString(id, key, nonce),
						}
					default:
						err = fmt.Errorf("%w: DV_IDENTIFIER id element %v", errors.ErrIncorrectFormat, id)
					}
				case "DV_MULTIMEDIA":
					switch uri := itemValue["uri"].(type) {
					case map[string]interface{}:
						switch value := uri["value"].(type) {
						case string:
							valueSet = map[string]interface{}{
								"uri": hm.EncryptString(value, key, nonce),
							}
						default:
							err = fmt.Errorf("%w: DV_MULTIMEDIA uri.value element %v", errors.ErrIncorrectFormat, value)
						}
					default:
						err = fmt.Errorf("%w: DV_MULTIMEDIA uri element %v", errors.ErrIncorrectFormat, uri)
					}
				case "DV_DATE_TIME":
					switch value := itemValue["value"].(type) {
					case string:
						var dateTime time.Time

						if dateTime, err = time.Parse(common.OpenEhrTimeFormat, value); err != nil {
							err = fmt.Errorf("%w: DV_DATE_TIME.value element %v", errors.ErrIncorrectFormat, value)
							break
						}

						valueSet = map[string]interface{}{
							"value": hm.EncryptInt(dateTime.Unix(), key),
						}
					default:
						err = fmt.Errorf("%w: DV_DATE_TIME.value element %v", errors.ErrIncorrectFormat, value)
					}
				case "DV_DATE":
					switch value := itemValue["value"].(type) {
					case string:
						var dateTime time.Time

						if dateTime, err = time.Parse("2006-01-02", value); err != nil {
							err = fmt.Errorf("%w: DV_DATE.value element %v", errors.ErrIncorrectFormat, value)
							break
						}

						valueSet = map[string]interface{}{
							"value": hm.EncryptInt(dateTime.Unix(), key),
						}
					default:
						err = fmt.Errorf("%w: DV_DATE.value element %v", errors.ErrIncorrectFormat, value)
					}
				case "DV_TIME":
					switch value := itemValue["value"].(type) {
					case string:
						var dateTime time.Time

						if dateTime, err = time.Parse("15:04:05.999", value); err != nil {
							err = fmt.Errorf("%w: DV_TIME.value element %v", errors.ErrIncorrectFormat, value)
							break
						}

						valueSet = map[string]interface{}{
							"value": hm.EncryptInt(dateTime.Unix(), key),
						}
					default:
						err = fmt.Errorf("%w: DV_TIME.value element %v", errors.ErrIncorrectFormat, value)
					}
				case "DV_QUANTITY":
					switch units := itemValue["units"].(type) {
					case string:
						valueSet = map[string]interface{}{
							"units": hm.EncryptString(itemValue["units"].(string), key, nonce),
						}
					default:
						err = fmt.Errorf("%w: DV_QUANTITY.units element %v", errors.ErrIncorrectFormat, units)
					}

					if err != nil {
						break
					}

					switch magnitude := itemValue["magnitude"].(type) {
					case float64:
						valueSet["magnitude"] = hm.EncryptFloat(magnitude, key)
					case int64:
						valueSet["magnitude"] = hm.EncryptInt(magnitude, key)
					default:
						err = fmt.Errorf("%w: DV_QUANTITY.magnitude element %v", errors.ErrIncorrectFormat, magnitude)
					}

					if err != nil {
						break
					}

					switch precision := itemValue["precision"].(type) {
					case float64:
						valueSet["precision"] = hm.EncryptFloat(precision, key)
					case int64:
						valueSet["precision"] = hm.EncryptInt(precision, key)
					}
				case "DV_COUNT":
					switch magnitude := itemValue["magnitude"].(type) {
					case float64:
						valueSet = map[string]interface{}{
							"magnitude": hm.EncryptFloat(magnitude, key),
						}
					case int64:
						valueSet = map[string]interface{}{
							"magnitude": hm.EncryptInt(magnitude, key),
						}
					default:
						err = fmt.Errorf("%w: DV_COUNT.magnitude element %v", errors.ErrIncorrectFormat, magnitude)
					}
				case "DV_PROPORTION":
					switch numerator := itemValue["numerator"].(type) {
					case float64:
						valueSet = map[string]interface{}{
							"numerator": hm.EncryptFloat(numerator, key),
						}
					default:
						err = fmt.Errorf("%w: DV_PROPORTION.numerator element %v", errors.ErrIncorrectFormat, numerator)
					}

					if err != nil {
						break
					}

					switch denominator := itemValue["denominator"].(type) {
					case float64:
						valueSet = map[string]interface{}{
							"denominator": hm.EncryptFloat(denominator, key),
						}
					default:
						err = fmt.Errorf("%w: DV_PROPORTION.denominator element %v", errors.ErrIncorrectFormat, denominator)
					}

					if err != nil {
						break
					}

					switch _type := itemValue["type"].(type) {
					case float64:
						valueSet = map[string]interface{}{
							"type": hm.EncryptFloat(_type, key),
						}
					case int64:
						valueSet = map[string]interface{}{
							"type": hm.EncryptInt(_type, key),
						}
					default:
						err = fmt.Errorf("%w: DV_PROPORTION.type element %v", errors.ErrIncorrectFormat, _type)
					}
				case "DV_URI":
					switch value := itemValue["value"].(type) {
					case string:
						valueSet = map[string]interface{}{
							"value": hm.EncryptString(value, key, nonce),
						}
					default:
						err = fmt.Errorf("%w: DV_URI.value element %v", errors.ErrIncorrectFormat, value)
					}
				case "DV_BOOLEAN":
					switch value := itemValue["value"].(type) {
					case bool:
						valueSet = map[string]interface{}{
							"value": hm.EncryptString(strconv.FormatBool(value), key, nonce),
						}
					default:
						err = fmt.Errorf("%w: DV_BOOLEAN.value element %v", errors.ErrIncorrectFormat, value)
					}
				case "DV_DURATION":
					// TODO make comparable duration
					switch value := itemValue["value"].(type) {
					case string:
						valueSet = map[string]interface{}{
							"value": hm.EncryptString(value, key, nonce),
						}
					default:
						err = fmt.Errorf("%w: DV_DURATION.value element %v", errors.ErrIncorrectFormat, value)
					}
				}

				if err != nil {
					log.Printf("Errors in item %v processing. Error: %v", item, err)
					continue
				}

				if node.Items == nil {
					node.Items = make(map[string]*Element)
				}

				element, ok := node.Items[itemNodeID]
				if !ok {
					element = &Element{
						ItemType:    _type,
						ElementType: hex.EncodeToString(hm.EncryptString(valueType, key, nonce)), // TODO make ElementType - []byte
						NodeID:      itemNodeID,
						Name:        hex.EncodeToString(hm.EncryptString(itemName["value"].(string), key, nonce)), // TODO make Name - []byte
						DataEntries: []*DataEntry{},
					}
					node.Items[itemNodeID] = element
				}

				dataEntry := &DataEntry{
					GroupID:       groupAccess.GroupUUID,
					ValueSet:      valueSet,
					DocStorIDEncr: docStorageIDEncrypted,
				}

				element.DataEntries = append(element.DataEntries, dataEntry)
			}
		}
	}

	iterate(content, node)

	node.dump()

	return i.index.Replace("INDEX", node)
}
