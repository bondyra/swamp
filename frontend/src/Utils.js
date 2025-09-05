export function randomString(length) {
    const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789";
    let result = "";
    for (let i = 0; i < length; i++) {
        result += chars.charAt(Math.floor(Math.random() * chars.length));
    }
    return result;
}


// todo move to backend (needs prefetched global cache)
export function getIconSrc(resourceType) {
    if(!resourceType)
        return "./icons/asset.svg"
    const [provider, resource] = resourceType.split(".")
    if(!provider)
        return "./icons/asset.svg"
    switch (provider) {
        case "aws":
            switch(resource){
                case "vpc":
                    return "https://icon.icepanel.io/AWS/svg/Networking-Content-Delivery/Virtual-Private-Cloud.svg"
                case "subnet":
                    return "https://icon.icepanel.io/AWS/svg/Networking-Content-Delivery/Virtual-Private-Cloud.svg"
                case "rtb":
                    return "https://icon.icepanel.io/AWS/svg/Networking-Content-Delivery/Virtual-Private-Cloud.svg"
                case "igw":
                    return "https://icon.icepanel.io/AWS/svg/Networking-Content-Delivery/Transit-Gateway.svg"
                case "sg":
                    return "https://icon.icepanel.io/AWS/svg/Security-Identity-Compliance/Network-Firewall.svg"
                case "nat":
                    return "https://icon.icepanel.io/AWS/svg/Networking-Content-Delivery/Transit-Gateway.svg"
                case "eip":
                    return "https://icon.icepanel.io/AWS/svg/Networking-Content-Delivery/Virtual-Private-Cloud.svg"
                case "eni":
                    return "https://icon.icepanel.io/AWS/svg/Networking-Content-Delivery/Virtual-Private-Cloud.svg"
                case "nacl":
                    return "https://icon.icepanel.io/AWS/svg/Security-Identity-Compliance/Network-Firewall.svg"
                default:
                    return "https://www.svgrepo.com/show/353443/aws.svg"
            }
        case "k8s":
            switch(resource){
                case "cm":
                    return "https://raw.githubusercontent.com/kubernetes/community/487f994c013ea61d92cf9a341af7620037abbce3/icons/svg/resources/unlabeled/cm.svg"
                case "ep":
                    return "https://raw.githubusercontent.com/kubernetes/community/19094aa6e60eb4a481650c4cbdb94badd9919b5b/icons/svg/resources/unlabeled/ep.svg"
                case "pod":
                    return "https://raw.githubusercontent.com/kubernetes/community/19094aa6e60eb4a481650c4cbdb94badd9919b5b/icons/svg/resources/unlabeled/pod.svg"
                case "pvc":
                    return "https://raw.githubusercontent.com/kubernetes/community/19094aa6e60eb4a481650c4cbdb94badd9919b5b/icons/svg/resources/unlabeled/pvc.svg"
                case "secret":
                    return "https://raw.githubusercontent.com/kubernetes/community/19094aa6e60eb4a481650c4cbdb94badd9919b5b/icons/svg/resources/unlabeled/secret.svg"
                case "sa":
                    return "https://raw.githubusercontent.com/kubernetes/community/19094aa6e60eb4a481650c4cbdb94badd9919b5b/icons/svg/resources/unlabeled/sa.svg"
                case "rs":
                    return "https://raw.githubusercontent.com/kubernetes/community/19094aa6e60eb4a481650c4cbdb94badd9919b5b/icons/svg/resources/unlabeled/rs.svg"
                case "deployment":
                    return "https://raw.githubusercontent.com/kubernetes/community/19094aa6e60eb4a481650c4cbdb94badd9919b5b/icons/svg/resources/unlabeled/deploy.svg"
                case "node":
                    return "https://raw.githubusercontent.com/kubernetes/community/refs/heads/master/icons/svg/infrastructure_components/unlabeled/node.svg"
                case "pv":
                    return "https://raw.githubusercontent.com/kubernetes/community/refs/heads/master/icons/svg/resources/unlabeled/pv.svg"
                default:
                    return "https://raw.githubusercontent.com/kubernetes/community/19094aa6e60eb4a481650c4cbdb94badd9919b5b/icons/svg/control_plane_components/labeled/api.svg"
            }
        default:
            return "./asset.svg"
    }
}
