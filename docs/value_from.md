# Dynamic Value

In cases where a rule needs a dynamic value, instead of specifying "value", "value_from" can be used instead. This allows an external source, such as an S3 object, or an HTTP endpoint to provider the values. Or the value can be provided on the command line when config-lint is invoked.


## Using an S3 Bucket


### Example:

```
...
  - id: VALUE_FROM_S3
    message: Instance type should be in list from S3 object
    resource: aws_instance
    assertions:
      - key: instance_type
        op: in
        value_from:
          url: s3://my-bucket/allowed-instance-types
...
```

## Using an HTTP endpoint

### Example:

```
...
  - id: VALUE_FROM_HTTPS
    message: Instance type should be in list from https endpoint
    resource: aws_instance
    assertions:
      - key: instance_type
        op: in
        value_from:
          url: https://my-api-endpoint/dev/instance_types
...
```

## Using a command line variable

### Example:

```
...
  - id: VALUE_FROM_COMMAND_LINE
    message: Instance type should be in list from https endpoint
    resource: aws_instance
    assertions:
      - key: instance_type
        op: in
        value_from:
          variable: instance_types
...
```

When invoking the config-lint, include the -var option, as in this example:

```
config-lint -rules <RULES_FILE> -var "instance_types=t2.small,c3.medium" <CONFIG_FILES>
```
