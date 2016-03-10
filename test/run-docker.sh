#!/bin/bash
#
# Build and run a docker image for Boulder. This is suitable for running
# repeatedly during development because Docker will cache the image it builds,
# and will only re-do the minimum necessary.
#
# NOTE: Currently we're not able to effectively cache the DB setup steps,
# because setting up the DB depends on source files in the Boulder repo. So any
# time source files change, Docker treats that as potentially invalidating the
# steps that came after the COPY. In theory we could add a step that copies only
# the files necessary to do the migrations, run them, and then copy the rest of
# the source.
set -o errexit
cd $(dirname $0)/..

# By default use host networking so non-docker clients running on the
# host can connect to the service using 127.0.0.1 and port
# 4000. NET_CONTAINER allows to replace that with a name or ID of
# another container where a client runs or that client uses for its
# networking. The latter is the case with Kubernetes.

docker_net=host
if [[ $NET_CONTAINER ]]; then
    docker_net="container:$NET_CONTAINER"
fi

# helper function to return the state of the container (true if running, false if not)
is_running(){
	local name=$1
	local state=$(docker inspect --format "{{.State.Running}}" $name 2>/dev/null)

	if [[ "$state" == "false" ]]; then
		# the container is up but not running
		# we should remove it so we can bring up another
		docker rm $name
	fi
	echo $state
}

if [[ "$(is_running boulder-mysql)" != "true" ]]; then
	# bring up mysql mariadb container - no need to specify 3306
	# port with host or container networking
	docker run -d \
		--net "$docker_net" \
		-e MYSQL_ALLOW_EMPTY_PASSWORD=yes \
		--name boulder-mysql \
		mariadb:10.0
fi

if [[ "$(is_running boulder-rabbitmq)" != "true" ]]; then
	# bring up rabbitmq container - no need to specify 5672 port
	# with host or container networking
	docker run -d \
		--net "$docker_net" \
		--name boulder-rabbitmq \
		rabbitmq:3
fi

# build the boulder docker image
docker build --rm --force-rm -t letsencrypt/boulder .

# In order to talk to a letsencrypt client running on the host, the fake DNS
# client used in Boulder's start.py needs to know what the host's IP is from the
# perspective of the container. The default value is 127.0.0.1. If you'd
# like your Boulder instance to always talk to some other host, you can set
# FAKE_DNS to that host's IP address.
fake_dns_args=()
if [[ $FAKE_DNS ]]; then
    fake_dns_args=(-e "FAKE_DNS=$FAKE_DNS")
fi

# run the boulder container
# The excluding `-d` command makes the instance interactive, so you can kill
# the boulder container with Ctrl-C.
docker run --rm -it \
	--net "$docker_net" \
	-e MYSQL_CONTAINER=yes \
	"${fake_dns_args[@]}" \
	--name boulder \
	letsencrypt/boulder
