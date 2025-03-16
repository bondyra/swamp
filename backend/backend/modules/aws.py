import configparser
import os
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
        if "metadata.profile" not in required_attrs:
            raise GenericQueryException("You need to provide __profile value to query AWS resource")
        if "metadata.region" not in required_attrs:
            raise GenericQueryException("You need to provide __region value to query AWS resource")
        profile, region = required_attrs["metadata.profile"], required_attrs["metadata.region"]
        async with aioboto3.Session(profile_name=profile, region_name=region).client(cls.boto_client_name) as client:
            async for i in cls._get(client):
                yield {
                    "metadata": {
                        "profile": profile,
                        "region": region,
                        **i["metadata"]
                    },
                    "data": i["data"]
                }

    @classmethod
    async def _get(self, client) -> AsyncGenerator[Dict, None]:
        raise NotImplementedError()

    @classmethod
    def attributes(cls) -> List[Attribute]:
        shp = boto3.client(cls.boto_client_name).meta.service_model.shape_for(cls.shape)
        return [
            Attribute(path="metadata.profile", description="AWS profile to use", query_required=True, allowed_values=_get_profiles()),
            Attribute(path="metadata.region", description="AWS region", query_required=True, allowed_values=_ALL_AWS_REGIONS),
            *cls._attributes_rec(shp)
        ]

    @classmethod
    def _attributes_rec(cls, obj, path=""):
        if "type_name" in dir(obj) and obj.type_name == "structure":
            for name, member in obj.members.items():
                new_path = f"{path}.{name}" if path else name
                yield from cls._attributes_rec(member, new_path)
        elif "type_name" in dir(obj) and obj.type_name == "list": # <- ignoring lists for now
            return
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
            yield {
                "metadata": {"id": item["VpcId"]},
                "data": item
            }


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
            yield {
                "metadata": {"id": item["SubnetId"]},
                "data": item
            }


class RouteTableHandler(LegacyAWSAPIHandler):
    boto_client_name = "ec2"
    shape = "RouteTable"

    @staticmethod
    def resource() -> str:
        return "rtb"

    @classmethod
    async def _get(cls, client) -> AsyncGenerator[Dict, None]:
        response = await client.describe_route_tables()
        for item in response["RouteTables"]:
            yield {
                "metadata": {"id": item["RouteTableId"]},
                "data": item
            }


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
            yield {
                "metadata": {"id": item["InternetGatewayId"]},
                "data": item
            }


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
            yield {
                "metadata": {"id": item["GroupId"]},
                "data": item
            }


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
            yield {
                "metadata": {"id": item["NatGatewayId"]},
                "data": item
            }


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
            yield {
                "metadata": {"id": item["AllocationId"]},
                "data": item
            }


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
            yield {
                "metadata": {"id": item["NetworkInterfaceId"]},
                "data": item
            }


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
            yield {
                "metadata": {"id": item["NetworkAclId"]},
                "data": item
            }


_PROFILES = []


def _get_profiles():
    global _PROFILES
    if _PROFILES:
        return _PROFILES
    credentials_parser = configparser.RawConfigParser()
    config_parser = configparser.RawConfigParser()
    if os.path.exists(f"{os.environ['HOME']}/.aws/credentials"):
        with open (f"{os.environ['HOME']}/.aws/credentials", "r") as f:
            credentials_parser.read_file(f)
    if os.path.exists(f"{os.environ['HOME']}/.aws/config"):
        with open (f"{os.environ['HOME']}/.aws/config", "r") as f:
            config_parser.read_file(f)
    _PROFILES = [
        *[section for section in credentials_parser.sections()],
        *["".join(section.split(" ")[1:]) for section in config_parser.sections() if section.startswith("profile ")]
    ]
    return _PROFILES


_ALL_AWS_REGIONS = [
    "us-east-2",
    "us-east-1",
    "us-west-1",
    "us-west-2",
    "af-south-1",
    "ap-east-1",
    "ap-south-2",
    "ap-southeast-3",
    "ap-southeast-5",
    "ap-southeast-4",
    "ap-south-1",
    "ap-northeast-3",
    "ap-northeast-2",
    "ap-southeast-1",
    "ap-southeast-2",
    "ap-southeast-7",
    "ap-northeast-1",
    "ca-central-1",
    "ca-west-1",
    "eu-central-1",
    "eu-west-1",
    "eu-west-2",
    "eu-south-1",
    "eu-west-3",
    "eu-south-2",
    "eu-north-1",
    "eu-central-2",
    "il-central-1",
    "mx-central-1",
    "me-south-1",
    "me-central-1",
    "sa-east-1",
    "us-gov-east-1",
    "us-gov-west-1"
]
