
all: staticpage pack

staticpage:
	echo web building...
	cd web && bower i

pack:
	echo 'package docs' > docs/docs.go
	bee pack -exs=".go:.DS_Store:.tmp:.log:.pid:.pprof:.memprof:Makefile:.bowerrc:.git:.gitignore:.travis.yml" -exp=".:swagger"

clean:
	rm -rf *.pprof *.memprof
	rm *.tar.gz || echo 'nothing'
	rm lastupdate.tmp || echo 'nothing'
	rm routers/commentsRouter_moduleab_server_controllers.go || echo 'nothing'
	rm moduleab_server || echo 'nothing'
	rm -rf web/app/bower_components
