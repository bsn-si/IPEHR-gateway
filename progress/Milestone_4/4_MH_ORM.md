## Integration with MH-ORM

At this stage, the information models of openEHR entities were partially implemented according to the specifications: [https://specifications.openehr.org/releases/RM/Release-1.1.0](https://specifications.openehr.org/releases/RM/Release-1.1.0):

- EHR Information Model: EHR, COMPOSITON, EVENT\_CONTEXT, CONTENT\_ITEM, SECTION, ENTRY, CARE\_ENTRY, OBSERVATION, EVALUATION, INSTRUCTION, ACTIVITY, ACTION, INSTRUCTION_DETAILS, ISM\_TRANSITION

- Common Information Model: PATHABLE, LOCATABLE, ARCHETYPED, LINK, FEEDER\_AUDIT, FEEDER\_AUDIT\_DETAILS, PARTY\_PROXY, PARTY\_SELF, PARTY\_IDENTIFIED, PARTY\_RELATED, PARTICIPATION, AUDIT\_DETAILS, CONTRIBUTION

- Data Structures Information Model: ITEM\_STRUCTURE, ITEM_SINGLE, ITEM\_LIST, ITEM\_TABLE, ITEM\_TREE, ITEM, CLUSTER, ELEMENT, HISTORY, EVENT, POINT\_EVENT, INTERVAL\_EVENT

- Data Types Information Model: DATA\_VALUE, DV\_BOOLEAN, DV\_STATE, DV\_IDENTIFIER, DV\_TEXT, TERM\_MAPPING, CODE\_PHRASE, DV\_CODED\_TEXT, DV\_PARAGRAPH, DV\_ORDERED, DV\_INTERVAL, REFERENCE\_RANGE, DV\_QUANTIFIED, DV\_AMOUNT, DV\_QUANTITY, DV\_COUNT, DV\_PROPORTION, DV\_TEMPORAL, DV\_DATE, DV\_TIME, DV\_DATE\_TIME, DV\_DURATION

Basic validation and connectivity checks are implemented when saving openEHR documents.
