from typing import AsyncGenerator, Dict

import aioboto3
import boto3

from overread.model import Module, Resource, Result


description = "Module for interacting with AWS resources"


class AWS(Module, name="aws"):
    pass


class AWSResource(Resource, module=AWS, name=""):
    boto_client_name = None
    shape_ls = None
    shape_get = None

    def description(cls) -> str:
        return boto3.client(cls.boto_client_name).meta.service_model.shape_for(cls.shape_ls).documentation

    async def ls(cls) -> AsyncGenerator[Result]:
        raise NotImplementedError()

    async def get(cls, resource_id) -> Result:
        raise NotImplementedError()

    def schema_ls(cls) -> Dict[str, str]:
        shp = boto3.client(cls.boto_client_name).meta.service_model.shape_for(cls.shape_ls)
        return dict(cls._schema_rec(shp))

    def schema_get(cls) -> Dict[str, str]:
        shp = boto3.client(cls.boto_client_name).meta.service_model.shape_for(cls.shape_get)
        return dict(cls._schema_rec(shp))

    def _schema_rec(cls, obj, path=""):
        if "type_name" in dir(obj) and obj.type_name == "structure":
            for name, member in obj.members.items():
                new_path = f"{path}.{name}" if path else name
                yield from cls._schema_rec(member, new_path)
        elif "type_name" in dir(obj) and obj.type_name == "list":
            for name, member in obj.member.members.items():
                new_path = f"{path}[*].{name}"
                yield from cls._schema_rec(member, new_path)
        else:
            yield (path, obj.documentation)


class Vpc(AWSResource, name="vpc"):
    boto_client_name = "ec2"
    shape_ls = "Vpc"
    shape_get = "Vpc"

    async def ls(cls) -> AsyncGenerator[Result]:
        async with aioboto3.client(cls.boto_client_name) as client:
            response = await client.describe_vpcs()
            for vpc in response["Vpcs"]:
                yield Result(vpc["VpcId"], vpc)

    async def get(cls, resource_id) -> Result:
        async with aioboto3.client(cls.boto_client_name) as client:
            response = await client.describe_vpcs(VpcIds=[resource_id])
            return Result(response["Vpcs"][0]["VpcId"], response["Vpcs"][0])

_config = {
    # "vpc": {
    #     "client": "ec2",
    #     "ls": {
    #         "request": lambda c: c.describe_vpcs(),
    #         "response": lambda r: ((i["VpcId"], i) for i in r["Vpcs"]),
    #         "shape": "Vpc"
    #     },
    #     "get": {
    #         "request": lambda c, i: c.describe_vpcs(VpcIds=[i]),
    #         "response": lambda r: (r["Vpcs"][0]["VpcId"], r["Vpcs"][0]),
    #         "shape": "Vpc"
    #     }
    # },
    "subnet": {
        "client": "ec2",
        "ls": {
            "request": lambda c: c.describe_subnets(),
            "response": lambda r: ((i["SubnetId"], i) for i in r["Subnets"]),
            "shape": "Subnet"
        },
        "get": {
            "request": lambda c, i: c.describe_subnets(SubnetIds=[i]),
            "response": lambda r: (r["Subnets"][0]["SubnetId"], r["Subnets"][0]),
            "shape": "Subnet"
        }
    },
    "rtb": {
        "client": "ec2",
        "ls": {  # todo - more complex than this
            "request": lambda c: c.describe_route_tables(),
            "response": lambda r: ((i["RouteTableId"], i) for i in r["RouteTables"]),
            "shape": "RouteTable"
        },
        "get": {
            "request": lambda c, i: c.describe_route_tables(RouteTableIds=[i]),
            "response": lambda r: (r["RouteTables"][0]["RouteTableId"], r["RouteTables"][0]),
            "shape": "RouteTable"
        }
    },
    "sg": {
        "client": "ec2",
        "ls": {  # todo - more complex than this
            "request": lambda c: c.describe_security_groups(),
            "response": lambda r: ((i["GroupId"], i) for i in r["SecurityGroups"]),
            "shape": "SecurityGroup"
        },
        "get": {
            "request": lambda c, i: c.describe_security_groups(GroupIds=[i]),
            "response": lambda r: (r["SecurityGroups"][0]["GroupId"], r["SecurityGroups"][0]),
            "shape": "SecurityGroup"
        }
    },
    "igw": {
        "client": "ec2",
        "ls": {
            "request": lambda c: c.describe_internet_gateways(),
            "response": lambda r: ((i["InternetGatewayId"], i) for i in r["InternetGateways"]),
            "shape": "InternetGateway"
        },
        "get": {
            "request": lambda c, i: c.describe_internet_gateways(InternetGatewayIds=[i]),
            "response": lambda r: (r["InternetGateways"][0]["InternetGatewayId"], r["InternetGateways"][0]),
            "shape": "InternetGateway"
        }
    },
    "nat": {
        "client": "ec2",
        "ls": {
            "request": lambda c: c.describe_nat_gateways(),
            "response": lambda r: ((i["NatGatewayId"], i) for i in r["NatGateways"]),
            "shape": "NatGateway"
        },
        "get": {
            "request": lambda c, i: c.describe_nat_gateways(NatGatewayIds=[i]),
            "response": lambda r: (r["NatGateways"][0]["NatGatewayId"], r["NatGateways"][0]),
            "shape": "NatGateway"
        }
    },
    "eip": {
        "client": "ec2",
        "ls": {
            "request": lambda c: c.describe_addresses(),
            "response": lambda r: ((i["AllocationId"], i) for i in r["Addresses"]),
            "shape": "Address"
        },
        "get": {
            "request": lambda c, i: c.describe_addresses(AllocationIds=[i]),
            "response": lambda r: (r["Addresses"][0]["AllocationId"], r["Addresses"][0]),
            "shape": "Address"
        }
    },
    "eni": {
        "client": "ec2",
        "ls": {
            "request": lambda c: c.describe_network_interfaces(),
            "response": lambda r: ((i["NetworkInterfaceId"], i) for i in r["NetworkInterfaces"]),
            "shape": "NetworkInterface"
        },
        "get": {
            "request": lambda c, i: c.describe_network_interfaces(NetworkInterfaceIds=[i]),
            "response": lambda r: (r["NetworkInterfaces"][0]["NetworkInterfaceId"], r["NetworkInterfaces"][0]),
            "shape": "NetworkInterface"
        }
    },
    "nacl": {
        "client": "ec2",
        "list": {
            "method": "describe_network_acls"
        },
        "ls": {
            "request": lambda c: c.describe_network_acls(),
            "response": lambda r: ((i["NetworkAclId"], i) for i in r["NetworkAcls"]),
            "shape": "NetworkAcl"
        },
        "get": {
            "request": lambda c, i: c.describe_network_acls(NetworkAclIds=[i]),
            "response": lambda r: (r["NetworkAcls"][0]["NetworkAclId"], r["NetworkAcls"][0]),
            "shape": "NetworkAcl"
        }
    }
}

# def ___resource_types():
#     clients = {c: boto3.client(c) for c in set(v["client"] for v in _config.values())}
#     return {r: clients[t["client"]].meta.service_model.shape_for(t["ls"]["shape"]).documentation for r, t in _config.items()}
