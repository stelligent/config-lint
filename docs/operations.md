## Assertion Operations

| Operation                         | Description    |
|-----------------------------------|----------------|
| [absent](#absent)                 | Absent         |
| [and](#and)                       | And            |
| [contains](#contains)             | Contains       |
| [does-not-contain](#does-not-contain)     | Does Not Contain   |
| [end-with](#end-with)             | Ends With      |
| [eq](#eq)                         | Equal          |
| [empty](#empty)                   | Empty          |
| [every](#every)                   | Every          |
| [has-properties](#has-properties) | Has Properties |
| [in](#in)                         | In             |
| [is-array](#is-array)             | Is Array       |
| [is-false](#is-false)             | Is False       |
| [is-not-array](#is-not-array)     | Is Not Array   |
| [is-true](#is-true)               | Is True        |
| [ne](#ne)                         | Not equal      |
| [none](#none)                     | None           |
| [not](#not)                       | Not            |
| [not-contains](#does-not-contain) | Does Not Contain |
| [not-empty](#not-empty)           | Not Empty      |
| [not-in](#not-in)                 | Not In         |
| [exactly-one](#exactly-one)       | Exactly One    |
| [or](#or)                         | Or             |
| [present](#present)               | Present        |
| [regex](#regex)                   | Regex          |
| [starts-with](#starts-with)       | Starts With    |
| [some](#some)                     | Some           |
| [xor](#xor)                       | Xor            |
| [is-subnet](#is-subnet)           | Is Subnet      |
| [is-private-ip](#is-private-ip)           | Is Private IP      |
| [exposed-hosts](#exposed-hosts)           | Number of hosts exposed to |

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

## does-not-contain

Attribute contains a substring, or array contains an element

You can also use the 'not-contains' operator to do the same thing.

## regex

Attribute matches a regular expression

See [here](https://github.com/google/re2/wiki/Syntax) for regular expression syntax.

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
          expressions:
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
          expressions:
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
          expressions:
            - key: "fromPort"
              op: eq
              value: 22
              value_type: integer
            - key: "ipRanges[]"
              op: contains
              value: 0.0.0.0/0
...
```

## exactly-one

Select an array from a resource, and run assertions against each element. Only one of the sub assertions should return true for the test to pass.
The key is a JMESPath expression that should return an array of objects. The key used in each sub assertion is relative to the selected objects.

This provides a simple looping mechanism that is easier to write and understand than a complex JMESPath expression.

Example:

```
  - id: ONLY_ONE_DEFAULT
    message:  Default should be true for only one element
    severity: FAILURE
    resource: sample
    assertions:
      - exactly-one:
          key: "items[]"
          expressions:
            - key: "default"
              op: is-true
```

## is-true

Check that the data has a true value. Shorthand for using { op: "eq" , "value": true }

Example:

```
...
  - id: ENCRYPTION_TRUE
    message: Encryption should be true
    severity: FAILURE
    resource: ebs_volume
    assertions:
       - key: encrypted
         op: is-true
...
```

## is-false

Check that the data has a false value. Shorthand for using { op: "eq" , "value": false }

Example:

```
...
  - id: ALLOW_PUBLIC_ACCESS
    message: Should not allow public access
    severity: FAILURE
    resource: some_resource
    assertions:
       - key: public_access
         op: is-false
...
```

## starts-with

Check that a string value starts with a value

Example:

```
...
  - id: NAME_PREFIX
    message: Name should have a certain prefix
    severity: FAILURE
    resource: some_resource
    assertions:
       - key: name
         op: starts-with
         value: Foo
...
```

## ends-with

Check that a string value ends with a value

Example:

```
...
  - id: NAME_SUFFIX
    message: Name should have a certain suffix
    severity: FAILURE
    resource: some_resource
    assertions:
       - key: name
         op: ends-with
         value: Resource
...
```

## is-array

Check that an attribute is an array

Example:

```
...
  - id: TAGS_ARRAY
    message: Tags should be an array
    severity: FAILURE
    resource: some_resource
    assertions:
       - key: Tags[]
         op: is-array
...
```

## is-not-array

Check that an attribute is not an array

Example:

```
...
  - id: DESCRIPTION_NOT_AN_ARRAY
    message: Description should not be an array
    severity: FAILURE
    resource: some_resource
    assertions:
       - key: Name
         op: is-not-array
...
```

## is-subnet

Check whether a given IP or CIDR block is a subnet of a larger CIDR  block

Example:

```
...
  - id: IP_IN_SUPERNET
    message: All ip_address values should be in the 10.0.0.0/8 supernet
    severity: FAILURE
    resource: some_resource
    assertions:
       - key: "ip_address"
         op: is-subnet
         value: "10.0.0.0/8"
...
```

## is-private-ip

Check whether a given IP is in RFC1918 address space.

Example:

```
...
  - id: IP_IN_PRIVATE_RANGE
    message: All ip_address values should be in private address space
    severity: FAILURE
    resource: some_resource
    assertions:
       - key: "ip_address"
         op: is-private-ip
...
```

## max-host-count

Checks how many hosts in a given CIDR range. This is useful for evaluating security group rules, for instance.

Example:

```
...
  - id: MAX_HOSTS_EXPOSED_PER_RULE
    message: All security group rules must expose less than 1016 hosts
    severity: FAILURE
    resource: aws_security_group_rule
    assertions:
      - every:
        key: "cidr_blocks"
        expressions:
          - key: "@"
            op: max-host-count
            value: 1016
...
```