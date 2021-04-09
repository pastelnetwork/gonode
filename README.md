# go-legroast

## GoLang CGo wrapper for LegRoast

Brings [LegRoast](https://github.com/pastelnetwork/legroast) C API into GoLang.

GoLang module contains pre-built LegRoast static library for Windows, Linux and MacOS platforms.

## MacOS Support

LegRoast expects pthread's stack size on MacOS bigger than its defaults.  
To prevent random crashes `/usr/local/go/src/runtime/cgo/gcc_darwin_amd64.c` should be patched manually on build machine as following:

```
40a41,45
> #ifdef CGO_STACK_SIZE
>     size = CGO_STACK_SIZE;
>     pthread_attr_setstacksize(&attr, size);
> #endif
>
$ CGO_CFLAGS="-DCGO_STACK_SIZE=0x200000" make
```


