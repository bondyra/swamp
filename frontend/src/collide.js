import { quadtree } from 'd3-quadtree';
 
export function collide() {
  let nodes = [];
  let force = (alpha) => {
    const tree = quadtree(
      nodes,
      (d) => d.x,
      (d) => d.y,
    );
 
    for (const node of nodes) {
      const w = node.measured.width / 2;
      const h = node.measured.height / 2;
      const nx1 = node.x - w;
      const nx2 = node.x + w;
      const ny1 = node.y - h;
      const ny2 = node.y + h;
 
      tree.visit((quad, x1, y1, x2, y2) => {
        if (!quad.length) {
          do {
            // quad.data is one of the nodes contained in this quad
            if (quad.data !== node) {
              const halfW_A = node.measured.width  / 2
              const halfH_A = node.measured.height / 2
              const halfW_B = quad.data.measured.width  / 2
              const halfH_B = quad.data.measured.height / 2

              // Vector from B to A
              const dx = node.x - quad.data.x
              const dy = node.y - quad.data.y

              // Overlap distances
              const overlapX = halfW_A + halfW_B - Math.abs(dx)
              const overlapY = halfH_A + halfH_B - Math.abs(dy)

              if (overlapX > 0 && overlapY > 0) {
                  // Resolve along the axis of least penetration
                  if (overlapX < overlapY){
                      // Push along X
                      const push = 100 * overlapX //* alpha
                      if (dx > 0) {
                          node.x += push
                          quad.data.x -= push
                      }
                      else {
                          node.x -= push
                          quad.data.x += push
                      }
                  }
                  else {
                    // Push along Y
                    const push = 100 * overlapY //* alpha
                    if (dy > 0){
                        node.y += push
                        quad.data.y -= push
                    }
                    else {
                        node.y -= push
                        quad.data.y += push
                    }
                  }
              }
            }
          } while ((quad = quad.next));
        }
        return x1 > nx2 || x2 < nx1 || y1 > ny2 || y2 < ny1; // stop signal - partition doesn't include the node
      });
    }
  };
 
  force.initialize = (newNodes) => (nodes = newNodes);
 
  return force;
}
 
export default collide;
