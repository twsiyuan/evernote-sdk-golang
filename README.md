# Evernote SDK Golang

This project was simple code generated from [Evernote-thrift](https://github.com/evernote/evernote-thrift)-1.28 using Apache Thrift.


# How to generate yourself code

* Download and install the newest [Apache Thrift](https://thrift.apache.org/). 
* Clone or download the official [evernote-thrift repository](https://github.com/evernote/evernote-thrift)
* Generate golang code with this command for each .thrift (Errors.thrift, Limits.thrift, NoteStore.thrift, Types.thrift, and UserStore.thrift):

    ```thrift -strict -nowarn --allow-64bit-consts --allow-neg-keys --gen go:package_prefix=github.com/your-github-id/evernote-sdk-golang/ evernote-thrift/src/UserStore.thrift```

* Modify import package path from 'git.apache.org/thrift.git/lib/go/thrift' to 'github.com/apache/thrift/lib/go/thrift'
* Install thrift golang package from github: ```go get github.com/apache/thrift/lib/go/thrift```
* Fix some compiler errors...(had some errors using thrift 0.10.0 generator)
* Enjoy!

# Example

see [utils/utils_test.go](/utils/utils_test.go) and [utils/oauth_test.go](/utils/oauth_test.go)

