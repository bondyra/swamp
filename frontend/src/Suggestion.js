export const suggestions = {
    "aws.vpc": {
        "aws.subnet": {fromAttr: ".VpcId", op: "=", toAttr: ".VpcId"},
        "aws.route_table": {fromAttr: ".VpcId", op: "=", toAttr: ".VpcId"},
        "aws.network_acl": {fromAttr: ".Associations[].SubnetId", op: "contains", toAttr: ".SubnetId"},

    },
    "aws.subnet": {
        "aws.network_acl": {fromAttr: ".SubnetId", op: "contains", toAttr: ".Associations[].SubnetId"}
    },
    "k8s.deployment": {
        "k8s.replica_set": {fromAttr: ".metadata.name", op: "contains", toAttr: `.metadata.ownerReferences[] | select(.kind="Deployment") | .name`}
    },
    "k8s.replica_set": {
        "k8s.pod": {fromAttr: ".metadata.name", op: "contains", toAttr: `.metadata.ownerReferences[] | select(.kind="ReplicaSet") | .name`}
    },
}
