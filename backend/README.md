docker build -t albus-auth:latest .

docker run -d --name albus-auth-test -p 6000:6000 --network albus-network albus-auth:latest