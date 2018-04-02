## Assertion Operations

| Operation                         | Description    |
|-----------------------------------|----------------|
| [and](#and)                       | And            |
| [absent](#absent)                 | Absent         |
| [contains](#contains)             | Contains       |
| [empty](#empty)                   | Empty          |
| [every](#every)                   | Every          |
| [eq](#eq)                         | Equal          |
| [in](#in)                         | In             |
| [has-properties](#has-properties) | Has Properties |
| [ne](#ne)                         | Not equal      |
| [none](#none)                     | None           |
| [not](#not)                       | Not            |
| [not-contains](#not-contains)     | Not Contains   |
| [not-empty](#not-empty)           | Not Empty      |
| [not-in](#not-in)                 | Not In         |
| [or](#or)                         | Or             |
| [present](#present)               | Present        |
| [regex](#regex)                   | Regex          |
| [some](#some)                     | Some           |
| [xor](#xor)                       | Xor            |

## eq

Equal

###Example:

```
...
  - id: VOLUME1
    resource: aws_ebs_volume
    message: EBS Volumes must be encrypted
    severity: FAILURE
    assertions:
      - key: encrypted
        op: eq
        value: true
...
```

## ne

Not Equal

Example:
```
...
  - id: SG1
    resource: aws_security_group
    message: Security group should not allow ingress from 0.0.0.0/0
    severity: FAILURE
    assertions:
      - key: "ingress[].cidr_blocks[] | [0]"
        op: ne
        value: "0.0.0.0/0"
...
```

## in

 In list of values

### Example:

```
...
  - id: R1
    message: Instance type should be t2.micro or m3.medium
    resource: aws_instance
    assertions:
      - key: instance_type
        op: in
        value: t2.micro,m3.medium
    severity: WARNING
...
```

## not-in

Not in list of values

## present

Attribute is present

###Example:

```
...
  - id: R6
    message: Department tag is required
    resource: aws_instance
    assertions:
      - key: "tags[].Department | [0]"
        op: present
    severity: FAILURE
...
```

## absent

Attribute is not present

## empty

Attribute is empty

## not-empty 

Attribute is not empty

## contains

Attribute contains a substring, or array contains an element

## not-contains

Attribute contains a substring, or array contains an element

## regex

Attribute matches a regular expression

## and

Logical and of a list of assertions

### Example:

```
...
  - id: ANDTEST
    resource: aws_instance
    message: Should have both Project and Department tags
    severity: WARNING
    assertions:
      - and:
        - key: "tags[].Department | [0]"
          op: present
        - key: "tags[].Project | [0]"
          op: present
    tags:
      - and-test
...
```

## or

Logical or of a list of assertions

### Example:

```
...
  - id: ORTEST
    resource: aws_instance
    message: Should have instance_type of t2.micro or m3.medium
    severity: WARNING
    assertions:
      - or:
        - key: instance_type
          op: eq
          value: t2.micro
        - key: instance_type
          op: eq
          value: m3.medium
...
```

## xor

Logical xor of a list of assertions. The assertion is true when exactly one test passes

### Example:

```
...
  - id: ORTEST
    resource: lint_rule
    message: Can have value or value_from, but not both
    severity: WARNING
    assertions:
      - xor:
        - key: value
          op: present
        - key: value_from
          op: present
...
```

## not

Logical not of an assertion

Example:

```
...
  - id: NOTTEST
    resource: aws_instance
    message: Should not have instance type of c4.large
    severity: WARNING
    assertions:
      - not:
        - key: instance_type
          op: eq
          value: c4.large
...
```

## has-properties

Checks for the present of every property in a comma separated list. This could also be done using the [and](#and) expression,
but this will often be more convenient.

Example:

```
...
  - id: VALID_ADDRESS
    message: Every address needs city, state and zip
    severity: FAILURE
    resource: address
    assertions:
      - key: address
        op: has-properties
        value: city,state,zip
...
```

## every

Select an array from a resource, and run assertions against each element. All of the sub assertions must pass for the test to pass.
The key is a JMESPath expression that should return an array of objects. The key used in each sub assertion is relative to the selected objects.

This provides a simple looping mechanism that is easier to write and understand than a complex JMESPath expression.

Example:

```
...
  - id: LOCATIONS_NEED_LAT_LONG
    message:  Every location requires a latitude and longitude
    severity: FAILURE
    resource: sample
    assertions:
      - every:
          key: Location
          assertions:
            - key: latitude
              op: present
            - key: longitude
              op: present
...
```

## some

Select an array from a resource, and run assertions against each element. At least one sub assertion must pass for the test to pass.
The key is a JMESPath expression that should return an array of objects. The key used in each sub assertion is relative to the selected objects.

This provides a simple looping mechanism that is easier to write and understand than a complex JMESPath expression.

Example:

```
...
  - id: LOCATION_REQUIRES_LAT_LONG
    message:  At least one location requires a latitude and longitude
    severity: FAILURE
    resource: sample
    assertions:
      - some:
          key: Location
          assertions:
            - key: latitude
              op: present
            - key: longitude
              op: present
...
```

## none

Select an array from a resource, and run assertions against each element. All of the sub assertions must fail for the test to pass.
The key is a JMESPath expression that should return an array of objects. The key used in each sub assertion is relative to the selected objects.

This provides a simple looping mechanism that is easier to write and understand than a complex JMESPath expression.

Example:

```
...
  - id: PORT_22_INGRESS
    message:  No ingress for port 22 should be open to the world
    severity: FAILURE
    resource: sample
    assertions:
      - none:
          key: "ipPermissions[]"
          assertions:
            - key: "fromPort"
              op: eq
              value: 22
              value_type: integer
            - key: "ipRanges[]"
              op: contains
              value: 0.0.0.0/0
...
```
