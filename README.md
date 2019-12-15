tfvalidate - A linter for Terraform Plans
==========

tfvalidate is a system for validating Terraform plans to ensure conformity in resource naming and tagging. It supports
multiple actions for both linting plans prior to apply, and determine mandatory reviewers based on resources changing
for CI pipelines.

Requires Terraform 0.12.9 or higher. Plans created prior to this version will likely fail to parse.

Installing
----------

```bash
go get github.com/justinm/tfvalidate
```
 
 
Usage
-----

### Linting plans

```bash
terraform plan -out=my.plan
tfvalidate -action lint my.plan
```

### Listing required approvers for a plan

```bash
terraform plan -out=my.plan
tfvalidate -action approvers my.plan
```


Capabilities
------------

tfvalidate works by ensuring attributes and their values meet certain specifications. For a full list of supported rules
and it's syntax, please see [.tfvalidate.yaml](.tfvalidate.yaml).

### Linting

Linting ensures planned resources follow policy, such as mandatory tagging and EC2 size limiting.

### Approvers

Approvers allows for a catered list of required users based on the resources being modified. Empower your CI by 
requiring additional review for specific, more sensitive resources.

### Machine Readable Output

Results are relayed by both specific exit codes and JSON data for easy machine readability.

#### Exit Codes
* 0 Success
* 1 Error
* 2 Validation Errors
* 3 Approvers Required

License
-------

This software uses compiles directly against the Terraform libraries. As such, this project will adopt the same licensing
as Terraform. Please see LICENSE for more information.
 