


export function getAllJSONPaths(obj, prefix = '') {
	let results = [];
	
	if (typeof obj === 'object' && obj !== null) {
		for (let key in obj) {
		  if (Array.isArray(obj[key])) {
            for(var i = 0; i < obj[key].length; i++) {
                let newPrefix = prefix === '' ? `${key}[${i}]` : `${prefix}.${key}[${i}]`;
                results = results.concat(getAllJSONPaths(obj[key][i], newPrefix));
            }
		  }
		  else if (typeof obj[key] === 'object'){
            let newPrefix = prefix === '' ? key : `${prefix}.${key}`;
            results = results.concat(getAllJSONPaths(obj[key], newPrefix));
		  } else {
            let newPrefix = prefix === '' ? key : `${prefix}.${key}`;
            results.push(newPrefix);
		  }
		}
	}
	return results;
}

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
                default:
                    return "https://www.svgrepo.com/show/353443/aws.svg"
            }
        default:
            return "./icons/asset.svg"
    }
}