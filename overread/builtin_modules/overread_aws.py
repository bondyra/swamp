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
        "list": {
            "method": "describe_subnets"
        },
        "shape": "Subnet",
        "field": "Subnets",
        "id": "SubnetId",
        "default_props": [
            "AvailabilityZone", "CidrBlock"
        ]
    },
    "rtb": {
        "client": "ec2",
        "list": {
            "method": "describe_route_tables"
        },
        "shape": "Dupa",
        "field": "RouteTables",
        "id": "RouteTableId",
        "default_props": [
            "Routes"
        ]
    },
    "sg": {
        "client": "ec2",
        "list": {
            "method": "describe_security_groups"
        },
        "shape": "Dupa",
        "field": "SecurityGroups",
        "id": "GroupId",
        "default_props": [
            "GroupName", "IpPermissions"
        ]
    },
    "igw": {
        "client": "ec2",
        "list": {
            "method": "describe_internet_gateways"
        },
        "shape": "Dupa",
        "field": "InternetGateways",
        "id": "InternetGatewayId",
        "default_props": [
            "Attachments"
        ]
    },
    "nat": {
        "client": "ec2",
        "list": {
            "method": "describe_nat_gateways"
        },
        "shape": "Dupa",
        "field": "NatGateways",
        "id": "NatGatewayId",
        "default_props": []
    },
    "eip": {
        "client": "ec2",
        "list": {
            "method": "describe_addresses"
        },
        "shape": "Dupa",
        "field": "Addresses",
        "id": "AllocationId",
        "default_props": [
            "PublicIp"
        ]
    },
    "eni": {
        "client": "ec2",
        "list": {
            "method": "describe_network_interfaces"
        },
        "shape": "Dupa",
        "field": "NetworkInterfaces",
        "id": "NetworkInterfaceId",
        "default_props": [
            "PrivateIpAddresses", "Association"
        ]
    },
    "nacl": {
        "client": "ec2",
        "list": {
            "method": "describe_network_acls"
        },
        "shape": "Dupa",
        "field": "NetworkAcls",
        "id": "NetworkAclId",
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
