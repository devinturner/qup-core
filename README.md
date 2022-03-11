# qup-core

qup-core is a library to help manage various cloud services. It serves as logic that can be used to simplify actions like creating an ec2 instance, or publishing a lambda in golang.

## services
- [] aws - ec2
- [] aws - lambda
- [] aws - dynamodb

## example

```golang
package main

import (
    ...
)

func main() {
    svc := ec2.New(session.New())
    client := cloud.NewClient(svc)

    opts := []cloud.InstanceOption{
        cloud.With
    }
    instance, err := client.CreateInstance()
    if err {
        ...
    }
    ...
}
```