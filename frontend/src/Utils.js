import {ICONS} from './Icons';

export function randomString(length) {
    const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789";
    let result = "";
    for (let i = 0; i < length; i++) {
        result += chars.charAt(Math.floor(Math.random() * chars.length));
    }
    return result;
}

export function getIconSrc(resourceType) {
    if(!resourceType)
        return "/icons/_DEFAULT_.svg"
    const [provider, resource] = resourceType.split(".")
    const iconPath = `/icons/${provider}/${resource}.svg`
    if (ICONS.has(iconPath))
        return iconPath
    const providerDefault = `/icons/${provider}/_DEFAULT_.svg`
    if (ICONS.has(providerDefault))
        return providerDefault
    return "/icons/_DEFAULT_.svg"
}
