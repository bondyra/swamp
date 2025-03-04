from typing import AsyncGenerator, Dict, List

import aioboto3
import boto3

from backend.model import Attribute, Handler, Provider, GenericQueryException


class AWS(Provider):
    @staticmethod
    def name() -> str:
        return "aws"

    @staticmethod
    def description() -> str:
        return "Provider for interacting with AWS resources via boto3"


class AWSHandler(Handler):
    @staticmethod
    def provider() -> str:
        return "aws"


class LegacyAWSAPIHandler(AWSHandler):
    boto_client_name = None
    shape = None

    @classmethod
    def description(cls) -> str:
        return boto3.client(cls.boto_client_name).meta.service_model.shape_for(cls.shape).documentation

    @classmethod
    async def get(cls, **required_attrs) -> AsyncGenerator[Dict, None]:
        if "__profile" not in required_attrs:
            raise GenericQueryException("You need to provide __profile value to query AWS resource")
        if "__region" not in required_attrs:
            raise GenericQueryException("You need to provide __region value to query AWS resource")
        profile, region = required_attrs["__profile"], required_attrs["__region"]
        async with aioboto3.Session(profile_name=profile, region_name=region).client(cls.boto_client_name) as client:
            async for r in cls._get(client):
                yield {
                    "__profile": profile,
                    "__region": region,
                    **r
                }

    @classmethod
    async def _get(self, client) -> AsyncGenerator[Dict, None]:
        raise NotImplementedError()

    @classmethod
    def attributes(cls) -> List[Attribute]:
        shp = boto3.client(cls.boto_client_name).meta.service_model.shape_for(cls.shape)
        return [
            Attribute(path="__profile", description="AWS profile to use", query_required=True),
            Attribute(path="__region", description="AWS region", query_required=True),
            *cls._attributes_rec(shp)
        ]

    @classmethod
    def _attributes_rec(cls, obj, path=""):
        if "type_name" in dir(obj) and obj.type_name == "structure":
            for name, member in obj.members.items():
                new_path = f"{path}.{name}" if path else name
                yield from cls._attributes_rec(member, new_path)
        # elif "type_name" in dir(obj) and obj.type_name == "list":  <- ignoring lists for now
        #     for name, member in obj.member.members.items():
        #         new_path = f"{path}[*].{name}"
        #         yield from cls._attributes_rec(member, new_path)
        else:
            yield Attribute(path=path, description=obj.documentation, query_required=False)


class VpcHandler(LegacyAWSAPIHandler):
    boto_client_name = "ec2"
    shape = "Vpc"

    @staticmethod
    def resource() -> str:
        return "vpc"

    @classmethod
    async def _get(cls, client) -> AsyncGenerator[Dict, None]:
        response = await client.describe_vpcs()
        for item in response["Vpcs"]:
            yield {"__id": item["VpcId"], **item}


class SubnetHandler(LegacyAWSAPIHandler):
    boto_client_name = "ec2"
    shape = "Subnet"

    @staticmethod
    def resource() -> str:
        return "subnet"

    @classmethod
    async def _get(cls, client) -> AsyncGenerator[Dict, None]:
        response = await client.describe_subnets()
        for item in response["Subnets"]:
            yield {"__id": item["SubnetId"], **item}


class RouteTableHandler(LegacyAWSAPIHandler):
    boto_client_name = "ec2"
    shape = "RouteTable"

    @staticmethod
    def resource() -> str:
        return "rtb"

    @classmethod
    async def _ls(cls, client) -> AsyncGenerator[Dict, None]:
        response = await client.describe_route_tables()
        for item in response["RouteTables"]:
            yield {"__id": item["RouteTableId"], **item}


class InternetGatewayHandler(LegacyAWSAPIHandler):
    boto_client_name = "ec2"
    shape = "InternetGateway"

    @staticmethod
    def resource() -> str:
        return "igw"

    @classmethod
    async def _get(cls, client) -> AsyncGenerator[Dict, None]:
        response = await client.describe_internet_gateways()
        for item in response["InternetGateways"]:
            yield {"__id": item["InternetGatewayId"], **item}


class SecurityGroupHandler(LegacyAWSAPIHandler):
    boto_client_name = "ec2"
    shape = "SecurityGroup"

    @staticmethod
    def resource() -> str:
        return "sg"

    @classmethod
    async def _get(cls, client) -> AsyncGenerator[Dict, None]:
        response = await client.describe_security_groups()
        for item in response["SecurityGroups"]:
            yield {"__id": item["GroupId"], **item}


class NATGatewayHandler(LegacyAWSAPIHandler):
    boto_client_name = "ec2"
    shape = "NatGateway"

    @staticmethod
    def resource() -> str:
        return "nat"

    @classmethod
    async def _get(cls, client) -> AsyncGenerator[Dict, None]:
        response = await client.describe_nat_gateways()
        for item in response["NatGateways"]:
            yield {"__id": item["NatGatewayId"], **item}


class ElasticIpHandler(LegacyAWSAPIHandler):
    boto_client_name = "ec2"
    shape = "Address"

    @staticmethod
    def resource() -> str:
        return "eip"

    @classmethod
    async def _get(cls, client) -> AsyncGenerator[Dict, None]:
        response = await client.describe_addresses()
        for item in response["Addresses"]:
            yield {"__id": item["AllocationId"], **item}


class NetworkInterfaceHandler(LegacyAWSAPIHandler):
    boto_client_name = "ec2"
    shape = "NetworkInterface"

    @staticmethod
    def resource() -> str:
        return "eni"

    @classmethod
    async def _get(cls, client) -> AsyncGenerator[Dict, None]:
        response = await client.describe_network_interfaces()
        for item in response["NetworkInterfaces"]:
            yield {"__id": item["NetworkInterfaceId"], **item}


class NetworkAclHandler(LegacyAWSAPIHandler):
    boto_client_name = "ec2"
    shape = "NetworkAcl"

    @staticmethod
    def resource() -> str:
        return "nacl"

    @classmethod
    async def _get(cls, client) -> AsyncGenerator[Dict, None]:
        response = await client.describe_network_acls()
        for item in response["NetworkAcls"]:
            yield {"__id": item["NetworkAclId"], **item}
