module github.com/elivoa/syd

go 1.14

require (
	github.com/axgle/mahonia v0.0.0-20180208002826-3358181d7394
	github.com/elivoa/got v0.0.0
	github.com/elivoa/gxl v0.0.0
	github.com/go-sql-driver/mysql v1.5.0
	github.com/golang/protobuf v1.4.2 // indirect
	golang.org/x/oauth2 v0.0.0-20200107190931-bf48bf16ab8d
	google.golang.org/appengine v1.6.6 // indirect
	google.golang.org/protobuf v1.24.0 // indirect
	gopkg.in/oauth2.v3 v3.12.0
	gopkg.in/session.v1 v1.0.1
)

replace (
	github.com/elivoa/got => ../got
	github.com/elivoa/gxl => ../gxl
)
