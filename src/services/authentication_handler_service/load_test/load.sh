#!/bin/bash

function RunAllLoadTests(){
	echo "running all load tests"
}

function RunLinearLoadScenario(){
	echo "running linear load scenario"
}

# https://linuxize.com/post/bash-functions/
function RunConcurrentLoadScenario(){
	echo "running concurrent load scenario"
}


function CreateAccountLoadTest(){
	echo "running create account api load test"
}

function GetAccountLoadTest(){
	echo "running create account api load test"
}

function UpdateAccountLoadTest(){
	echo "running create account api load test"
}

function LockAndUnlockAccountLoadTest(){
	echo "running create account api load test"
}

function ArchiveAccountLoadTest(){
	echo "running create account api load test"
}

function AuthenticateAccountLoadTest(){
	echo "running create account api load test"
}

#ghz --insecure --async --proto ../proto/authentication_handler_service.proto \
#		--call proto.authentication_handler_serviceApi/CreateAccount \
#  	-d '{"email":"yoan@gmail.com", "password":"yoanyomba"}' 0.0.0.0:9897

# linear load testing
# ============================================== #
ghz --insecure --async --proto ../proto/authentication_handler_service.proto \
		--call proto.authentication_handler_serviceApi/GetAccount -c 10 -n 10000 --rps 200 \
		-d '{"id":1}' 0.0.0.0:9897 #-O html

# Performs step load starting at 50 RPS and inscreasing by 10 RPS every 5s until we reach 10000 total requests.
# The RPS load is distributed among the 10 workers, all sharing 1 connection
ghz --insecure --async --proto ../proto/authentication_handler_service.proto \
		--call proto.authentication_handler_serviceApi/GetAccount -n 10000 -c 10 --load-schedule=step \
		--load-start=50 --load-step=10 --load-step-duration=5s \
		-d '{"id":1}' 0.0.0.0:9897 #-O html

# Performs step load starting at 50 RPS and inscreasing by 10 RPS every 5s until we reach 150 RPS at which point
# the load is sustained at constant RPS rate until we reach 10000 total requests. The RPS load is
# distributed among the 10 workers, all sharing 1 connection.
ghz --insecure --async --proto ../proto/authentication_handler_service.proto \
		--call proto.authentication_handler_serviceApi/GetAccount -n 10000 -c 10 --load-schedule=step \
		--load-start=50 --load-end=150 --load-step=10 \
		--load-step-duration=5s \
		-d '{"id":1}' 0.0.0.0:9897 #-O html


# Performs step load starting at 50 RPS and inscreasing by 10 RPS every 5s until 60s has elapsed at which
# point the load is sustained at that RPS rate until we reach 10000 total requests. The RPS load is
# distributed among the 10 workers, all sharing 1 connection.
ghz --insecure --async --proto ../proto/authentication_handler_service.proto \
		--call proto.authentication_handler_serviceApi/GetAccount -n 10000 -c 10 --load-schedule=step --load-start=50 --load-step=10 \
		--load-step-duration=5s --load-max-duration=60s \
		--load-step-duration=5s \
		-d '{"id":1}' 0.0.0.0:9897 #-O html


# Performs linear load starting at 200 RPS and decreasing by 2 RPS every 1s until 20 RPS has been reached,
# at which point the load is sustained at that RPS rate until we reach 10000 total requests. The RPS
# load is distributed among the 10 workers, all sharing 1 connection.
ghz --insecure --async --proto ../proto/authentication_handler_service.proto \
		--call proto.authentication_handler_serviceApi/GetAccount -n 10000 -c 10 --load-schedule=line \
		 --load-start=200 --load-step=-2 --load-end=50 \
		-d '{"id":1}' 0.0.0.0:9897 #-O html

# Concurrent Load Testing
# ======================= #
# Performs RPS load of 200 RPS. The number of concurrent workers starts at 5 and is increased by 5 every 5s until
# we reach 50 workers. At that point we keep the sustained 200 RPS load spread over the 50 workers until total
# of 10000 requests is reached. That means as we increase the number of total concurrent workers, their share of RPS load decreases.
ghz --insecure --async --proto ../proto/authentication_handler_service.proto \
		--call proto.authentication_handler_serviceApi/GetAccount -n 100000 --rps 200 --concurrency-schedule=step \
		--concurrency-start=5 --concurrency-step=5 --concurrency-end=50 --concurrency-step-duration=5s \
		-d '{"id":1}' 0.0.0.0:9897 #-O html

# Performs RPS load of 200 RPS. The number of concurrent workers starts at 10 and is increased by 10 every 5s until 60s has elapsed.
# At that point we keep the sustained 200 RPS load spread over the same number of workers until total of 20000 requests is reached.
ghz --insecure --async --proto ../proto/authentication_handler_service.proto \
		--call proto.authentication_handler_serviceApi/GetAccount -n 20000 -rps 200 --concurrency-schedule=step \
		--concurrency-start=10 --concurrency-step=10 --concurrency-step-duration=5s --concurrency-max-duration=60s \
		-d '{"id":1}' 0.0.0.0:9897 #-O html


# Performs RPS load of 200 RPS. The number of concurrent workers starts at 200 and is decreased linearly by 2 every 1s
# until we are at 20 concurrent workers. At that point we keep the sustained 200 RPS load spread over the same number
# of workers until total of 10000 requests is reached. As total number of active concurrent workers decreases,
# their share of RPS load increases.
ghz --insecure --async --proto ../proto/authentication_handler_service.proto \
		--call proto.authentication_handler_serviceApi/GetAccount -n 10000 --rps 200 --concurrency-schedule=line \
		--concurrency-start=200 --concurrency-step=-2 --concurrency-end=20 \
		-d '{"id":1}' 0.0.0.0:9897 #-O html
