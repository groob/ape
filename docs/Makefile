all: middleman docker

middleman: 
	bundle exec middleman build --clean
docker:
	docker build -t groob/ape-docs .
