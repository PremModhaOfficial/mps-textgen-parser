import xml.etree.ElementTree as ET
import random
import string
import time

def generate_mps_id():
    # MPS IDs are typically base62 strings.
    # Let's use timestamp + random to ensure global uniqueness.
    chars = string.ascii_letters + string.digits
    rand_part = ''.join(random.choice(chars) for _ in range(6))
    return f'1{int(time.time()*1000) % 1000000}{rand_part}'

tree = ET.parse('/home/prem-modha/MPSProjects/UserManagement/languages/UserManagement/models/UserManagement.textGen.mps')
root = tree.getroot()

# Build a mapping of old_id -> new_id for the duplicates
# But wait, within a single WtQ9Q tree, internal references use these IDs!
# (e.g. ForEachVariable has an ID, and ForEachVariableReference uses node="THAT_ID")
# So we must remap consistently WITHIN each component, but uniquely across components.

# Actually, the safest way is:
# We know the first component is the original (it has original IDs).
# For all SUBSEQUENT components (the 4 we appended), we should rewrite ALL their IDs 
# and update any internal references.

components = root.findall('.//{*}node[@concept="WtQ9Q"]')
print(f'Found {len(components)} components.')

# The first one is original, leave it alone.
for comp in components[1:]:
    id_map = {}
    
    # 1. First pass: assign new ID to every node in this component
    # Also include the component root node itself
    for n in comp.iter('node'):
        if 'id' in n.attrib:
            old_id = n.attrib['id']
            new_id = generate_mps_id()
            id_map[old_id] = new_id
            n.attrib['id'] = new_id
            
    # 2. Second pass: update all internal references that point to these old IDs
    for ref in comp.iter('ref'):
        # Some refs use 'node="id"' format for internal variables
        if 'node' in ref.attrib:
            old_ref = ref.attrib['node']
            if old_ref in id_map:
                ref.attrib['node'] = id_map[old_ref]

# Save the fixed file
tree.write('/home/prem-modha/MPSProjects/UserManagement/languages/UserManagement/models/UserManagement.textGen.mps', encoding='utf-8', xml_declaration=True)

# Verify uniqueness now
ids = []
for n in root.iter('node'):
    if 'id' in n.attrib:
        ids.append(n.attrib['id'])
print(f'Total node IDs: {len(ids)}, Unique: {len(set(ids))}')
if len(ids) == len(set(ids)):
    print('✅ ALL IDs ARE NOW UNIQUE!')
