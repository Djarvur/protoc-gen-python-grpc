# protoc-gen-python-grpc

This is a drop-in replacement for the same name [protoc](https://protobuf.dev/reference/python/python-generated/) plugin re-written in Golang.

Golang-written programs are easy to compile to any supported platform, easy to maintain
and better suite the nowadays requirements.

## Input and output

Input is exactly the same as any other protoc plugin.

Output is exactly the same as the [good old protoc-gen-python-grpc](https://github.com/grpc/grpc/blob/master/src/compiler/python_generator.cc) provides.

Please rise the issue in case of any difference found in these plugins output. 

## How it works

### Install

```
go install github.com/Djarvur/protoc-gen-python-grpc/cmd/protoc-gen-python-grpc@v0.0.1
```

Make sure your [protoc](https://grpc.io/docs/protoc-installation/)/[buf](https://buf.build/docs/installation) compiler can see the `protoc-gen-python-grpc` in the path.

### Run

Exactly the same as any other protoc plugin.

#### Parameters


### Generic template support

There is the embedded template for the `_pb2_grpc.py` files in the program.

Of course, you can provide your own in Go [text/template](https://pkg.go.dev/text/template) format utilising the following placeholders:

For each file:
- `{{.Name}}` - name of the source file as it is provided by protoc/buf
- `{{.Package}}` - proto package name
- `{{.Services}}` - list of the services defined in the proto file (see below)

For each service:
- `{{.Name}}` - name of the service
- `{{.Comment}}` - service-related comment, might be multi-line
- `{{.Methods}}` - list of the methods defined for this service in proto file (see below)

For each method:
- `{{.Name}}` - name of the method
- `{{.Comment}}` - method-related comment, might be multi-line
- `{{.Request}}` - name of the method request message        
- `{{.Response}}` - name of the method response message       
- `{{.ClientStreaming}}` - boolean indicates is this method client-streaming
- `{{.ServerStreaming}}` - boolean indicates is this method server-streaming

Also, the following functions are available to be used as template pipelines:
- `{{trimSuffix separator value}}` - returns a substring with the after-the-last-separator part removed
- `{{baseName separator value}}` - returns a substring after-the-last-separator
- `{{replace from to value}}` - returns a string with all the occurrences of `from` replaced to `to`
- `{{split separator value}}` - splitting value by separator and returns a list of substrings 
- `{{join separator ...values}}` - joining a values list to one string with values divided by separator

## Limitations and caveats

This program based on the [github.com/pseudomuto/protokit](https://pkg.go.dev/github.com/pseudomuto/protokit?utm_source=godoc) package,
which is a bit outdated and based on the [github.com/golang/protobuf](https://pkg.go.dev/github.com/golang/protobuf) package
which is deprecated. Still works though.
