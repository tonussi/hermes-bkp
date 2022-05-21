build_debug_hermes:
	docker-compose -f debug-hermes.docker-compose.yml up --build

build_hermes:
	docker-compose -f hermes.docker-compose.yml up --build

run_hermes:
	docker-compose -f hermes.docker-compose.yml up
