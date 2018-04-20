# Digital Lumens Devops Test Exercise

In this exercise, we need to take a trivial Go server, containerize it, and deploy it in a Docker host.
The tricky part is that we want to use TLS certificates for authentication, so the deployed container must be configured to support this.

We provide the simple server and a client to test with, as well as a head start on generating a certificate authority and certificates for both the server and client. The server should be deployable to both an EC2 instance (which we will provide) and a local Docker host (like Docker for Mac). Make a private fork of this repo and submit your answer in the form of a PR against your private fork.

You will also be invited to our Slack server so that you can ask any questions you might have.

### Base Exercise

  1. Containerize the server to run in Docker, listening on port 8080
  2. Deploy to a local Docker host (https://localhost:8080) and verify that the client can connect and authenticate.
  3. Deploy to a vanilla Docker host in EC2 (which we will provide).

### Additional Challenges

  1. Create a Docker image that does not contain any secrets or credentials (in this case, certs and keys)
  2. Create a Docker image that contains only the server binary and not the sources (e.g. server.go)
