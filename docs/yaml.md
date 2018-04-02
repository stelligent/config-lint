# Parsing an Arbitrary YAML file

Set the Type to YAML, and provide a `resources` section to describe how to discover the resources in a YAML file.
The other parsers have code that expects a specific file format, and it converts those into a collection of
resource objects that will be linted.

For arbitrary YAML files, the resources describes how resources should be loaded. Each resources has these attributes:

* type: identifies the type of resource. Rules with a matching `resource` will be applied to these resources.
* key: JMESPath expression that should return an array of objects from the YAML file.
* id: a JMESPATH expression applied to each resource to extract its unique identifier.

The rules section works the same as in any other linter.

Example:

```
version: 1
description: Rules for generic YAML file
type: YAML
files:
  - "*.config"

# For generic YAML linting, we need a list of resources
# Each entry in the list describes the resource type, how to discover it in the file, and how to get its ID
# The key attribute is a JMESPath expression that should return an array

resources:
  - type: widget
    key: widgets[]
    id: id
  - type: gadget
    key: gadgets[]
    id: name
  - type: contraption
    key: other_stuff.contraptions[]
    id: ids.serial_number
  # include the root document in a single element array with a literal id
  - type: document
    key: '[@]'
    id: '`"Document"`'

rules:

  - id: DOCUMENT_KEYS
    message: Unexpected document key
    severity: FAILURE
    resource: document
    assertions:
      - every:
          key: "keys(@)"
          assertions:
            - key: "@"
              op: in
              value: widgets,gadgets,other_stuff

  - id: WIDGET_NAME
    message: Widget needs a name
    severity: FAILURE
    resource: widget
    assertions:
      - key: name
        op: present

  - id: GADGET_COLOR
    message: Gadget has missing or invalid color
    severity: FAILURE
    resource: gadget
    assertions:
      - key: color
        op: in
        value: red,blue,green

  - id: GADGET_PROPERTIES
    message: Gadget has missing properties
    severity: FAILURE
    resource: gadget
    assertions:
      - key: "@"
        op: has-properties
        value: name,color

  - id: CONTRAPTION_SIZE
    message: Contraption size should be less than 1000
    resource: contraption
    severity: FAILURE
    assertions:
      - key: size
        op: lt
        value: 1000
        value_type: integer

  - id: CONTRAPTION_LOCATIONS
    message: Contraption location must have city
    resource: contraption
    severity: FAILURE
    assertions:
      - every:
          key: locations
          assertions:
            - key: city
              op: present
```

There is an example file [here](example-files/config/generic.config) that this rule file can scan.
