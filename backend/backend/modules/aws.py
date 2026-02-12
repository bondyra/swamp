import configparser
import os
from typing import AsyncGenerator, Dict, List

import aioboto3
import boto3
from itertools import product

from backend.model import Attribute, Label, Provider, GenericQueryException
from backend.utils import get_matches


class AWS(Provider):
    @staticmethod
    def provider_name() -> str:
        return "aws"

    @staticmethod
    def provider_description() -> str:
        return "Provider for interacting with AWS resources via boto3"

    @staticmethod
    def resources() -> List[str]:
        return list(_resources.keys())

    @classmethod
    def description(cls, r: str) -> str:
        client, shape = _resources[r]["client"], _resources[r]["shape"]
        return boto3.client(client).meta.service_model.shape_for(shape).documentation
    
    @classmethod
    def icon(cls, r: str) -> str:
        return _resources[r].get("icon", "todo-needs-prefetch-cache-on-frontend")

    @classmethod
    async def get(cls, r: str, labels: Dict[str, Label]) -> AsyncGenerator[Dict, None]:
        if "_aws_profile" not in labels:
            raise GenericQueryException("You need to provide _aws_profile value to query AWS resource")
        if "_aws_region" not in labels:
            raise GenericQueryException("You need to provide _aws_region value to query AWS resource")
        profile_label, region_label = labels["_aws_profile"], labels["_aws_region"]
        client = _resources[r]["client"]
        profiles = get_matches(profile_label,  _get_profiles())
        regions = get_matches(region_label, _ALL_AWS_REGIONS)
        for profile, region in product(profiles, regions):
            async for x in cls._single_get(client, r, profile, region):
                yield x

    @classmethod
    async def _single_get(cls, client, r: str, profile: str, region: str) -> AsyncGenerator[Dict, None]:
        async with aioboto3.Session(profile_name=profile, region_name=region).client(client) as client:
            response = await _resources[r]["api_call"](client)
            for item in _resources[r]["iter_items"](response):
                yield {
                    **{"_id": _resources[r]["get_id"](item)},
                    "_aws_profile": profile,
                    "_aws_region": region,
                    **item
                }

    @classmethod
    async def attributes(cls, r: str) -> List[Attribute]:
        return [
            Attribute(path="_aws_profile", description="AWS profile to use", allowed_values=_get_profiles()),
            Attribute(path="_aws_region", description="AWS region", allowed_values=_ALL_AWS_REGIONS)
        ]

    @classmethod
    async def example(cls, r: str) -> Dict:
        client, shape = _resources[r]["client"],_resources[r]["shape"]
        shp = boto3.client(client).meta.service_model.shape_for(shape)
        return cls._example_rec(shp, f"{cls.provider_name()}.{r}")

    @classmethod
    def _example_rec(cls, obj, n):
        if "type_name" in dir(obj) and obj.type_name == "structure":
            return {
                name: cls._example_rec(member, name)
                for name, member in obj.members.items()
            }
        elif "type_name" in dir(obj) and obj.type_name == "list":
            el = cls._example_rec(obj.member, obj.name)
            return [el, el]
        elif "type_name" in dir(obj) and obj.type_name in {"string", "boolean", "integer"}:
            return f"{n}_VALUE"
        else:
            return None


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


_resources = {
    "vpc": {
        "client": "ec2",
        "shape": "Vpc",
        "api_call": lambda c: c.describe_vpcs(),
        "iter_items": lambda r: r["Vpcs"],
        "get_id": lambda i: i["VpcId"]
    },
    "instance": {
        "client": "ec2",
        "shape": "Instance",
        "api_call": lambda c: c.describe_instances(),
        "iter_items": lambda r: [inst for r2 in r.get("Reservations", []) for inst in r2.get("Instances")],
        "get_id": lambda i: i["InstanceId"]
    },
    "subnet": {
        "client": "ec2",
        "shape": "Subnet",
        "api_call": lambda c: c.describe_subnets(),
        "iter_items": lambda r: r["Subnets"],
        "get_id": lambda i: i["SubnetId"]
    },
    "route_table": {
        "client": "ec2",
        "shape": "RouteTable",
        "api_call": lambda c: c.describe_route_tables(),
        "iter_items": lambda r: r["RouteTables"],
        "get_id": lambda i: i["RouteTableId"]
    },
    "internet_gateway": {
        "client": "ec2",
        "shape": "InternetGateway",
        "api_call": lambda c: c.describe_internet_gateways(),
        "iter_items": lambda r: r["InternetGateways"],
        "get_id": lambda i: i["InternetGatewayId"]
    },
    "security_group": {
        "client": "ec2",
        "shape": "SecurityGroup",
        "api_call": lambda c: c.describe_security_groups(),
        "iter_items": lambda r: r["SecurityGroups"],
        "get_id": lambda i: i["GroupId"]
    },
    "nat_gateway": {
        "client": "ec2",
        "shape": "NatGateway",
        "api_call": lambda c: c.describe_nat_gateways(),
        "iter_items": lambda r: r["NatGateways"],
        "get_id": lambda i: i["NatGatewayId"]
    },
    "elastic_ip": {
        "client": "ec2",
        "shape": "Address",
        "api_call": lambda c: c.describe_addresses(),
        "iter_items": lambda r: r["Addresses"],
        "get_id": lambda i: i["AllocationId"]
    },
    "eni": {
        "client": "ec2",
        "shape": "NetworkInterface",
        "api_call": lambda c: c.describe_network_interfaces(),
        "iter_items": lambda r: r["NetworkInterfaces"],
        "get_id": lambda i: i["NetworkInterfaceId"]
    },
    "network_acl": {
        "client": "ec2",
        "shape": "NetworkAcl",
        "api_call": lambda c: c.describe_network_acls(),
        "iter_items": lambda r: r["NetworkAcls"],
        "get_id": lambda i: i["NetworkAclId"]
    },
    "ami": {
        "client": "ec2",
        "shape": "Image",
        "api_call": lambda c: c.describe_images(),
        "iter_items": lambda r: r["Images"],
        "get_id": lambda i: i["ImageId"],
    },
    "launch_template": {
        "client": "ec2",
        "shape": "PrefixList",
        "api_call": lambda c: c.describe_launch_templates(),
        "iter_items": lambda r: r["LaunchTemplates"],
        "get_id": lambda i: i["LaunchTemplateId"],
    },
    "prefix_list": {
        "client": "ec2",
        "shape": "",
        "api_call": lambda c: c.describe_prefix_lists(),
        "iter_items": lambda r: r["PrefixLists"],
        "get_id": lambda i: i["PrefixListId"],
    },
    "reserved_instances": {
        "client": "ec2",
        "shape": "ReservedInstances",
        "api_call": lambda c: c.describe_reserved_instances(),
        "iter_items": lambda r: r["ReservedInstances"],
        "get_id": lambda i: i["ReservedInstancesId"],
    },
    "snapshot": {
        "client": "ec2",
        "shape": "Snapshot",
        "api_call": lambda c: c.describe_snapshots(),
        "iter_items": lambda r: r["Snapshots"],
        "get_id": lambda i: i["SnapshotId"],
    },
    "tgw_attachment": {
        "client": "ec2",
        "shape": "TransitGatewayAttachment",
        "api_call": lambda c: c.describe_transit_gateway_attachments(),
        "iter_items": lambda r: r["TransitGatewayAttachments"],
        "get_id": lambda i: i["TransitGatewayAttachmentId"],
    },
    "tgw_route_table": {
        "client": "ec2",
        "shape": "TransitGatewayRouteTable",
        "api_call": lambda c: c.describe_transit_gateway_route_tables(),
        "iter_items": lambda r: r["TransitGatewayRouteTables"],
        "get_id": lambda i: i["TransitGatewayRouteTableId"],
    },
    "transit_gateway": {
        "client": "ec2",
        "shape": "TransitGateway",
        "api_call": lambda c: c.describe_transit_gateways(),
        "iter_items": lambda r: r["TransitGateways"],
        "get_id": lambda i: i["TransitGatewayId"],
    },
    "volume": {
        "client": "ec2",
        "shape": "Volume",
        "api_call": lambda c: c.describe_volumes(),
        "iter_items": lambda r: r["Volumes"],
        "get_id": lambda i: i["VolumeId"],
    },
    "vpc_endpoint": {
        "client": "ec2",
        "shape": "VpcEndpoint",
        "api_call": lambda c: c.describe_vpc_endpoints(),
        "iter_items": lambda r: r["VpcEndpoints"],
        "get_id": lambda i: i["VpcEndpointId"],
    },
    "vpc_peering_connection": {
        "client": "ec2",
        "shape": "VpcPeeringConnection",
        "api_call": lambda c: c.describe_vpc_peering_connections(),
        "iter_items": lambda r: r["VpcPeeringConnections"],
        "get_id": lambda i: i["VpcPeeringConnectionId"],
    },
    "dx_connection": {
        "client": "directconnect",
        "shape": "Connection",
        "api_call": lambda c: c.describe_connections(),
        "iter_items": lambda r: r["connections"],
        "get_id": lambda i: i["connectionId"]
    },
    "dx_lag": {
        "client": "directconnect",
        "shape": "Lag",
        "api_call": lambda c: c.describe_lags(),
        "iter_items": lambda r: r["lags"],
        "get_id": lambda i: i["lagId"]
    },
    "dx_virtual_gateway": {
        "client": "directconnect",
        "shape": "VirtualGateway",
        "api_call": lambda c: c.describe_virtual_gateways(),
        "iter_items": lambda r: r["virtualGateways"],
        "get_id": lambda i: i["virtualGatewayId"]
    },
    "dx_virtual_interface": {
        "client": "directconnect",
        "shape": "VirtualInterface",
        "api_call": lambda c: c.describe_virtual_interfaces(),
        "iter_items": lambda r: r["virtualInterfaces"],
        "get_id": lambda i: i["virtualInterfaceId"]
    },
    "resolver_endpoint": {
        "client": "route53resolver",
        "shape": "ResolverEndpoint",
        "api_call": lambda c: c.list_resolver_endpoints(),
        "iter_items": lambda r: r["ResolverEndpoints"],
        "get_id": lambda i: i["Id"]
    },
    "resolver_rule": {
        "client": "route53resolver",
        "shape": "ResolverRule",
        "api_call": lambda c: c.list_resolver_rules(),
        "iter_items": lambda r: r["ResolverRules"],
        "get_id": lambda i: i["Id"]
    },
    "autoscaling_group": {
        "client": "autoscaling",
        "shape": "AutoScalingGroup",
        "api_call": lambda c: c.describe_auto_scaling_groups(),
        "iter_items": lambda r: r["AutoScalingGroups"],
        "get_id": lambda i: i["AutoScalingGroupName"]
    },
    "load_balancer": {
        "client": "elbv2",
        "shape": "LoadBalancer",
        "api_call": lambda c: c.describe_load_balancers(),
        "iter_items": lambda r: r["LoadBalancers"],
        "get_id": lambda i: i["LoadBalancerArn"]
    },
    "target_group": {
        "client": "elbv2",
        "shape": "TargetGroup",
        "api_call": lambda c: c.describe_target_groups(),
        "iter_items": lambda r: r["TargetGroups"],
        "get_id": lambda i: i["TargetGroupArn"]
    },
}
