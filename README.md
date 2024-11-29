# Tool that tries to show you your cloud stuff
It contains two tools - CLI and UI.

## CLI
### What it does
It does two operations - `ls` and `get`. Former lists given resource types, latter provides more details about specific resource.

It essentialy tries to do something similar that `kubectl` does for K8S, but also for AWS, and with deep autocompletion support.

**! AWS only for now**

### Setup
1. Install wheel from current GitHub release
3. Enable autocomplete:
```
activate-global-python-argcomplete --user
```
On MacOS, you need to exec `compinit`:
```
autoload -U +X compinit 
compinit
```
(You might want to add those two lines to your `~/.zprofile`)

### How to use
Tool is bound to default AWS profile - so if you want to query other account, just change AWS_PROFILE.

List all VPCs:
```
swamp ls aws vpc
```

List all VPCs, but embed some field you want (JSON path compliant):
```
swamp ls aws vpc CidrBlock
```

(you can also add additional such paths as next positional arguments)

Get full JSON of a specific VPC:
```
swamp get aws vpc vpc-xxx
```

### Autocompletion
Please press tabs often - this way you can see what resource types are supported, what fields you could query, what resource IDs you can get, and so on.

## UI
I'm also trying to visualize resources and build some sensible graphs out of them

### How to use
For the time being I've bundled everything in a single docker image that you should run locally.

You'll need to expose two ports and mount your AWS creds for the container (I know it sounds stupid, but a normal support is on the way):
```
docker run -p 8000:8000 -p 3000:3000 -v $HOME/.aws:/root/.aws ghcr.io/bondyra/swamp:latest
```

This image spins up python backend server (that uses the same library CLI above uses) and react (reactflow based) frontend UI you access.

Go to `localhost:8000`, click fetch button to load the resources. Then you can organize stuff with another button, or change the theme.

That's essentially it for now.

