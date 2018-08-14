build:
	docker-compose run app ./codeship/build.sh

slsdeploy:
	docker-compose run app sls deploy

deploy: dep build slsdeploy

destroy:
	docker-compose run app sls remove

dep:
	docker-compose run app dep ensure

test:
	docker-compose run test ./codeship/test.sh

codeshipsetup: dep build

dynamo:
	docker-compose up -d dynamo
	sleep 5
	docker-compose run init-dynamo

clean:
	docker-compose kill
	docker-compose rm -f