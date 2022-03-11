# qup-core

qup-core is a library to help manage various cloud services. It serves as logic that can be used to simplify actions like creating an ec2 instance, or publishing a lambda in golang.

## example

```golang
package main

import (
    ...
)

func main() {
    svc := ec2.New(session.New())
    client := cloud.NewClient(svc)

    //create an instance with some defaults
    instance1, _ := client.CreateInstance()

    // create an instance with a custom image
    instance2, _ := client.CreateInstance(cloud.WithImageId("ami-12345"))

    // create an instance with a custom image and custom tags
    instance3, _ := client.CreateInstance(
        cloud.WithImageId("ami-12345"),
        cloud.WithTags(map[string]string{
            "Name": "instance3",
            "env": "test"
        })
    )
}
```