import aioboto3


_config = {
    "vpc": {
        "client": "ec2",
        "ls": {
            "request": lambda c: c.describe_vpcs(),
            "response": lambda r: ((i["VpcId"], i) for i in r["Vpcs"])
        },
        "get": {
            "request": lambda c, i: c.describe_vpcs(VpcIds=[i]),
            "response": lambda r: (r["Vpcs"][0]["VpcId"], r["Vpcs"][0])
        },
        "shape": "Vpc",
        "default_props": [
            "CidrBlock"
        ]
    },
    "subnet": {
        "client": "ec2",
        "ls": {
            "request": lambda c: c.describe_subnets(),
            "response": lambda r: ((i["SubnetId"], i) for i in r["Subnets"])
        },
        "get": {
            "request": lambda c, i: c.describe_subnets(SubnetIds=[i]),
            "response": lambda r: (r["Subnets"][0]["SubnetId"], r["Subnets"][0])
        },
        "shape": "Subnet",
        "id": "SubnetId",
        "default_props": [
            "VpcId", "AvailabilityZone", "CidrBlock"
        ]
    },
    "rtb": {
        "client": "ec2",
        "ls": {  # todo - more complex than this
            "request": lambda c: c.describe_route_tables(),
            "response": lambda r: ((i["RouteTableId"], i) for i in r["RouteTables"])
        },
        "get": {
            "request": lambda c, i: c.describe_route_tables(RouteTableIds=[i]),
            "response": lambda r: (r["RouteTables"][0]["RouteTableId"], r["RouteTables"][0])
        },
        "shape": "RouteTable",
        "default_props": [
            "Routes"
        ]
    },
    "sg": {
        "client": "ec2",
        "ls": {  # todo - more complex than this
            "request": lambda c: c.describe_security_groups(),
            "response": lambda r: ((i["GroupId"], i) for i in r["SecurityGroups"])
        },
        "get": {
            "request": lambda c, i: c.describe_security_groups(GroupIds=[i]),
            "response": lambda r: (r["SecurityGroups"][0]["GroupId"], r["SecurityGroups"][0])
        },
        "shape": "SecurityGroup",
        "default_props": [
            "GroupName", "IpPermissions"
        ]
    },
    "igw": {
        "client": "ec2",
        "ls": {
            "request": lambda c: c.describe_internet_gateways(),
            "response": lambda r: ((i["InternetGatewayId"], i) for i in r["InternetGateways"])
        },
        "get": {
            "request": lambda c, i: c.describe_internet_gateways(InternetGatewayIds=[i]),
            "response": lambda r: (r["InternetGateways"][0]["InternetGatewayId"], r["InternetGateways"][0])
        },
        "shape": "InternetGateway",
        "default_props": [
            "Attachments"
        ]
    },
    "nat": {
        "client": "ec2",
        "ls": {
            "request": lambda c: c.describe_nat_gateways(),
            "response": lambda r: ((i["NatGatewayId"], i) for i in r["NatGateways"])
        },
        "get": {
            "request": lambda c, i: c.describe_nat_gateways(NatGatewayIds=[i]),
            "response": lambda r: (r["NatGateways"][0]["NatGatewayId"], r["NatGateways"][0])
        },
        "shape": "NatGateway",
        "default_props": []
    },
    "eip": {
        "client": "ec2",
        "ls": {
            "request": lambda c: c.describe_addresses(),
            "response": lambda r: ((i["AllocationId"], i) for i in r["Addresses"])
        },
        "get": {
            "request": lambda c, i: c.describe_addresses(AllocationIds=[i]),
            "response": lambda r: (r["Addresses"][0]["AllocationId"], r["Addresses"][0])
        },
        "shape": "Address",
        "default_props": [
            "PublicIp"
        ]
    },
    "eni": {
        "client": "ec2",
        "ls": {
            "request": lambda c: c.describe_network_interfaces(),
            "response": lambda r: ((i["NetworkInterfaceId"], i) for i in r["NetworkInterfaces"])
        },
        "get": {
            "request": lambda c, i: c.describe_network_interfaces(NetworkInterfaceIds=[i]),
            "response": lambda r: (r["NetworkInterfaces"][0]["NetworkInterfaceId"], r["NetworkInterfaces"][0])
        },
        "shape": "NetworkInterface",
        "default_props": [
            "PrivateIpAddresses", "Association"
        ]
    },
    "nacl": {
        "client": "ec2",
        "list": {
            "method": "describe_network_acls"
        },
        "ls": {
            "request": lambda c: c.describe_network_acls(),
            "response": lambda r: ((i["NetworkAclId"], i) for i in r["NetworkAcls"])
        },
        "get": {
            "request": lambda c, i: c.describe_network_acls(NetworkAclIds=[i]),
            "response": lambda r: (r["NetworkAcls"][0]["NetworkAclId"], r["NetworkAcls"][0])
        },
        "shape": "NetworkAcl",
        "default_props": [
            "Entries", "IsDefault"
        ]
    }
}


async def ls(thing_type):
    t = _config[thing_type]
    async with aioboto3.Session().client(t["client"]) as c:
        response = await t["ls"]["request"](c)
    for _ in t["ls"]["response"](response):
        yield _


async def get(thing_type, id):
    t = _config[thing_type]
    async with aioboto3.Session().client(t["client"]) as c:
        response = await t["get"]["request"](c, id)
        return t["get"]["response"](response)


def thing_types():
    return list(_config.keys())


def default_props(thing_type):
    return _config[thing_type].get("default_props", []) if thing_type in _config else []


async def schema_ls(thing_type):
    t = _config[thing_type]
    async with aioboto3.Session().client(t["client"]) as c:
        shp = c.meta.service_model.shape_for(t["shape"])
    return list(json_paths(shp))


schema_get = schema_ls


def json_paths(obj, path=""):
    if "type_name" in dir(obj) and obj.type_name == "structure":
        for name, member in obj.members.items():
            new_path = f"{path}.{name}" if path else name
            yield from json_paths(member, new_path)
    elif "type_name" in dir(obj) and obj.type_name == "list":
        for name, member in obj.member.members.items():
            new_path = f"{path}[*].{name}"
            yield from json_paths(member, new_path)
    else:
        yield path
