# crawlers

A collection of crawlers used to crawl SEC data into SQL databases

Setting up Live-Crawler:
To build the live crawler docker image from the root directory (The file is named Dockerfile.live_crawler anticipating the creation of other dockerfiles):
docker build -f Dockerfile.live_crawler --platform linux/amd64 -t docker-image:test .

Create an ECR_REPO in AWS then run these docker commands to populate it:
docker tag docker-image:test <ECR_REPO_URI>:latest
docker push <ECR_REPO_URI>:latest

Create a lambda function from the container image you created in ECR
