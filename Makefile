docker: proto docker-treeservice docker-treecli

proto:
	cd messages && make regenerate-docker

docker-treeservice:
	docker build -f treeservice.dockerfile -t treeservice .

docker-treecli:
	docker build -f treecli.dockerfile -t treecli .