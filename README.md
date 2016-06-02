# notifier
A library for sending native desktop notifications from Go. This uses
platform-specific helper libraries as follows:

-OSX: [Terminal Notifier](https://github.com/julienXX/terminal-notifier)
-Windows: [notifu](https://www.paralint.com/projects/notifu/)

Those libraries are embedded directly in Go in this library, so there are no external dependencies or expectations for installations on the user's system. These libraries were also chosen for their small size, particularly in the case of notifu, which is far smaller than things like Growl or Snarl.

To generate updated embedded binaries, you can run, for example:

```
go get -u github.com/jteeuwen/go-bindata/...
cd platform/terminal-notifier-1.6.3
go-bindata terminal-notifier.app/...
mv bindata.go ../../osx
```

You then have to manually change the package in `bindata.go` to `osx` instead of `main` in that
case.

This is currently a work in progress and only runs on OSX and Windows and embeds
all binaries for all platforms instead of dynamically only including the
required platform at build time, for example.

The platform directory is only here to serve as a reference for the native binaries
being used.
