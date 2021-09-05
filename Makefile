
NAME=terraform-provider-unifi
PLUGIN_PATH=$(HOME)/.terraform.d/plugins

all: build

build:
	go build -o $(NAME)


install: build
	install -d $(PLUGIN_PATH)
	install -m 775 $(NAME) $(PLUGIN_PATH)/

test:
	./controller.sh update
	./controller.sh start
	./controller.sh test || echo "Failed"
	./controller.sh stop
	
