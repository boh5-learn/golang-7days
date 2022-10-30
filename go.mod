module golang-7days

go 1.19

require (
	gee v0.0.0
	geecache v0.0.0
)

replace (
	gee => ./gee_web/gee
	geecache => ./gee_cache/geecache
)
