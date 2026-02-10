import { BaseEdge, EdgeLabelRenderer, getStraightPath } from '@xyflow/react';
 
export function SwampEdge({ id, selected, sourceX, sourceY, targetX, targetY, data}) {
  const [edgePath, labelX, labelY] = getStraightPath({
    sourceX,
    sourceY,
    targetX,
    targetY,
  });
 
  return (
    <>
      <BaseEdge id={id} path={edgePath} style={{ strokeWidth: 5 }}/>
      {
        selected &&
        <EdgeLabelRenderer>
            <div
            style={{
                position: "absolute",
                transform: `translate(-50%, -50%) translate(${labelX}px, ${labelY}px)`,
                background: "#1E1E1E",
                padding: "2px 6px",
                borderRadius: 4,
                fontSize: 8,
                fontFamily: "monospace",
                pointerEvents: "all",
            }}
            >
            {data?.label ?? "My label"}
            </div>
        </EdgeLabelRenderer>
      }
    </>
  );
}
