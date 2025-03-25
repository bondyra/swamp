import { useReactFlow } from '@xyflow/react';


export default class StateManager {
    constructor() {
        this.reactFlow = useReactFlow();
    }

    createRandomString(length) {
        const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789";
        let result = "";
        for (let i = 0; i < length; i++) {
            result += chars.charAt(Math.floor(Math.random() * chars.length));
        }
        return result;
    }
  
    onQueryResult(result){
        var newNodes = [];
        var newEdges = [];
        result.forEach(item => {
            const newNodeId = `${item.resourceType}.${item.result._id}`;
            newNodes.push({
                id: newNodeId,
                position: { x: 0, y: 0 },
                type: 'resource',
                data: {
                    resourceType: item.resourceType,
                    result: item.result
                },
            });
            newEdges.push({id: `${id}-${newNodeId}`, source: id, target: newNodeId, style: {strokeWidth: 5} });
        });// todo: check for duplicates
        this.reactFlow.addNodes(newNodes);
        this.reactFlow.addEdges(newEdges);
    }

    modifyNode(nodeId, {resourceData, nestedData, nestedIndex}){
        this.reactFlow.updateNodeData(nodeId, (node) => {
            if (nestedData){
                assert(nestedIndex)
                return {
                    ...node.data,
                    nested: node.data.nested.map((n, index)=>{
                        if(index == nestedIndex)
                            return nestedData
                        return n
                    })
                }
            }
            assert(resourceData)
            return {
              ...node.data,
              ...resourceData
            };
          })
    }

    modifyQuery(nodeId, data){
        this.reactFlow.updateNodeData(nodeId, (node) => {
            return {...node.data, ...data};
        })
    }

    initialNodes(){
        return [
            {
              id: this.createRandomString(),
              position: { x: 0, y: 0 },
              type: 'query',
              data: {labels: []},
            }
        ]
    }

    initialEdges(){
        return []
    }

    newQuery(sourceNode){
        const newNodeId = this.createRandomString()
        const newNode = {
            id: newNodeId,
            position: { x: 0, y: 0 },
            type: 'query',
            data: {resourceType: null, labels: [], parentResourceType: sourceNode.data.resourceType,  parent: sourceNode.data.result},
            origin: [0.5, 0.0],
        };
  
      reactFlow.setNodes((nds) => nds.concat(newNode));
      reactFlow.setEdges((eds) =>
        eds.concat({ id: `${sourceNode.id}-${newNodeId}`, source: sourceNode.id, target: newNodeId, style: {strokeWidth: 5} }),
      );
    }

    deleteQuery(nodeId, nodeData){
        assert(!nodeData.hasRun)
        const newEdges = this.reactFlow.getEdges().filter(e=> e.target !== nodeId)
        const newNodes = this.reactFlow.getNodes().filter(n=> n.id !== nodeId)
        this.reactFlow.setEdges(newEdges)
        this.reactFlow.setNodes(newNodes)
    }
}
