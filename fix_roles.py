import re

with open('/home/prem-modha/MPSProjects/UserManagement/languages/UserManagement/models/UserManagement.textGen.mps', 'r') as f:
    text = f.read()

# Fix WtQ9Q (ConceptTextGenDeclaration) to include the filename child (29tGrW)
# We need to find the <concept id="..." index="WtQ9Q"> block and inject the child.
wtq_match = re.search(r'(<concept[^>]+index="WtQ9Q">)', text)
if wtq_match:
    pos = wtq_match.end()
    if 'index="29tGrW"' not in text[pos:pos+500]: # Check if already in this concept block
        addition = '\n        <child id="1233920508237" name="filename" index="29tGrW" />'
        text = text[:pos] + addition + text[pos:]
        print('✅ Added 29tGrW to WtQ9Q')

# Let's check for any other used roles not registered
used_roles = set(re.findall(r'<node[^>]+role="([^"]+)"', text))
used_refs = set(re.findall(r'<ref[^>]+role="([^"]+)"', text))
registered_children = set(re.findall(r'<child[^>]+index="([^"]+)"', text))
registered_refs = set(re.findall(r'<reference[^>]+index="([^"]+)"', text))
registered_props = set(re.findall(r'<property[^>]+index="([^"]+)"', text))

# Properties used as roles (for simple values)
used_props = set(re.findall(r'<property[^>]+role="([^"]+)"', text))

all_used_roles = used_roles | used_refs | used_props
all_registered = registered_children | registered_refs | registered_props

missing = all_used_roles - all_registered
print(f'Missing roles before fix: {all_used_roles - all_registered}')

with open('/home/prem-modha/MPSProjects/UserManagement/languages/UserManagement/models/UserManagement.textGen.mps', 'w') as f:
    f.write(text)

