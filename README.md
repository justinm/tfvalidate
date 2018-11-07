tfvalidate - A linter for Terraform Plans
==========

tfvalidate is a system for validating Terraform plans to ensure conformity in resource naming and tagging.

Installing
----------

```bash
go get github.com/justinm/tfvalidate
```
 
Usage
-----

```bash
terraform plan -out=my.plan
tfvalidate --plan my.plan
```

Capabilities
------------
tfvalidate works by ensuring attributes and their values meet certain specifications. For a full list of supported rules
and it's syntax, please see [.tfvalidate.json](.tfvalidate.json).
