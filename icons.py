import os
import json

def build_map():
    for x, _, files in os.walk("frontend/public/icons"):
        for f in files:
            rel = os.path.relpath(x, "frontend/public").replace("\\", "/")
            yield f"/{rel}/{f}"


if __name__ == "__main__":
    # Format object into JS code
    js = "export const ICONS = new Set(" + json.dumps(list(build_map()), indent=2) + ");\n"

    with open("frontend/src/Icons.js", "w", encoding="utf8") as out:
        out.write("// GENERATED WITH icons.py\n")
        out.write("export default {};\n\n")
        out.write(js)
