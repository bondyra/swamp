from typing import AsyncGenerator, Dict

import aioboto3
import boto3

from swamp.model import Handler, Module, Result


class AWS(Module):
    @staticmethod
    def name() -> str:
        return "aws"

    @staticmethod
    def description() -> str:
        return "Module for interacting with AWS resources"


class AWSHandler(Handler):
    @staticmethod
    def module() -> str:
        return "aws"


class LegacyAWSAPIHandler(AWSHandler):
    boto_client_name = None
    shape_ls = None
    shape_get = None

    @classmethod
    def description(cls) -> str:
        return boto3.client(cls.boto_client_name).meta.service_model.shape_for(cls.shape_ls).documentation

    @classmethod
    async def ls(cls) -> AsyncGenerator[Result, None]:
        async with aioboto3.Session().client(cls.boto_client_name) as client:
            async for r in cls._ls(client):
                yield r

    @classmethod
    async def _ls(self, client) -> AsyncGenerator[Result, None]:
        raise NotImplementedError()

    @classmethod
    async def get(cls, resource_id: str) -> Result:
        async with aioboto3.Session().client(cls.boto_client_name) as client:
            return await cls._get(client, resource_id)

    @classmethod
    async def _get(self) -> Result:
        raise NotImplementedError()

    @classmethod
    def schema_ls(cls) -> Dict[str, str]:
        shp = boto3.client(cls.boto_client_name).meta.service_model.shape_for(cls.shape_ls)
        return dict(cls._schema_rec(shp))

    @classmethod
    def schema_get(cls) -> Dict[str, str]:
        shp = boto3.client(cls.boto_client_name).meta.service_model.shape_for(cls.shape_get)
        return dict(cls._schema_rec(shp))

    @classmethod
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


class VpcHandler(LegacyAWSAPIHandler):
    boto_client_name = "ec2"
    shape_ls = "Vpc"
    shape_get = "Vpc"

    @staticmethod
    def resource_type() -> str:
        return "vpc"

    @classmethod
    async def _ls(cls, client) -> AsyncGenerator[Result, None]:
        response = await client.describe_vpcs()
        for item in response["Vpcs"]:
            yield Result(item["VpcId"], item)

    @classmethod
    async def _get(cls, client, resource_id: str) -> Result:
        response = await client.describe_vpcs(VpcIds=[resource_id])
        return Result(response["Vpcs"][0]["VpcId"], response["Vpcs"][0])


class SubnetHandler(LegacyAWSAPIHandler):
    boto_client_name = "ec2"
    shape_ls = "Subnet"
    shape_get = "Subnet"

    @staticmethod
    def resource_type() -> str:
        return "subnet"

    @classmethod
    async def _ls(cls, client) -> AsyncGenerator[Result, None]:
        response = await client.describe_subnets()
        for item in response["Subnets"]:
            yield Result(item["SubnetId"], item)

    @classmethod
    async def _get(cls, client, resource_id: str) -> Result:
        response = await client.describe_subnets(SubnetIds=[resource_id])
        return Result(response["Subnets"][0]["SubnetId"], response["Subnets"][0])


class RouteTableHandler(LegacyAWSAPIHandler):
    boto_client_name = "ec2"
    shape_ls = "RouteTable"
    shape_get = "RouteTable"

    @staticmethod
    def resource_type() -> str:
        return "rtb"

    @classmethod
    async def _ls(cls, client) -> AsyncGenerator[Result, None]:
        response = await client.describe_route_tables()
        for item in response["RouteTables"]:
            yield Result(item["RouteTableId"], item)

    @classmethod
    async def _get(cls, client, resource_id: str) -> Result:
        response = await client.describe_route_tables(RouteTableIds=[resource_id])
        return Result(response["RouteTables"][0]["RouteTableId"], response["RouteTables"][0])


class InternetGatewayHandler(LegacyAWSAPIHandler):
    boto_client_name = "ec2"
    shape_ls = "InternetGateway"
    shape_get = "InternetGateway"

    @staticmethod
    def resource_type() -> str:
        return "igw"

    @classmethod
    async def _ls(cls, client) -> AsyncGenerator[Result, None]:
        response = await client.describe_internet_gateways()
        for item in response["InternetGateways"]:
            yield Result(item["InternetGatewayId"], item)

    @classmethod
    async def _get(cls, client, resource_id: str) -> Result:
        response = await client.describe_internet_gateways(InternetGatewayIds=[resource_id])
        return Result(response["InternetGateways"][0]["InternetGatewayId"], response["InternetGateways"][0])


class SecurityGroupHandler(LegacyAWSAPIHandler):
    boto_client_name = "ec2"
    shape_ls = "SecurityGroup"
    shape_get = "SecurityGroup"

    @staticmethod
    def resource_type() -> str:
        return "sg"

    @classmethod
    async def _ls(cls, client) -> AsyncGenerator[Result, None]:
        response = await client.describe_security_groups()
        for item in response["SecurityGroups"]:
            yield Result(item["GroupId"], item)

    @classmethod
    async def _get(cls, client, resource_id: str) -> Result:
        response = await client.describe_security_groups(GroupIds=[resource_id])
        return Result(response["SecurityGroups"][0]["GroupId"], response["SecurityGroups"][0])


class NATGatewayHandler(LegacyAWSAPIHandler):
    boto_client_name = "ec2"
    shape_ls = "NatGateway"
    shape_get = "NatGateway"

    @staticmethod
    def resource_type() -> str:
        return "nat"

    @classmethod
    async def _ls(cls, client) -> AsyncGenerator[Result, None]:
        response = await client.describe_nat_gateways()
        for item in response["NatGateways"]:
            yield Result(item["NatGatewayId"], item)

    @classmethod
    async def _get(cls, client, resource_id: str) -> Result:
        response = await client.describe_nat_gateways(NatGatewayIds=[resource_id])
        return Result(response["NatGateways"][0]["NatGatewayId"], response["NatGateways"][0])


class NATGatewayHandler(LegacyAWSAPIHandler):
    boto_client_name = "ec2"
    shape_ls = "Address"
    shape_get = "Address"

    @staticmethod
    def resource_type() -> str:
        return "eip"

    @classmethod
    async def _ls(cls, client) -> AsyncGenerator[Result, None]:
        response = await client.describe_addresses()
        for item in response["Addresses"]:
            yield Result(item["AllocationId"], item)

    @classmethod
    async def _get(cls, client, resource_id: str) -> Result:
        response = await client.describe_addresses(AllocationIds=[resource_id])
        return Result(response["Addresses"][0]["AllocationId"], response["Addresses"][0])


class NATGatewayHandler(LegacyAWSAPIHandler):
    boto_client_name = "ec2"
    shape_ls = "NatGateway"
    shape_get = "NatGateway"

    @staticmethod
    def resource_type() -> str:
        return "nat"

    @classmethod
    async def _ls(cls, client) -> AsyncGenerator[Result, None]:
        response = await client.describe_addresses()
        for item in response["Addresses"]:
            yield Result(item["AllocationId"], item)

    @classmethod
    async def _get(cls, client, resource_id: str) -> Result:
        response = await client.describe_addresses(AllocationIds=[resource_id])
        return Result(response["Addresses"][0]["AllocationId"], response["Addresses"][0])


class NetworkInterfaceHandler(LegacyAWSAPIHandler):
    boto_client_name = "ec2"
    shape_ls = "NetworkInterface"
    shape_get = "NetworkInterface"

    @staticmethod
    def resource_type() -> str:
        return "eni"

    @classmethod
    async def _ls(cls, client) -> AsyncGenerator[Result, None]:
        response = await client.describe_network_interfaces()
        for item in response["NetworkInterfaces"]:
            yield Result(item["NetworkInterfaceId"], item)

    @classmethod
    async def _get(cls, client, resource_id: str) -> Result:
        response = await client.describe_network_interfaces(NetworkInterfaceIds=[resource_id])
        return Result(response["NetworkInterfaces"][0]["NetworkInterfaceId"], response["NetworkInterfaces"][0])


class NetworkAclHandler(LegacyAWSAPIHandler):
    boto_client_name = "ec2"
    shape_ls = "NetworkAcl"
    shape_get = "NetworkAcl"

    @staticmethod
    def resource_type() -> str:
        return "nacl"

    @classmethod
    async def _ls(cls, client) -> AsyncGenerator[Result, None]:
        response = await client.describe_network_acls()
        for item in response["NetworkAcls"]:
            yield Result(item["NetworkAclId"], item)

    @classmethod
    async def _get(cls, client, resource_id: str) -> Result:
        response = await client.describe_network_acls(NetworkAclIds=[resource_id])
        return Result(response["NetworkAcls"][0]["NetworkAclId"], response["NetworkAcls"][0])
