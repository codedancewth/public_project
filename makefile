export GO111MODULE=on
PWD := $(shell echo `pwd`)

.PHONY: account agent relation room_pk wish_gift base_proto chat\
 	heartbeat_game msgdispatcher recommend_tagging strategic user\
 	community greedy_yaplay recommend_app onlyfans onlyfans_offline activity_center

all: account agent relation room_pk wish_gift base_proto chat\
 	heartbeat_game msgdispatcher recommend_tagging strategic user\
 	community greedy_yaplay recommend_app onlyfans onlyfans_offline welfare_center activity_center


publicproject:
	powerproto build -r -p $(PWD)/proto/public_project
