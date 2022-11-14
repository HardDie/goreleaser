# goreleaser

## How to use:
```
goreleaser build --name 'MyApp' \
	--company 'org.myorg.myapp' \
	--image 'path/to/image.png' \
	--license 'Licensed under GPLv3.' \
	--version 'v1.2.3' \
	--ldflags '-X main.Version=v1.2.3' \
	--path 'cmd/myapp/main.go'
```
As a result, you will get a 'release' folder with upload-ready archives for linux, darwin (MacOS) and windows platforms.
