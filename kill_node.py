import re

with open('/home/prem-modha/MPSProjects/UserManagement/languages/UserManagement/models/UserManagement.textGen.mps', 'r') as f:
    text = f.read()

start_match = re.search(r'<node concept="WtQ9Q" id="6DJmAW\$krO0"', text)
if not start_match:
    print("Node not found!")
else:
    start = start_match.start()
    pos = start
    depth = 0
    while pos < len(text):
        m = re.search(r'<(/?)node\b([^>]*)>', text[pos:])
        if not m: break
        
        tag_pos = pos + m.start()
        is_close = m.group(1) == '/'
        is_self = m.group(2).rstrip().endswith('/')
        
        if is_close:
            depth -= 1
            if depth == 0:
                end = tag_pos + len(m.group(0))
                # Delete this block from text
                text = text[:start] + text[end:]
                print("✅ Successfully deleted node 6DJmAW$krO0")
                break
            pos = tag_pos + len(m.group(0))
        elif is_self:
            if depth == 0:
                end = tag_pos + len(m.group(0))
                text = text[:start] + text[end:]
                break
            pos = tag_pos + len(m.group(0))
        else:
            depth += 1
            pos = tag_pos + len(m.group(0))

    with open('/home/prem-modha/MPSProjects/UserManagement/languages/UserManagement/models/UserManagement.textGen.mps', 'w') as f:
        f.write(text)

import xml.etree.ElementTree as ET
root = ET.parse('/home/prem-modha/MPSProjects/UserManagement/languages/UserManagement/models/UserManagement.textGen.mps').getroot()
for n in root.findall('.//node[@concept="WtQ9Q"]'):
    ref = n.find('ref[@role="WuzLi"]')
    print(f'Remaining: {ref.attrib.get("resolve")}')
