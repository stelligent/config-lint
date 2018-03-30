## Assertion Operations

| Operation               | Description |
|-------------------------|-------------|
| [eq](#eq)               | Equal       |
| [ne](#ne)               | Not equal   |
| [in](#in)               | In          |
| [not-in](#not-in)       | Not In      | 
| [present](#present)     | Present     |
| [absent](#absent)       | Absent      |
| [empty](#empty)         | Empty       |
| [not-empty](#not-empty) | Not Empty   |
| [contains](#contains)   | Contains    |
| [regex](#regex)         | Regex       |
| [and](#and)             | And         |
| [or](#or)               | Or          |
| [not](#not)             | Not         |

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
