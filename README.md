# Graph explorer of cloud resources

## How to use
For the time being I've bundled everything in a single docker image that you should run locally.

You'll need to expose two ports (some port, e.g. 3000 for UI & REQUIRED 8000 for backend) and mount your AWS creds & kubeconfig in the container:
```
docker run -p 3000:80 -p 8000:8000 -v $HOME/.aws:/root/.aws -v $HOME/.kube:/root/.kube ghcr.io/bondyra/swamp:latest
```

This image spins up python backend server (simple layer over boto3 & kubernetes python SDK) and react (reactflow) frontend UI.

Go to `localhost:3000` to use the UI.

That's essentially it for now.
