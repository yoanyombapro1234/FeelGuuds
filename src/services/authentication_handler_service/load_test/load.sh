#!/bin/bash
# script takes as input the rpc method name, its the request input in json form, and the output type
# load ${operation} '{"id":1}' (html|json)

svcApiName=authentication_handler_serviceApi
svcProtoPath=../proto/authentication_handler_service.proto

function RunAllLoadTests(){
	echo "running all load tests"
	RunLoadTest CreateAccount '{"email":"yassir1234@gmail.com", "password":"honeybooboo17"}' html
	RunLoadTest UpdateAccount '{"id":1, "email":"yassir@gmail.com"}' html
	RunLoadTest LockAccount '{"id":1}' html
	RunLoadTest UnLockAccount '{"id":1}' html
	RunLoadTest GetAccount '{"id":1}' html
	RunLoadTest AuthenticateAccount '{"email":"yassir@gmail.com", "password":"honeybooboo17"}' html
	RunLoadTest LogoutAccount '{"id":1}' html
	RunLoadTest DeleteAccount '{"id":1}' html
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

function RunLoadTest(){
	echo "running create account api load test"
	operation=${1:-GetAccount}
	request=${2:-'{\"id\":1}}'}
	outputType=${3:-html}
	outputPath=""
	var flags=


	if [ $outputType == "html" ];
	then
			outputPath=".result/html/result.html"
			# shellcheck disable=SC2037
			flags="-O $outputType -o $outputPath"
	elif [ $outputType == "json" ]; then
			outputPath=".result/json/result.json"
			# shellcheck disable=SC2037
			flags="-O $outputType -o $outputPath"
	else
		  # shellcheck disable=SC2037
		  flags="-O pretty"
	fi

	echo $operation $request $outputType $flags

	# run linear load test
	ghz --insecure --async --proto "${svcProtoPath}" \
		--call proto.authentication_handler_serviceApi/"${operation}"  -c 10 -n 10000 --rps 200 \
		-d "${request}" 0.0.0.0:9897 $flags

	open "$outputPath"

	# Performs step load starting at 50 RPS and inscreasing by 10 RPS every 5s until we reach 10000 total requests.
	# The RPS load is distributed among the 10 workers, all sharing 1 connection
	ghz --insecure --async --proto "${svcProtoPath}" \
			--call proto.${svcApiName}/"${operation}" -n 10000 -c 10 --load-schedule=step \
			--load-start=50 --load-step=10 --load-step-duration=5s \
			-d "${request}" 0.0.0.0:9897 $flags
	open "$outputPath"

	# Performs step load starting at 50 RPS and inscreasing by 10 RPS every 5s until we reach 150 RPS at which point
	# the load is sustained at constant RPS rate until we reach 10000 total requests. The RPS load is
	# distributed among the 10 workers, all sharing 1 connection.
	ghz --insecure --async --proto "${svcProtoPath}" \
			--call proto.${svcApiName}/"${operation}" -n 10000 -c 10 --load-schedule=step \
			--load-start=50 --load-end=150 --load-step=10 \
			--load-step-duration=5s \
			-d "${request}" 0.0.0.0:9897 $flags
	open "$outputPath"

	# Performs step load starting at 50 RPS and inscreasing by 10 RPS every 5s until 60s has elapsed at which
	# point the load is sustained at that RPS rate until we reach 10000 total requests. The RPS load is
	# distributed among the 10 workers, all sharing 1 connection.
	ghz --insecure --async --proto "${svcProtoPath}" \
			--call proto.${svcApiName}/"${operation}" -n 10000 -c 10 --load-schedule=step --load-start=50 --load-step=10 \
			--load-step-duration=5s --load-max-duration=60s \
			--load-step-duration=5s \
			-d "${request}" 0.0.0.0:9897 $flags
	open "$outputPath"

	# Performs linear load starting at 200 RPS and decreasing by 2 RPS every 1s until 20 RPS has been reached,
	# at which point the load is sustained at that RPS rate until we reach 10000 total requests. The RPS
	# load is distributed among the 10 workers, all sharing 1 connection.
	ghz --insecure --async --proto "${svcProtoPath}" \
			--call proto.${svcApiName}/"${operation}" -n 10000 -c 10 --load-schedule=line \
			 --load-start=200 --load-step=-2 --load-end=50 \
			-d "${request}" 0.0.0.0:9897 $flags
	open "$outputPath"

	# Concurrent Load Testing
	# ======================= #
	# Performs RPS load of 200 RPS. The number of concurrent workers starts at 5 and is increased by 5 every 5s until
	# we reach 50 workers. At that point we keep the sustained 200 RPS load spread over the 50 workers until total
	# of 10000 requests is reached. That means as we increase the number of total concurrent workers, their share of RPS load decreases.
	ghz --insecure --async --proto "${svcProtoPath}" \
			--call proto.${svcApiName}/"${operation}" -n 10000 --rps 200 --concurrency-schedule=step \
			--concurrency-start=5 --concurrency-step=5 --concurrency-end=50 --concurrency-step-duration=5s \
			-d "${request}" 0.0.0.0:9897 $flags
	open "$outputPath"

	# Performs RPS load of 200 RPS. The number of concurrent workers starts at 10 and is increased by 10 every 5s until 60s has elapsed.
	# At that point we keep the sustained 200 RPS load spread over the same number of workers until total of 200000 requests is reached.
	ghz --insecure --async --proto "${svcProtoPath}" \
			--call proto.${svcApiName}/"${operation}" -n 200000 -rps 200 --concurrency-schedule=step \
			--concurrency-start=10 --concurrency-step=10 --concurrency-step-duration=5s --concurrency-max-duration=60s \
			-d "${request}" 0.0.0.0:9897 $flags
	open "$outputPath"

	# Performs RPS load of 200 RPS. The number of concurrent workers starts at 200 and is decreased linearly by 2 every 1s
	# until we are at 20 concurrent workers. At that point we keep the sustained 200 RPS load spread over the same number
	# of workers until total of 10000 requests is reached. As total number of active concurrent workers decreases,
	# their share of RPS load increases.
	ghz --insecure --async --proto "${svcProtoPath}" \
			--call proto.${svcApiName}/"${operation}" -n 10000 --rps 200 --concurrency-schedule=line \
			--concurrency-start=200 --concurrency-step=-2 --concurrency-end=20 \
			-d "${request}" 0.0.0.0:9897 $flags
	open "$outputPath"
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

RunAllLoadTests
